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
	defer f.Close()

	data, _ := json.MarshalIndent(s.All(), "", "\t")

	_, err = f.WriteAt(data, 0)
	if err != nil {
		return err
	}

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
