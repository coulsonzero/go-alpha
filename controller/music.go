package controller

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func MusicList(c *gin.Context) {
	script, _ := filepath.Abs("cmd/music.py")
	cmd := exec.Command("python3", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch music list",
			"details": string(output),
		})
		return
	}

	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	code := 200
	if c, ok := result["code"].(float64); ok {
		code = int(c)
	}
	c.JSON(code, result)
}
