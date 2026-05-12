package types

import (
	"bytes"
	"compress/gzip"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

var fileDescriptorTx []byte

func init() {
	registerProtoFileDescriptors()
}

func registerProtoFileDescriptors() {
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	typeString := descriptorpb.FieldDescriptorProto_TYPE_STRING
	typeBytes := descriptorpb.FieldDescriptorProto_TYPE_BYTES
	typeUint64 := descriptorpb.FieldDescriptorProto_TYPE_UINT64
	typeInt64 := descriptorpb.FieldDescriptorProto_TYPE_INT64
	typeBool := descriptorpb.FieldDescriptorProto_TYPE_BOOL

	fd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/dex/tx.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.dex"),
		MessageType: []*descriptorpb.DescriptorProto{
			{ // 0: MsgCreatePool
				Name: strp("MsgCreatePool"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("denom_a"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("denomA")},
					{Name: strp("denom_b"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("denomB")},
					{Name: strp("amount_a"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("amountA")},
					{Name: strp("amount_b"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("amountB")},
				},
			},
			{ // 1: MsgCreatePoolResponse
				Name: strp("MsgCreatePoolResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{ // 2: MsgAddLiquidity
				Name: strp("MsgAddLiquidity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("amount_a"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amountA")},
					{Name: strp("amount_b"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("amountB")},
					{Name: strp("min_shares_out"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("minSharesOut")},
				},
			},
			{ // 3: MsgAddLiquidityResponse
				Name: strp("MsgAddLiquidityResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("shares_minted"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("sharesMinted")},
				},
			},
			{ // 4: MsgRemoveLiquidity
				Name: strp("MsgRemoveLiquidity"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("shares_in"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("sharesIn")},
					{Name: strp("min_amount_a_out"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("minAmountAOut")},
					{Name: strp("min_amount_b_out"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("minAmountBOut")},
				},
			},
			{ // 5: MsgRemoveLiquidityResponse
				Name: strp("MsgRemoveLiquidityResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("amount_a"), Number: int32p(1), Label: &label, Type: &typeBytes, JsonName: strp("amountA")},
					{Name: strp("amount_b"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("amountB")},
				},
			},
			{ // 6: MsgSwap
				Name: strp("MsgSwap"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("token_in_denom"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("tokenInDenom")},
					{Name: strp("token_in_amount"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("tokenInAmount")},
					{Name: strp("min_token_out"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("minTokenOut")},
				},
			},
			{ // 7: MsgSwapResponse
				Name: strp("MsgSwapResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("token_out_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("tokenOutDenom")},
					{Name: strp("token_out_amount"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("tokenOutAmount")},
				},
			},
			{ // 8: MsgCreateReferralCode
				Name: strp("MsgCreateReferralCode"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("creator"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("creator")},
				},
			},
			{ // 9: MsgCreateReferralCodeResponse
				Name: strp("MsgCreateReferralCodeResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("referral_code"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("referralCode")},
				},
			},
			{ // 10: MsgRegisterReferral
				Name: strp("MsgRegisterReferral"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("user"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("user")},
					{Name: strp("referral_code"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("referralCode")},
				},
			},
			{ // 11: MsgRegisterReferralResponse
				Name: strp("MsgRegisterReferralResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("referrer"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("referrer")},
				},
			},
			{ // 12: MsgPlaceOrder
				Name: strp("MsgPlaceOrder"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("creator"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("creator")},
					{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("side"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("side")},
					{Name: strp("order_type"), Number: int32p(4), Label: &label, Type: &typeString, JsonName: strp("orderType")},
					{Name: strp("price"), Number: int32p(5), Label: &label, Type: &typeString, JsonName: strp("price")},
					{Name: strp("quantity"), Number: int32p(6), Label: &label, Type: &typeBytes, JsonName: strp("quantity")},
					{Name: strp("time_in_force"), Number: int32p(7), Label: &label, Type: &typeString, JsonName: strp("timeInForce")},
					{Name: strp("trigger_price"), Number: int32p(8), Label: &label, Type: &typeString, JsonName: strp("triggerPrice")},
				},
			},
			{ // 13: MsgPlaceOrderResponse
				Name: strp("MsgPlaceOrderResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("order_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("orderId")},
					{Name: strp("status"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("status")},
				},
			},
			{ // 14: MsgCancelOrder
				Name: strp("MsgCancelOrder"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("creator"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("creator")},
					{Name: strp("order_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("orderId")},
				},
			},
			{ // 15: MsgCancelOrderResponse
				Name: strp("MsgCancelOrderResponse"),
			},
			{ // 16: MsgModifyOrder
				Name: strp("MsgModifyOrder"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("creator"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("creator")},
					{Name: strp("order_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("orderId")},
					{Name: strp("new_price"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("newPrice")},
					{Name: strp("new_quantity"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("newQuantity")},
				},
			},
			{ // 17: MsgModifyOrderResponse
				Name: strp("MsgModifyOrderResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("new_order_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("newOrderId")},
				},
			},
			{ // 18: MsgFollowTrader
				Name: strp("MsgFollowTrader"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("follower"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("follower")},
					{Name: strp("trader"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("trader")},
					{Name: strp("max_per_trade"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("maxPerTrade")},
					{Name: strp("total_budget"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("totalBudget")},
					{Name: strp("copy_ratio"), Number: int32p(5), Label: &label, Type: &typeInt64, JsonName: strp("copyRatio")},
				},
			},
			{ // 19: MsgFollowTraderResponse
				Name: strp("MsgFollowTraderResponse"),
			},
			{ // 20: MsgUnfollowTrader
				Name: strp("MsgUnfollowTrader"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("follower"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("follower")},
					{Name: strp("trader"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("trader")},
				},
			},
			{ // 21: MsgUnfollowTraderResponse
				Name: strp("MsgUnfollowTraderResponse"),
			},
			{ // 22: MsgUpdateCopySettings
				Name: strp("MsgUpdateCopySettings"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("follower"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("follower")},
					{Name: strp("trader"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("trader")},
					{Name: strp("max_per_trade"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("maxPerTrade")},
					{Name: strp("total_budget"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("totalBudget")},
					{Name: strp("copy_ratio"), Number: int32p(5), Label: &label, Type: &typeInt64, JsonName: strp("copyRatio")},
				},
			},
			{ // 23: MsgUpdateCopySettingsResponse
				Name: strp("MsgUpdateCopySettingsResponse"),
			},
			{ // 24: MsgClaimReferralRewards
				Name: strp("MsgClaimReferralRewards"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("referrer"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("referrer")},
				},
			},
			{ // 25: MsgClaimReferralRewardsResponse
				Name: strp("MsgClaimReferralRewardsResponse"),
			},
			{ // 26: MsgMultiHopSwap
				Name: strp("MsgMultiHopSwap"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("route"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("route")},
					{Name: strp("token_in_denom"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("tokenInDenom")},
					{Name: strp("token_in_amount"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("tokenInAmount")},
					{Name: strp("min_token_out"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("minTokenOut")},
				},
			},
			{ // 27: MsgMultiHopSwapResponse
				Name: strp("MsgMultiHopSwapResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("token_out_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("tokenOutDenom")},
					{Name: strp("token_out_amount"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("tokenOutAmount")},
				},
			},
			{ // 28: MsgSetPoolFeeConfig
				Name: strp("MsgSetPoolFeeConfig"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("sender"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("sender")},
					{Name: strp("pool_id"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("volatility_fee_enabled"), Number: int32p(3), Label: &label, Type: &typeBool, JsonName: strp("volatilityFeeEnabled")},
					{Name: strp("base_fee"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("baseFee")},
					{Name: strp("max_fee"), Number: int32p(5), Label: &label, Type: &typeBytes, JsonName: strp("maxFee")},
					{Name: strp("volatility_multiplier"), Number: int32p(6), Label: &label, Type: &typeBytes, JsonName: strp("volatilityMultiplier")},
				},
			},
			{Name: strp("MsgSetPoolFeeConfigResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("CreatePool"), InputType: strp(".syreen.dex.MsgCreatePool"), OutputType: strp(".syreen.dex.MsgCreatePoolResponse")},
					{Name: strp("AddLiquidity"), InputType: strp(".syreen.dex.MsgAddLiquidity"), OutputType: strp(".syreen.dex.MsgAddLiquidityResponse")},
					{Name: strp("RemoveLiquidity"), InputType: strp(".syreen.dex.MsgRemoveLiquidity"), OutputType: strp(".syreen.dex.MsgRemoveLiquidityResponse")},
					{Name: strp("Swap"), InputType: strp(".syreen.dex.MsgSwap"), OutputType: strp(".syreen.dex.MsgSwapResponse")},
					{Name: strp("CreateReferralCode"), InputType: strp(".syreen.dex.MsgCreateReferralCode"), OutputType: strp(".syreen.dex.MsgCreateReferralCodeResponse")},
					{Name: strp("RegisterReferral"), InputType: strp(".syreen.dex.MsgRegisterReferral"), OutputType: strp(".syreen.dex.MsgRegisterReferralResponse")},
					{Name: strp("PlaceOrder"), InputType: strp(".syreen.dex.MsgPlaceOrder"), OutputType: strp(".syreen.dex.MsgPlaceOrderResponse")},
					{Name: strp("CancelOrder"), InputType: strp(".syreen.dex.MsgCancelOrder"), OutputType: strp(".syreen.dex.MsgCancelOrderResponse")},
					{Name: strp("ModifyOrder"), InputType: strp(".syreen.dex.MsgModifyOrder"), OutputType: strp(".syreen.dex.MsgModifyOrderResponse")},
					{Name: strp("FollowTrader"), InputType: strp(".syreen.dex.MsgFollowTrader"), OutputType: strp(".syreen.dex.MsgFollowTraderResponse")},
					{Name: strp("UnfollowTrader"), InputType: strp(".syreen.dex.MsgUnfollowTrader"), OutputType: strp(".syreen.dex.MsgUnfollowTraderResponse")},
					{Name: strp("UpdateCopySettings"), InputType: strp(".syreen.dex.MsgUpdateCopySettings"), OutputType: strp(".syreen.dex.MsgUpdateCopySettingsResponse")},
					{Name: strp("ClaimReferralRewards"), InputType: strp(".syreen.dex.MsgClaimReferralRewards"), OutputType: strp(".syreen.dex.MsgClaimReferralRewardsResponse")},
					{Name: strp("MultiHopSwap"), InputType: strp(".syreen.dex.MsgMultiHopSwap"), OutputType: strp(".syreen.dex.MsgMultiHopSwapResponse")},
					{Name: strp("SetPoolFeeConfig"), InputType: strp(".syreen.dex.MsgSetPoolFeeConfig"), OutputType: strp(".syreen.dex.MsgSetPoolFeeConfigResponse")},
				},
			},
		},
	}

	// Serialize and gzip the file descriptor for Descriptor() methods
	rawDesc, err := proto.Marshal(fd)
	if err != nil {
		panic("dex: failed to marshal file descriptor: " + err.Error())
	}
	var buf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	gz.Write(rawDesc)
	gz.Close()
	fileDescriptorTx = buf.Bytes()

	// Register with nil resolver (no dependencies needed)
	file, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("dex: failed to create proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(file); err != nil {
		// Already registered is fine
	}

	// Register query proto
	qfd := &descriptorpb.FileDescriptorProto{
		Name:    strp("syreen/dex/query.proto"),
		Syntax:  strp("proto3"),
		Package: strp("syreen.dex"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: strp("QueryParamsRequest")},
			{Name: strp("QueryParamsResponse")},
			{
				Name: strp("QueryPoolRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QueryPoolResponse")},
			{Name: strp("QueryPoolsRequest")},
			{Name: strp("QueryPoolsResponse")},
			{
				Name: strp("QueryQuoteRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("token_in_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("tokenInDenom")},
					{Name: strp("token_in_amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("tokenInAmount")},
				},
			},
			{
				Name: strp("QueryQuoteResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("token_out_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("tokenOutDenom")},
					{Name: strp("token_out_amount"), Number: int32p(2), Label: &label, Type: &typeBytes, JsonName: strp("tokenOutAmount")},
					{Name: strp("price_impact"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("priceImpact")},
				},
			},
			{
				Name: strp("QuerySpotPriceRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{
				Name: strp("QuerySpotPriceResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("price_a_b"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("priceAB")},
					{Name: strp("price_b_a"), Number: int32p(3), Label: &label, Type: &typeString, JsonName: strp("priceBA")},
					{Name: strp("denom_a"), Number: int32p(4), Label: &label, Type: &typeString, JsonName: strp("denomA")},
					{Name: strp("denom_b"), Number: int32p(5), Label: &label, Type: &typeString, JsonName: strp("denomB")},
				},
			},
			{
				Name: strp("QueryOptimalRouteRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("input_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("inputDenom")},
					{Name: strp("output_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("outputDenom")},
					{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
				},
			},
			{
				Name: strp("QueryOptimalRouteResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("input_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("inputDenom")},
					{Name: strp("output_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("outputDenom")},
					{Name: strp("input_amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("inputAmount")},
					{Name: strp("output_amount"), Number: int32p(4), Label: &label, Type: &typeBytes, JsonName: strp("outputAmount")},
					{Name: strp("route_type"), Number: int32p(5), Label: &label, Type: &typeString, JsonName: strp("routeType")},
				},
			},
			{
				Name: strp("QueryRoutesRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("input_denom"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("inputDenom")},
					{Name: strp("output_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("outputDenom")},
					{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
				},
			},
			{Name: strp("QueryRoutesResponse")},
			// Risk Score query types
			{
				Name: strp("QueryRiskScoreRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QueryRiskScoreResponse")},
			{
				Name: strp("QueryTradeRiskRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("input_denom"), Number: int32p(2), Label: &label, Type: &typeString, JsonName: strp("inputDenom")},
					{Name: strp("amount"), Number: int32p(3), Label: &label, Type: &typeBytes, JsonName: strp("amount")},
				},
			},
			{Name: strp("QueryTradeRiskResponse")},
			{
				Name: strp("QueryRiskHistoryRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QueryRiskHistoryResponse")},
			// Signals query types
			{
				Name: strp("QuerySignalsRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QuerySignalsResponse")},
			{
				Name: strp("QuerySignalHistoryRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QuerySignalHistoryResponse")},
			// Sentiment Oracle query types
			{
				Name: strp("QuerySentimentRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QuerySentimentResponse")},
			{
				Name: strp("QuerySentimentHistoryRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("limit"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("limit")},
				},
			},
			{Name: strp("QuerySentimentHistoryResponse")},
			{
				Name: strp("QuerySentimentAlertsRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QuerySentimentAlertsResponse")},
			// Referral query types
			{
				Name: strp("QueryReferralInfoRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryReferralInfoResponse")},
			{
				Name: strp("QueryReferralListRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryReferralListResponse")},
			{
				Name: strp("QueryReferralCodeLookupRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("code"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("code")},
				},
			},
			{Name: strp("QueryReferralCodeLookupResponse")},
			{
				Name: strp("QueryReferralEarningsRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryReferralEarningsResponse")},
			// Order Book query types
			{
				Name: strp("QueryOrderBookRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QueryOrderBookResponse")},
			{
				Name: strp("QueryOrderBookDepthRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
				},
			},
			{Name: strp("QueryOrderBookDepthResponse")},
			{
				Name: strp("QueryUserOrdersRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryUserOrdersResponse")},
			{
				Name: strp("QueryOrderDetailRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("order_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("orderId")},
				},
			},
			{Name: strp("QueryOrderDetailResponse")},
			{
				Name: strp("QueryTradesRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("pool_id"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("poolId")},
					{Name: strp("limit"), Number: int32p(2), Label: &label, Type: &typeUint64, JsonName: strp("limit")},
				},
			},
			{Name: strp("QueryTradesResponse")},
			// Copy Trading query types
			{
				Name: strp("QueryLeaderboardRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("limit"), Number: int32p(1), Label: &label, Type: &typeUint64, JsonName: strp("limit")},
				},
			},
			{Name: strp("QueryLeaderboardResponse")},
			{
				Name: strp("QueryTraderStatsRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryTraderStatsResponse")},
			{
				Name: strp("QueryTraderFollowersRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryTraderFollowersResponse")},
			{
				Name: strp("QueryCopyFollowingRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryCopyFollowingResponse")},
			{
				Name: strp("QueryCopyHistoryRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: strp("address"), Number: int32p(1), Label: &label, Type: &typeString, JsonName: strp("address")},
				},
			},
			{Name: strp("QueryCopyHistoryResponse")},
			// Next-gen query types
			{Name: strp("QueryOracleCompositeRequest")},
			{Name: strp("QueryOracleCompositeResponse")},
			{Name: strp("QueryWhaleAlertsRequest")},
			{Name: strp("QueryWhaleAlertsResponse")},
			{Name: strp("QueryTimeWeightedPowerRequest")},
			{Name: strp("QueryTimeWeightedPowerResponse")},
			{Name: strp("QueryTraderRecordRequest")},
			{Name: strp("QueryTraderRecordResponse")},
			{Name: strp("QueryTraderLeaderboardRequest")},
			{Name: strp("QueryTraderLeaderboardResponse")},
			{Name: strp("QueryIBCOrdersRequest")},
			{Name: strp("QueryIBCOrdersResponse")},
			{Name: strp("QueryPoolFeeRequest")},
			{Name: strp("QueryPoolFeeResponse")},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strp("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{Name: strp("Params"), InputType: strp(".syreen.dex.QueryParamsRequest"), OutputType: strp(".syreen.dex.QueryParamsResponse")},
					{Name: strp("Pool"), InputType: strp(".syreen.dex.QueryPoolRequest"), OutputType: strp(".syreen.dex.QueryPoolResponse")},
					{Name: strp("Pools"), InputType: strp(".syreen.dex.QueryPoolsRequest"), OutputType: strp(".syreen.dex.QueryPoolsResponse")},
					{Name: strp("Quote"), InputType: strp(".syreen.dex.QueryQuoteRequest"), OutputType: strp(".syreen.dex.QueryQuoteResponse")},
					{Name: strp("SpotPrice"), InputType: strp(".syreen.dex.QuerySpotPriceRequest"), OutputType: strp(".syreen.dex.QuerySpotPriceResponse")},
					{Name: strp("OptimalRoute"), InputType: strp(".syreen.dex.QueryOptimalRouteRequest"), OutputType: strp(".syreen.dex.QueryOptimalRouteResponse")},
					{Name: strp("Routes"), InputType: strp(".syreen.dex.QueryRoutesRequest"), OutputType: strp(".syreen.dex.QueryRoutesResponse")},
					{Name: strp("Signals"), InputType: strp(".syreen.dex.QuerySignalsRequest"), OutputType: strp(".syreen.dex.QuerySignalsResponse")},
					{Name: strp("SignalHistory"), InputType: strp(".syreen.dex.QuerySignalHistoryRequest"), OutputType: strp(".syreen.dex.QuerySignalHistoryResponse")},
					{Name: strp("RiskScore"), InputType: strp(".syreen.dex.QueryRiskScoreRequest"), OutputType: strp(".syreen.dex.QueryRiskScoreResponse")},
					{Name: strp("TradeRisk"), InputType: strp(".syreen.dex.QueryTradeRiskRequest"), OutputType: strp(".syreen.dex.QueryTradeRiskResponse")},
					{Name: strp("RiskHistory"), InputType: strp(".syreen.dex.QueryRiskHistoryRequest"), OutputType: strp(".syreen.dex.QueryRiskHistoryResponse")},
					{Name: strp("Sentiment"), InputType: strp(".syreen.dex.QuerySentimentRequest"), OutputType: strp(".syreen.dex.QuerySentimentResponse")},
					{Name: strp("SentimentHistory"), InputType: strp(".syreen.dex.QuerySentimentHistoryRequest"), OutputType: strp(".syreen.dex.QuerySentimentHistoryResponse")},
					{Name: strp("SentimentAlerts"), InputType: strp(".syreen.dex.QuerySentimentAlertsRequest"), OutputType: strp(".syreen.dex.QuerySentimentAlertsResponse")},
					{Name: strp("ReferralInfo"), InputType: strp(".syreen.dex.QueryReferralInfoRequest"), OutputType: strp(".syreen.dex.QueryReferralInfoResponse")},
					{Name: strp("ReferralList"), InputType: strp(".syreen.dex.QueryReferralListRequest"), OutputType: strp(".syreen.dex.QueryReferralListResponse")},
					{Name: strp("ReferralCodeLookup"), InputType: strp(".syreen.dex.QueryReferralCodeLookupRequest"), OutputType: strp(".syreen.dex.QueryReferralCodeLookupResponse")},
					{Name: strp("ReferralEarnings"), InputType: strp(".syreen.dex.QueryReferralEarningsRequest"), OutputType: strp(".syreen.dex.QueryReferralEarningsResponse")},
					// Order Book queries
					{Name: strp("OrderBook"), InputType: strp(".syreen.dex.QueryOrderBookRequest"), OutputType: strp(".syreen.dex.QueryOrderBookResponse")},
					{Name: strp("OrderBookDepth"), InputType: strp(".syreen.dex.QueryOrderBookDepthRequest"), OutputType: strp(".syreen.dex.QueryOrderBookDepthResponse")},
					{Name: strp("UserOrders"), InputType: strp(".syreen.dex.QueryUserOrdersRequest"), OutputType: strp(".syreen.dex.QueryUserOrdersResponse")},
					{Name: strp("OrderDetail"), InputType: strp(".syreen.dex.QueryOrderDetailRequest"), OutputType: strp(".syreen.dex.QueryOrderDetailResponse")},
					{Name: strp("Trades"), InputType: strp(".syreen.dex.QueryTradesRequest"), OutputType: strp(".syreen.dex.QueryTradesResponse")},
					// Copy Trading queries
					{Name: strp("Leaderboard"), InputType: strp(".syreen.dex.QueryLeaderboardRequest"), OutputType: strp(".syreen.dex.QueryLeaderboardResponse")},
					{Name: strp("TraderStats"), InputType: strp(".syreen.dex.QueryTraderStatsRequest"), OutputType: strp(".syreen.dex.QueryTraderStatsResponse")},
					{Name: strp("TraderFollowers"), InputType: strp(".syreen.dex.QueryTraderFollowersRequest"), OutputType: strp(".syreen.dex.QueryTraderFollowersResponse")},
					{Name: strp("CopyFollowing"), InputType: strp(".syreen.dex.QueryCopyFollowingRequest"), OutputType: strp(".syreen.dex.QueryCopyFollowingResponse")},
					{Name: strp("CopyHistory"), InputType: strp(".syreen.dex.QueryCopyHistoryRequest"), OutputType: strp(".syreen.dex.QueryCopyHistoryResponse")},
					// Next-gen queries
					{Name: strp("OracleComposite"), InputType: strp(".syreen.dex.QueryOracleCompositeRequest"), OutputType: strp(".syreen.dex.QueryOracleCompositeResponse")},
					{Name: strp("WhaleAlerts"), InputType: strp(".syreen.dex.QueryWhaleAlertsRequest"), OutputType: strp(".syreen.dex.QueryWhaleAlertsResponse")},
					{Name: strp("TimeWeightedPower"), InputType: strp(".syreen.dex.QueryTimeWeightedPowerRequest"), OutputType: strp(".syreen.dex.QueryTimeWeightedPowerResponse")},
					{Name: strp("TraderRecord"), InputType: strp(".syreen.dex.QueryTraderRecordRequest"), OutputType: strp(".syreen.dex.QueryTraderRecordResponse")},
					{Name: strp("TraderLeaderboard"), InputType: strp(".syreen.dex.QueryTraderLeaderboardRequest"), OutputType: strp(".syreen.dex.QueryTraderLeaderboardResponse")},
					{Name: strp("IBCOrders"), InputType: strp(".syreen.dex.QueryIBCOrdersRequest"), OutputType: strp(".syreen.dex.QueryIBCOrdersResponse")},
					{Name: strp("PoolFee"), InputType: strp(".syreen.dex.QueryPoolFeeRequest"), OutputType: strp(".syreen.dex.QueryPoolFeeResponse")},
				},
			},
		},
	}
	qfile, err := protodesc.NewFile(qfd, nil)
	if err != nil {
		panic("dex: failed to create query proto file descriptor: " + err.Error())
	}
	if err := protoregistry.GlobalFiles.RegisterFile(qfile); err != nil {
		// Already registered is fine
	}
}

func strp(s string) *string  { return &s }
func int32p(i int32) *int32 { return &i }
