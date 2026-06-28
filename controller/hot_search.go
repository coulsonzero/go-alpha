package controller

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HotSearch(c *gin.Context) {
	script, _ := filepath.Abs("cmd/weibo.py")
	date := c.Query("date") // MM-DD 格式，为空则爬实时

	args := []string{script, "--json"}
	if date != "" {
		args = append(args, "--date", date)
	}

	cmd := exec.Command("python3", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to fetch hot search",
			"details": string(output),
		})
		return
	}

	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	// 透传 Python 返回的状态码
	code := 200
	if c, ok := result["code"].(float64); ok {
		code = int(c)
	}
	c.JSON(code, result)
}
