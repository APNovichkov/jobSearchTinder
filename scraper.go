package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type JobListing struct{
	Title string	`json:"posting_name"`
	Company string  `json:"company"`
	Location string `json:"location"`
	Url string		`json:"posting_url"`
	Age string		`json:"posting_age"`
	IntAge int		`json:"posting_age_numerical"`
	Origin string   `json:"origin"`
}



func RunScraper() ([]JobListing){
	log.Info("Running Scraper Module")

	
	totalJobListings := []JobListing{}
	jobListingsChan := make(chan JobListing)
	

	go getYCJobListings(4, jobListingsChan)

	// ycJobListings, err := getYCJobListings(4)
	// if err != nil{
	// 	log.Panic("There was an error getting YC Job Listings")
	// }

	for jobListing := range jobListingsChan{
		totalJobListings = append(totalJobListings, jobListing)
		fmt.Println("Recieved a new job listing")
	}

	// totalJobListings = append(totalJobListings, ycJobListings...)

	// indeedJobListings, err := getIndeedJobListings(1)
	// if err != nil{
	// 	log.Panic("There was an error getting Indeed Job Listing")
	// }

	// totalJobListings = append(totalJobListings, indeedJobListings...)

	return totalJobListings
}	

// Indeed Handler
func getIndeedJobListings(numPages int) ([]JobListing, error){
	// create context
	log.Info("Initializing Context")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Get Job Listings from LinkedIn Job Board

	const IndeedOrigin string = "https://www.indeed.com/jobs?q=Software+Engineer&l=San+Francisco+Bay+Area%2C+CA"

	// Define output array
	jobListings := []JobListing{}

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Navigate to page
	if err := chromedp.Run(ctx, chromedp.Navigate(IndeedOrigin)); err != nil {
		return nil, fmt.Errorf("Error navigating to indeed url")
	}

	log.Info("Looking at Indeed page #1")

	var postingTitles []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(`//*[contains(concat( " ", @class, " " ), concat( " ", "jobtitle", " " ))]`, &postingTitles)); err != nil {
		return nil, fmt.Errorf("Error getting job titles from indeed")
	}

	var postingCompanies []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(`.company .turnstileLink , .company`, &postingCompanies)); err != nil {
		return nil, fmt.Errorf("Error getting company name from indeed")
	}

	var postingDates []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".date", &postingDates)); err != nil {
		return nil, fmt.Errorf("Error getting posting date from indeed")
	}

	var postingLocation []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".accessible-contrast-color-location", &postingLocation)); err != nil {
		return nil, fmt.Errorf("Error getting posting location from indeed")
	}

	for i := 0; i < len(postingTitles); i++{
		titleVal, _ := postingTitles[i].Attribute("title")
		locationVal := ""
		dateVal := ""

		if (postingLocation[i].Children != nil){
			locationVal =  postingLocation[i].Children[0].NodeValue
		}else{
			locationVal = "NA"
		}

		if (postingDates[i].Children != nil){
			dateVal = postingDates[i].Children[0].NodeValue
		}else{
			dateVal = "NA"
		}

		newListing := JobListing{
			Title: titleVal,
			Company: postingCompanies[i].Children[0].NodeValue,
			Location: locationVal,
			Url: "NA",
			Age: dateVal,
			Origin: IndeedOrigin,
		}
		
		jobListings = append(jobListings, newListing)
	}

	return jobListings, nil
}

