package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

// RegisterStandardCosmosRoutes registers standard Cosmos SDK REST API endpoints
// that proxy to working keeper-based queries, bypassing the broken IAVL
// versioned store in SDK v0.50.x. This makes the chain compatible with Keplr,
// Mintscan, and other standard Cosmos ecosystem tools.
func (app *SyreenApp) RegisterStandardCosmosRoutes(router *mux.Router) {
	// Helper: create SDK context from committed state
	newCtx := func() sdk.Context {
		cms := app.CommitMultiStore().CacheMultiStore()
		return sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
	}

	// Helper: write JSON response
	writeJSON := func(w http.ResponseWriter, v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}

	// Helper: write error response (matches Cosmos SDK error format)
	writeErr := func(w http.ResponseWriter, code int, msg string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    code,
			"message": msg,
		})
	}

	// Helper: empty pagination (single page)
	paginationResp := func(total int) map[string]interface{} {
		return map[string]interface{}{
			"next_key": nil,
			"total":    fmt.Sprintf("%d", total),
		}
	}

	// ─── Bank: Balances ───────────────────────────────────────────────
	// GET /cosmos/bank/v1beta1/balances/{address}
	router.HandleFunc("/cosmos/bank/v1beta1/balances/{address}", func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["address"]
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			writeErr(w, 400, "invalid address: "+err.Error())
			return
		}

		ctx := newCtx()
		coins := app.BankKeeper.SpendableCoins(ctx, accAddr)

		balances := make([]map[string]string, 0, len(coins))
		for _, c := range coins {
			balances = append(balances, map[string]string{
				"denom":  c.Denom,
				"amount": c.Amount.String(),
			})
		}
		if len(balances) == 0 {
			balances = []map[string]string{}
		}

		writeJSON(w, map[string]interface{}{
			"balances":   balances,
			"pagination": paginationResp(len(balances)),
		})
	}).Methods("GET")

	// ─── Bank: Total Supply ───────────────────────────────────────────
	// GET /cosmos/bank/v1beta1/supply
	router.HandleFunc("/cosmos/bank/v1beta1/supply", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()

		var supply []map[string]string
		app.BankKeeper.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
			supply = append(supply, map[string]string{
				"denom":  coin.Denom,
				"amount": coin.Amount.String(),
			})
			return false
		})
		if supply == nil {
			supply = []map[string]string{}
		}

		writeJSON(w, map[string]interface{}{
			"supply":     supply,
			"pagination": paginationResp(len(supply)),
		})
	}).Methods("GET")

	// ─── Bank: Supply by Denom ────────────────────────────────────────
	// GET /cosmos/bank/v1beta1/supply/by_denom?denom=usyreen
	router.HandleFunc("/cosmos/bank/v1beta1/supply/by_denom", func(w http.ResponseWriter, r *http.Request) {
		denom := r.URL.Query().Get("denom")
		if denom == "" {
			writeErr(w, 400, "denom query parameter required")
			return
		}

		ctx := newCtx()
		coin := app.BankKeeper.GetSupply(ctx, denom)

		writeJSON(w, map[string]interface{}{
			"amount": map[string]string{
				"denom":  coin.Denom,
				"amount": coin.Amount.String(),
			},
		})
	}).Methods("GET")

	// ─── Auth: Account ────────────────────────────────────────────────
	// GET /cosmos/auth/v1beta1/accounts/{address}
	router.HandleFunc("/cosmos/auth/v1beta1/accounts/{address}", func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["address"]
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			writeErr(w, 400, "invalid address: "+err.Error())
			return
		}

		ctx := newCtx()
		acc := app.AccountKeeper.GetAccount(ctx, accAddr)
		if acc == nil {
			writeErr(w, 404, "account not found")
			return
		}

		// Use codec to marshal with proper @type and Any handling
		baseAcct, ok := acc.(*authtypes.BaseAccount)
		if !ok {
			writeErr(w, 400, "account is not a BaseAccount")
			return
		}
		any, err := codectypes.NewAnyWithValue(baseAcct)
		if err != nil {
			writeErr(w, 500, "failed to pack account: "+err.Error())
			return
		}

		accJSON, err := app.appCodec.MarshalJSON(any)
		if err != nil {
			writeErr(w, 500, "failed to marshal account: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"account":`))
		w.Write(accJSON)
		w.Write([]byte(`}`))
	}).Methods("GET")

	// ─── Staking: All Validators ──────────────────────────────────────
	// GET /cosmos/staking/v1beta1/validators?status=BOND_STATUS_BONDED
	router.HandleFunc("/cosmos/staking/v1beta1/validators", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		status := r.URL.Query().Get("status")

		var validators []stakingtypes.Validator
		var err error

		if status == "BOND_STATUS_BONDED" {
			validators, err = app.StakingKeeper.GetBondedValidatorsByPower(ctx)
		} else {
			validators, err = app.StakingKeeper.GetAllValidators(ctx)
		}
		if err != nil {
			writeErr(w, 500, "failed to get validators: "+err.Error())
			return
		}

		// Filter by status if specified (for non-bonded statuses)
		if status != "" && status != "BOND_STATUS_BONDED" {
			var filtered []stakingtypes.Validator
			for _, val := range validators {
				if val.Status.String() == status {
					filtered = append(filtered, val)
				}
			}
			validators = filtered
		}

		// Marshal each validator with codec for proper proto JSON
		valJSONs := make([]json.RawMessage, 0, len(validators))
		for i := range validators {
			bz, err := app.appCodec.MarshalJSON(&validators[i])
			if err != nil {
				continue
			}
			valJSONs = append(valJSONs, json.RawMessage(bz))
		}

		writeJSON(w, map[string]interface{}{
			"validators": valJSONs,
			"pagination": paginationResp(len(valJSONs)),
		})
	}).Methods("GET")

	// ─── Staking: Single Validator ────────────────────────────────────
	// GET /cosmos/staking/v1beta1/validators/{validator_addr}
	router.HandleFunc("/cosmos/staking/v1beta1/validators/{validator_addr}", func(w http.ResponseWriter, r *http.Request) {
		valAddrStr := mux.Vars(r)["validator_addr"]
		valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			writeErr(w, 400, "invalid validator address: "+err.Error())
			return
		}

		ctx := newCtx()
		val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			writeErr(w, 404, "validator not found")
			return
		}

		valJSON, err := app.appCodec.MarshalJSON(&val)
		if err != nil {
			writeErr(w, 500, "failed to marshal validator: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"validator":`))
		w.Write(valJSON)
		w.Write([]byte(`}`))
	}).Methods("GET")

	// ─── Staking: Delegations ─────────────────────────────────────────
	// GET /cosmos/staking/v1beta1/delegations/{delegator_addr}
	router.HandleFunc("/cosmos/staking/v1beta1/delegations/{delegator_addr}", func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["delegator_addr"]
		delAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			writeErr(w, 400, "invalid delegator address: "+err.Error())
			return
		}

		ctx := newCtx()
		delegations, err := app.StakingKeeper.GetDelegatorDelegations(ctx, delAddr, 200)
		if err != nil {
			writeErr(w, 500, "failed to get delegations: "+err.Error())
			return
		}

		responses := make([]map[string]interface{}, 0, len(delegations))
		for _, del := range delegations {
			valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
			if err != nil {
				continue
			}

			val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
			if err != nil {
				continue
			}

			tokens := val.TokensFromSharesTruncated(del.Shares).TruncateInt()

			responses = append(responses, map[string]interface{}{
				"delegation": map[string]string{
					"delegator_address": del.DelegatorAddress,
					"validator_address": del.ValidatorAddress,
					"shares":            del.Shares.String(),
				},
				"balance": map[string]string{
					"denom":  "usyreen",
					"amount": tokens.String(),
				},
			})
		}

		writeJSON(w, map[string]interface{}{
			"delegation_responses": responses,
			"pagination":           paginationResp(len(responses)),
		})
	}).Methods("GET")

	// ─── Staking: Validator Delegations ───────────────────────────────
	// GET /cosmos/staking/v1beta1/validators/{validator_addr}/delegations
	router.HandleFunc("/cosmos/staking/v1beta1/validators/{validator_addr}/delegations", func(w http.ResponseWriter, r *http.Request) {
		valAddrStr := mux.Vars(r)["validator_addr"]
		valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
		if err != nil {
			writeErr(w, 400, "invalid validator address: "+err.Error())
			return
		}

		ctx := newCtx()
		val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			writeErr(w, 404, "validator not found")
			return
		}

		delegations, err := app.StakingKeeper.GetValidatorDelegations(ctx, valAddr)
		if err != nil {
			writeErr(w, 500, "failed to get delegations: "+err.Error())
			return
		}

		responses := make([]map[string]interface{}, 0, len(delegations))
		for _, del := range delegations {
			tokens := val.TokensFromSharesTruncated(del.Shares).TruncateInt()
			responses = append(responses, map[string]interface{}{
				"delegation": map[string]string{
					"delegator_address": del.DelegatorAddress,
					"validator_address": del.ValidatorAddress,
					"shares":            del.Shares.String(),
				},
				"balance": map[string]string{
					"denom":  "usyreen",
					"amount": tokens.String(),
				},
			})
		}

		writeJSON(w, map[string]interface{}{
			"delegation_responses": responses,
			"pagination":           paginationResp(len(responses)),
		})
	}).Methods("GET")

	// ─── Staking: Params ──────────────────────────────────────────────
	// GET /cosmos/staking/v1beta1/params
	router.HandleFunc("/cosmos/staking/v1beta1/params", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		params, err := app.StakingKeeper.GetParams(ctx)
		if err != nil {
			writeErr(w, 500, "failed to get staking params: "+err.Error())
			return
		}

		paramsJSON, err := app.appCodec.MarshalJSON(&params)
		if err != nil {
			writeErr(w, 500, "failed to marshal params: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"params":`))
		w.Write(paramsJSON)
		w.Write([]byte(`}`))
	}).Methods("GET")

	// ─── Staking: Pool ────────────────────────────────────────────────
	// GET /cosmos/staking/v1beta1/pool
	router.HandleFunc("/cosmos/staking/v1beta1/pool", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()

		bondedPoolAddr := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
		notBondedPoolAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

		bondedBal := app.BankKeeper.GetBalance(ctx, bondedPoolAddr, "usyreen")
		notBondedBal := app.BankKeeper.GetBalance(ctx, notBondedPoolAddr, "usyreen")

		writeJSON(w, map[string]interface{}{
			"pool": map[string]string{
				"bonded_tokens":     bondedBal.Amount.String(),
				"not_bonded_tokens": notBondedBal.Amount.String(),
			},
		})
	}).Methods("GET")

	// ─── Staking: Unbonding Delegations ───────────────────────────────
	// GET /cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations
	router.HandleFunc("/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations", func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["delegator_addr"]
		delAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			writeErr(w, 400, "invalid delegator address: "+err.Error())
			return
		}

		ctx := newCtx()
		unbondings, err := app.StakingKeeper.GetUnbondingDelegations(ctx, delAddr, 200)
		if err != nil {
			writeErr(w, 500, "failed to get unbonding delegations: "+err.Error())
			return
		}

		ubJSONs := make([]json.RawMessage, 0, len(unbondings))
		for i := range unbondings {
			bz, err := app.appCodec.MarshalJSON(&unbondings[i])
			if err != nil {
				continue
			}
			ubJSONs = append(ubJSONs, json.RawMessage(bz))
		}

		writeJSON(w, map[string]interface{}{
			"unbonding_responses": ubJSONs,
			"pagination":          paginationResp(len(ubJSONs)),
		})
	}).Methods("GET")

	// ─── Distribution: Delegation Rewards ─────────────────────────────
	// GET /cosmos/distribution/v1beta1/delegators/{delegator_addr}/rewards
	router.HandleFunc("/cosmos/distribution/v1beta1/delegators/{delegator_addr}/rewards", func(w http.ResponseWriter, r *http.Request) {
		addr := mux.Vars(r)["delegator_addr"]
		delAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			writeErr(w, 400, "invalid delegator address: "+err.Error())
			return
		}

		ctx := newCtx()

		// Get all delegations for this delegator
		delegations, err := app.StakingKeeper.GetDelegatorDelegations(ctx, delAddr, 200)
		if err != nil {
			writeErr(w, 500, "failed to get delegations: "+err.Error())
			return
		}

		type rewardEntry struct {
			ValidatorAddress string              `json:"validator_address"`
			Reward           []map[string]string `json:"reward"`
		}

		var rewards []rewardEntry
		totalRewards := sdk.DecCoins{}

		for _, del := range delegations {
			valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
			if err != nil {
				continue
			}

			val, err := app.StakingKeeper.GetValidator(ctx, valAddr)
			if err != nil {
				continue
			}

			endingPeriod, err := app.DistrKeeper.IncrementValidatorPeriod(ctx, val)
			if err != nil {
				continue
			}

			delRewards, err := app.DistrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
			if err != nil {
				continue
			}

			rewardCoins := make([]map[string]string, 0, len(delRewards))
			for _, r := range delRewards {
				rewardCoins = append(rewardCoins, map[string]string{
					"denom":  r.Denom,
					"amount": r.Amount.String(),
				})
			}
			if len(rewardCoins) == 0 {
				rewardCoins = []map[string]string{}
			}

			rewards = append(rewards, rewardEntry{
				ValidatorAddress: del.ValidatorAddress,
				Reward:           rewardCoins,
			})

			totalRewards = totalRewards.Add(delRewards...)
		}

		totalCoins := make([]map[string]string, 0, len(totalRewards))
		for _, t := range totalRewards {
			totalCoins = append(totalCoins, map[string]string{
				"denom":  t.Denom,
				"amount": t.Amount.String(),
			})
		}
		if len(totalCoins) == 0 {
			totalCoins = []map[string]string{}
		}
		if rewards == nil {
			rewards = []rewardEntry{}
		}

		writeJSON(w, map[string]interface{}{
			"rewards": rewards,
			"total":   totalCoins,
		})
	}).Methods("GET")

	// ─── Mint: Params ─────────────────────────────────────────────────
	// GET /cosmos/mint/v1beta1/params
	router.HandleFunc("/cosmos/mint/v1beta1/params", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		params, err := app.MintKeeper.Params.Get(ctx)
		if err != nil {
			writeErr(w, 500, "failed to get mint params: "+err.Error())
			return
		}

		paramsJSON, err := app.appCodec.MarshalJSON(&params)
		if err != nil {
			writeErr(w, 500, "failed to marshal params: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"params":`))
		w.Write(paramsJSON)
		w.Write([]byte(`}`))
	}).Methods("GET")

	// ─── Mint: Inflation ──────────────────────────────────────────────
	// GET /cosmos/mint/v1beta1/inflation
	router.HandleFunc("/cosmos/mint/v1beta1/inflation", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		minter, err := app.MintKeeper.Minter.Get(ctx)
		if err != nil {
			writeErr(w, 500, "failed to get minter: "+err.Error())
			return
		}

		writeJSON(w, map[string]string{
			"inflation": minter.Inflation.String(),
		})
	}).Methods("GET")

	// ─── Node Info ────────────────────────────────────────────────────
	// GET /cosmos/base/tendermint/v1beta1/node_info
	router.HandleFunc("/cosmos/base/tendermint/v1beta1/node_info", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]interface{}{
			"default_node_info": map[string]interface{}{
				"protocol_version": map[string]string{
					"p2p":   "8",
					"block": "11",
					"app":   "0",
				},
				"network": "syreen-1",
				"version": "0.38.12",
				"moniker": "syreen-validator",
				"other": map[string]string{
					"tx_index":    "on",
					// M6: Do not leak the internal CometBFT bind address.
					// The public RPC is served via nginx; redact the
					// internal listener so node_info doesn't advertise
					// 0.0.0.0:26657 to the world.
					"rpc_address": "tcp://localhost:26657",
				},
			},
			"application_version": map[string]interface{}{
				"name":              "SyreenChain",
				"app_name":          "syreend",
				"version":           "1.0.0",
				"cosmos_sdk_version": "v0.50.11",
				"go_version":        "go1.21",
			},
		})
	}).Methods("GET")

	// ─── Gov: Params ──────────────────────────────────────────────────
	// GET /cosmos/gov/v1beta1/params/{params_type}
	router.HandleFunc("/cosmos/gov/v1beta1/params/{params_type}", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		govParams, err := app.GovKeeper.Params.Get(ctx)
		if err != nil {
			writeErr(w, 500, "failed to get gov params: "+err.Error())
			return
		}

		paramsJSON, err := app.appCodec.MarshalJSON(&govParams)
		if err != nil {
			writeErr(w, 500, "failed to marshal params: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(paramsJSON)
	}).Methods("GET")

	// ─── Slashing: Params ─────────────────────────────────────────────
	// GET /cosmos/slashing/v1beta1/params
	router.HandleFunc("/cosmos/slashing/v1beta1/params", func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtx()
		params, err := app.SlashingKeeper.GetParams(ctx)
		if err != nil {
			writeErr(w, 500, "failed to get slashing params: "+err.Error())
			return
		}

		paramsJSON, err := app.appCodec.MarshalJSON(&params)
		if err != nil {
			writeErr(w, 500, "failed to marshal params: "+err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"params":`))
		w.Write(paramsJSON)
		w.Write([]byte(`}`))
	}).Methods("GET")
}
