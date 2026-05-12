package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"syreen/x/dex/types"
)

// ---------------------------------------------------------------------------
// Referral data types (stored as JSON in KV store)
// ---------------------------------------------------------------------------

// ReferrerStats tracks per-referrer statistics.
type ReferrerStats struct {
	Address        string   `json:"address"`
	ReferralCode   string   `json:"referral_code"`
	TotalReferrals int64    `json:"total_referrals"`
	TotalVolume    math.Int `json:"total_volume"`
	TotalEarned    math.Int `json:"total_earned"`
	Tier           string   `json:"tier"`
}

// ReferredUserStats tracks per-referred-user statistics.
type ReferredUserStats struct {
	Address     string   `json:"address"`
	Referrer    string   `json:"referrer"`
	TotalVolume math.Int `json:"total_volume"`
	TotalSaved  math.Int `json:"total_saved"`
}

// ReferralGlobalStats tracks global referral statistics.
type ReferralGlobalStats struct {
	TotalReferrals      int64    `json:"total_referrals"`
	TotalFeesDistributed math.Int `json:"total_fees_distributed"`
}

// ReferralEarnings stores a list of earning entries for a referrer.
type ReferralEarnings struct {
	Entries []ReferralEarningEntry `json:"entries"`
}

// ReferralEarningEntry is a single earning event.
type ReferralEarningEntry struct {
	Height int64    `json:"height"`
	Trader string   `json:"trader"`
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
	PoolID uint64   `json:"pool_id"`
}

// ---------------------------------------------------------------------------
// Tier definitions
// ---------------------------------------------------------------------------

const (
	TierBronze   = "bronze"
	TierSilver   = "silver"
	TierGold     = "gold"
	TierPlatinum = "platinum"
)

// GetTier returns the tier name based on total referral count.
func GetTier(totalReferrals int64) string {
	switch {
	case totalReferrals > 200:
		return TierPlatinum
	case totalReferrals > 50:
		return TierGold
	case totalReferrals > 10:
		return TierSilver
	default:
		return TierBronze
	}
}

// GetTierReferrerPctBps returns the referrer fee share in basis points of the total fee.
// e.g., 500 = 5% of the 0.3% fee.
func GetTierReferrerPctBps(tier string) int64 {
	switch tier {
	case TierPlatinum:
		return 1200 // 12%
	case TierGold:
		return 1000 // 10%
	case TierSilver:
		return 700 // 7%
	default:
		return 500 // 5%
	}
}

// GetTierDiscountPctBps returns the trader discount in basis points of the total fee.
func GetTierDiscountPctBps(tier string) int64 {
	switch tier {
	case TierPlatinum:
		return 1000 // 10%
	case TierGold:
		return 700 // 7%
	default:
		return 500 // 5%
	}
}

// ---------------------------------------------------------------------------
// Referral Code Generation
// ---------------------------------------------------------------------------

// GenerateReferralCode derives a deterministic referral code from an address.
// Uses first 8 chars of SHA-256 hash of the address, uppercase.
func GenerateReferralCode(address string) string {
	hash := sha256.Sum256([]byte(address))
	return strings.ToUpper(hex.EncodeToString(hash[:4]))
}

// ---------------------------------------------------------------------------
// KV Store: Referral Code <-> Address mapping
// ---------------------------------------------------------------------------

// SetReferralCode stores the mapping from referral code to referrer address,
// and from address to referral code.
func (k Keeper) SetReferralCode(ctx context.Context, code, address string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set([]byte(types.ReferralCodePrefix+code), []byte(address))
	kvStore.Set([]byte(types.ReferralByAddrPrefix+address), []byte(code))
}

// GetReferrerByCode returns the referrer address for a given referral code.
func (k Keeper) GetReferrerByCode(ctx context.Context, code string) (string, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralCodePrefix + code))
	if err != nil || bz == nil {
		return "", false
	}
	return string(bz), true
}

