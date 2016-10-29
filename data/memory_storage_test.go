package data

import "testing"

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	checkExists := func(key, value string) {
		value0, _ := storage.Get(key)
		value1 := string(value0)
		if value != value1 {
			t.Errorf("storage[%s] expected to be %s, was %s.", key, value, value1)
		}
	}

	checkNotExists := func(key string) {
		if value, _ := storage.Get(key); len(value) != 0 {
			t.Errorf("storage[%s] expected to be empty, was %s.", key, string(value))
		}
	}

	storage.Set("a", []byte("1"))
	storage.Set("b", []byte("2"))
	storage.Set("c", []byte("3"))

	checkExists("a", "1")
	checkExists("b", "2")
	checkExists("c", "3")
	checkNotExists("d")

	storage.Remove("b")
	checkNotExists("b")

	storage.Set("b", []byte("20"))
	storage.Set("c", []byte("30"))
	storage.Set("d", []byte("40"))
	storage.Set("e", []byte("50"))
	checkExists("c", "30")

	keys, _ := storage.GetKeyRange("b", "d")
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d.", len(keys))
		return
	}
	if keys[0] != "b" {
		t.Errorf("Expected keys[0] to be b, got %s.", keys[0])
		return
	}
	if keys[1] != "c" {
		t.Errorf("Expected keys[1] to be c, got %s.", keys[0])
		return
	}
}
