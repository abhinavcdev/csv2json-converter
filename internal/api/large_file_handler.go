package api

import (
	"csv2json/internal/converter"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// uploadHandlerLarge handles file upload and conversion with large file support
func uploadHandlerLarge(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
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

	// For large files (>10MB), trigger automatic download instead of temp storage
	if len(jsonData) > 10*1024*1024 {
		// Calculate processing stats - count JSON objects properly
		jsonStr := string(jsonData)
		recordCount := strings.Count(jsonStr, `},{`) + 1 // Count object separators + 1
		if strings.HasPrefix(strings.TrimSpace(jsonStr), "[") {
			// Array format - count objects more accurately
			recordCount = strings.Count(jsonStr, `"Region":`) // Count by a field that should be in every record
			if recordCount == 0 {
				recordCount = strings.Count(jsonStr, `{`) // Fallback to brace count
			}
		}

		// Set headers for direct download
		filename := strings.TrimSuffix(header.Filename, ".csv") + ".json"
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/json")
		c.Header("Content-Length", fmt.Sprintf("%d", len(jsonData)))
		
		// Send file directly
		c.Data(http.StatusOK, "application/json", jsonData)
		return
	}

	// For smaller files, return inline JSON
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

// downloadHandler serves large converted files for download
func downloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "/tmp/" + filename
	
	// Security check - only allow csv2json files
	if !strings.HasPrefix(filename, "csv2json_") || !strings.HasSuffix(filename, ".json") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	
	// This function is no longer needed since we do direct downloads
	c.JSON(http.StatusNotFound, gin.H{"error": "Direct download method used"})
	
	// Set headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	
	c.File(filePath)
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes int) string {
	if bytes == 0 {
		return "0 B"
	}
	k := 1024
	sizes := []string{"B", "KB", "MB", "GB"}
	i := 0
	for bytes >= k && i < len(sizes)-1 {
		bytes /= k
		i++
	}
	return fmt.Sprintf("%.1f %s", float64(bytes), sizes[i])
}
