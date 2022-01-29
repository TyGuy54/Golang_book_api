// testing api GET: curl localhost:8080/books
// testing api POST:  curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
// testing api PATCH: curl localhost:8080/checkout?id=2 --request "PATCH"
// testing api PATCH: curl localhost:8080/return?id=2 --request "PATCH"

package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

// look up what a struck is
// im useing this struct kinda like json so it can be returned as json
// definitly need to study this concept more
type book struct {
	ID		string `json: "id"`
	Title	string `json: "title"`
	Author	string `json: "author"`
	Quantity	int `json: "quantity"`
}

// data structure for books
// slice = array
var books = []book{
	{ID: "1", Title: "Cool Book for Cool Dudes", Author: "Adam Demamf", Quantity: 2},
	{ID: "2", Title: "Magic the Gathering for Dummies", Author: "Tyler Ross", Quantity: 9},
	{ID: "3", Title: "Treating Cats 101", Author: "Jobie Donarski", Quantity: 5},
}

func get_books(c *gin.Context) {
	// *gin.Context is all the information abput the request( the * is a pointer)
	// it stores all the information related to a request
	// alows me to return a response

	// IndentedJSON means im going to get nicley formated JSON
	// the data we are sending is books
	c.IndentedJSON(http.StatusOK, books)
}

func book_by_id(c *gin.Context){
	id := c.Param("id")
	book, err := get_book_by_id(id)

	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return 
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkout_book(c *gin.Context){
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing quary parameter"})
		return
	}

	book, err := get_book_by_id(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not found"})
		return
	}

	 book.Quantity -= 1
	 c.IndentedJSON(http.StatusOK, book)
}

func return_book(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing quary parameter"})
		return
	}

	book, err := get_book_by_id(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// this is a helper function that 
// will return a book based on its id
// the (*book. error) is the return type
func get_book_by_id(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func create_book(c *gin.Context) {
	var new_book book

	// binding the JSON data to the new book struct
	if err := c.BindJSON(&new_book); err != nil {
		return  // return a status error as a response from .BindJSON
	}
	// if everyting is ok we append the new book to the books slice(array)
	books = append(books, new_book)
	c.IndentedJSON(http.StatusCreated, new_book)
}

func main() {
	// setting up a router with gin
	router := gin.Default()
	// handling route /books
	router.GET("/books", get_books) // GET: Getting information 
	router.GET("/books/:id", book_by_id) // path parameters
	router.POST("/books", create_book) // POST: Adding information 
	router.PATCH("/checkout", checkout_book) // PATCH: Updating information
	router.PATCH("/return", return_book)
	// when I go to local host 8080 and go to /books itll call the 
	// function get_books()
	// router.Run() runs the web server
	router.Run("localhost:8080")
}
