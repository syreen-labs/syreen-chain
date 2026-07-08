package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"syreen/x/intent/types"
)

// MsgRouter defines the interface for routing sdk.Msg to their handlers.
type MsgRouter interface {
	Handler(msg sdk.Msg) baseapp.MsgServiceHandler
}

// executingKey is a context key used for the reentrancy guard.
// Using a context value instead of a struct field is thread-safe
// for Block-STM parallel execution.
type executingKey struct{}

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	dexKeeper      types.DexKeeper
	transferKeeper types.TransferKeeper
	mevKeeper      types.MEVKeeper
	msgRouter      MsgRouter
	authority     string
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	msgRouter MsgRouter,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		msgRouter:     msgRouter,
		authority:     authority,
	}
}

// GetAuthority returns the module's governance authority address (the account
// permitted to execute MsgUpdateParams). Set from the gov module address at
// keeper construction.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// SetDexKeeper sets the DEX keeper for trading intent execution.
func (k *Keeper) SetDexKeeper(dexKeeper types.DexKeeper) {
	k.dexKeeper = dexKeeper
}

// SetTransferKeeper sets the IBC transfer keeper for cross-chain intents.
func (k *Keeper) SetTransferKeeper(transferKeeper types.TransferKeeper) {
	k.transferKeeper = transferKeeper
}

// SetMEVKeeper wires the Fairness Engine "Return" sink (fairness pool + user
// rebate ledger). Optional; when nil, slashed stake falls back to the community
// pool and no rebate beneficiary is recorded.
func (k *Keeper) SetMEVKeeper(mevKeeper types.MEVKeeper) {
	k.mevKeeper = mevKeeper
}

// routeSlashToFairnessPool moves confiscated solver stake into the Fairness
// Engine pool as REAL coins (the user-rebate funding source) when the MEV keeper
// is wired; otherwise it falls back to the community pool. Previously the coins
// were sent to the distribution *account* without FundCommunityPool, orphaning
// them — routing to the fairness pool both fixes that and funds the rebate.
func (k Keeper) routeSlashToFairnessPool(ctx context.Context, slashCoins sdk.Coins) error {
	if k.mevKeeper != nil {
		return k.mevKeeper.CreditFairnessPool(ctx, types.ModuleName, slashCoins, "")
	}
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, "distribution", slashCoins)
}

// isExecuting checks the context-based reentrancy guard.
// Thread-safe for Block-STM parallel execution since each
// goroutine has its own context.
func (k Keeper) isExecuting(ctx context.Context) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Value(executingKey{}) != nil
}

// withExecuting returns a new context with the reentrancy guard set.
func (k Keeper) withExecuting(ctx context.Context) context.Context {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.WithValue(executingKey{}, true)
}

// SetExecuting is a test helper that returns a new context with
// the reentrancy guard set (v=true) or returns ctx unchanged (v=false).
func (k *Keeper) SetExecuting(ctx context.Context, v bool) context.Context {
	if v {
		return k.withExecuting(ctx)
	}
	return ctx
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// MaxIntentsPerUserPerBlock is the maximum number of intents a single user
// can submit in a single block, preventing spam and resource exhaustion.
const MaxIntentsPerUserPerBlock = 10

// SubmitIntent creates a new intent and locks the maxFee from the creator
// validateIntentBody performs defense-in-depth validation of an intent body at
// submission time. It guarantees that any body which later reaches the trading
// execution paths (extractTradingInputCoins at submit, and executeTradingSwap /
// tryDCA in BeginBlock) carries non-nil, positive amounts and valid denoms — so
// downstream sdk.NewCoin / QuoRaw calls cannot panic on malformed input. Swap
// intents keep their historical, looser rules (input is delivered by the solver,
// not locked here): only a parseable body and a non-negative min floor.
func validateIntentBody(intentType string, body json.RawMessage) error {
	requirePositive := func(field string, amt math.Int) error {
		if amt.IsNil() {
			return fmt.Errorf("%s must be set", field)
		}
		if !amt.IsPositive() {
			return fmt.Errorf("%s must be positive", field)
		}
		return nil
	}
	requireDenom := func(field, denom string) error {
		if err := sdk.ValidateDenom(denom); err != nil {
			return fmt.Errorf("invalid %s %q: %w", field, denom, err)
		}
		return nil
	}
	// An optional min-output floor: nil (omitted) is allowed and treated as zero
	// downstream; a set value must be >= 0.
	requireNonNegMin := func(amt math.Int) error {
		if !amt.IsNil() && amt.IsNegative() {
			return fmt.Errorf("min_output_amount must not be negative")
		}
		return nil
	}

	switch intentType {
	case types.IntentTypeSwap:
		var b types.SwapIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid swap intent body: %w", err)
		}
		return requireNonNegMin(b.MinOutputAmount)

	case types.IntentTypeLimitBuy:
		var b types.LimitBuyIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid limit_buy intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("input_amount", b.InputAmount); err != nil {
			return err
		}
		return requireNonNegMin(b.MinOutputAmount)

	case types.IntentTypeLimitSell:
		var b types.LimitSellIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid limit_sell intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("input_amount", b.InputAmount); err != nil {
			return err
		}
		return requireNonNegMin(b.MinOutputAmount)

	case types.IntentTypeStopLoss:
		var b types.StopLossIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid stop_loss intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("input_amount", b.InputAmount); err != nil {
			return err
		}
		return requireNonNegMin(b.MinOutputAmount)

	case types.IntentTypeTakeProfit:
		var b types.TakeProfitIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid take_profit intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("input_amount", b.InputAmount); err != nil {
			return err
		}
		return requireNonNegMin(b.MinOutputAmount)

	case types.IntentTypeDCA, types.IntentTypeTWAP:
		var b types.DCAIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid dca/twap intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("total_amount", b.TotalAmount); err != nil {
			return err
		}
		if b.NumExecutions == 0 {
			return fmt.Errorf("num_executions must be greater than zero")
		}
		return nil

	case types.IntentTypeCrossChainSwap:
		var b types.CrossChainSwapIntent
		if err := json.Unmarshal(body, &b); err != nil {
			return fmt.Errorf("invalid cross_chain_swap intent body: %w", err)
		}
		if err := requireDenom("input_denom", b.InputDenom); err != nil {
			return err
		}
		if err := requireDenom("output_denom", b.OutputDenom); err != nil {
			return err
		}
		if err := requirePositive("input_amount", b.InputAmount); err != nil {
			return err
		}
		if b.Receiver == "" {
			return fmt.Errorf("receiver must be set for cross_chain_swap")
		}
		if b.IBCSourceChannel == "" {
			return fmt.Errorf("ibc_source_channel must be set for cross_chain_swap")
		}
		return requireNonNegMin(b.MinOutputAmount)

	default:
		// Non-trading / unknown intent types are not amount-bearing here; leave
		// their handling to their own execution paths.
		return nil
	}
}

