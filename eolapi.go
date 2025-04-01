package main

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// Helper for optional values
func some[T any](t T) *T {
	return &t
}

type ProductResponse []string

type ProductDetails struct {
	Cycle        interface{} `json:"cycle"`        // Can be a number or string
	ReleaseDate  string      `json:"releaseDate"`  // Date in string format (YYYY-MM-DD)
	EOL          interface{} `json:"eol"`          // Can be a string or boolean
	Latest       string      `json:"latest"`       // Latest release in this cycle
	Link         *string     `json:"link"`         // Link to changelog or null
	LTS          interface{} `json:"lts"`          // Can be a string or boolean
	Support      interface{} `json:"support"`      // Date in string format or boolean
	Discontinued interface{} `json:"discontinued"` // Date in string format or boolean
}

func getProducts() (*ProductResponse, error) {
	// Make API request to get available products
	url := "https://endoflife.date/api/all.json"
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf(
			"product API error: %d - %s",
			resp.Status(),
			string(resp.Body()),
		)
	}

	var products ProductResponse
	if err := json.Unmarshal(resp.Body(), &products); err != nil {
		return nil, fmt.Errorf("failed to parse product data: %v", err)
	}

	return &products, nil
}

func getProductCycles(productName string) (*[]ProductDetails, error) {
	// Make API request to get product cycles
	url := fmt.Sprintf("https://endoflife.date/api/%s.json", productName)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf(
			"product cycle API error: %d - %s",
			resp.Status(),
			string(resp.Body()),
		)
	}

	var productDetails []ProductDetails
	if err := json.Unmarshal(resp.Body(), &productDetails); err != nil {
		return nil, fmt.Errorf("failed to parse product cycle data: %v", err)
	}

	return &productDetails, nil
}

func getCycleDetails(productName, cycle string) (*ProductDetails, error) {
	// Make API request to get cycle details
	url := fmt.Sprintf("https://endoflife.date/api/%s/%s.json", productName, cycle)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf(
			"cycle details API error: %d - %s",
			resp.Status(),
			string(resp.Body()),
		)
	}

	var productDetails ProductDetails
	if err := json.Unmarshal(resp.Body(), &productDetails); err != nil {
		return nil, fmt.Errorf("failed to parse cycle detail data: %v", err)
	}

	return &productDetails, nil
}
