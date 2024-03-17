package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	redisCache    *redis.Client
	mySqlDb       *sql.DB
	keysFromCache int
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
	err := mySqlDb.QueryRow("SELECT first_name FROM employees WHERE emp_no = ?", key).Scan(&val)
	if err != nil {
		//fmt.Println("Error querying database:", err)
		return ""
	}

	return val

	// 10001 -> 499999
}

func getKeyValue(key string) string {
	val, err := redisCache.Get(context.Background(), key).Result()

	if err != nil {
		//fmt.Println("Error getting value for key from cache: ", err)
		//fmt.Println("Getting value for key from the database")

		val = getKeyValueFromDB(key)

		if val == "" {
			return ""
		}

		//fmt.Println("Setting value for key in cache")
		err = redisCache.Set(context.Background(), key, val, 0).Err()

		if err != nil {
			//fmt.Println("Error setting key: ", err)
		}
	} else {
		keysFromCache++
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

	rand.Seed(time.Now().UnixNano())
	rangeMin := 10001
	rangeMax := 100001

	iterations := 1000000
	iterations_jump := iterations / 10
	iterations_chunk_count := -1
	keysFromCache = 0

	isLoadLeakEnabled := true
	loadLeakPercentage := 20.0
	loadLeakJump := int(math.Round(float64(100.0 / loadLeakPercentage)))

	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		randomIndex := rand.Intn(rangeMax-rangeMin+1) + rangeMin
		getKeyValue(fmt.Sprintf("%d", randomIndex))

		if isLoadLeakEnabled && (i%loadLeakJump == 0) {
			go getKeyValueFromDB(fmt.Sprintf("%d", randomIndex))
		}

		if i%iterations_jump == 0 {
			iterations_chunk_count++
			fmt.Println("Completed", iterations_chunk_count*10, "% of the total iterations")
		}

		//val := getKeyValue("key_" + fmt.Sprintf("%d", i))

		/*if val == "" {
			fmt.Println("Key not found, or some error occurred.")
		} else {
			fmt.Println("Value for key: "+"key_"+fmt.Sprintf("%d", i)+" =", val)
		}

		fmt.Println("-------------------------")*/
	}

	duration := time.Since(startTime)

	fmt.Println("Total queries:", iterations)
	fmt.Println("Queries answered from cache:", keysFromCache)
	fmt.Println("Queries answered from database:", iterations-keysFromCache)
	fmt.Println("Load leaking enabled:", isLoadLeakEnabled)
	if isLoadLeakEnabled {
		fmt.Println("Load leaking percentage:", loadLeakPercentage)
	}
	fmt.Println("Total time taken by all the queries:", duration.Microseconds(), "microseconds")
	fmt.Println("Average time taken by one single query:", float64(duration.Microseconds())/float64(iterations), "microseconds")

	fmt.Println("-------------------------")
	fmt.Println("Ending operations...")
	fmt.Println("Bye bye!")
	fmt.Println("-------------------------")
}
