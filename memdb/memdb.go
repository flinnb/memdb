package memdb

import "fmt"

type MemDB struct {
	parent     *MemDB
	store      map[string]string
	deleted    map[string]string
	valueCount map[string]int
}

var activeDB *MemDB

func Init() {
	db := MemDB{
		parent:     nil,
		store:      map[string]string{},
		valueCount: map[string]int{},
	}
	activeDB = &db
}

func Set(k, v string) {
	activeDB.set(k, v)
}

func Get(k string) string {
	v, _ := activeDB.get(k)
	return v
}

func Delete(k string) {
	activeDB.deleteKey(k)
}

func Count(v string) int {
	return activeDB.getCount(v)
}

func Begin() {
	db := MemDB{
		parent:     activeDB,
		store:      map[string]string{},
		deleted:    map[string]string{},
		valueCount: map[string]int{},
	}
	activeDB = &db
}

func Rollback() error {
	if activeDB.parent != nil {
		activeDB = activeDB.parent
	} else {
		return fmt.Errorf("TRANSACTION NOT FOUND")
	}
	return nil
}

func Commit() {
	activeDB.commit()
	resetActiveDB()
}

func resetActiveDB() {
	if activeDB.parent != nil {
		p := activeDB.parent
		activeDB = p
		resetActiveDB()
	}
}

func (m *MemDB) set(k, v string) {
	oldV, ok := m.get(k)
	if ok {
		oldCount := m.getCount(oldV)
		// If we already have a zero count, we don't need to do anything.
		// This shouldn't ever happen if we actually have an old value,
		// but we should check anyway
		if oldCount > 0 {
			m.valueCount[oldV] = oldCount - 1
		}
	}
	// If we're seeting a value on a key we've previously deleted, we need
	// to remove it from the deleted map. `delete()` is a no-op if the key
	// isn't found, so we can be lazy here
	if m.deleted != nil {
		delete(m.deleted, k)
	}
	m.store[k] = v
	vCount := m.getCount(v)
	m.valueCount[v] = vCount + 1
}

func (m *MemDB) get(k string) (string, bool) {
	v, ok := m.store[k]

	if !ok {
		// If we can't find key while in transaction, look for it
		// in parent db, unless we've marked it as deleted in
		// current transaction
		if m.parent != nil {
			_, deleted := m.getDeleted(k)
			if deleted {
				v = "NULL"
			} else {
				return m.parent.get(k)
			}
		} else {
			v = "NULL"
		}
	}
	return v, ok
}

func (m *MemDB) getDeleted(k string) (string, bool) {
	if m.deleted != nil {
		v, deleted := m.deleted[k]

		if !deleted {
			// If we can't find key, look for it in parent db
			if m.parent != nil {
				return m.parent.getDeleted(k)
			}
		} else {
			return v, deleted
		}

	}
	return "", false
}

func (m *MemDB) deleteKey(k string) {
	v, ok := m.get(k)
	if ok {
		oldCount := m.getCount(v)
		// If we already have a zero count, we don't need to do anything.
		// This shouldn't ever happen if we actually have a  value, but
		// we should check anyway
		if oldCount > 0 {
			m.valueCount[v] = oldCount - 1
		}
	}
	if m.parent != nil {
		m.deleted[k] = v
	}
	delete(m.store, k)
}

func (m *MemDB) getCount(v string) int {
	c, ok := m.valueCount[v]
	if !ok {
		if m.parent != nil {
			return m.parent.getCount(v)
		} else {
			c = 0
		}
	}
	return c
}

func (m *MemDB) commit() {
	if m.parent == nil {
		// Not in a transaction, so this is no-op
		return
	}
	for k, v := range m.store {
		m.parent.set(k, v)
	}
	for k, _ := range m.deleted {
		m.parent.deleteKey(k)
	}
	m.parent.commit()
}
