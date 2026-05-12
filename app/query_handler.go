package app

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	intenttypes "syreen/x/intent/types"
)

// RegisterCustomQueryRoutes adds direct-state query endpoints that bypass
// the broken CacheMultiStoreWithVersion in SDK v0.50.x IAVL integration.
// These read from the current committed state without version loading.
func (app *SyreenApp) RegisterCustomQueryRoutes(router *mux.Router, _ client.Context) {
	router.HandleFunc("/syreen/v1/balances/{address}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]

		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Create context from committed multistore (no version loading needed)
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		coins := app.BankKeeper.SpendableCoins(ctx, accAddr)

		type coinJSON struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		}
		result := make([]coinJSON, 0, len(coins))
		for _, c := range coins {
			result = append(result, coinJSON{Denom: c.Denom, Amount: c.Amount.String()})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"balances": result})
	}).Methods("GET")

	router.HandleFunc("/syreen/v1/supply", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		supply := app.BankKeeper.GetSupply(ctx, "usyreen")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"denom":  supply.Denom,
			"amount": supply.Amount.String(),
		})
	}).Methods("GET")

	// --- Custom Module Param Endpoints ---

	router.HandleFunc("/syreen/feemarket/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.FeeMarketKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/tokenfactory/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.TokenfactoryKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/mevprotection/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.MEVProtectionKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.IntentKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/abstractaccount/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.AbstractAccountKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/compute/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.ComputeKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/dex/v1/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.DexKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	router.HandleFunc("/syreen/dex/v1/pools", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		pools := app.DexKeeper.GetAllPools(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"pools": pools})
	}).Methods("GET")

	router.HandleFunc("/syreen/v1/evm/params", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		params := app.EVMKeeper.GetParams(ctx)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"params": params})
	}).Methods("GET")

	// --- Trading Intent Endpoints ---

	router.HandleFunc("/syreen/dex/v1/pool/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolIDStr := vars["pool_id"]

		poolID, parseErr := strconv.ParseUint(poolIDStr, 10, 64)
		if parseErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		pool, found := app.DexKeeper.GetPool(ctx, poolID)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "pool not found"})
			return
		}

		// Calculate spot prices
		priceAB := "0"
		priceBA := "0"
		if !pool.ReserveA.IsZero() {
			p, _ := app.DexKeeper.GetSpotPrice(ctx, poolID, pool.DenomA, pool.DenomB)
			priceAB = p.String()
		}
		if !pool.ReserveB.IsZero() {
			p, _ := app.DexKeeper.GetSpotPrice(ctx, poolID, pool.DenomB, pool.DenomA)
			priceBA = p.String()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pool":     pool,
			"price_ab": priceAB,
			"price_ba": priceBA,
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/dex/v1/spot_price/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolIDStr := vars["pool_id"]

		poolID, parseErr := strconv.ParseUint(poolIDStr, 10, 64)
		if parseErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		pool, found := app.DexKeeper.GetPool(ctx, poolID)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "pool not found"})
			return
		}

		priceAB := "0"
		priceBA := "0"
		if !pool.ReserveA.IsZero() {
			p, _ := app.DexKeeper.GetSpotPrice(ctx, poolID, pool.DenomA, pool.DenomB)
			priceAB = p.String()
		}
		if !pool.ReserveB.IsZero() {
			p, _ := app.DexKeeper.GetSpotPrice(ctx, poolID, pool.DenomB, pool.DenomA)
			priceBA = p.String()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pool_id":  poolID,
			"denom_a":  pool.DenomA,
			"denom_b":  pool.DenomB,
			"price_ab": priceAB,
			"price_ba": priceBA,
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/intents", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		intents := app.IntentKeeper.GetAllIntents(ctx)

		w.Header().Set("Content-Type", "application/json")
		if intents == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"intents": []interface{}{},
				"count":   0,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"intents": intents,
				"count":   len(intents),
			})
		}
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/intents/{type}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		intentType := vars["type"]

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		intents := app.IntentKeeper.GetIntentsByType(ctx, intentType)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"intents": intents,
			"count":   len(intents),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/intent/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		intentID := vars["id"]

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		intent, found := app.IntentKeeper.GetIntent(ctx, intentID)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "intent not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"intent": intent})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/pending", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		intents := app.IntentKeeper.GetPendingIntents(ctx)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"intents": intents,
			"count":   len(intents),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/dex/v1/risk/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, parseErr := strconv.ParseUint(vars["pool_id"], 10, 64)
		if parseErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		halted, reason := app.DexKeeper.IsPoolHalted(ctx, poolID)
		riskState, found := app.DexKeeper.GetRiskState(ctx, poolID)

		result := map[string]interface{}{
			"pool_id": poolID,
			"halted":  halted,
			"reason":  reason,
			"found":   found,
		}
		if found {
			result["price_change_pct"] = riskState.PriceChangePct.String()
			result["snapshot_height"] = riskState.SnapshotHeight
			result["halted_until"] = riskState.HaltedUntil
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/orders/{address}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		intents := app.IntentKeeper.GetIntentsByCreator(ctx, address)

		// Separate into active vs completed
		var active, completed []interface{}
		for _, intent := range intents {
			if intent.Status == "pending" || intent.Status == "solving" {
				active = append(active, intent)
			} else {
				completed = append(completed, intent)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active":    active,
			"completed": completed,
			"total":     len(intents),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/dex/v1/quote/{pool_id}/{denom}/{amount}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		poolID, parseErr := strconv.ParseUint(vars["pool_id"], 10, 64)
		if parseErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}

		denom := vars["denom"]
		amountStr := vars["amount"]
		amount, ok := math.NewIntFromString(amountStr)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid amount"})
			return
		}

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		tokenIn := sdk.NewCoin(denom, amount)
		tokenOut, priceImpact, err := app.DexKeeper.GetQuote(ctx, poolID, tokenIn)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token_in":     tokenIn.String(),
			"token_out":    tokenOut.String(),
			"out_amount":   tokenOut.Amount.String(),
			"out_denom":    tokenOut.Denom,
			"price_impact": priceImpact.String(),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/solvers", func(w http.ResponseWriter, r *http.Request) {
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		solvers := app.IntentKeeper.GetAllSolvers(ctx)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"solvers": solvers,
			"count":   len(solvers),
		})
	}).Methods("GET")

	// --- EVM Endpoints ---

	// Get EVM contract code at an address
	router.HandleFunc("/syreen/v1/evm/code/{address}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addrHex := vars["address"]

		if !common.IsHexAddress(addrHex) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid hex address"})
			return
		}

		addr := common.HexToAddress(addrHex)
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		code := app.EVMKeeper.GetCode(ctx, addr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"address": addr.Hex(),
			"code":    "0x" + hex.EncodeToString(code),
		})
	}).Methods("GET")

	// Get EVM storage value at address + slot
	router.HandleFunc("/syreen/v1/evm/storage/{address}/{slot}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addrHex := vars["address"]
		slotHex := vars["slot"]

		if !common.IsHexAddress(addrHex) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid hex address"})
			return
		}

		// H4: Validate slot is valid hex
		cleanSlot := slotHex
		if len(cleanSlot) >= 2 && cleanSlot[:2] == "0x" {
			cleanSlot = cleanSlot[2:]
		}
		if len(cleanSlot) == 0 || len(cleanSlot) > 64 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid slot: must be 1-64 hex characters"})
			return
		}
		for _, c := range cleanSlot {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid slot: must be hex"})
				return
			}
		}

		addr := common.HexToAddress(addrHex)
		slot := common.HexToHash(slotHex)
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		value := app.EVMKeeper.GetStorageAt(ctx, addr, slot)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"address": addr.Hex(),
			"slot":    slot.Hex(),
			"value":   value.Hex(),
		})
	}).Methods("GET")

	// Get EVM nonce for an address
	router.HandleFunc("/syreen/v1/evm/nonce/{address}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addrHex := vars["address"]

		if !common.IsHexAddress(addrHex) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid hex address"})
			return
		}

		addr := common.HexToAddress(addrHex)
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		nonce := app.EVMKeeper.GetEVMNonce(ctx, addr)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"address": addr.Hex(),
			"nonce":   nonce,
		})
	}).Methods("GET")

	// --- Strategy Endpoints (12 AI templates) ---

	router.HandleFunc("/syreen/intent/v1/strategy_templates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"templates": intenttypes.GetAvailableTemplates(),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/strategies/{address}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid address: " + err.Error()})
			return
		}

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		strategies := app.IntentKeeper.GetStrategiesByCreator(ctx, address)
		if strategies == nil {
			strategies = []intenttypes.Strategy{}
		}

		// Separate into active vs completed for UI convenience
		var active, completed []intenttypes.Strategy
		for _, s := range strategies {
			if s.Status == intenttypes.StrategyStatusActive {
				active = append(active, s)
			} else {
				completed = append(completed, s)
			}
		}
		if active == nil {
			active = []intenttypes.Strategy{}
		}
		if completed == nil {
			completed = []intenttypes.Strategy{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"strategies": strategies,
			"active":     active,
			"completed":  completed,
			"total":      len(strategies),
		})
	}).Methods("GET")

	router.HandleFunc("/syreen/intent/v1/strategy/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strategyID := vars["id"]

		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())

		strategy, found := app.IntentKeeper.GetStrategy(ctx, strategyID)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]string{"error": "strategy not found"})
			return
		}

		chain, chainFound := app.IntentKeeper.GetChain(ctx, strategy.ChainID)

		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"strategy": strategy,
			"found":    true,
		}
		if chainFound {
			resp["chain"] = chain
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	// --- DEX Signals ---
	router.HandleFunc("/syreen/dex/v1/signals/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{Height: app.LastBlockHeight()}, false, app.Logger())
		signals := app.DexKeeper.ComputePoolSignals(ctx, poolID)
		w.Header().Set("Content-Type", "application/json")
		if signals == nil {
			// Fallback: compute live from pool reserves
			pool, found := app.DexKeeper.GetPool(ctx, poolID)
			if found && !pool.ReserveA.IsZero() && !pool.ReserveB.IsZero() {
				price := math.LegacyNewDecFromInt(pool.ReserveB).Quo(math.LegacyNewDecFromInt(pool.ReserveA))
				json.NewEncoder(w).Encode(map[string]interface{}{
					"signals": map[string]interface{}{
						"pool_id":          poolID,
						"height":           app.LastBlockHeight(),
						"current_price":    price.String(),
						"signal":           "NEUTRAL",
						"composite_score":  50,
						"volatility_score": "0",
						"moving_averages": map[string]string{
							"sma_10": price.String(),
							"sma_20": price.String(),
							"sma_50": price.String(),
							"ema_10": price.String(),
							"ema_20": price.String(),
							"ema_50": price.String(),
						},
						"momentum": map[string]interface{}{
							"rsi":             "50",
							"price_roc":       "0",
							"volume_momentum": "0",
						},
					},
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"signals": nil})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"signals": signals})
	}).Methods("GET")

	// --- DEX Signal History ---
	router.HandleFunc("/syreen/dex/v1/signal_history/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		history, found := app.DexKeeper.GetSignalHistory(ctx, poolID)
		w.Header().Set("Content-Type", "application/json")
		if !found {
			json.NewEncoder(w).Encode(map[string]interface{}{"history": nil})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"history": history})
	}).Methods("GET")

	// --- DEX Sentiment ---
	router.HandleFunc("/syreen/dex/v1/sentiment/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		state, found := app.DexKeeper.GetSentimentState(ctx, poolID)
		w.Header().Set("Content-Type", "application/json")
		if !found {
			json.NewEncoder(w).Encode(map[string]interface{}{"sentiment": nil})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"sentiment": state})
	}).Methods("GET")

	// --- DEX Trades ---
	router.HandleFunc("/syreen/dex/v1/trades/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		limit := 100
		if l := r.URL.Query().Get("limit"); l != "" {
			if n, e := strconv.Atoi(l); e == nil && n > 0 && n <= 500 {
				limit = n
			}
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		trades := app.DexKeeper.GetTradesByPool(ctx, poolID, limit)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"trades": trades})
	}).Methods("GET")

	// --- DEX Order Book ---
	router.HandleFunc("/syreen/dex/v1/orderbook/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		ob := app.DexKeeper.GetOrderBook(ctx, poolID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"orderbook": ob})
	}).Methods("GET")

	// --- DEX Sentiment History ---
	router.HandleFunc("/syreen/dex/v1/sentiment_history/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		limit := int64(100)
		if l := r.URL.Query().Get("limit"); l != "" {
			if n, e := strconv.ParseInt(l, 10, 64); e == nil && n > 0 && n <= 500 {
				limit = n
			}
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		history := app.DexKeeper.GetSentimentHistory(ctx, poolID, limit)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"history": history})
	}).Methods("GET")

	// --- DEX Sentiment Alerts ---
	router.HandleFunc("/syreen/dex/v1/sentiment_alerts/{pool_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolID, err := strconv.ParseUint(vars["pool_id"], 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid pool_id"})
			return
		}
		cms := app.CommitMultiStore().CacheMultiStore()
		ctx := sdk.NewContext(cms, cmtproto.Header{}, false, app.Logger())
		alerts := app.DexKeeper.GetSentimentAlerts(ctx, poolID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"alerts": alerts})
	}).Methods("GET")
}
