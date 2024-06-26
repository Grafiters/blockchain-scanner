package main

import (
	"fmt"
	"goblock/db"
	"goblock/db/models"
	"goblock/services/blockchain"
	"log"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

var (
	run    bool = true
	params string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	connectionDb, err := db.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Printf("Connected to Database")

	log.Printf("Try to get all data blockchain")

	for run {
		runtime.GC()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		blockchains, err := models.GetAllBlockchainWitHDecrypt(connectionDb, "enabled")
		if err != nil {
			log.Fatalf("Error retrieving blockchain data: %v", err)
		}
		for _, v := range blockchains {
			fmt.Printf("Start fetching %s on block %d", v.Key, v.Height)
			fmt.Println("")

			if v.Server == "" {
				fmt.Printf("Server is null on blockchain currencies %s", v.Key)
				fmt.Println(err)
				continue
			}

			bc_service := blockchain.NewBlockchainService(v, connectionDb)
			fmt.Printf("Fetching block %s on height %d", v.Client, v.Height)
			fmt.Println("")
			lastBlock, err := bc_service.LatestBlockNumber()
			if err != nil {
				fmt.Println("Error:", err)
			}

			if int64(v.Height) < lastBlock-3 {
				bc_service.Fetch(v.Height)
				fmt.Printf("Done fetching %s on block %d", v.Key, v.Height)
				fmt.Println("")
				models.UpdateHeight(connectionDb, int(v.ID), int(v.Height+1))
			} else if lastBlock-3 == int64(v.Height) {
				time.Sleep(5 * time.Second)
			}
			fmt.Println("")
		}
		// time.Sleep(1 * time.Second)
	}

	if !run {
		fmt.Printf("Blockchain Service has Stopped")
	}
}
