package entities

import "time"

// RobotsCheck: type กลางที่ทุกแพ็กเกจใช้ “ตัวเดียวกัน”
type RobotsCheck struct {
	InputURL   string         `json:"input_url"`
	RobotsURL  string         `json:"robots_url"`
	StatusCode int            `json:"status_code"`
	UserAgent  string         `json:"user_agent"`
	TestedPath string         `json:"tested_path"`

	Allowed    bool           `json:"allowed"`
	Reason     string         `json:"reason"`
	CrawlDelay time.Duration  `json:"crawl_delay,omitempty"` // ใช้ value (ถ้าอยาก pointer เปลี่ยนได้)
	Sitemaps   []string       `json:"sitemaps,omitempty"`
	FetchedAt  time.Time      `json:"fetched_at"`
	RawRobots  string         `json:"raw_robots,omitempty"`
}