func (k Keeper) SubmitIntent(ctx context.Context, msg *types.MsgSubmitIntent) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	if !params.EnableIntents {
		return "", fmt.Errorf("intents are currently disabled")
	}

	// Rate limit: count how many intents this creator already submitted in the current block
	currentBlock := sdkCtx.BlockHeight()
	existingIntents := k.GetIntentsByCreator(ctx, msg.Creator)
	countThisBlock := 0
	for _, intent := range existingIntents {
		if intent.CreatedAt == currentBlock {
			countThisBlock++
		}
	}
	if countThisBlock >= MaxIntentsPerUserPerBlock {
		return "", fmt.Errorf("rate limit exceeded: maximum %d intents per user per block", MaxIntentsPerUserPerBlock)
	}

	// Defense in depth: reject malformed intent bodies at submission so they never
	// enter state and reach the BeginBlock execution paths. A nil/zero input amount
	// on a trading body would otherwise panic in sdk.NewCoin (extractTradingInputCoins
	// / executeTradingSwap); an empty denom would produce invalid coins. Validating
	// here keeps the BeginBlock choke points (guarded separately) from ever seeing
	// obviously-broken bodies.
	if err := validateIntentBody(msg.IntentType, msg.Body); err != nil {
		return "", err
	}

	// Generate intent ID
	intentID := k.nextIntentID(ctx)

	// Calculate expiry block height
	expiryBlocks := msg.ExpiryBlocks
	if expiryBlocks > params.IntentExpiryBlocks {
		expiryBlocks = params.IntentExpiryBlocks
	}
	expiry := sdkCtx.BlockHeight() + int64(expiryBlocks)

	// The solver auction runs at created + SolvingWindow, and FulfillIntent
	// rejects an intent whose hard Expiry has passed. If Expiry <= the solving
	// deadline (e.g. the CLI's default ExpiryBlocks happens to equal SolvingWindow),
	// the intent hard-expires on the exact block the auction fires and can never be
	// fulfilled. Guarantee a margin past the solving deadline regardless of the
	// requested expiry so every intent gets a real chance to be settled.
	if minExpiry := sdkCtx.BlockHeight() + int64(params.SolvingWindow) + 10; expiry < minExpiry {
		expiry = minExpiry
	}

	intent := types.Intent{
		ID:         intentID,
		Creator:    msg.Creator,
		IntentType: msg.IntentType,
		Body:       msg.Body,
		MaxFee:     msg.MaxFee,
		Tip:        msg.Tip,
		Expiry:     expiry,
		Status:     types.StatusPending,
		CreatedAt:  sdkCtx.BlockHeight(),
		SolverAddr: "",
	}

	// Lock maxFee + tip from creator into module account
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return "", fmt.Errorf("invalid creator address: %w", err)
	}

	totalLock := msg.MaxFee.Add(msg.Tip...)
	if totalLock.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, totalLock); err != nil {
			return "", fmt.Errorf("failed to lock fees: %w", err)
		}
	}

	// For trading intents, also lock the input tokens that will be swapped
	tradingCoins := k.extractTradingInputCoins(msg.IntentType, msg.Body)
	if tradingCoins.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, tradingCoins); err != nil {
			return "", fmt.Errorf("failed to lock trading input tokens: %w", err)
		}
	}

	k.SetIntent(ctx, intent)
	k.Logger(ctx).Info("intent submitted", "id", intentID, "creator", msg.Creator, "type", msg.IntentType, "expiry", expiry)

	return intentID, nil
}

// RegisterSolver registers a new solver and locks their stake
func (k Keeper) RegisterSolver(ctx context.Context, msg *types.MsgRegisterSolver) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	// Check if solver already registered
	_, found := k.GetSolver(ctx, msg.Address)
	if found {
		return types.ErrSolverAlreadyRegistered
	}

	// Check minimum stake
	if msg.StakeAmount.IsLT(params.MinSolverStake) {
		return types.ErrInsufficientStake
	}

	// Lock stake into module account
	solverAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return fmt.Errorf("invalid solver address: %w", err)
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, solverAddr, types.ModuleName, sdk.NewCoins(msg.StakeAmount)); err != nil {
		return fmt.Errorf("failed to lock solver stake: %w", err)
	}

	solver := types.Solver{
		Address:         msg.Address,
		Moniker:         msg.Moniker,
		StakedAmount:    msg.StakeAmount,
		ReputationScore: 100, // starting reputation
		TotalSolved:     0,
		TotalFailed:     0,
		Active:          true,
		JoinedAt:        sdkCtx.BlockHeight(),
	}

	k.SetSolver(ctx, solver)
	k.Logger(ctx).Info("solver registered", "address", msg.Address, "moniker", msg.Moniker, "stake", msg.StakeAmount)

	return nil
}

// DeregisterSolver marks a solver as inactive and starts the unbonding period.
// The stake is NOT returned immediately — it remains locked for SolverUnbondingBlocks
// to allow slashing of misbehavior discovered after deregistration.
func (k Keeper) DeregisterSolver(ctx context.Context, msg *types.MsgDeregisterSolver) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	solver, found := k.GetSolver(ctx, msg.Address)
	if !found {
		return types.ErrSolverNotRegistered
	}

	if !solver.Active {
		return fmt.Errorf("solver is already deregistered and unbonding")
	}

	params := k.GetParams(ctx)

	// Mark as inactive with an unbonding completion height
	solver.Active = false
	solver.UnbondingHeight = sdkCtx.BlockHeight() + int64(params.SolverUnbondingBlocks)
	k.SetSolver(ctx, solver)

	k.Logger(ctx).Info("solver deregistered, unbonding started",
		"address", msg.Address,
		"unbonding_complete_at", solver.UnbondingHeight,
	)
	return nil
}

