package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/apod/mcp-server/config"
	"github.com/apod/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Get_apodHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["date"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("date=%v", val))
		}
		if val, ok := args["hd"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("hd=%v", val))
		}
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("api_key=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/apod%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			// API key already added to query string
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result []interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGet_apodTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_apod",
		mcp.WithDescription("Returns images"),
		mcp.WithString("date", mcp.Description("The date of the APOD image to retrieve")),
		mcp.WithBoolean("hd", mcp.Description("Retrieve the URL for the high resolution image")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Get_apodHandler(cfg),
	}
}
