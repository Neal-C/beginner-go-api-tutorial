package main

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
);

type book struct {
	ID	string `json:"id"`
	Title string	`json:"title"`
	Author string	`json:"author"`
	Quantity int	`json:"quantity"`
};

var books = []book{
{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
};

//*CREATE
func createBook(context *gin.Context){
	var newBook book;

	if err := context.BindJSON(&newBook); err != nil {
		return;
	};

	books = append(books, newBook);
	context.IndentedJSON(http.StatusCreated, newBook);
};

//*READ
func getBooks(context *gin.Context){
	context.IndentedJSON(http.StatusOK, books);
};

func findBookById(id string)(*book, error){
	for i , book := range books {
		if(book.ID == id){
			return &books[i], nil
		};
	}

	return nil, errors.New("book not found 😫")
};

func getBookById(context *gin.Context){
	id, ok := context.GetQuery("id");

	if ok == false {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request 😮"})
		return;
	}

	book, err := findBookById(id);

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found 😣"});
		return;
	}

	context.IndentedJSON(http.StatusOK, book);
};

//*UPDATE

func checkoutBook(context *gin.Context){
	id, ok := context.GetQuery("id");

	if ok == false {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request, missing 'id' parameter 😵"});
		return;
	};

	book, err := findBookById(id);

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, 
			gin.H{"message": "Book not found 🥺"});
		return;
	}

	if( book.Quantity <= 0){
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available 😥"});
		return;
	}

	(*book).Quantity -= 1;
	context.IndentedJSON(http.StatusOK, book);

};

func returnABookById(context *gin.Context){
	id, ok := context.GetQuery("id");

	if ok == false {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request, missing 'id' parameter 😵"});
		return;
	}

	book, err := findBookById(id);

	if err != nil {
		context.IndentedJSON(http.StatusForbidden, gin.H{"message": "That book doesn't come from here, did you mean to do a donation? 😅"});
		return;
	}

	(*book).Quantity += 1;

	context.IndentedJSON(http.StatusAccepted, gin.H{"message": "thanks for being an upstanding citizen 😘"});
}

//*DELETE
//We don't destroy books


func main(){
	router := gin.Default();
	router.GET("/books", getBooks);
	router.GET("/books/by-id",getBookById);
	router.PATCH("/checkout", checkoutBook);
	router.PATCH("/return", returnABookById);
	router.POST("/books", createBook);
	router.Run("localhost:8080");
}