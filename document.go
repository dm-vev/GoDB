package main

type Document struct {
	ID   string
	Data string
}

func NewDocument(id, data string) *Document {
	return &Document{ID: id, Data: data}
}