// CompleteSolverUnbonding returns stake to solvers whose unbonding period has ended.
// Called in BeginBlocker.
func (k Keeper) CompleteSolverUnbonding(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)

	prefix := []byte(types.SolverPrefixKey)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return
	}
	// Collect keys to delete while iterating, then delete after the iterator is
	// closed. Deleting under a live iterator over the same range is unsafe.
	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		var solver types.Solver
		if err := json.Unmarshal(iter.Value(), &solver); err != nil {
			continue
		}

		// Check if this solver is unbonding and the period has passed
		if !solver.Active && solver.UnbondingHeight > 0 && sdkCtx.BlockHeight() >= solver.UnbondingHeight {
			solverAddr, err := sdk.AccAddressFromBech32(solver.Address)
			if err != nil {
				continue
			}

			if solver.StakedAmount.IsPositive() {
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, solverAddr, sdk.NewCoins(solver.StakedAmount)); err != nil {
					k.Logger(ctx).Error("failed to return solver stake", "solver", solver.Address, "error", err)
					// Fall through: still delete solver to keep state consistent across all
					// validators. Stake remains in module account, recoverable via governance.
				}
			}

			// Queue solver for removal — must happen regardless of bank send result
			// to ensure deterministic state across all validators. Copy the key since
			// the iterator may reuse the underlying buffer.
			key := append([]byte(nil), iter.Key()...)
			keysToDelete = append(keysToDelete, key)
			sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
				"solver_unbonding_complete",
				sdk.NewAttribute("address", solver.Address),
			))
			k.Logger(ctx).Info("solver unbonding complete", "address", solver.Address)
		}
	}
	iter.Close()

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}
}

// prefixEndBytes returns the end key for a prefix scan.
func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}

// SubmitSolution stores a solution for an intent
func (k Keeper) SubmitSolution(ctx context.Context, msg *types.MsgSubmitSolution) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	// Verify solver is registered and active
	solver, found := k.GetSolver(ctx, msg.SolverAddr)
	if !found || !solver.Active {
		return types.ErrSolverNotRegistered
	}

	// Verify intent exists and is in a solvable state
	intent, found := k.GetIntent(ctx, msg.IntentID)
	if !found {
		return types.ErrIntentNotFound
	}

	// Prevent self-solving: creator cannot be the solver (prevents reputation farming)
	if intent.Creator == msg.SolverAddr {
		return fmt.Errorf("intent creator cannot solve their own intent")
	}

	if intent.Status != types.StatusPending && intent.Status != types.StatusSolving {
		return types.ErrIntentAlreadyFulfilled
	}

	// Check if intent has expired
	if sdkCtx.BlockHeight() > intent.Expiry {
		return types.ErrIntentExpired
	}

	// Check solving window (solutions must be submitted within SolvingWindow blocks of intent creation)
	if sdkCtx.BlockHeight() > intent.CreatedAt+int64(params.SolvingWindow) {
		return types.ErrSolvingWindowClosed
	}

	// Check max solutions limit
	solutions := k.GetSolutionsForIntent(ctx, msg.IntentID)
	if uint64(len(solutions)) >= params.MaxSolutionsPerIntent {
		return types.ErrMaxSolutionsReached
	}

	solution := types.Solution{
		IntentID:        msg.IntentID,
		SolverAddr:      msg.SolverAddr,
		ExecutionMsgs:   msg.ExecutionMsgs,
		ExpectedOutcome: msg.ExpectedOutcome,
		GasEstimate:     0,
		SubmittedAt:     sdkCtx.BlockHeight(),
	}

	k.SetSolution(ctx, solution)

	// Update intent status to solving if still pending
	if intent.Status == types.StatusPending {
		intent.Status = types.StatusSolving
		k.SetIntent(ctx, intent)
	}

	k.Logger(ctx).Info("solution submitted", "intent_id", msg.IntentID, "solver", msg.SolverAddr)
	return nil
}

