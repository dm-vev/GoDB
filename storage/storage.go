package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type Storage struct {
	data map[string]map[string]string
	mu   sync.RWMutex
}

func New(filename string) (*Storage, error) {
	s := &Storage{
		data: make(map[string]map[string]string),
	}

	err := s.loadFromFile(filename)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) CreateTable(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[name]; ok {
		return errors.New("table already exists")
	}

	s.data[name] = make(map[string]string)
	return s.saveToFile("data.json")
}

func (s *Storage) DeleteTable(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[name]; !ok {
		return errors.New("table not found")
	}

	delete(s.data, name)
	return s.saveToFile("data.json")
}

func (s *Storage) Set(name, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.data[name]
	if !ok {
		return errors.New("table not found")
	}

	table[key] = value
	return s.saveToFile("data.json")
}

func (s *Storage) Get(name, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	table, ok := s.data[name]
	if !ok {
		return "", errors.New("table not found")
	}

	value, ok := table[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return value, nil
}

func (s *Storage) Delete(name, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	table, ok := s.data[name]
	if !ok {
		return errors.New("table not found")
	}

	if _, ok := table[key]; !ok {
		return errors.New("key not found")
	}

	delete(table, key)
	return s.saveToFile("data.json")
}

func (s *Storage) ExportToFile(filename string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.saveToFile(filename)
}

func (s *Storage) ImportFromFile(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadFromFile(filename)
}

func (s *Storage) saveToFile(filename string) error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (s *Storage) loadFromFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.data)
}
