package main

import (
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

	log.Info("Getting Job Listings... Be patient")
	jobListings := RunScraper()

	log.Info("Converting Job Listings to JSON")
	fmt.Println("------------------------------------------------")
	// jobListingsString, _ := json.Marshal(jobListings)

	// log.Info(fmt.Sprintf("Got Result: %v", string(jobListingsString)))

	likedJobListings := []JobListing{}
	dislikedJobListings := []JobListing{}

	var maxLikedListings int
	fmt.Printf("Enter number of jobs you want to apply for: ")
	fmt.Scan(&maxLikedListings)

	listingIndex := 0
	for len(likedJobListings) < maxLikedListings && jobListings != nil {
		fmt.Printf("Do you like the following listing (Y,N):\n%v\n", jobListings[listingIndex])
		var yesOrNo string
		fmt.Scan(&yesOrNo)

		if yesOrNo == "Y" {
			likedJobListings = append(likedJobListings, jobListings[listingIndex])
		}else{
			dislikedJobListings = append(dislikedJobListings, jobListings[listingIndex])
		}

		removeElementFromSlice(jobListings, listingIndex)
	}

	fmt.Printf("CONGRATS! You are done for the day!\nHere are the job listings that you liked!\n")

	for i := 0; i < len(likedJobListings); i++ {
		fmt.Printf("%v. %v\n", i, likedJobListings[i])
	}
}

func removeElementFromSlice(arr []JobListing, index int) ([]JobListing){
	if len(arr) == 1 {
		return nil
	}

	return append(arr[:index], arr[index+1:]...)
}







