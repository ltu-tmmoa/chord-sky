package data

import (
	"sort"
	"testing"
)

type keyList []*ID

func (lst keyList) Len() int {
	return len(lst)
}

func (lst keyList) Less(i, j int) bool {
	return lst[i].Cmp(lst[j]) < 0
}

func (lst keyList) Swap(i, j int) {
	lst[i], lst[j] = lst[j], lst[i]
}

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	id := func(i int64) *ID {
		return newID64(i, 3)
	}

	checkExists := func(key *ID, value string) {
		value0, _ := storage.Get(key)
		value1 := string(value0)
		if value != value1 {
			t.Errorf("storage[%s] expected to be %s, was %s.", key, value, value1)
		}
	}

	checkNotExists := func(key *ID) {
		if value, _ := storage.Get(key); len(value) != 0 {
			t.Errorf("storage[%s] expected to be empty, was %s.", key, string(value))
		}
	}

	storage.Set(id(1), []byte("1"))
	storage.Set(id(2), []byte("2"))
	storage.Set(id(3), []byte("3"))

	checkExists(id(1), "1")
	checkExists(id(2), "2")
	checkExists(id(3), "3")
	checkNotExists(id(4))

	storage.Remove(id(2))
	checkNotExists(id(2))

	storage.Set(id(0), []byte("00"))
	storage.Set(id(1), []byte("10"))
	storage.Set(id(2), []byte("20"))
	storage.Set(id(3), []byte("30"))
	storage.Set(id(4), []byte("40"))
	storage.Set(id(5), []byte("50"))
	storage.Set(id(6), []byte("60"))
	storage.Set(id(7), []byte("70"))
	checkExists(id(3), "30")

	storage.Set(id(10), []byte("22"))
	checkExists(id(2), "22")

	{
		keys, _ := storage.GetKeyRange(id(2), id(4))
		sort.Sort(keyList(keys))
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d.", len(keys))
			return
		}
		if !keys[0].Eq(id(2)) {
			t.Errorf("Expected keys[0] to be 2, got %s.", keys[0])
		}
		if !keys[1].Eq(id(3)) {
			t.Errorf("Expected keys[1] to be 3, got %s.", keys[1])
		}
	}
	{
		keys, _ := storage.GetKeyRange(id(7), id(1))
		sort.Sort(keyList(keys))
		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d.", len(keys))
			return
		}
		if !keys[0].Eq(id(0)) {
			t.Errorf("Expected keys[0] to be 0, got %s.", keys[0])
		}
		if !keys[1].Eq(id(7)) {
			t.Errorf("Expected keys[1] to be 7, got %s.", keys[1])
		}
	}
}
