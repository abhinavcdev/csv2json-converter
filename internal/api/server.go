package api

import (
	"csv2json/internal/converter"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ConvertRequest represents the API request for CSV conversion
type ConvertRequest struct {
	CSVData   string                        `json:"csv_data" binding:"required"`
	Options   converter.ConversionOptions   `json:"options"`
}

// ConvertResponse represents the API response
type ConvertResponse struct {
	Success     bool        `json:"success"`
	Data        interface{} `json:"data,omitempty"`
	Error       string      `json:"error,omitempty"`
	Message     string      `json:"message,omitempty"`
	DownloadURL string      `json:"download_url,omitempty"`
	FileSize    int         `json:"file_size,omitempty"`
}

// StartServer initializes and starts the API server
func StartServer() {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "csv2json-api",
		})
	})

	// Convert CSV to JSON endpoint
	r.POST("/convert", convertHandler)

	// File upload endpoint
	r.POST("/upload", uploadHandlerLarge)

	// Download endpoint for large files
	r.GET("/download/:filename", downloadHandler)

	// Serve static files for the frontend
	r.Static("/static", "./frontend/dist")
	r.StaticFile("/", "./frontend/dist/landing.html")
	r.StaticFile("/app.html", "./frontend/dist/app.html")
	r.StaticFile("/docs.html", "./frontend/dist/docs.html")
	r.StaticFile("/landing.html", "./frontend/dist/landing.html")

	r.Run(":8080")
}

// convertHandler handles CSV to JSON conversion via JSON API
func convertHandler(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ConvertResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Use default options if not provided
	if req.Options.Delimiter == 0 {
		req.Options = converter.DefaultOptions()
	}

	// Convert CSV to JSON
	reader := strings.NewReader(req.CSVData)
	jsonData, err := converter.ConvertCSVToJSON(reader, req.Options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Success: false,
			Error:   "Conversion failed: " + err.Error(),
		})
		return
	}

	// Parse JSON data to return as proper JSON response
	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Success: false,
			Error:   "Failed to parse converted JSON: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{
		Success: true,
		Data:    result,
	})
}

// uploadHandler handles file upload and conversion
func uploadHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ConvertResponse{
			Success: false,
			Error:   "No file uploaded: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Parse options from form data
	options := converter.DefaultOptions()
	
	if delimiter := c.PostForm("delimiter"); delimiter != "" {
		if len(delimiter) == 1 {
			options.Delimiter = rune(delimiter[0])
		}
	}
	
	if hasHeader := c.PostForm("has_header"); hasHeader != "" {
		if val, err := strconv.ParseBool(hasHeader); err == nil {
			options.HasHeader = val
		}
	}
	
	if outputFormat := c.PostForm("output_format"); outputFormat != "" {
		options.OutputFormat = outputFormat
	}
	
	if prettyPrint := c.PostForm("pretty_print"); prettyPrint != "" {
		if val, err := strconv.ParseBool(prettyPrint); err == nil {
			options.PrettyPrint = val
		}
	}
	
	if inferTypes := c.PostForm("infer_types"); inferTypes != "" {
		if val, err := strconv.ParseBool(inferTypes); err == nil {
			options.InferTypes = val
		}
	}

	// Convert CSV to JSON
	jsonData, err := converter.ConvertCSVToJSON(file, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Success: false,
			Error:   "Conversion failed: " + err.Error(),
		})
		return
	}

	// Parse JSON data to return as proper JSON response
	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		c.JSON(http.StatusInternalServerError, ConvertResponse{
			Success: false,
			Error:   "Failed to parse converted JSON: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ConvertResponse{
		Success: true,
		Data:    result,
	})
}