// FulfillIntent picks the best solution, executes it, pays the solver tip, and updates reputation.
// If msg.SolverAddr is empty, auto-selects the best solver via SelectWinningSolver auction.
func (k *Keeper) FulfillIntent(ctx context.Context, msg *types.MsgFulfillIntent) error {
	// C-01: Reentrancy guard — prevent a solution from calling FulfillIntent recursively.
	// Uses context-based guard which is thread-safe for Block-STM parallel execution.
	if k.isExecuting(ctx) {
		return types.ErrReentrant
	}
	ctx = k.withExecuting(ctx)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	intent, found := k.GetIntent(ctx, msg.IntentID)
	if !found {
		return types.ErrIntentNotFound
	}

	if intent.Status == types.StatusFulfilled {
		return types.ErrIntentAlreadyFulfilled
	}

	if intent.Status == types.StatusExpired || intent.Status == types.StatusFailed {
		return fmt.Errorf("intent is in terminal state: %s", intent.Status)
	}

	// Defense in depth: check expiry by block height, not just status
	if sdkCtx.BlockHeight() > intent.Expiry {
		return types.ErrIntentExpired
	}

	solverAddrStr := msg.SolverAddr

	// Verify solver address matches sender when explicitly specified
	if solverAddrStr != "" && msg.Sender != "" && msg.Sender != solverAddrStr {
		return fmt.Errorf("sender %s does not match solver address %s", msg.Sender, solverAddrStr)
	}

	// Auto-select solver via auction if SolverAddr is empty
	if solverAddrStr == "" {
		// H-14: When auto-selecting, verify the sender is the intent creator
		// to prevent arbitrary third parties from triggering fulfillment.
		// The module's own AutoFulfillIntents (BeginBlock) passes Sender="" which
		// is allowed since it constructs the msg internally.
		if msg.Sender != "" && msg.Sender != intent.Creator {
			return fmt.Errorf("only the intent creator can trigger auto-fulfillment (sender %s != creator %s)", msg.Sender, intent.Creator)
		}

		winningSol, err := k.SelectWinningSolver(ctx, msg.IntentID)
		if err != nil {
			return fmt.Errorf("auction auto-selection failed: %w", err)
		}
		solverAddrStr = winningSol.SolverAddr
	}

	// Verify the solver submitted a solution for this intent
	solutions := k.GetSolutionsForIntent(ctx, msg.IntentID)
	var winningSolution *types.Solution
	for _, sol := range solutions {
		if sol.SolverAddr == solverAddrStr {
			s := sol
			winningSolution = &s
			break
		}
	}

	if winningSolution == nil {
		return types.ErrInvalidSolution
	}

	// Lifecycle event: the auction winner has been selected for this intent.
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"auction_won",
		sdk.NewAttribute("intent_id", msg.IntentID),
		sdk.NewAttribute("solver_addr", solverAddrStr),
	))

	// Snapshot the creator's output-denom balance BEFORE the solution executes so
	// that outcome verification can measure the delta actually delivered by the
	// solver rather than the creator's pre-existing holdings.
	var swapOutputBefore math.Int
	if intent.IntentType == types.IntentTypeSwap {
		swapOutputBefore = k.swapOutputBalanceBefore(ctx, intent)
	}

	// Execute the solution's messages
	if err := k.executeSolutionMsgs(ctx, winningSolution); err != nil {
		// Mark intent as failed, not fulfilled
		intent.Status = types.StatusFailed
		k.SetIntent(ctx, intent)

		// Update solver failure stats and apply REAL slashing
		solver, solverFound := k.GetSolver(ctx, solverAddrStr)
		if solverFound {
			solver.TotalFailed++
			if solver.ReputationScore >= 5 {
				solver.ReputationScore -= 5
			} else {
				solver.ReputationScore = 0
			}

			// Real slashing: confiscate a fraction of staked tokens
			slashAmount := params.SolverSlashFraction.MulInt(solver.StakedAmount.Amount).TruncateInt()
			if slashAmount.IsPositive() {
				slashCoins := sdk.NewCoins(sdk.NewCoin(solver.StakedAmount.Denom, slashAmount))

				// Send slashed coins from intent module to distribution (community pool)
				if sendErr := k.routeSlashToFairnessPool(ctx, slashCoins); sendErr != nil {
					k.Logger(ctx).Error("failed to slash solver stake to community pool",
						"solver", solverAddrStr, "amount", slashCoins, "error", sendErr)
				} else {
					// Reduce the solver's recorded stake
					solver.StakedAmount.Amount = solver.StakedAmount.Amount.Sub(slashAmount)

					sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
						"solver_slashed",
						sdk.NewAttribute("solver_addr", solverAddrStr),
						sdk.NewAttribute("intent_id", msg.IntentID),
						sdk.NewAttribute("slashed", slashCoins.String()),
						sdk.NewAttribute("reason", "execution_failed"),
					))

					k.Logger(ctx).Info("solver slashed",
						"solver", solverAddrStr, "slashed", slashCoins, "remaining_stake", solver.StakedAmount)
				}
			}

			// Deactivate solver if stake drops below minimum
			if solver.StakedAmount.IsLT(params.MinSolverStake) {
				solver.Active = false
				sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
					"solver_deactivated",
					sdk.NewAttribute("solver_addr", solverAddrStr),
					sdk.NewAttribute("reason", "insufficient_stake"),
				))
				k.Logger(ctx).Info("solver deactivated due to insufficient stake after slashing",
					"solver", solverAddrStr, "remaining_stake", solver.StakedAmount)
			}

			k.SetSolver(ctx, solver)
		}

		// C-09: Refund creator's locked MaxFee + Tip on failed solution execution.
		// This happens AFTER solver slashing so slashed coins come from solver's stake,
		// not from the locked fees.
		creatorAddr, addrErr := sdk.AccAddressFromBech32(intent.Creator)
		if addrErr == nil {
			totalRefund := intent.MaxFee.Add(intent.Tip...)
			if totalRefund.IsAllPositive() {
				if refundErr := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalRefund); refundErr != nil {
					k.Logger(ctx).Error("failed to refund creator after failed solution",
						"intent_id", msg.IntentID, "creator", intent.Creator, "error", refundErr)
				}
			}
		}

		return fmt.Errorf("solution execution failed: %w", err)
	}

	// Outcome verification: for swap intents, verify the creator's balance changed
	// in the expected direction. If not, mark intent as failed and slash the solver.
	if intent.IntentType == types.IntentTypeSwap {
		// Bind the winning solver to the output it declared in the auction (best
		// execution). declaredOutput is zero if the solver made no declaration, in
		// which case verifySwapOutcome falls back to the intent's MinOutputAmount.
		declaredOutput, _ := winningSolution.DeclaredOutput()
		if verifyErr := k.verifySwapOutcome(ctx, intent, solverAddrStr, params, swapOutputBefore, declaredOutput); verifyErr != nil {
			k.Logger(ctx).Error("outcome verification failed",
				"intent_id", msg.IntentID, "solver", solverAddrStr, "error", verifyErr)
			// Mark as failed, slash solver
			intent.Status = types.StatusFailed
			k.SetIntent(ctx, intent)

			solver, solverFound := k.GetSolver(ctx, solverAddrStr)
			if solverFound {
				solver.TotalFailed++
				if solver.ReputationScore >= 10 {
					solver.ReputationScore -= 10
				} else {
					solver.ReputationScore = 0
				}
				slashAmount := params.SolverSlashFraction.MulInt(solver.StakedAmount.Amount).TruncateInt()
				if slashAmount.IsPositive() {
					slashCoins := sdk.NewCoins(sdk.NewCoin(solver.StakedAmount.Denom, slashAmount))
					if sendErr := k.routeSlashToFairnessPool(ctx, slashCoins); sendErr == nil {
						solver.StakedAmount.Amount = solver.StakedAmount.Amount.Sub(slashAmount)
						sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
							"solver_slashed",
							sdk.NewAttribute("solver_addr", solverAddrStr),
							sdk.NewAttribute("intent_id", msg.IntentID),
							sdk.NewAttribute("slashed", slashCoins.String()),
							sdk.NewAttribute("reason", "outcome_verification_failed"),
						))
					}
				}
				if solver.StakedAmount.IsLT(params.MinSolverStake) {
					solver.Active = false
					sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
						"solver_deactivated",
						sdk.NewAttribute("solver_addr", solverAddrStr),
						sdk.NewAttribute("reason", "insufficient_stake"),
					))
				}
				k.SetSolver(ctx, solver)
			}

			return fmt.Errorf("outcome verification failed: %w", verifyErr)
		}
	}

	// Mark intent as fulfilled
	intent.Status = types.StatusFulfilled
	intent.SolverAddr = solverAddrStr
	k.SetIntent(ctx, intent)

	// Fairness Engine · Return: mark the trading user as a first-priority rebate
	// beneficiary for the next fairness-pool distribution, weighted by their tip
	// (a trade-size proxy). Nil-safe when the MEV keeper isn't wired.
	if k.mevKeeper != nil {
		weight := intent.Tip.AmountOf(sdk.DefaultBondDenom).Add(intent.MaxFee.AmountOf(sdk.DefaultBondDenom))
		k.mevKeeper.RecordRebateBeneficiary(ctx, intent.Creator, weight)
	}

	// Pay solver the tip from module account
	solverAddr, err := sdk.AccAddressFromBech32(solverAddrStr)
	if err != nil {
		return fmt.Errorf("invalid solver address: %w", err)
	}

	if intent.Tip.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, solverAddr, intent.Tip); err != nil {
			return fmt.Errorf("failed to pay solver tip: %w", err)
		}
	}

	// Refund remaining maxFee to creator
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	if intent.MaxFee.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, intent.MaxFee); err != nil {
			return fmt.Errorf("failed to refund max fee: %w", err)
		}
	}

	// Update solver reputation
	solver, solverFound := k.GetSolver(ctx, solverAddrStr)
	if solverFound {
		solver.TotalSolved++
		solver.ReputationScore += 10
		if solver.ReputationScore > 1000 {
			solver.ReputationScore = 1000
		}
		k.SetSolver(ctx, solver)
	}

	k.Logger(ctx).Info("intent fulfilled", "intent_id", msg.IntentID, "solver", solverAddrStr)
	return nil
}

