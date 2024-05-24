package main

import (
	"fmt"
	"os"
	db "search/internal/mango/db"
	localsource "search/internal/mango/madokami"
	"time"
)

func main() {

	// web.Serve()
	connectionString := fmt.Sprintf("mongodb+srv://ductran:%s@inkling-cluster.jnpkxro.mongodb.net/?tls=true", os.Getenv("atlaspwd"))

	client, err := db.Connect(connectionString)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("Database connected")
	}

	database := client.Database("madokami")

	start := time.Now()
	localsource.HashDirectory("/home/noel/Madokami/", database)
	fmt.Println("Hashing took", time.Since(start))
}
