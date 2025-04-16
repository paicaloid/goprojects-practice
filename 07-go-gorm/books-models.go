package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func CreateBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Book created successfully")
	return nil
}

func GetBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	return &book
}

func UpdateBook(db *gorm.DB, book *Book) error {
	result := db.Model(&book).Updates(book)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("Book updated successfully")
	return nil
}

func DeleteBook(db *gorm.DB, id int) {
	var book Book
	result := db.Delete(&book, id)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
	}
	fmt.Println("Book deleted successfully")
}

func GetBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	return books
}