// executeSolutionMsgs decodes and routes each execution message through the SDK message router.
// All messages are executed against a CacheContext; state is only committed if ALL succeed.
func (k *Keeper) executeSolutionMsgs(ctx context.Context, solution *types.Solution) error {
	// C-01 defense in depth: verify reentrancy guard is set by the caller (FulfillIntent).
	// This prevents bypasses where executeSolutionMsgs could be reached without the guard.
	if !k.isExecuting(ctx) {
		return fmt.Errorf("intent: executeSolutionMsgs called without reentrancy guard")
	}

	if len(solution.ExecutionMsgs) == 0 {
		return fmt.Errorf("solution has no execution messages")
	}
	// C-3 defense in depth: enforce upper bound on execution msgs at runtime
	if len(solution.ExecutionMsgs) > types.MaxExecutionMsgs {
		return fmt.Errorf("solution execution messages count %d exceeds maximum %d", len(solution.ExecutionMsgs), types.MaxExecutionMsgs)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Execute all messages in a cache context so partial failures are rolled back.
	cacheCtx, write := sdkCtx.CacheContext()

	// C-11: Build allowed message type whitelist from params
	params := k.GetParams(ctx)
	allowedTypes := make(map[string]bool, len(params.AllowedMsgTypes))
	for _, t := range params.AllowedMsgTypes {
		allowedTypes[t] = true
	}

	// SECURITY (fund-safety invariants). The signer-based authorization below is
	// necessary but NOT sufficient: solution messages are executed here in
	// BeginBlock without any real signature, so a solver could craft a message
	// that spends the module account's pooled escrow (all users' fees/tips/inputs
	// + every solver's stake) or drains the creator's wallet. We defend by
	// snapshotting balances and enforcing, after execution, that:
	//   (1) the intent MODULE ACCOUNT balance does not DECREASE in any denom
	//       (solutions must never move protocol escrow — keeper-driven paths do
	//        that separately), and
	//   (2) the CREATOR's balance only decreases by at most the intent's declared
	//       input (InputAmount of InputDenom); no other creator denom may drop.
	// A violation fails the whole solution (cacheCtx is never written).
	moduleAccAddr := authtypes.NewModuleAddress(types.ModuleName)
	intent, intentFound := k.GetIntent(ctx, solution.IntentID)
	var intentCreatorAddr sdk.AccAddress
	if intentFound {
		intentCreatorAddr, _ = sdk.AccAddressFromBech32(intent.Creator)
	}
	// Authorized creator input spend (denom + max amount); zero for intents that
	// do not declare a swap input, which fully protects every creator denom.
	authorizedInputDenom := ""
	authorizedInputAmount := math.ZeroInt()
	if intentFound {
		var sb types.SwapIntent
		if err := json.Unmarshal(intent.Body, &sb); err == nil && sb.InputDenom != "" && !sb.InputAmount.IsNil() && sb.InputAmount.IsPositive() {
			authorizedInputDenom = sb.InputDenom
			authorizedInputAmount = sb.InputAmount
		}
	}
	moduleBefore := k.bankKeeper.GetAllBalances(cacheCtx, moduleAccAddr)
	var creatorBefore sdk.Coins
	if intentCreatorAddr != nil {
		creatorBefore = k.bankKeeper.GetAllBalances(cacheCtx, intentCreatorAddr)
	}

	for i, rawMsg := range solution.ExecutionMsgs {
		// Each execution msg is a protojson-encoded sdk.Msg using the standard
		// "@type" discriminator, e.g. {"@type":"/cosmos.bank.v1beta1.MsgSend", ...}.
		// This MUST be decoded through the codec's interface-JSON path, which
		// resolves "@type" against the interface registry. A plain json.Unmarshal
		// into a codectypes.Any does NOT map "@type" (the Any field is tagged
		// "type_url"), leaving TypeUrl empty so the message resolves to type "/"
		// and is rejected — which silently broke ALL solution execution.
		var sdkMsg sdk.Msg
		if err := k.cdc.UnmarshalInterfaceJSON(rawMsg, &sdkMsg); err != nil {
			// Fallback for HAND-ROLLED custom types (e.g. /syreen.dex.MsgSwap) that
			// lack gogoproto jsonpb support and so can't be decoded by the interface
			// JSON path above. Resolve the concrete type from the registry by its
			// "@type" and json.Unmarshal via its struct tags. Standard cosmos types
			// still take the jsonpb path. Signer resolution + whitelist + the balance
			// invariants below all run identically on the result, so this widens what
			// can be decoded WITHOUT weakening any security check.
			var disc struct {
				Type string `json:"@type"`
			}
			if jerr := json.Unmarshal(rawMsg, &disc); jerr != nil || disc.Type == "" {
				return fmt.Errorf("failed to decode execution msg %d: %w", i, err)
			}
			resolved, rerr := k.cdc.InterfaceRegistry().Resolve(disc.Type)
			if rerr != nil {
				return fmt.Errorf("failed to decode execution msg %d (type %s): %w", i, disc.Type, rerr)
			}
			if jerr := json.Unmarshal(rawMsg, resolved); jerr != nil {
				return fmt.Errorf("failed to decode execution msg %d (type %s): %w", i, disc.Type, jerr)
			}
			m, ok := resolved.(sdk.Msg)
			if !ok {
				return fmt.Errorf("execution msg %d (type %s) is not an sdk.Msg", i, disc.Type)
			}
			sdkMsg = m
		}

		// C-11: Check message type against whitelist
		msgTypeURL := sdk.MsgTypeURL(sdkMsg)
		if !allowedTypes[msgTypeURL] {
			return fmt.Errorf("%w: %s", types.ErrDisallowedMsgType, msgTypeURL)
		}

		// C-01b: Verify all signers match either the intent module account or the intent creator.
		// (First line of defense; the balance invariants enforced after the loop are
		// what actually bound fund movement — see the SECURITY note above.)
		// SECURITY: The solver address is NOT an allowed signer — only the module account and
		// intent creator are authorized. This prevents solver self-solving via signer manipulation.
		//
		// SDK v0.53 removed the legacy sdk.Msg.GetSigners() method, so a type
		// assertion on `interface{ GetSigners() []sdk.AccAddress }` silently fails
		// (ok=false) for bank/IBC/compute messages, which would skip authorization
		// entirely and let a solver spend arbitrary accounts' funds. Resolve signers
		// through the codec's protobuf signing context instead, and FAIL CLOSED: any
		// error or a message with no resolvable signers is rejected, never executed.
		signerBzs, _, err := k.cdc.GetMsgV1Signers(sdkMsg)
		if err != nil {
			return fmt.Errorf("execution msg %d (type %s): failed to resolve signers: %w", i, msgTypeURL, err)
		}
		if len(signerBzs) == 0 {
			return fmt.Errorf("execution msg %d (type %s) has no resolvable signers; refusing to execute", i, msgTypeURL)
		}
		for _, signerBz := range signerBzs {
			signer := sdk.AccAddress(signerBz)
			if signer.Equals(moduleAccAddr) {
				continue
			}
			if intentCreatorAddr != nil && signer.Equals(intentCreatorAddr) {
				continue
			}
			return fmt.Errorf("execution msg %d has unauthorized signer %s (must be module account or intent creator)", i, signer)
		}

		// Route and execute the message against cacheCtx
		handler := k.msgRouter.Handler(sdkMsg)
		if handler == nil {
			return fmt.Errorf("no handler found for execution msg %d (type %s)", i, msgTypeURL)
		}

		if _, err := handler(cacheCtx, sdkMsg); err != nil {
			return fmt.Errorf("execution msg %d (type %s) failed: %w", i, msgTypeURL, err)
		}
	}

	// SECURITY invariant (1): the module account (pooled escrow) must not lose
	// funds to a solution. Reject any denom whose module balance decreased.
	moduleAfter := k.bankKeeper.GetAllBalances(cacheCtx, moduleAccAddr)
	for _, before := range moduleBefore {
		if moduleAfter.AmountOf(before.Denom).LT(before.Amount) {
			return fmt.Errorf("solution would reduce intent module escrow of %s (%s -> %s); refusing",
				before.Denom, before.Amount, moduleAfter.AmountOf(before.Denom))
		}
	}

	// SECURITY invariant (2): the creator may only be debited up to the intent's
	// declared input (InputAmount of InputDenom). Every other denom must not drop,
	// and InputDenom must not drop by more than the authorized amount.
	if intentCreatorAddr != nil {
		creatorAfter := k.bankKeeper.GetAllBalances(cacheCtx, intentCreatorAddr)
		for _, before := range creatorBefore {
			spent := before.Amount.Sub(creatorAfter.AmountOf(before.Denom))
			if !spent.IsPositive() {
				continue
			}
			allowed := math.ZeroInt()
			if before.Denom == authorizedInputDenom {
				allowed = authorizedInputAmount
			}
			if spent.GT(allowed) {
				return fmt.Errorf("solution would spend %s%s of creator funds, exceeding authorized input (%s%s); refusing",
					spent, before.Denom, allowed, before.Denom)
			}
		}
	}

	// All messages succeeded and fund-safety invariants held — commit.
	write()
	return nil
}

// GetIntent retrieves an intent by ID
func (k Keeper) GetIntent(ctx context.Context, intentID string) (types.Intent, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.IntentKey(intentID))
	if err != nil || bz == nil {
		return types.Intent{}, false
	}
	var intent types.Intent
	if err := json.Unmarshal(bz, &intent); err != nil {
		return types.Intent{}, false
	}
	return intent, true
}

