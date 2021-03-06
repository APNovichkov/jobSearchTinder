package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)



func main() {
	fmt.Println("Welcome to Job Search Tinder!")

	// Initialize logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Run the scraper
	log.Info("Getting Job Listings... Be patient")
	jobListings := RunScraper()

	// Convert to JSON and write output.json file
	log.Info("Converting Job Listings to JSON")
	fmt.Println("------------------------------------------------")
	jobListingsString, _ := json.Marshal(jobListings)
	filename := "output.json"
	ioutil.WriteFile(filename, jobListingsString, os.ModePerm)

	
	// Define liked and disliked arrays of Job Listings and prompt user for input to simulate JobSearchTinder
	likedJobListings := []JobListing{}
	dislikedJobListings := []JobListing{}

	var maxLikedListings int
	fmt.Printf("Enter number of jobs you want to apply for: ")
	fmt.Scan(&maxLikedListings)
	fmt.Println("------------------------------------------------")

	listingIndex := 0
	for len(likedJobListings) < maxLikedListings && jobListings != nil {
		lst := jobListings[listingIndex]
		fmt.Printf("Do you like the following listing (Y,N):\n\n")
		printListing(lst)
		var yesOrNo string
		fmt.Scan(&yesOrNo)
		fmt.Println("------------------------------------------------")

		if yesOrNo == "Y" {
			likedJobListings = append(likedJobListings, jobListings[listingIndex])
		}else{
			dislikedJobListings = append(dislikedJobListings, jobListings[listingIndex])
		}

		removeElementFromSlice(jobListings, listingIndex)
	}

	fmt.Printf("CONGRATS! You are done for the day!\nHere are the job listings that you liked!\n\n")

	for i := 0; i < len(likedJobListings); i++ {
		fmt.Printf("%v. ", i)
		printListing(likedJobListings[i])
		fmt.Println("")
	}
}

func printListing(lst JobListing){
	// Print Job Listing

	fmt.Printf("Posting Title: %v\nCompany Name: %v\nLocation: %v\nURL of posting: %v\nPosting Date: %v\n", lst.Title, lst.Company, lst.Location, lst.Url, lst.Age)
}

func removeElementFromSlice(arr []JobListing, index int) ([]JobListing){
	// Remove element at given index from given slice

	return append(arr[:index], arr[index+1:]...)
}







