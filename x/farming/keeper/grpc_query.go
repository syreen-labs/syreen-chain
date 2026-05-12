package keeper

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterRESTRoutes registers REST query routes for the farming module
func (k Keeper) RegisterRESTRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/syreen/farming/v1/farms", k.handleGetFarms)
	mux.HandleFunc("/syreen/farming/v1/farm/", k.handleGetFarm)
	mux.HandleFunc("/syreen/farming/v1/positions/", k.handleGetPositions)
	mux.HandleFunc("/syreen/farming/v1/pending/", k.handleGetPendingReward)
}

type FarmResponse struct {
	PoolID         uint64 `json:"pool_id"`
	LPDenom        string `json:"lp_denom"`
	RewardDenom    string `json:"reward_denom"`
	RewardPerBlock string `json:"reward_per_block"`
	TotalStaked    string `json:"total_staked"`
	APR            string `json:"apr"`
	Active         bool   `json:"active"`
	StartBlock     int64  `json:"start_block"`
	EndBlock       int64  `json:"end_block"`
}

type PositionResponse struct {
	Address       string `json:"address"`
	PoolID        uint64 `json:"pool_id"`
	Amount        string `json:"amount"`
	PendingReward string `json:"pending_reward"`
	StakedAt      int64  `json:"staked_at"`
}

func (k Keeper) handleGetFarms(w http.ResponseWriter, r *http.Request) {
	sdkCtx := r.Context().Value("sdkCtx")
	if sdkCtx == nil {
		http.Error(w, "no context", 500)
		return
	}
	ctx := sdkCtx.(sdk.Context)

	farms := k.GetAllFarms(ctx)
	var resp []FarmResponse
	for _, f := range farms {
		resp = append(resp, FarmResponse{
			PoolID:         f.PoolID,
			LPDenom:        f.LPDenom,
			RewardDenom:    f.RewardDenom,
			RewardPerBlock: f.RewardPerBlock.String(),
			TotalStaked:    f.TotalStaked.String(),
			APR:            k.GetFarmAPR(ctx, f),
			Active:         f.Active,
			StartBlock:     f.StartBlock,
			EndBlock:       f.EndBlock,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "farms": resp})
}

func (k Keeper) handleGetFarm(w http.ResponseWriter, r *http.Request) {
	sdkCtx := r.Context().Value("sdkCtx")
	if sdkCtx == nil {
		http.Error(w, "no context", 500)
		return
	}
	ctx := sdkCtx.(sdk.Context)

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/syreen/farming/v1/farm/"), "/")
	if len(parts) < 1 {
		http.Error(w, "pool_id required", 400)
		return
	}
	poolID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid pool_id", 400)
		return
	}

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "farm not found"})
		return
	}

	resp := FarmResponse{
		PoolID:         farm.PoolID,
		LPDenom:        farm.LPDenom,
		RewardDenom:    farm.RewardDenom,
		RewardPerBlock: farm.RewardPerBlock.String(),
		TotalStaked:    farm.TotalStaked.String(),
		APR:            k.GetFarmAPR(ctx, farm),
		Active:         farm.Active,
		StartBlock:     farm.StartBlock,
		EndBlock:       farm.EndBlock,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": resp})
}

func (k Keeper) handleGetPositions(w http.ResponseWriter, r *http.Request) {
	sdkCtx := r.Context().Value("sdkCtx")
	if sdkCtx == nil {
		http.Error(w, "no context", 500)
		return
	}
	ctx := sdkCtx.(sdk.Context)

	address := strings.TrimPrefix(r.URL.Path, "/syreen/farming/v1/positions/")
	if address == "" {
		http.Error(w, "address required", 400)
		return
	}

	positions := k.GetPositionsByAddress(ctx, address)
	var resp []PositionResponse
	for _, p := range positions {
		pending := p.PendingReward
		if farm, ok := k.GetFarm(ctx, p.PoolID); ok {
			pending = k.CalculatePendingReward(farm, p)
		}
		resp = append(resp, PositionResponse{
			Address:       p.Address,
			PoolID:        p.PoolID,
			Amount:        p.Amount.String(),
			PendingReward: pending.String(),
			StakedAt:      p.StakedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "positions": resp})
}

func (k Keeper) handleGetPendingReward(w http.ResponseWriter, r *http.Request) {
	sdkCtx := r.Context().Value("sdkCtx")
	if sdkCtx == nil {
		http.Error(w, "no context", 500)
		return
	}
	ctx := sdkCtx.(sdk.Context)

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/syreen/farming/v1/pending/"), "/")
	if len(parts) < 2 {
		http.Error(w, "pool_id/address required", 400)
		return
	}
	poolID, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid pool_id", 400)
		return
	}
	address := parts[1]

	farm, ok := k.GetFarm(ctx, poolID)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "farm not found"})
		return
	}

	pos, exists := k.GetPosition(ctx, poolID, address)
	pending := math.ZeroInt()
	staked := math.ZeroInt()
	if exists {
		pending = k.CalculatePendingReward(farm, pos)
		staked = pos.Amount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"pending_reward": pending.String(),
		"staked":         staked.String(),
	})
}
