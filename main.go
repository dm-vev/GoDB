package main

import "fmt"

func main() {
	db := NewDB()

	// Создаем коллекцию "users"
	err := db.CreateCollection("users")
	if err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}

	// Получаем коллекцию "users"
	col, err := db.GetCollection("users")
	if err != nil {
		fmt.Println("Error getting collection:", err)
		return
	}

	// Создаем документ с ID "user1" и данными {"name": "John Doe", "age": 30}
	doc := NewDocument("user1", `{"name": "John Doe", "age": 30}`)

	// Добавляем документ в коллекцию "users"
	err = col.InsertDocumentAsync(doc)
	if err != nil {
		fmt.Println("Error inserting document:", err)
		return
	}

	// Получаем документ с ID "user1" из коллекции "users"
	doc, err = col.GetDocument("user1")
	if err != nil {
		fmt.Println("Error getting document:", err)
		return
	}

	// Выводим данные документа
	fmt.Println("Data:", doc.Data)
}
