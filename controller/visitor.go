package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"go-alpha/models"
	"go-alpha/response"
)

const visitorCookieKey = "_user_uuid"

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	// Format as hex string (32 chars)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// RecordVisit checks Cookie to identify new vs returning visitors.
func RecordVisit(c *gin.Context) {
	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	cookieUUID, err := c.Cookie(visitorCookieKey)
	isNewVisitor := (err != nil || cookieUUID == "")

	if isNewVisitor {
		cookieUUID = generateUUID()
		// Set cookie, expires in 1 year
		c.SetCookie(visitorCookieKey, cookieUUID, 365*24*3600, "/", "", false, true)

		// New visitor: count in both total and today
		models.RDB.PFAdd(ctx, "visit:uv:total", cookieUUID)
	}

	// Count into today's UV (also covers new visitor's first day)
	models.RDB.PFAdd(ctx, fmt.Sprintf("visit:uv:%s", today), cookieUUID)

	// PV — always increment
	models.RDB.Incr(ctx, fmt.Sprintf("visit:pv:%s", today))

	totalUV, _ := models.RDB.PFCount(ctx, "visit:uv:total").Result()

	response.Success("记录成功", gin.H{
		"is_new_visitor": isNewVisitor,
		"total_visitors": totalUV,
	}, c)
}

// GetVisitorStats returns today/ weekly/ all-time UV and PV,
// and syncs today's data from Redis into MySQL.
func GetVisitorStats(c *gin.Context) {
	ctx := context.Background()
	today := time.Now().Format("2006-01-02")

	todayUV, _ := models.RDB.PFCount(ctx, fmt.Sprintf("visit:uv:%s", today)).Result()
	todayPV, _ := models.RDB.Get(ctx, fmt.Sprintf("visit:pv:%s", today)).Int64()

	// Sync today's data to MySQL
	models.VisitorDaily{}.Upsert(today, todayUV, todayPV)

	weekKeys := make([]string, 7)
	for i := 0; i < 7; i++ {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		weekKeys[i] = fmt.Sprintf("visit:uv:%s", d)
	}
	weekUV, _ := models.RDB.PFCount(ctx, weekKeys...).Result()

	totalUV, _ := models.RDB.PFCount(ctx, "visit:uv:total").Result()

	response.Success("获取统计成功", gin.H{
		"today_uv":  todayUV,
		"today_pv":  todayPV,
		"weekly_uv": weekUV,
		"total_uv":  totalUV,
	}, c)
}
