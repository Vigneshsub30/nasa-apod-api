package main

import (
	"github.com/apod/mcp-server/config"
	"github.com/apod/mcp-server/models"
	tools_request_tag "github.com/apod/mcp-server/tools/request_tag"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_request_tag.CreateGet_apodTool(cfg),
	}
}
