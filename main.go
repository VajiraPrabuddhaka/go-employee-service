package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var db *sql.DB

func initDB() {
	var err error
	connString := fmt.Sprintf("sqlserver://sa:VFhWccmYunx9yyD0@%s?database=employee", os.Getenv("TS_SERVICE_URL"))
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Error creating connection pool: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}

	fmt.Println("Connected to database!")
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/employee", createEmployee)
	r.GET("/employee/:id", getEmployee)

	r.Run(":8080")
}

func createEmployee(c *gin.Context) {
	empName := c.Query("emp_name")
	if empName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "emp_name query parameter is required"})
		return
	}

	id := uuid.New().String()

	query := "INSERT INTO employee (id, name) VALUES (@p1, @p2)"
	_, err := db.Exec(query, id, empName)
	if err != nil {
		log.Printf("Error inserting new employee: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting new employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": empName})
}

func getEmployee(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	var name string
	query := "SELECT name FROM employee WHERE id = @p1"
	err := db.QueryRow(query, id).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			log.Printf("Error retrieving employee: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving employee"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": name})
}
