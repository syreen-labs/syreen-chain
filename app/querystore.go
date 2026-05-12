package app

import (
	"io"

	"cosmossdk.io/store/types"
)

// FallbackQueryMultiStore wraps a CommitMultiStore and overrides
// CacheMultiStoreWithVersion to fall back to the latest committed state
// when IAVL version loading fails (SDK v0.50.x bug).
//
// This fixes:
//   - All standard gRPC/REST query endpoints (/cosmos/bank/v1beta1/balances, etc.)
//   - Transaction simulation (required for normal CLI and wallet usage)
//   - Account lookups needed for signing
//
// Ref: https://github.com/cosmos/cosmos-sdk/issues/13317
type FallbackQueryMultiStore struct {
	inner types.CommitMultiStore
}

var _ types.MultiStore = (*FallbackQueryMultiStore)(nil)

func NewFallbackQueryMultiStore(cms types.CommitMultiStore) *FallbackQueryMultiStore {
	return &FallbackQueryMultiStore{inner: cms}
}

// CacheMultiStoreWithVersion tries the versioned load first. If it fails,
// falls back to CacheMultiStore() which reads from the latest committed state.
func (f *FallbackQueryMultiStore) CacheMultiStoreWithVersion(version int64) (types.CacheMultiStore, error) {
	cms, err := f.inner.CacheMultiStoreWithVersion(version)
	if err == nil {
		return cms, nil
	}
	return f.inner.CacheMultiStore(), nil
}

func (f *FallbackQueryMultiStore) CacheMultiStore() types.CacheMultiStore {
	return f.inner.CacheMultiStore()
}

func (f *FallbackQueryMultiStore) GetStore(key types.StoreKey) types.Store {
	return f.inner.GetStore(key)
}

func (f *FallbackQueryMultiStore) GetKVStore(key types.StoreKey) types.KVStore {
	return f.inner.GetKVStore(key)
}

func (f *FallbackQueryMultiStore) TracingEnabled() bool {
	return f.inner.TracingEnabled()
}

func (f *FallbackQueryMultiStore) SetTracer(w io.Writer) types.MultiStore {
	f.inner.SetTracer(w)
	return f
}

func (f *FallbackQueryMultiStore) SetTracingContext(tc types.TraceContext) types.MultiStore {
	f.inner.SetTracingContext(tc)
	return f
}

func (f *FallbackQueryMultiStore) LatestVersion() int64 {
	return f.inner.LatestVersion()
}

func (f *FallbackQueryMultiStore) GetStoreType() types.StoreType {
	return f.inner.GetStoreType()
}

func (f *FallbackQueryMultiStore) CacheWrap() types.CacheWrap {
	return f.inner.CacheWrap()
}

func (f *FallbackQueryMultiStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return f.inner.CacheWrapWithTrace(w, tc)
}