// GetReferralCodeByAddress returns the referral code for a given address.
func (k Keeper) GetReferralCodeByAddress(ctx context.Context, address string) (string, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralByAddrPrefix + address))
	if err != nil || bz == nil {
		return "", false
	}
	return string(bz), true
}

// ---------------------------------------------------------------------------
// KV Store: Referral Registration (who referred whom)
// ---------------------------------------------------------------------------

// SetReferralRegistration records that `user` was referred by `referrer`.
func (k Keeper) SetReferralRegistration(ctx context.Context, user, referrer string) {
	kvStore := k.storeService.OpenKVStore(ctx)
	kvStore.Set([]byte(types.ReferralRegistrationPrefix+user), []byte(referrer))
}

// GetReferrer returns the referrer for a given user, if any.
func (k Keeper) GetReferrer(ctx context.Context, user string) (string, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralRegistrationPrefix + user))
	if err != nil || bz == nil {
		return "", false
	}
	return string(bz), true
}

// ---------------------------------------------------------------------------
// KV Store: Referrer Stats
// ---------------------------------------------------------------------------

func (k Keeper) GetReferrerStats(ctx context.Context, address string) (ReferrerStats, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralStatsPrefix + address))
	if err != nil || bz == nil {
		return ReferrerStats{}, false
	}
	var stats ReferrerStats
	if err := json.Unmarshal(bz, &stats); err != nil {
		return ReferrerStats{}, false
	}
	return stats, true
}

func (k Keeper) SetReferrerStats(ctx context.Context, stats ReferrerStats) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(stats)
	kvStore.Set([]byte(types.ReferralStatsPrefix+stats.Address), bz)
}

// ---------------------------------------------------------------------------
// KV Store: Referred User Stats
// ---------------------------------------------------------------------------

func (k Keeper) GetReferredUserStats(ctx context.Context, address string) (ReferredUserStats, bool) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralUserStatsPrefix + address))
	if err != nil || bz == nil {
		return ReferredUserStats{}, false
	}
	var stats ReferredUserStats
	if err := json.Unmarshal(bz, &stats); err != nil {
		return ReferredUserStats{}, false
	}
	return stats, true
}

func (k Keeper) SetReferredUserStats(ctx context.Context, stats ReferredUserStats) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(stats)
	kvStore.Set([]byte(types.ReferralUserStatsPrefix+stats.Address), bz)
}

// ---------------------------------------------------------------------------
// KV Store: Global Stats
// ---------------------------------------------------------------------------

func (k Keeper) GetReferralGlobalStats(ctx context.Context) ReferralGlobalStats {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralGlobalPrefix))
	if err != nil || bz == nil {
		return ReferralGlobalStats{
			TotalReferrals:       0,
			TotalFeesDistributed: math.ZeroInt(),
		}
	}
	var stats ReferralGlobalStats
	if err := json.Unmarshal(bz, &stats); err != nil {
		return ReferralGlobalStats{
			TotalReferrals:       0,
			TotalFeesDistributed: math.ZeroInt(),
		}
	}
	return stats
}

func (k Keeper) SetReferralGlobalStats(ctx context.Context, stats ReferralGlobalStats) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(stats)
	kvStore.Set([]byte(types.ReferralGlobalPrefix), bz)
}

// ---------------------------------------------------------------------------
// KV Store: Earnings History
// ---------------------------------------------------------------------------

func (k Keeper) GetReferralEarnings(ctx context.Context, address string) ReferralEarnings {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get([]byte(types.ReferralEarningsPrefix + address))
	if err != nil || bz == nil {
		return ReferralEarnings{Entries: []ReferralEarningEntry{}}
	}
	var earnings ReferralEarnings
	if err := json.Unmarshal(bz, &earnings); err != nil {
		return ReferralEarnings{Entries: []ReferralEarningEntry{}}
	}
	return earnings
}

