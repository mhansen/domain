package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mhansen/domain"
)

var (
	apiKey = flag.String("api_key", "", "Domain API Key")
)

func main() {
	flag.Parse()

	var c http.Client
	dc := domain.NewClient(&c, *apiKey)

	listings, err := dc.SearchResidential(
		domain.ResidentialSearchRequest{
			ListingType: "Rent",
			Locations: []domain.LocationFilter{
				{
					State:                     "NSW",
					Suburb:                    "Pyrmont",
					IncludeSurroundingSuburbs: false,
				},
			},
		})
	if err != nil {
		log.Fatalf("error searching: %v", err)
	}
	for _, l := range listings {
		fmt.Println(l.Listing.ID)
	}
	// fmt.Printf("%+v\n", listings)
}
