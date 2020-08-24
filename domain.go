package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	c      *http.Client
	apiKey string
}

func NewClient(c *http.Client, apiKey string) *Client {
	return &Client{c, apiKey}
}

func (dc Client) SearchResidentialPage(rsr ResidentialSearchRequest) ([]SearchResult, error) {
	rsrJSON, err := json.Marshal(rsr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.domain.com.au/v1/listings/residential/_search", bytes.NewBuffer(rsrJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Api-Key", dc.apiKey)
	req.Header.Add("accept", "application/json")
	log.Printf("making request for page #%v: %v, %+v", rsr.PageNumber, req.URL, rsr)
	resp, err := dc.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %v failed: %v", req.URL.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		log.Print(string(b))
		return nil, fmt.Errorf("got non-200 code: %v, %v", resp.StatusCode, resp.Status)
	}

	listingsPage := []SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(&listingsPage)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse json: %v", err)
	}
	log.Printf("got %v listings", len(listingsPage))
	return listingsPage, nil
}

func (dc Client) SearchResidential(rsr ResidentialSearchRequest) ([]SearchResult, error) {
	rsr.PageSize = 200
	rsr.PageNumber = 0
	listings := []SearchResult{}

	// Domain returns an error: "Cannot page beyond 1000 records" if you try to.
	for len(listings) < 1000 {
		listingsPage, err := dc.SearchResidentialPage(rsr)
		if err != nil {
			return nil, err
		}
		if len(listingsPage) == 0 {
			break
		}
		listings = append(listings, listingsPage...)
		rsr.PageNumber++
	}
	return listings, nil
}

type LocationFilter struct {
	State                     string `json:"state"`
	Region                    string `json:"region"`
	Area                      string `json:"area"`
	Suburb                    string `json:"suburb"`
	PostCode                  string `json:"postCode"`
	IncludeSurroundingSuburbs bool   `json:"includeSurroundingSuburbs"`
}

type ResidentialSearchRequest struct {
	ListingType  string `json:"listingType"`
	MinBedrooms  int    `json:"minBedrooms"`
	MinBathrooms int    `json:"minBathrooms"`
	MinCarspaces int    `json:"minCarspaces"`
	PageSize     int    `json:"pageSize"`
	PageNumber   int    `json:"pageNumber"`
	Locations    []LocationFilter
}

// Listing https://developer.domain.com.au/docs/latest/apis/pkg_agents_listings/references/listings_detailedresidentialsearch
type SearchResult struct {
	Listing PropertyListing `json:"listing"`
}

type PropertyListing struct {
	PropertyDetails PropertyDetails `json:"propertyDetails"`
}

type PropertyDetails struct {
	State        string  `json:"state"`
	PropertyType string  `json:"propertyType"`
	Bathrooms    float32 `json:"bathrooms"`
	Bedrooms     float32 `json:"bedrooms"`
	CarSpaces    int32   `json:"carspaces"`
	Suburb       string  `json:"suburb"`
	Postcode     string  `json:"postcode"`
}