func (k Keeper) SetReferralEarnings(ctx context.Context, address string, earnings ReferralEarnings) {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(earnings)
	kvStore.Set([]byte(types.ReferralEarningsPrefix+address), bz)
}

// ---------------------------------------------------------------------------
// Core Logic: CreateReferralCode
// ---------------------------------------------------------------------------

// CreateReferralCode generates a referral code for the given address.
// If customCode is non-empty, it is used instead of the auto-generated code.
func (k Keeper) CreateReferralCode(ctx context.Context, creator string, customCode string) (string, error) {
	_, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return "", types.ErrInvalidSender
	}

	// Check if user already has a code
	if _, exists := k.GetReferralCodeByAddress(ctx, creator); exists {
		return "", types.ErrReferralCodeExists
	}

	var code string
	if customCode != "" {
		// Validate custom code: alphanumeric, 3-20 chars
		if len(customCode) < 3 || len(customCode) > 20 {
			return "", fmt.Errorf("custom referral code must be 3-20 characters")
		}
		code = customCode
		// Check if custom code is already taken
		if _, exists := k.GetReferrerByCode(ctx, code); exists {
			return "", fmt.Errorf("referral code %s is already taken", code)
		}
	} else {
		code = GenerateReferralCode(creator)
	}

	// Handle collisions: if the code already exists, append a counter
	baseCode := code
	for i := 1; i <= 100; i++ {
		if _, exists := k.GetReferrerByCode(ctx, code); !exists {
			break // unique code found
		}
		code = baseCode + fmt.Sprintf("%02d", i)
	}
	// Final collision check after exhausting retries
	if _, exists := k.GetReferrerByCode(ctx, code); exists {
		return "", fmt.Errorf("unable to generate unique referral code")
	}

	// Store the mapping
	k.SetReferralCode(ctx, code, creator)

	// Initialize referrer stats
	k.SetReferrerStats(ctx, ReferrerStats{
		Address:        creator,
		ReferralCode:   code,
		TotalReferrals: 0,
		TotalVolume:    math.ZeroInt(),
		TotalEarned:    math.ZeroInt(),
		Tier:           TierBronze,
	})

	k.Logger(ctx).Info("referral code created",
		"creator", creator,
		"code", code,
	)

	return code, nil
}

// ---------------------------------------------------------------------------
// Core Logic: RegisterReferral
// ---------------------------------------------------------------------------

// RegisterReferral registers a user under a referral code.
func (k Keeper) RegisterReferral(ctx context.Context, user, referralCode string) (string, error) {
	_, err := sdk.AccAddressFromBech32(user)
	if err != nil {
		return "", types.ErrInvalidSender
	}

	// Check if user already has a referrer
	if _, exists := k.GetReferrer(ctx, user); exists {
		return "", types.ErrAlreadyReferred
	}

	// Look up referral code
	referrer, found := k.GetReferrerByCode(ctx, referralCode)
	if !found {
		return "", types.ErrReferralCodeNotFound
	}

	// Self-referral check
	if user == referrer {
		return "", types.ErrSelfReferral
	}

	// Circular referral check: walk up the referral chain up to 10 levels
	current := referrer
	for i := 0; i < 10; i++ {
		ancestor, exists := k.GetReferrer(ctx, current)
		if !exists {
			break
		}
		if ancestor == user {
			return "", types.ErrCircularReferral
		}
		current = ancestor
	}

	// Register the referral
	k.SetReferralRegistration(ctx, user, referrer)

	// Initialize referred user stats
	k.SetReferredUserStats(ctx, ReferredUserStats{
		Address:     user,
		Referrer:    referrer,
		TotalVolume: math.ZeroInt(),
		TotalSaved:  math.ZeroInt(),
	})

	// Update referrer stats
	stats, _ := k.GetReferrerStats(ctx, referrer)
	stats.TotalReferrals++
	stats.Tier = GetTier(stats.TotalReferrals)
	k.SetReferrerStats(ctx, stats)

	// Update global stats
	global := k.GetReferralGlobalStats(ctx)
	global.TotalReferrals++
	k.SetReferralGlobalStats(ctx, global)

	k.Logger(ctx).Info("referral registered",
		"user", user,
		"referrer", referrer,
		"code", referralCode,
	)

	return referrer, nil
}

