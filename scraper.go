package main

import (
	"context"
	"fmt"
	"sync"

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
	
	var wg sync.WaitGroup

	wg.Add(1)
	go getYCJobListings(4, jobListingsChan, wg)
	wg.Add(1)
	go getIndeedJobListings(1, jobListingsChan, wg)

	for jobListing := range jobListingsChan{
		totalJobListings = append(totalJobListings, jobListing)
	}

	// wg.Wait()
	return totalJobListings
}	

// Indeed Handler
func getIndeedJobListings(numPages int, jobListingsChan chan JobListing, wg sync.WaitGroup){
	// Get Job Listings from LinkedIn Job Board 
	defer wg.Done()

	// create context
	log.Info("Initializing Context for Indeed scraper")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()

	const IndeedOrigin string = "https://www.indeed.com/jobs?q=Software+Engineer&l=San+Francisco+Bay+Area%2C+CA"

	// Navigate to page
	if err := chromedp.Run(ctx, chromedp.Navigate(IndeedOrigin)); err != nil {
		panic("Error navigating to indeed url")
	}

	log.Info("Looking at Indeed page #1")

	var postingTitles []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(`//*[contains(concat( " ", @class, " " ), concat( " ", "jobtitle", " " ))]`, &postingTitles)); err != nil {
		panic("Error getting job titles from indeed")
	}

	var postingCompanies []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(`.company .turnstileLink , .company`, &postingCompanies)); err != nil {
		panic("Error getting company name from indeed")
	}

	var postingDates []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".date", &postingDates)); err != nil {
		panic("Error getting posting date from indeed")
	}

	var postingLocation []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".accessible-contrast-color-location", &postingLocation)); err != nil {
		panic("Error getting posting location from indeed")
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
		
		jobListingsChan <- newListing
	}

	

	cancel()
}

// YC Handler
func getYCJobListings(numPages int, jobListingsChan chan JobListing, wg sync.WaitGroup){
	// Get Job Listings from the YCombinator hacker news jobs page

	defer wg.Done()

	// Define Context
	log.Info("Initializing context for YC scraper")
	ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

	const YcOrigin string = "https://news.ycombinator.com/jobs"

	// Add timeout to context
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()

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
		}

		// Click on More link
		log.Info(fmt.Sprintf("Clicking on More link"))
		if err := chromedp.Run(ctx, chromedp.Click(`.morelink`, chromedp.NodeVisible)); err != nil {
			panic("Error clicking on 'More' link")
		}
	}

	close(jobListingsChan)

	cancel()
}