// SetIntent stores an intent
func (k Keeper) SetIntent(ctx context.Context, intent types.Intent) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(intent)
	kvStore.Set(types.IntentKey(intent.ID), bz)
}

// GetSolver retrieves a solver by address
func (k Keeper) GetSolver(ctx context.Context, address string) (types.Solver, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.SolverKey(address))
	if err != nil || bz == nil {
		return types.Solver{}, false
	}
	var solver types.Solver
	if err := json.Unmarshal(bz, &solver); err != nil {
		return types.Solver{}, false
	}
	return solver, true
}

// SetSolver stores a solver
func (k Keeper) SetSolver(ctx context.Context, solver types.Solver) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(solver)
	kvStore.Set(types.SolverKey(solver.Address), bz)
}

// SetSolution stores a solution
func (k Keeper) SetSolution(ctx context.Context, solution types.Solution) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(solution)
	kvStore.Set(types.SolutionKey(solution.IntentID, solution.SolverAddr), bz)
}

// GetSolutionsForIntent returns all solutions for a given intent
func (k Keeper) GetSolutionsForIntent(ctx context.Context, intentID string) []types.Solution {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.SolutionsByIntentPrefix(intentID)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var solutions []types.Solution
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var solution types.Solution
		if err := json.Unmarshal(iter.Value(), &solution); err != nil {
			continue
		}
		solutions = append(solutions, solution)
	}

	// Explicit deterministic sort by solver address to guarantee identical
	// results regardless of KV store backend iteration order.
	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].SolverAddr < solutions[j].SolverAddr
	})

	return solutions
}

// deleteSolutionsForIntent removes all solutions associated with an intent (M-04).
func (k Keeper) deleteSolutionsForIntent(ctx context.Context, intentID string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := types.SolutionsByIntentPrefix(intentID)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return
	}
	defer iter.Close()

	var keysToDelete [][]byte
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		keysToDelete = append(keysToDelete, key)
	}

	for _, key := range keysToDelete {
		kvStore.Delete(key)
	}
}

// GetAllIntents returns all intents regardless of status.
func (k Keeper) GetAllIntents(ctx context.Context) []types.Intent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var intents []types.Intent
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		intents = append(intents, intent)
	}

	sort.Slice(intents, func(i, j int) bool {
		return intents[i].ID < intents[j].ID
	})

	return intents
}

