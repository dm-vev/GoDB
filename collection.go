package main

import "fmt"

type Collection struct {
	documents map[string]*Document
}

func NewCollection() *Collection {
	return &Collection{documents: make(map[string]*Document)}
}

func (col *Collection) InsertDocument(doc *Document) error {
	if _, ok := col.documents[doc.ID]; ok {
		return fmt.Errorf("document already exists: %s", doc.ID)
	}
	col.documents[doc.ID] = doc
	return nil
}

func (col *Collection) GetDocument(id string) (*Document, error) {
	doc, ok := col.documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	return doc, nil
}
