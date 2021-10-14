package main

import (
	"reflect"
	"sync"
	"testing"
)

func Test_store_NewStore(t *testing.T) {
	got := NewStore()
	want := &store{
		lock: sync.RWMutex{},
		data: map[string]string{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("failed to create a new store, got(%v), want(%v)", got, want)
	}
}

func Test_store_Set(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")

	if s.Count() < 2 {
		t.Fatalf("failed to set keys")
	}

	val, found := s.Get("foo")
	if !found {
		t.Fatalf("key %s set but it's not exist", "foo")
	}

	if val != "bar" {
		t.Fatalf("key %s set as %s but it's %s now", "foo", "bar", val)
	}
}

func Test_store_All(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")

	data := s.All()
	c := len(data)
	if c != 2 {
		t.Fatalf("there should be 2 key but there are %d keys", c)
	}
}

func Test_store_Get(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")

	val, found := s.Get("foo")
	if !found {
		t.Fatalf("key %s set but it's not exist", "foo")
	}

	if val != "bar" {
		t.Fatalf("key %s set as %s but it's %s now", "foo", "bar", val)
	}
}

func Test_store_Delete(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")
	s.Delete("foo")

	if _, found := s.Get("foo"); found {
		t.Fatalf("key %s deleted but it exists", "foo")
	}

	c := s.Count()
	if c > 1 {
		t.Fatalf("there should be 2 key but there are %d keys", c)
	}
}

func Test_store_Flush(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")
	s.Flush()

	if _, found := s.Get("foo"); found {
		t.Fatalf("all keys deleted but key %s is exists", "foo")
	}

	c := s.Count()
	if c != 0 {
		t.Fatalf("all keys deleted but there are %d keys", c)
	}
}

func Test_store_Count(t *testing.T) {
	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")

	c1 := s.Count()
	if c1 != 2 {
		t.Fatalf("2 keys set but there are %d keys", c1)
	}

	s.Delete("bar")

	c2 := s.Count()
	if c2 != 1 {
		t.Fatalf("there should be 1 key but there are %d keys", c2)
	}

	s.Flush()

	c3 := s.Count()
	if c3 != 0 {
		t.Fatalf("there should be 0 key but there are %d keys", c3)
	}
}