// GetIntentsByType returns all pending intents matching the given intent type.
// This is useful for the executor/solver service to query active limit orders,
// stop-loss orders, DCA intents, etc.
func (k Keeper) GetIntentsByType(ctx context.Context, intentType string) []types.Intent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var intents []types.Intent
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		if intent.IntentType == intentType && (intent.Status == types.StatusPending || intent.Status == types.StatusSolving) {
			intents = append(intents, intent)
		}
	}

	sort.Slice(intents, func(i, j int) bool {
		return intents[i].ID < intents[j].ID
	})

	return intents
}

// GetPendingIntents returns all intents with StatusPending
func (k Keeper) GetPendingIntents(ctx context.Context) []types.Intent {
	return k.getIntentsByStatus(ctx, types.StatusPending)
}

// GetIntentsByCreator returns all intents created by the given address.
func (k Keeper) GetIntentsByCreator(ctx context.Context, creator string) []types.Intent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var intents []types.Intent
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		if intent.Creator == creator {
			intents = append(intents, intent)
		}
	}

	sort.Slice(intents, func(i, j int) bool {
		return intents[i].ID < intents[j].ID
	})

	return intents
}

// GetAllSolvers returns all registered solvers.
func (k Keeper) GetAllSolvers(ctx context.Context) []types.Solver {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.SolverPrefixKey)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var solvers []types.Solver
	for ; iter.Valid(); iter.Next() {
		var solver types.Solver
		if err := json.Unmarshal(iter.Value(), &solver); err != nil {
			continue
		}
		solvers = append(solvers, solver)
	}

	sort.Slice(solvers, func(i, j int) bool {
		return solvers[i].Address < solvers[j].Address
	})

	return solvers
}

// getIntentsByStatus returns all intents matching the given status
func (k Keeper) getIntentsByStatus(ctx context.Context, status types.IntentStatus) []types.Intent {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return nil
	}
	defer iter.Close()

	var intents []types.Intent
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		if intent.Status == status {
			intents = append(intents, intent)
		}
	}

	// Explicit deterministic sort by intent ID to guarantee identical
	// results regardless of KV store backend iteration order.
	sort.Slice(intents, func(i, j int) bool {
		return intents[i].ID < intents[j].ID
	})

	return intents
}

// ExpireIntents finds expired intents and refunds creators
func (k Keeper) ExpireIntents(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.IntentPrefix)

	iter, err := kvStore.Iterator(prefix, nil)
	if err != nil {
		return
	}
	defer iter.Close()

	var expiredIntents []types.Intent
	var pruneIntentIDs []string
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		var intent types.Intent
		if err := json.Unmarshal(iter.Value(), &intent); err != nil {
			continue
		}
		if (intent.Status == types.StatusPending || intent.Status == types.StatusSolving) && sdkCtx.BlockHeight() > intent.Expiry {
			expiredIntents = append(expiredIntents, intent)
			continue
		}
		// LIVENESS: terminal intents (fulfilled/failed/expired) are never otherwise
		// deleted, so every BeginBlock re-scans and re-unmarshals the entire intent
		// history — an unbounded cost an attacker can grow with cheap zero-fee
		// intents until block production stalls. Prune terminal intents once they
		// are older than the retention window so the working set stays bounded.
		if intent.IsTerminal() && sdkCtx.BlockHeight()-intent.CreatedAt > types.IntentPruneRetentionBlocks {
			pruneIntentIDs = append(pruneIntentIDs, intent.ID)
		}
	}

	for _, id := range pruneIntentIDs {
		kvStore.Delete(types.IntentKey(id))
		k.deleteSolutionsForIntent(ctx, id)
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"intent_pruned",
			sdk.NewAttribute("intent_id", id),
		))
	}

	for _, intent := range expiredIntents {
		// Refund locked funds to creator
		creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
		if err != nil {
			// Still mark as expired to keep state deterministic across validators.
			k.Logger(ctx).Error("invalid creator address in expired intent", "intent_id", intent.ID, "error", err)
		} else {
			totalLock := intent.MaxFee.Add(intent.Tip...)
			if totalLock.IsAllPositive() {
				if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalLock); err != nil {
					k.Logger(ctx).Error("failed to refund expired intent", "intent_id", intent.ID, "error", err)
					// Fall through: still expire intent for deterministic state.
				}
			}

			// Refund locked trading input tokens for trading intents
			k.refundTradingTokens(ctx, intent)
		}

		// Always mark as expired and clean up — deterministic across all validators.
		intent.Status = types.StatusExpired
		k.SetIntent(ctx, intent)

		// M-04: Delete orphaned solutions for expired intents to prevent state bloat
		k.deleteSolutionsForIntent(ctx, intent.ID)

		sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
			"intent_expired",
			sdk.NewAttribute("intent_id", intent.ID),
			sdk.NewAttribute("creator", intent.Creator),
		))

		k.Logger(ctx).Info("intent expired", "id", intent.ID, "creator", intent.Creator)
	}
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx context.Context) types.Params {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ParamsKey))
	if err != nil || bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}
	return params
}

// SetParams stores the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return kvStore.Set([]byte(types.ParamsKey), bz)
}

// GetSolvingIntents returns all intents with StatusSolving
func (k Keeper) GetSolvingIntents(ctx context.Context) []types.Intent {
	return k.getIntentsByStatus(ctx, types.StatusSolving)
}

// AutoFulfillIntents iterates intents in "solving" state whose solving window has expired
// and auto-fulfills them using the auction. If no solutions exist, the intent is expired and refunded.
func (k Keeper) AutoFulfillIntents(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := k.GetParams(ctx)

	solvingIntents := k.GetSolvingIntents(ctx)
	for _, intent := range solvingIntents {
		// Check if the solving window has expired
		solvingDeadline := intent.CreatedAt + int64(params.SolvingWindow)
		if sdkCtx.BlockHeight() <= solvingDeadline {
			continue // Still within solving window, skip
		}

		// Solving window expired -- try to auto-fulfill via auction
		solutions := k.GetSolutionsForIntent(ctx, intent.ID)
		if len(solutions) == 0 {
			// No solutions: expire the intent and refund
			k.expireAndRefundIntent(ctx, intent)
			continue
		}

		// Auto-select and fulfill
		err := k.FulfillIntent(ctx, &types.MsgFulfillIntent{
			SolverAddr: "", // empty triggers auto-selection
			IntentID:   intent.ID,
		})
		if err != nil {
			k.Logger(ctx).Error("auto-fulfill failed for intent",
				"intent_id", intent.ID, "error", err)
			// If fulfillment failed (e.g., execution error), the intent is already marked failed
			// by FulfillIntent, so no further action needed.
		}
	}
}

