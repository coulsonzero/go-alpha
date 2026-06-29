package models

import "time"

type VisitorDaily struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Date      string    `gorm:"type:date;uniqueIndex;not null" json:"date"` // 统计日期
	UV        int64     `gorm:"default:0" json:"uv"`                        // 真实访客人数cookie + uuid唯一标识
	PV        int64     `gorm:"default:0" json:"pv"`                        // 总访问量
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (VisitorDaily) TableName() string {
	return "vistor"
}

func (VisitorDaily) Upsert(date string, uv, pv int64) {
	// Insert or update if date already exists
	DB.Where("date = ?", date).Assign(VisitorDaily{UV: uv, PV: pv}).FirstOrCreate(&VisitorDaily{
		Date: date,
		UV:   uv,
		PV:   pv,
	})
}

func (VisitorDaily) GetStats() *VisitorDaily {
	var stats VisitorDaily
	DB.Where("date = ?", time.Now().Format("2006-01-02")).First(&stats)
	return &stats
}
