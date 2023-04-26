package main

import "fmt"

import (
	"sync"
	"time"
)

type Collection struct {
	documents map[string]*Document
	mutex     sync.RWMutex
}

func (col *Collection) InsertDocumentAsync(doc *Document) error {
	// Check if document already exists
	col.mutex.RLock()
	if _, ok := col.documents[doc.ID]; ok {
		col.mutex.RUnlock()
		return fmt.Errorf("document already exists: %s", doc.ID)
	}
	col.mutex.RUnlock()

	// Create a channel to receive the result of the insert operation
	resultChan := make(chan error)

	// Insert the document asynchronously
	go func() {
		// Acquire a write lock on the collection to perform the insertion
		col.mutex.Lock()
		defer col.mutex.Unlock()

		// Check if document already exists (again, in case it was inserted while we were waiting for the lock)
		if _, ok := col.documents[doc.ID]; ok {
			resultChan <- fmt.Errorf("document already exists: %s", doc.ID)
			return
		}

		// Insert the document
		col.documents[doc.ID] = doc

		// Schedule a garbage collection for the document after a delay
		go func() {
			time.Sleep(5 * time.Minute)
			col.mutex.Lock()
			defer col.mutex.Unlock()
			delete(col.documents, doc.ID)
		}()

		// Send a success signal on the result channel
		resultChan <- nil
	}()

	// Wait for the insert operation to complete or fail
	err := <-resultChan
	if err != nil {
		return err
	}

	// Document was inserted successfully
	return nil
}

func (col *Collection) SearchByValueAsync(value string) []string {
	resultChan := make(chan []string)

	go func() {
		col.mutex.RLock()
		defer col.mutex.RUnlock()

		var results []string
		for _, doc := range col.documents {
			if doc.Data == value {
				results = append(results, doc.ID)
			}
		}

		resultChan <- results
	}()

	return <-resultChan
}

func (col *Collection) InsertDocument(doc *Document) error {
	return col.InsertDocumentAsync(doc)
}

func NewCollection() *Collection {
	return &Collection{documents: make(map[string]*Document)}
}

func (col *Collection) GetDocument(id string) (*Document, error) {
	doc, ok := col.documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	return doc, nil
}
