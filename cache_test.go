package cache

import (
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	mc := New()

	// test non-existing get
	k := "foo"
	a, exist := mc.Get(k)
	if exist {
		t.Errorf("Expected to not find key={%v}, but got value={%v}", k, a)
	}

	// test set
	ev := "bar"
	mc.Set(k, ev)
	b, exist := mc.Get(k)
	if !exist {
		t.Errorf("For key={%v}, expected value={%v} but key not found", k, ev)
	}
	if b != ev {
		t.Errorf("For key={%v}, expected value={%v} but got value={%v}", k, ev, b)
	}

	// test set with overwrite
	nev := "baz"
	mc.Set(k, nev)
	c, exist := mc.Get(k)
	if !exist {
		t.Errorf("For key={%v}, expected value={%v} but key not found", k, nev)
	}
	if c != nev {
		t.Errorf("For key={%v}, expected value={%v} but got value={%v}", k, nev, c)
	}
}

func TestMemCacheExpiration(t *testing.T) {
	ttl := 10 * time.Millisecond
	mc := NewWith(ttl, 20*time.Millisecond)

	// test expiration
	k := "foo"
	ev := "bar"
	mc.Set(k, ev)
	<-time.After(15 * time.Millisecond)
	a, exist := mc.Get(k)
	if exist {
		t.Errorf("For key={%v}, expected eviction after ttl={%v} but got value={%v}", k, ttl, a)
	}

	// test expiration reset after overwrite
	nev := "baz"
	mc.Set(k, ev)
	<-time.After(7 * time.Millisecond)
	mc.Set(k, nev)
	<-time.After(5 * time.Millisecond)
	_, exist = mc.Get(k)
	if !exist {
		t.Errorf("For key={%v}, expected value={%v} but key not found before ttl={%v}", k, nev, ttl)
	}
}
