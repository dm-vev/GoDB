package main

import "fmt"

type DB struct {
	collections map[string]*Collection
}

func NewDB() *DB {
	return &DB{collections: make(map[string]*Collection)}
}

func (db *DB) CreateCollection(name string) error {
	if _, ok := db.collections[name]; ok {
		return fmt.Errorf("collection already exists: %s", name)
	}
	db.collections[name] = NewCollection()
	return nil
}

func (db *DB) GetCollection(name string) (*Collection, error) {
	col, ok := db.collections[name]
	if !ok {
		return nil, fmt.Errorf("collection not found: %s", name)
	}
	return col, nil
}
