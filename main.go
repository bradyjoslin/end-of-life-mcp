package main

import (
	"fmt"
	"strings"
)

// Main call handler for the servlet
func Call(input CallToolRequest) (CallToolResult, error) {

	args := input.Params.Arguments.(map[string]interface{})

	switch input.Params.Name {
	case "list_available_products":
		return handleGetProducts()
	case "get_product_cycles":
		return handleGetProductCycles(args)
	case "get_cycle_details":
		return handleGetCycleDetails(args)
	default:
		return CallToolResult{}, fmt.Errorf("unknown tool: %s", input.Params.Name)
	}
}

func handleGetProducts() (CallToolResult, error) {
	products, err := getProducts()
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some(err.Error()),
			}},
		}, nil
	}

	response := fmt.Sprintf("Available products: %v", products)

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(response),
		}},
	}, nil
}

func handleGetProductCycles(args map[string]interface{}) (CallToolResult, error) {
	product, ok := args["product_name"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("product_name parameter is required"),
			}},
		}, nil
	}

	productDetails, err := getProductCycles(product)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Failed to get product cycles"),
			}},
		}, nil
	}

	response := make([]string, len(*productDetails))
	for i, cycleDetails := range *productDetails {
		response[i] = formatCycleDetails(product, cycleDetails)
	}

	responseText := strings.Join(response, "\n")

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(responseText),
		}},
	}, nil
}

func handleGetCycleDetails(args map[string]interface{}) (CallToolResult, error) {
	product, ok := args["product_name"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("product_name parameter is required"),
			}},
		}, nil
	}

	cycle, ok := args["cycle_name"].(string)
	if !ok {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("cycle_name parameter is required"),
			}},
		}, nil
	}

	cycleDetails, err := getCycleDetails(product, cycle)
	if err != nil {
		return CallToolResult{
			IsError: some(true),
			Content: []Content{{
				Type: ContentTypeText,
				Text: some("Failed to get cycle details"),
			}},
		}, nil
	}

	details := formatCycleDetails(product, *cycleDetails)

	return CallToolResult{
		Content: []Content{{
			Type: ContentTypeText,
			Text: some(details),
		}},
	}, nil
}

func formatCycleDetails(product string, cycleDetails ProductDetails) string {
	link := "N/A"
	if cycleDetails.Link != nil {
		link = *cycleDetails.Link
	}

	lts := fmt.Sprintf("%v", cycleDetails.LTS)

	discontinued := "N/A"
	if cycleDetails.Discontinued != nil {
		discontinued = fmt.Sprintf("%v", cycleDetails.Discontinued)
	}

	support := "N/A"
	if cycleDetails.Support != nil {
		support = fmt.Sprintf("%v", cycleDetails.Support)
	}

	var eol string
	switch v := cycleDetails.EOL.(type) {
	case string:
		eol = v
	case bool:
		eol = fmt.Sprintf("%t", v)
	default:
		eol = "N/A"
	}

	return fmt.Sprintf(
		"Details for product %s cycle %s details: Release Date: %s, EOL: %s, Latest: %s, Link: %s, LTS: %s, Support: %s, Discontinued: %s",
		product,
		cycleDetails.Cycle,
		cycleDetails.ReleaseDate,
		eol,
		cycleDetails.Latest,
		link,
		lts,
		support,
		discontinued,
	)
}

// Tool description for mcp.run
func Describe() (ListToolsResult, error) {
	return ListToolsResult{
		Tools: []ToolDescription{
			{
				Name:        "list_available_products",
				Description: "Lists all software, operating systems, and devices tracked by endoflife.date. Use this to discover available products before checking specific EOL information.",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
			{
				Name:        "get_product_cycles",
				Description: "Retrieves all release cycles for a specific product with their end-of-life dates. Use this to check EOL status for products like operating systems, programming languages, or software.",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"product_name"},
					"properties": map[string]interface{}{
						"product_name": map[string]interface{}{
							"type":        "string",
							"description": "Name of one of the available products, e.g., 'ubuntu', 'php', 'windows'",
						},
					},
				},
			},
			{
				Name:        "get_cycle_details",
				Description: "Retrieves detailed EOL information about a specific release cycle of a product. Use this to get support dates, LTS status, and latest release information for a particular version.",
				InputSchema: map[string]interface{}{
					"type":     "object",
					"required": []string{"product_name", "cycle_name"},
					"properties": map[string]interface{}{
						"product_name": map[string]interface{}{
							"type":        "string",
							"description": "Name of one of the available products, e.g., 'ubuntu', 'php', 'windows'",
						},
						"cycle_name": map[string]interface{}{
							"type":        "string",
							"description": "Version or release cycle identifier, e.g., '22.04', '8.1', '11'",
						},
					},
				},
			},
		},
	}, nil
}
