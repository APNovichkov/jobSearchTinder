package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)



func main() {
	fmt.Println("Welcome to Job Search Tinder!")

	// Initialize logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// create context
	log.Info("Initializing Context")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	log.Info("Getting Job Listings")
	jobListings, _ := getYCJobListings(ctx, 3)

	log.Info("Converting Job Listings to JSON")
	jobListingsString, _ := json.Marshal(jobListings)

	log.Info(fmt.Sprintf("Got Result: %v", string(jobListingsString)))
}