// ---------------------------------------------------------------------------
// Fee Sharing: Applied during Swap
// ---------------------------------------------------------------------------

// ReferralFeeResult holds the calculated fee distribution for a referral swap.
type ReferralFeeResult struct {
	HasReferral      bool
	ReferrerAddr     sdk.AccAddress
	ReferrerAmount   math.Int // amount to send to referrer (in input token)
	TraderDiscount   math.Int // amount to return to trader (in input token)
	ReferrerAddrStr  string
}

// CalculateReferralFees calculates the referral fee split for a swap.
// Returns the amounts to send to referrer and discount for trader.
// The amounts are taken from the fee portion (feeBps applied to tokenIn).
func (k Keeper) CalculateReferralFees(ctx context.Context, trader string, tokenInAmount math.Int, feeBps int64) ReferralFeeResult {
	result := ReferralFeeResult{HasReferral: false}

	// Check if trader has a referrer
	referrer, exists := k.GetReferrer(ctx, trader)
	if !exists {
		return result
	}

	referrerAddr, err := sdk.AccAddressFromBech32(referrer)
	if err != nil {
		return result
	}

	// Get referrer stats for tier
	stats, found := k.GetReferrerStats(ctx, referrer)
	if !found {
		return result
	}

	// Calculate total fee amount: tokenInAmount * feeBps / 10000
	totalFee := tokenInAmount.MulRaw(feeBps).QuoRaw(10000)
	if totalFee.IsZero() {
		return result
	}

	tier := stats.Tier
	referrerPctBps := GetTierReferrerPctBps(tier)
	discountPctBps := GetTierDiscountPctBps(tier)

	// referrerAmount = totalFee * referrerPctBps / 10000
	referrerAmount := totalFee.MulRaw(referrerPctBps).QuoRaw(10000)
	// discountAmount = totalFee * discountPctBps / 10000
	discountAmount := totalFee.MulRaw(discountPctBps).QuoRaw(10000)

	result.HasReferral = true
	result.ReferrerAddr = referrerAddr
	result.ReferrerAmount = referrerAmount
	result.TraderDiscount = discountAmount
	result.ReferrerAddrStr = referrer

	return result
}

