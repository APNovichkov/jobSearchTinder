package main

import (
	"testing"
)

// For Dani - Can't really run a table test in this example
func TestRemoveElementFromSlice(t *testing.T) {
	testSlice := []JobListing{}
	
	testJobListing1 := JobListing{
		Title: "test1",
		Company: "testCompany1",
	}
	testJobListing2 := JobListing {
		Title: "test2",
		Company: "testCompany2",
	}

	testSlice = append(testSlice, testJobListing1)
	testSlice = append(testSlice, testJobListing2)


	if removeElementFromSlice(testSlice, 0)[0].Company != "testCompany2" {
        t.Error("Expected Company to equal to equal testCompany2")
    }
}