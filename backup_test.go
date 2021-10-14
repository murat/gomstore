package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path"
	"reflect"
	"testing"
)

func Test_backup_Save(t *testing.T) {
	filePath := path.Join(t.TempDir(), "test.json")

	s := NewStore()
	s.Set("foo", "bar")
	s.Set("bar", "foo")

	err := s.Save(filePath)
	if err != nil {
		t.Fatalf("failed to save backup file, %v", err)
	}

	f, _ := ioutil.ReadFile(filePath)
	data, _ := json.MarshalIndent(s.All(), "", "\t")

	if !bytes.Equal(f, data) {
		t.Fatalf("backup file(%d bytes) does not contain stored data(%d bytes)", len(f), len(data))
	}
}

func Test_backup_Load(t *testing.T) {
	filePath := path.Join(t.TempDir(), "test.json")

	s1 := NewStore()

	s1.Set("foo", "bar")
	s1.Set("bar", "foo")
	_ = s1.Save(filePath)

	s2 := NewStore()
	_ = s2.Load(filePath)

	s1data, s2data := s1.All(), s2.All()
	if !reflect.DeepEqual(s1data, s2data) {
		t.Fatalf("failed to load from file, saved %v, loaded %v", s1data, s2data)
	}

	val, found := s2.Get("foo")
	if !found || val != "bar" {
		t.Fatalf("comparison failed saved and loaded data")
	}
}
