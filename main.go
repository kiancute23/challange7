package main

import (
	"challange7/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "skripsikia23"
	dbname   = "challange7"
)

var (
	db  *sql.DB
	err error
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database")

	router := gin.Default()

	router.GET("/books", GetBook)
	router.GET("/books/:id", GetBookByID)
	router.POST("/books", CreateBook)
	router.PUT("/books/:id", UpdateBook)
	router.DELETE("/books/:id", DeleteBook)

	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}

}

func CreateBook(c *gin.Context) {
	var book model.Book

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input!",
		})
	}

	sqlStatement := `INSERT INTO books (title, author, description) VALUES ($1, $2, $3)	`

	_, err := db.Exec(sqlStatement, &book.Title, &book.Author, &book.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error insert data!",
		})

		panic(err)
	} else {
		c.JSON(200,
			model.Response{
				200,
				"Success Insert Book",
				nil,
			})
	}

	fmt.Printf("New Book Data Inserted : %+v\n", book)
}

func GetBook(c *gin.Context) {
	var result = []model.Book{}

	sqlStatement := `SELECT title, author, description from Books`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var book model.Book

		if err = rows.Scan(
			&book.Title,
			&book.Author,
			&book.Description); err != nil {
			log.Fatal("Error binding to struct")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error Get Data!",
			})
			panic(err)
		} else {
			result = append(result, book)
		}
	}

	c.JSON(200,
		model.Response{
			200,
			"Success",
			result,
		})

	fmt.Printf("Book data : %v\n", result)
}

func GetBookByID(c *gin.Context) {
	var arr_book []model.Book
	var book model.Book

	bookId := c.Param("id")

	sqlStatement := `SELECT title, author, description FROM books WHERE id = $1;`

	rows := db.QueryRow(sqlStatement, bookId)
	switch err := rows.Scan(&book.Title,
		&book.Author,
		&book.Description); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		c.JSON(200,
			model.Response{
				200,
				"Success",
				append(arr_book, book),
			})
	default:
		panic(err)
	}

	fmt.Printf("Book data by id: %v\n", arr_book)
}

func UpdateBook(c *gin.Context) {
	var book model.Book

	bookId := c.Param("id")

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input request!",
		})
	}

	sqlStatement := `UPDATE books SET title = $1, author = $2, description =$3 WHERE id = $4;`
	res, err := db.Exec(sqlStatement, &book.Title, &book.Author, &book.Description, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error Update Book!",
		})
		panic(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	} else {
		c.JSON(200, model.Response{
			200,
			"Successfully update book",
			nil,
		})
	}

	fmt.Printf("Updated data book %d \n", count)
}

func DeleteBook(c *gin.Context) {
	bookId := c.Param("id")

	sqlStatement := `DELETE FROM books	WHERE id = $1;`
	res, err := db.Exec(sqlStatement, bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error delete book!",
		})
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	} else {
		c.JSON(200, model.Response{
			200,
			"Successfully delete book",
			nil,
		})
	}

	fmt.Printf("Deleted data book : %d", count)
}