// ApplyReferralFees processes the referral fee distribution during a swap.
//
// H5: The referrer portion is NOT sent directly to the referrer account. It is
// credited to a claimable balance held in the dex module account, which the
// referrer must explicitly withdraw via MsgClaimReferralRewards. This prevents
// a referrer with a blocked / non-existent / module account from bricking every
// downstream swap (an attacker could otherwise register a poisoned referrer and
// halt trading for everyone they refer).
//
// The trader discount is still sent immediately because the trader is the
// transaction sender — their account is guaranteed to exist and accept funds.
func (k Keeper) ApplyReferralFees(ctx context.Context, trader string, poolID uint64, tokenInDenom string, feeResult ReferralFeeResult, swapVolume ...math.Int) error {
	if !feeResult.HasReferral {
		return nil
	}

	traderAddr, err := sdk.AccAddressFromBech32(trader)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// H5: Credit referrer portion to a claimable balance instead of bank-sending.
	if feeResult.ReferrerAmount.IsPositive() {
		k.AddClaimableReferralRewards(ctx, feeResult.ReferrerAddrStr, tokenInDenom, feeResult.ReferrerAmount)
	}

	// Send discount back to trader from module
	if feeResult.TraderDiscount.IsPositive() {
		discountCoins := sdk.NewCoins(sdk.NewCoin(tokenInDenom, feeResult.TraderDiscount))
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, traderAddr, discountCoins); err != nil {
			return fmt.Errorf("failed to send trader discount: %w", err)
		}
	}

	// Determine actual swap volume for stats tracking. If the caller passed
	// the real swap volume use that; otherwise fall back to the fee amounts
	// (legacy behavior, should not happen in practice).
	actualVolume := feeResult.ReferrerAmount.Add(feeResult.TraderDiscount)
	if len(swapVolume) > 0 && !swapVolume[0].IsNil() && swapVolume[0].IsPositive() {
		actualVolume = swapVolume[0]
	}

	// Update referrer stats
	stats, _ := k.GetReferrerStats(ctx, feeResult.ReferrerAddrStr)
	stats.TotalVolume = stats.TotalVolume.Add(actualVolume)
	stats.TotalEarned = stats.TotalEarned.Add(feeResult.ReferrerAmount)
	k.SetReferrerStats(ctx, stats)

	// Update referred user stats
	userStats, _ := k.GetReferredUserStats(ctx, trader)
	userStats.TotalVolume = userStats.TotalVolume.Add(actualVolume)
	userStats.TotalSaved = userStats.TotalSaved.Add(feeResult.TraderDiscount)
	k.SetReferredUserStats(ctx, userStats)

	// Update global stats
	global := k.GetReferralGlobalStats(ctx)
	global.TotalFeesDistributed = global.TotalFeesDistributed.Add(feeResult.ReferrerAmount)
	k.SetReferralGlobalStats(ctx, global)

	// Record earning entry (keep last 100 entries max)
	earnings := k.GetReferralEarnings(ctx, feeResult.ReferrerAddrStr)
	entry := ReferralEarningEntry{
		Height: sdkCtx.BlockHeight(),
		Trader: trader,
		Denom:  tokenInDenom,
		Amount: feeResult.ReferrerAmount,
		PoolID: poolID,
	}
	earnings.Entries = append(earnings.Entries, entry)
	if len(earnings.Entries) > 100 {
		earnings.Entries = earnings.Entries[len(earnings.Entries)-100:]
	}
	k.SetReferralEarnings(ctx, feeResult.ReferrerAddrStr, earnings)

	k.Logger(ctx).Info("referral fees applied",
		"trader", trader,
		"referrer", feeResult.ReferrerAddrStr,
		"referrer_amount", feeResult.ReferrerAmount,
		"trader_discount", feeResult.TraderDiscount,
		"denom", tokenInDenom,
		"pool_id", poolID,
	)

	return nil
}

// ---------------------------------------------------------------------------
// H5: Claimable referral rewards (withdrawal pattern)
// ---------------------------------------------------------------------------

// claimableReferralKey returns the KV store key for a referrer/denom pair.
func claimableReferralKey(referrer, denom string) []byte {
	return []byte(types.ReferralClaimablePrefix + referrer + "/" + denom)
}

// claimableReferralPrefix returns the prefix for all claimable balances of a
// given referrer (across all denoms).
func claimableReferralPrefix(referrer string) []byte {
	return []byte(types.ReferralClaimablePrefix + referrer + "/")
}

// GetClaimableReferralRewards returns the currently claimable amount for a
// referrer/denom pair (zero if none).
func (k Keeper) GetClaimableReferralRewards(ctx context.Context, referrer, denom string) math.Int {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz, err := kvStore.Get(claimableReferralKey(referrer, denom))
	if err != nil || bz == nil {
		return math.ZeroInt()
	}
	var amt math.Int
	if err := amt.Unmarshal(bz); err != nil {
		return math.ZeroInt()
	}
	return amt
}

