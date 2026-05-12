package statedb

import "github.com/ethereum/go-ethereum/common"

// accessList implements EIP-2929 access list tracking
type accessList struct {
	addresses map[common.Address]int
	slots     []map[common.Hash]struct{}
}

func newAccessList() *accessList {
	return &accessList{
		addresses: make(map[common.Address]int),
	}
}

func (al *accessList) ContainsAddress(addr common.Address) bool {
	_, ok := al.addresses[addr]
	return ok
}

func (al *accessList) Contains(addr common.Address, slot common.Hash) (bool, bool) {
	idx, ok := al.addresses[addr]
	if !ok {
		return false, false
	}
	if idx == -1 {
		return true, false
	}
	_, slotOk := al.slots[idx][slot]
	return true, slotOk
}

func (al *accessList) AddAddress(addr common.Address) bool {
	if _, present := al.addresses[addr]; present {
		return false
	}
	al.addresses[addr] = -1
	return true
}

// copy returns a deep copy of the access list for snapshot/revert support.
func (al *accessList) copy() *accessList {
	cp := newAccessList()
	for addr, idx := range al.addresses {
		cp.addresses[addr] = idx
	}
	cp.slots = make([]map[common.Hash]struct{}, len(al.slots))
	for i, slotMap := range al.slots {
		cpSlotMap := make(map[common.Hash]struct{}, len(slotMap))
		for k, v := range slotMap {
			cpSlotMap[k] = v
		}
		cp.slots[i] = cpSlotMap
	}
	return cp
}

func (al *accessList) AddSlot(addr common.Address, slot common.Hash) (bool, bool) {
	idx, addrPresent := al.addresses[addr]
	if !addrPresent || idx == -1 {
		al.addresses[addr] = len(al.slots)
		slotmap := map[common.Hash]struct{}{slot: {}}
		al.slots = append(al.slots, slotmap)
		return !addrPresent, true
	}
	_, slotPresent := al.slots[idx][slot]
	if !slotPresent {
		al.slots[idx][slot] = struct{}{}
	}
	return false, !slotPresent
}

// TODO: transientStorage is defined but currently unused by StateDB.
// It should be wired into StateDB and exposed via GetTransientState/SetTransientState
// methods once EIP-1153 (TSTORE/TLOAD) opcodes are enabled.
//
// transientStorage implements EIP-1153 transient storage
type transientStorage struct {
	storage map[common.Address]map[common.Hash]common.Hash
}

func newTransientStorage() transientStorage {
	return transientStorage{
		storage: make(map[common.Address]map[common.Hash]common.Hash),
	}
}

func (t transientStorage) Get(addr common.Address, key common.Hash) common.Hash {
	val, ok := t.storage[addr]
	if !ok {
		return common.Hash{}
	}
	return val[key]
}

func (t transientStorage) Set(addr common.Address, key, value common.Hash) {
	if _, ok := t.storage[addr]; !ok {
		t.storage[addr] = make(map[common.Hash]common.Hash)
	}
	t.storage[addr][key] = value
}
