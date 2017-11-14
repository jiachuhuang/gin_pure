package cache

import (
	"testing"
	"strconv"
)

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

func BenchmarkMemoryCache_Set(b *testing.B) {
	mc := NewMemoryCache()
	mc.Init("1024")

	b.StartTimer()
	var s string
	for i := 0; i < 100000; i++ {
		s = strconv.Itoa(i)
		mc.Set(s,s,0)
	}
	b.StopTimer()

}