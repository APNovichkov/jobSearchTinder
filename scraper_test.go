package main

import "testing"

// For Dani - Can't really do a table test here either, result is always going to be the same (unpredictable), no input values
func TestRunScraper(t *testing.T){

	jobListings := RunScraper()

	if (len(jobListings) == 0){
		t.Error("Scraper returned an empty job listings array, something is wrong")
	}
}

func BenchmarkRunScraper(b *testing.B){
	for i := 0; i < 5; i++{
		RunScraper()
	}
}