package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	pageSize = 200
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
	rsr.PageSize = int32(pageSize)
	rsr.PageNumber = 0
	listings := []SearchResult{}

	// Domain returns an error: "Cannot page beyond 1000 records" if you try to.
	for len(listings) < 1000 {
		listingsPage, err := dc.SearchResidentialPage(rsr)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listingsPage...)
		if len(listingsPage) < pageSize {
			break
		}
		rsr.PageNumber++
	}
	return listings, nil
}

// LocationFilter is Domain.SearchService.v2.Model.DomainSearchWebApiV2ModelsSearchLocation
type LocationFilter struct {
	// [ ACT, NSW, QLD, VIC, SA, WA, NT, TAS ]
	State                     string `json:"state"`
	Region                    string `json:"region"`
	Area                      string `json:"area"`
	Suburb                    string `json:"suburb"`
	PostCode                  string `json:"postCode"`
	IncludeSurroundingSuburbs bool   `json:"includeSurroundingSuburbs"`
	SurroundingRadiusInMeters *int32 `json:"surroundingRadiusInMeters"`
}

// ResidentialSearchRequest is Domain.SearchService.v2.Model.DomainSearchWebApiV2ModelsSearchParameters
type ResidentialSearchRequest struct {
	ListingType  string           `json:"listingType"`
	MinBedrooms  *float32         `json:"minBedrooms"`
	MaxBedrooms  *float32         `json:"maxBedrooms"`
	MinBathrooms *float32         `json:"minBathrooms"`
	MaxBathrooms *float32         `json:"maxBathrooms"`
	MinCarspaces *int32           `json:"minCarspaces"`
	MaxCarspaces *int32           `json:"maxCarspaces"`
	MinPrice     *int32           `json:"minPrice"`
	MaxPrice     *int32           `json:"maxPrice"`
	PageSize     int32            `json:"pageSize"`
	PageNumber   int32            `json:"pageNumber"`
	Locations    []LocationFilter `json:"locations"`
	UpdatedSince string           `json:"updatedSince"`
	ListedSince  string           `json:"listedSince"`
}

// SearchResult is Domain.SearchService.v2.Model.DomainSearchContractsV2SearchResult
type SearchResult struct {
	Type    string          `json:"type"`
	Listing PropertyListing `json:"listing"`
}

// PriceDetails is Domain.SearchService.v2.Model.DomainSearchContractsV2PriceDetails
type PriceDetails struct {
	Price        int32  `json:"price"`
	PriceFrom    int32  `json:"priceFrom"`
	PriceTo      int32  `json:"priceTo"`
	DisplayPrice string `json:"displayPrice"`
}

// PropertyListing is Domain.SearchService.v2.Model.DomainSearchContractsV2PropertyListing
// https://production-api.domain.com.au/swagger/index.html?urls.primaryName=API%20Version%201#model-Domain.SearchService.v2.Model.DomainSearchContractsV2PropertyListing
type PropertyListing struct {
	ID int32 `json:"id"`
	// Sale, Rent, Share, Sold, NewHomes
	ListingType        string          `json:"listingType"`
	Headline           string          `json:"headline"`
	SummaryDescription string          `json:"summaryDescription"`
	HasFloorplan       bool            `json:"hasFloorplan"`
	Labels             []string        `json:"labels"`
	ListingSlug        string          `json:"listingSlug"`
	PropertyDetails    PropertyDetails `json:"propertyDetails"`
	PriceDetails       PriceDetails    `json:"priceDetails"`
	DateAvailable      string          `json:"dateAvailable"`
	DateListed         string          `json:"dateListed"`
}

// AuctionSchedule is Domain.SearchService.v2.Model.DomainSearchContractsV2AuctionSchedule
type AuctionSchedule struct {
	Time            string `json:"time"`
	AuctionLocation string `json:"auctionLocation"`
}

// PropertyDetails is Domain.SearchService.v2.Model.DomainSearchContractsV2PropertyDetails
type PropertyDetails struct {
	State              string   `json:"state"`
	PropertyType       string   `json:"propertyType"`
	Bathrooms          float32  `json:"bathrooms"`
	Bedrooms           float32  `json:"bedrooms"`
	CarSpaces          int32    `json:"carspaces"`
	Features           []string `json:"featuress"`
	AllPropertyTypes   []string `json:"allPropertyTypes"`
	UnitNumber         string   `json:"unitNumber"`
	StreetNumber       string   `json:"streetNumber"`
	Street             string   `json:"street"`
	Area               string   `json:"area"`
	Region             string   `json:"region"`
	Suburb             string   `json:"suburb"`
	SuburbID           int32    `json:"suburbId"`
	Postcode           string   `json:"postcode"`
	DisplayableAddress string   `json:"displayableAddress"`
	Latitude           float32  `json:"latitude"`
	Longitude          float32  `json:"longitude"`
	MapCertainty       int32    `json:"mapCertainty"`
	LandArea           float64  `json:"landArea"`
	BuildingArea       float64  `json:"buildingArea"`
	OnlyShowProperties []string `json:"onlyShowProperties"`
	DisplayAddressType string   `json:"displayAddressType"`
	IsRural            bool     `json:"isRural"`
	TopSpotKeywords    []string `json:"topSpotKeywords"`
	IsNew              bool     `json:"isNew"`
	Tags               []string `json:"tags"`
}
