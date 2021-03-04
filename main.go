package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)



func main() {
	fmt.Println("Welcome to Job Search Tinder!")

	// Initialize logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.Info("Getting Job Listings")
	jobListings := RunScraper()

	log.Info("Converting Job Listings to JSON")
	jobListingsString, _ := json.Marshal(jobListings)

	log.Info(fmt.Sprintf("Got Result: %v", string(jobListingsString)))
}







