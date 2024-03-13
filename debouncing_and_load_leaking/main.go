package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

func connectMySQL() {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "dllexpuser:dllexppwd@tcp(127.0.0.1:3306)/dll_experiment")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Ping the database to check if it's reachable
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	fmt.Println("Connected to MySQL database")

	// Perform database operations
	// For example, querying data from a table
	rows, err := db.Query("SELECT * FROM keyvalue")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	// Iterate over the rows and print the results
	for rows.Next() {
		var key string
		var value string
		err = rows.Scan(&key, &value)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return
	}
}

func connectRedis() {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Assuming Redis is running on the default port
		Password: "",               // No password set
		DB:       0,                // Use the default database
	})

	// Ping the Redis server to check if it's reachable
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	// Use the Redis client to perform operations
	err = rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}

	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil {
		fmt.Println("Error getting value for key:", err)
		return
	}
	fmt.Println("Value for key:", val)
}

func main() {
	fmt.Println("Hello World!")
	connectRedis()
	fmt.Println("Connection to Redis worked")
	connectMySQL()
	fmt.Println("Connection to MySQL worked")
}