// YC Handler
func getYCJobListings(numPages int, jobListingsChan chan<- JobListing){
	// Get Job Listings from the YCombinator hacker news jobs page

	// create context
	log.Info("Initializing Context")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	const YcOrigin string = "https://news.ycombinator.com/jobs"

	// Define output array
	// jobListings := []JobListing{}

	// Add timeout to context
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Navigate to page
	if err := chromedp.Run(ctx, chromedp.Navigate(YcOrigin)); err != nil {
		panic("Error getting to yc url")
	}

	for i := 0; i < numPages; i++ {
		log.Info(fmt.Sprintf("Looking at page #%v", i+1))

		// Scrape Posting titles data
		var postingTitles []*cdp.Node
		if err := chromedp.Run(ctx, chromedp.Nodes(`.storylink`, &postingTitles)); err != nil {
			panic("Error getting to job posting object: %v")
		}

		// Scrape Posting dates data
		var postingDates []*cdp.Node
		if err := chromedp.Run(ctx, chromedp.Nodes(`.age a`, &postingDates)); err != nil {
			panic("Error getting jon posting dates: %v")
		}

		// Check if lengths of these two are the same
		if len(postingTitles) != len(postingDates) {
			panic("Length of posting titles and dates do not align!!")
		}

		// Parse data into a new struct and append to output array
		for i := 0; i < len(postingTitles); i++ {
			newListing := JobListing{
				Title: postingTitles[i].Children[0].NodeValue,
				Company: "NA",
				Location: "NA",
				Url: postingTitles[i].AttributeValue("href"),
				Age: postingDates[i].Children[0].NodeValue,
				IntAge: ConvertStringDateToInt(postingDates[i].Children[0].NodeValue),
				Origin: YcOrigin,
			}

			jobListingsChan <- newListing
			// jobListings = append(jobListings, newListing)
		}

		// Click on More link
		log.Info(fmt.Sprintf("Clicking on More link"))
		if err := chromedp.Run(ctx, chromedp.Click(`.morelink`, chromedp.NodeVisible)); err != nil {
			panic("Error clicking on 'More' link")
		}
	}

	close(jobListingsChan)

}


// YC Handler
func getYCJobListings__(numPages int) ([]JobListing, error){
	// Get Job Listings from the YCombinator hacker news jobs page

	// create context
	log.Info("Initializing Context")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	const YcOrigin string = "https://news.ycombinator.com/jobs"

	// Define output array
	jobListings := []JobListing{}

	// Add timeout to context
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Navigate to page
	if err := chromedp.Run(ctx, chromedp.Navigate(YcOrigin)); err != nil {
		return nil, fmt.Errorf("Error getting to yc url")
	}

	for i := 0; i < numPages; i++ {
		log.Info(fmt.Sprintf("Looking at page #%v", i+1))

		// Scrape Posting titles data
		var postingTitles []*cdp.Node
		if err := chromedp.Run(ctx, chromedp.Nodes(`.storylink`, &postingTitles)); err != nil {
			return nil, fmt.Errorf("Error getting to job posting object: %v", err)
		}

		// Scrape Posting dates data
		var postingDates []*cdp.Node
		if err := chromedp.Run(ctx, chromedp.Nodes(`.age a`, &postingDates)); err != nil {
			return nil, fmt.Errorf("Error getting jon posting dates: %v", err)
		}

		// Check if lengths of these two are the same
		if len(postingTitles) != len(postingDates) {
			panic("Length of posting titles and dates do not align!!")
		}

		// Parse data into a new struct and append to output array
		for i := 0; i < len(postingTitles); i++ {
			newListing := JobListing{
				Title: postingTitles[i].Children[0].NodeValue,
				Company: "NA",
				Location: "NA",
				Url: postingTitles[i].AttributeValue("href"),
				Age: postingDates[i].Children[0].NodeValue,
				IntAge: ConvertStringDateToInt(postingDates[i].Children[0].NodeValue),
				Origin: YcOrigin,
			}

			jobListings = append(jobListings, newListing)
		}

		// Click on More link
		log.Info(fmt.Sprintf("Clicking on More link"))
		if err := chromedp.Run(ctx, chromedp.Click(`.morelink`, chromedp.NodeVisible)); err != nil {
			return nil, fmt.Errorf("Error clicking on 'More' link")
		}
	}

	return jobListings, nil
}