// setClaimableReferralRewards stores the claimable amount, deleting the entry
// if it drops to zero.
func (k Keeper) setClaimableReferralRewards(ctx context.Context, referrer, denom string, amount math.Int) {
	kvStore := k.storeService.OpenKVStore(ctx)
	key := claimableReferralKey(referrer, denom)
	if amount.IsNil() || !amount.IsPositive() {
		kvStore.Delete(key)
		return
	}
	bz, err := amount.Marshal()
	if err != nil {
		return
	}
	kvStore.Set(key, bz)
}

// AddClaimableReferralRewards credits `amount` of `denom` to a referrer's
// claimable balance. Funds are expected to already be held in the dex module
// account (they were never bank-sent out during the originating swap).
func (k Keeper) AddClaimableReferralRewards(ctx context.Context, referrer, denom string, amount math.Int) {
	if amount.IsNil() || !amount.IsPositive() {
		return
	}
	current := k.GetClaimableReferralRewards(ctx, referrer, denom)
	k.setClaimableReferralRewards(ctx, referrer, denom, current.Add(amount))
}

// ClaimableBalance is a single denom/amount entry returned to callers.
type ClaimableBalance struct {
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
}

// GetAllClaimableReferralRewards returns every (denom, amount) entry the
// referrer can currently claim.
func (k Keeper) GetAllClaimableReferralRewards(ctx context.Context, referrer string) []ClaimableBalance {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := claimableReferralPrefix(referrer)
	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var out []ClaimableBalance
	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key())
		denom := key[len(prefix):]
		var amt math.Int
		if err := amt.Unmarshal(iter.Value()); err != nil {
			continue
		}
		if amt.IsPositive() {
			out = append(out, ClaimableBalance{Denom: denom, Amount: amt})
		}
	}
	return out
}

// ClaimReferralRewards withdraws all claimable referral rewards for the given
// referrer (across every denom) from the dex module account to their address.
// Returns the list of coins paid out.
func (k Keeper) ClaimReferralRewards(ctx context.Context, referrer string) (sdk.Coins, error) {
	addr, err := sdk.AccAddressFromBech32(referrer)
	if err != nil {
		return nil, types.ErrInvalidSender
	}

	balances := k.GetAllClaimableReferralRewards(ctx, referrer)
	if len(balances) == 0 {
		return sdk.Coins{}, fmt.Errorf("no claimable referral rewards")
	}

	coins := sdk.NewCoins()
	for _, b := range balances {
		coins = coins.Add(sdk.NewCoin(b.Denom, b.Amount))
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
		return nil, fmt.Errorf("failed to pay claimable referral rewards: %w", err)
	}

	// Zero-out all entries on success
	for _, b := range balances {
		k.setClaimableReferralRewards(ctx, referrer, b.Denom, math.ZeroInt())
	}

	k.Logger(ctx).Info("referral rewards claimed",
		"referrer", referrer,
		"coins", coins,
	)

	return coins, nil
}

// ---------------------------------------------------------------------------
// GetAllReferredUsers returns all users referred by a given referrer.
// Uses prefix scan on referral_reg/ to find all registrations, then filters.
// ---------------------------------------------------------------------------

func (k Keeper) GetAllReferredUsers(ctx context.Context, referrer string) []ReferredUserStats {
	kvStore := k.storeService.OpenKVStore(ctx)
	prefix := []byte(types.ReferralRegistrationPrefix)

	iter, err := kvStore.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil
	}
	defer iter.Close()

	var results []ReferredUserStats
	for ; iter.Valid(); iter.Next() {
		storedReferrer := string(iter.Value())
		if storedReferrer != referrer {
			continue
		}
		// Extract user address from key
		key := string(iter.Key())
		userAddr := key[len(types.ReferralRegistrationPrefix):]

		userStats, found := k.GetReferredUserStats(ctx, userAddr)
		if !found {
			userStats = ReferredUserStats{
				Address:     userAddr,
				Referrer:    referrer,
				TotalVolume: math.ZeroInt(),
				TotalSaved:  math.ZeroInt(),
			}
		}
		results = append(results, userStats)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Address < results[j].Address
	})
	return results
}
