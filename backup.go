package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func (s *store) Save(filePath string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	data, _ := json.MarshalIndent(s.All(), "", "\t")

	_, _ = f.WriteAt(data, 0)

	return nil
}

func (s *store) Load(filePath string) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	d := map[string]string{}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &d)
	if err != nil {
		return err
	}

	s.data = d

	return nil
}
