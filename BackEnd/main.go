package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//Define Meal type
type Meal struct {
	Id          int
	Name        string
	Description string
	Price       string
}

type Order struct {
	MealID int    `json:"id"`
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Price  string `json:"price"`
}

func main() {
	//Create DB connection
	var (
		host     = "localhost"
		dbName   = "recipes"
		port     = "3306"
		user     = "root"
		password = "test"
	)

	connection := connectDB(user, password, host, port, dbName)
	defer connection.Close()
	requestHandler(connection)
}

//Handle GET requests for Meals
func requestHandler(db *sql.DB) {
	router := gin.Default()
	router.Use(cors.Default())

	//Endpoint gets products from MYSQL
	router.GET("/get-meals", func(c *gin.Context) {
		var rows *sql.Rows
		var mealList []Meal
		query := "SELECT * FROM meals"

		//Query data from database
		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}
		//Close this db connection when handler finishes
		defer rows.Close()
		for rows.Next() {
			var p Meal
			err = rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price)
			if err != nil {
				panic(err)
			}
			mealList = append(mealList, p)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}

		//Map structs into response for front-end
		c.JSON(200, mealList)

	})
	router.POST("/place-order", func(c *gin.Context) {
		fmt.Println("order placed")

		// var order Order
		jsonString, err := ioutil.ReadAll(c.Request.Body)
		// if err != nil {
		// 	fmt.Println("Error reading JSON string")
		// }
		var order []Order
		err = json.Unmarshal([]byte(jsonString), &order)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("order: %+v\n", order)
		// json.Unmarshal([]byte(jsonString), &order)

	})
	router.Run(":3000")
}

func connectDB(user, password, host, port, dbName string) *sql.DB {
	fmt.Println(user, password)
	connectionString := user + ":" + password + "@" + "tcp(" + host + ":" + port + ")/" + dbName
	fmt.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	} else {
		println("Database connection succseful")
	}

	return db
}