// expireAndRefundIntent marks an intent as expired and refunds locked funds to creator.
func (k Keeper) expireAndRefundIntent(ctx context.Context, intent types.Intent) {
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		k.Logger(ctx).Error("invalid creator address during auto-expire", "intent_id", intent.ID, "error", err)
		return
	}

	totalLock := intent.MaxFee.Add(intent.Tip...)
	if totalLock.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalLock); err != nil {
			k.Logger(ctx).Error("failed to refund auto-expired intent", "intent_id", intent.ID, "error", err)
			// Fall through to mark intent as expired anyway to prevent infinite retry
		}
	}

	// Refund locked trading input tokens for trading intents
	k.refundTradingTokens(ctx, intent)

	intent.Status = types.StatusExpired
	k.SetIntent(ctx, intent)
	k.Logger(ctx).Info("intent auto-expired and refunded (solving window elapsed, no fulfillment)", "id", intent.ID)
}

// CancelIntent lets the intent creator cancel a still-active (Pending/Solving)
// intent and reclaim all locked funds (MaxFee + Tip and any locked trading
// input), mirroring the refund path used by ExpireIntents/expireAndRefundIntent.
// Only the creator may cancel. Chain-linked intents cannot be cancelled directly
// (their funds are escrowed at the chain level) — cancel the chain instead.
func (k *Keeper) CancelIntent(ctx context.Context, msg *types.MsgCancelIntent) error {
	intent, found := k.GetIntent(ctx, msg.IntentID)
	if !found {
		return types.ErrIntentNotFound
	}

	if intent.Creator != msg.Creator {
		return types.ErrIntentNotCreator
	}

	if intent.Status != types.StatusPending && intent.Status != types.StatusSolving {
		return types.ErrIntentNotCancellable
	}

	// A chain step's tokens are locked/refunded at the chain level; cancelling the
	// step intent directly would double-refund. Require chain cancellation instead.
	if k.IsChainLinkedIntent(ctx, msg.IntentID) {
		return types.ErrIntentNotCancellable
	}

	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Refund locked MaxFee + Tip (same as the expiry refund path).
	totalLock := intent.MaxFee.Add(intent.Tip...)
	if totalLock.IsAllPositive() {
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, totalLock); err != nil {
			return fmt.Errorf("failed to refund locked fees: %w", err)
		}
	}

	// Refund locked trading input tokens for trading intents.
	k.refundTradingTokens(ctx, intent)

	// Clean up any submitted solutions.
	k.deleteSolutionsForIntent(ctx, msg.IntentID)

	// Move to a terminal state. Reuse StatusExpired so existing terminal-state
	// checks (FulfillIntent guard, IsTerminal pruning) treat it correctly.
	intent.Status = types.StatusExpired
	k.SetIntent(ctx, intent)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		"cancel_intent",
		sdk.NewAttribute("creator", intent.Creator),
		sdk.NewAttribute("intent_id", intent.ID),
	))

	k.Logger(ctx).Info("intent cancelled by creator", "id", intent.ID, "creator", intent.Creator)
	return nil
}

// nextIntentID generates the next sequential intent ID
func (k Keeper) nextIntentID(ctx context.Context) string {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(types.IntentCounterKey())
	var counter uint64
	if err == nil && bz != nil {
		counter = types.BytesToUint64(bz)
	}
	counter++
	kvStore.Set(types.IntentCounterKey(), types.Uint64ToBytes(counter))
	return strconv.FormatUint(counter, 10)
}

// swapOutputBalanceBefore snapshots the creator's balance of the swap intent's
// output denom prior to executing the solution messages. The returned value is
// compared against the post-execution balance in verifySwapOutcome so that the
// check measures the DELTA delivered by the solver rather than the creator's
// absolute holdings (which could already exceed MinOutputAmount).
func (k *Keeper) swapOutputBalanceBefore(ctx context.Context, intent types.Intent) math.Int {
	var swapBody types.SwapIntent
	if err := json.Unmarshal(intent.Body, &swapBody); err != nil {
		// Non-swap body format — verification will be skipped, return zero.
		return math.ZeroInt()
	}
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return math.ZeroInt()
	}
	return k.bankKeeper.GetBalance(ctx, creatorAddr, swapBody.OutputDenom).Amount
}

// verifySwapOutcome checks that a swap intent's solution actually produced
// the expected outcome — the creator must have RECEIVED at least MinOutputAmount
// of the output denom as a result of the solution executing. It compares the
// post-execution balance against the pre-execution snapshot (beforeBalance) and
// requires (after - before) >= MinOutputAmount. This prevents solvers from
// submitting solutions that appear to succeed but don't actually deliver the
// expected tokens, including the case where the creator already held >=
// MinOutputAmount before the solution ran.
func (k *Keeper) verifySwapOutcome(ctx context.Context, intent types.Intent, solverAddr string, params types.Params, beforeBalance math.Int, declaredOutput math.Int) error {
	var swapBody types.SwapIntent
	if err := json.Unmarshal(intent.Body, &swapBody); err != nil {
		// If we can't parse the body, skip verification (non-swap format)
		return nil
	}

	// MinOutputAmount can be nil if the swap body omitted it. A nil math.Int
	// wraps a nil big.Int, so any comparison (.GT/.LT) or String() below would
	// panic — and this runs inside AutoFulfillIntents in BeginBlock, so a panic
	// would HALT THE CHAIN. Normalize nil to zero ("no minimum floor").
	if swapBody.MinOutputAmount.IsNil() {
		swapBody.MinOutputAmount = math.ZeroInt()
	}

	// Check the creator's balance of the output denom
	creatorAddr, err := sdk.AccAddressFromBech32(intent.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	afterBalance := k.bankKeeper.GetBalance(ctx, creatorAddr, swapBody.OutputDenom).Amount

	// Fairness Engine · best execution: the creator must receive at least the
	// intent's MinOutputAmount AND at least what the winning solver DECLARED in
	// the auction. Binding the solver to their declaration is what stops them
	// winning the best-execution auction by over-promising and delivering only
	// the floor — under-delivery fails here and slashes the solver.
	required := swapBody.MinOutputAmount
	if !declaredOutput.IsNil() && declaredOutput.GT(required) {
		required = declaredOutput
	}

	// The creator should have received at least `required` as a delta.
	delta := afterBalance.Sub(beforeBalance)
	if delta.LT(required) {
		return fmt.Errorf("creator received %s %s, less than required %s (min %s, solver promised %s)",
			delta, swapBody.OutputDenom, required, swapBody.MinOutputAmount, declaredOutput)
	}

	return nil
}
