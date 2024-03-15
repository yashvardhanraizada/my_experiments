package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	redisCache *redis.Client
	mySqlDb    *sql.DB
)

func connectMySQL() bool {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "dllexpuser:dllexppwd@tcp(127.0.0.1:3306)/dll_experiment")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return false
	}
	//defer db.Close()

	// Ping the database to check if it's reachable
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return false
	}

	fmt.Println("Connected to MySQL database")
	mySqlDb = db

	// Perform database operations
	// For example, querying data from a table
	/*rows, err := db.Query("SELECT * FROM keyvalue")
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
	}*/

	/*for i := 100; i < 10000; i++ {
		// Insert data into the table
		_, err = db.Exec("INSERT INTO keyvalue (`key`, `value`) VALUES (?, ?)", fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
		if err != nil {
			fmt.Println("Error inserting into database:", err)
			return
		}
	}*/

	return true
}

func connectRedis() bool {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Assuming Redis is running on the default port
		Password: "",               // No password set
		DB:       0,                // Use the default database
	})

	// Ping the Redis server to check if it's reachable
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return false
	}

	fmt.Println("Connected to Redis cache")
	redisCache = rdb

	// Use the Redis client to perform operations
	/*err = rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}

	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil {
		fmt.Println("Error getting value for key:", err)
		return
	}
	fmt.Println("Value for key:", val)*/

	return true
}

func getKeyValueFromDB(key string) string {
	var val string
	err := mySqlDb.QueryRow("SELECT `value` FROM keyvalue WHERE `key` = ?", key).Scan(&val)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return ""
	}

	return val
}

func getKeyValue(key string) string {
	val, err := redisCache.Get(context.Background(), key).Result()

	if err != nil {
		fmt.Println("Error getting value for key from cache: ", err)
		fmt.Println("Getting value for key from the database")

		val = getKeyValueFromDB(key)

		if val == "" {
			return ""
		}

		fmt.Println("Setting value for key in cache")
		err = redisCache.Set(context.Background(), key, val, 0).Err()

		if err != nil {
			fmt.Println("Error setting key: ", err)
		}
	}

	return val
}

func main() {
	fmt.Println("-------------------------")
	fmt.Println("Welcome to Load Leaking experiment!")
	fmt.Println("-------------------------")

	if !connectRedis() || !connectMySQL() {
		fmt.Println("Startup failed. Exiting...")
		return
	}

	defer mySqlDb.Close()

	fmt.Println("-------------------------")
	fmt.Println("Starting operations...")
	fmt.Println("-------------------------")

	for i := 5000; i < 5100; i++ {
		val := getKeyValue("key_" + fmt.Sprintf("%d", i))

		if val == "" {
			fmt.Println("Key not found, or some error occurred.")
		} else {
			fmt.Println("Value for key: "+"key_"+fmt.Sprintf("%d", i)+" =", val)
		}

		fmt.Println("-------------------------")
	}

	fmt.Println("Ending operations...")
	fmt.Println("Bye bye!")
	fmt.Println("-------------------------")
}