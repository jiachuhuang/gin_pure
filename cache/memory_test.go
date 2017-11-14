package cache

import "testing"

func TestNewMemoryCache(t *testing.T) {
	mc := NewMemoryCache()
	if mc == nil {
		t.Errorf("NewMemoryCache Error")
	} else {
		t.Logf("NewMemoryCache OK")
	}
}

func TestMemoryCache_Init(t *testing.T) {
	mc := NewMemoryCache()
	if mc == nil {
		t.Errorf("NewMemoryCache Error")
	} else {
		t.Logf("NewMemoryCache OK")
	}

	err := mc.Init("1024")
	if err != nil {
		t.Errorf("NewMemoryCache Error %s", err)
	} else {
		t.Logf("MemoryCache Init OK")
	}
}

func TestMemoryCache_Set(t *testing.T) {
	mc := NewMemoryCache()
	if mc == nil {
		t.Errorf("NewMemoryCache Error")
	} else {
		t.Logf("NewMemoryCache OK")
	}

	err := mc.Init("1024")
	if err != nil {
		t.Errorf("NewMemoryCache Error %s", err)
	} else {
		t.Logf("MemoryCache Init OK")
	}

	mc.Set("a", "a", 0)
	mc.Delete("a")
	mc.Set("b", "b", 0)
	v := mc.Get("b")
	t.Logf("%s", v.(string))
	mc.Set("a", "a", 0)
	v = mc.Get("a")
	t.Logf("%s", v.(string))
}