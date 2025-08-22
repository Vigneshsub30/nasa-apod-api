"""
MCP Server - Python Implementation
"""

import os
import json
import requests
from pathlib import Path
from typing import Annotated
from pydantic import Field
from mcp.server.fastmcp import FastMCP

# Create MCP server instance
mcp = FastMCP("MCP Server")

def get_config():
    """Get configuration from environment or config file."""
    class Config:
        def __init__(self):
            self.base_url = os.getenv("API_BASE_URL")
            self.bearer_token = os.getenv("API_BEARER_TOKEN")
            
            # Try to load from config file if env vars not set
            if not self.base_url or not self.bearer_token:
                config_path = Path.home() / ".api" / "config.json"
                if config_path.exists():
                    with open(config_path, 'r') as f:
                        config_data = json.load(f)
                        self.base_url = self.base_url or config_data.get("baseURL")
                        self.bearer_token = self.bearer_token or config_data.get("bearerToken")
    
    return Config()

# Add configuration resource
@mcp.resource("config://settings")
def get_config_resource() -> str:
    """Get current configuration settings."""
    config = get_config()
    return json.dumps({
        "base_url": config.base_url,
        "bearer_token": "***" if config.bearer_token else None
    }, indent=2)

# Tool functions
@mcp.tool()
def get_apod(date: Annotated[str, Field(description="The date of the APOD image to retrieve")], hd: Annotated[str, Field(description="Retrieve the URL for the high resolution image")]) -> str:
    """Returns images"""
    try:
        config = get_config()
        
        if not config.base_url or not config.bearer_token:
            return "Error: Missing API configuration. Please set API_BASE_URL and API_BEARER_TOKEN environment variables."
        
        # Build request parameters
        params = {}
        if date: params["date"] = date
        if hd: params["hd"] = hd
        
        # Make API call
        url = f"{config.base_url}/api/unknown"
        
        headers = {
            "Authorization": f"Bearer {config.bearer_token}",
            "Accept": "application/json",
            "Content-Type": "application/json"
        }
        
        response = requests.get(url, headers=headers, params=params)
        
        # Handle HTTP errors
        if response.status_code >= 400:
            try:
                error_data = response.json()
                return f"Failed to format JSON: {json.dumps(error_data, indent=2)}"
            except json.JSONDecodeError:
                return f"Failed to format JSON: {response.text}"
        
        # Parse response
        try:
            result = response.json()
            return json.dumps(result, indent=2)
        except json.JSONDecodeError:
            # Fallback to raw text if JSON parsing fails
            return response.text
            
    except requests.exceptions.ConnectionError as e:
        return f"Request failed: Connection error - {str(e)}"
    except requests.exceptions.Timeout as e:
        return f"Request failed: Request timeout - {str(e)}"
    except requests.exceptions.RequestException as e:
        return f"Request failed: {str(e)}"
    except Exception as e:
        return f"Unexpected error: {str(e)}"


if __name__ == "__main__":
    mcp.run()
