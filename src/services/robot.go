// reunderstand

package services

import (
	"context"
	"fmt"
	"go-fiber-template/domain/entities" // สำคัญ: path ต้องตรงกับ go.mod + โฟลเดอร์จริง
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/temoto/robotstxt"
)

// CheckRobot: ใช้ entities.RobotsCheck ที่ประกาศไว้แล้ว (อย่าประกาศซ้ำใน package นี้)
func CheckRobot(link, userAgent string) (entities.RobotsCheck, error) {
	if strings.TrimSpace(userAgent) == "" {
		userAgent = "ai-tutor-bot"
	}

	u, err := url.Parse(strings.TrimSpace(link))
	if err != nil {
		return entities.RobotsCheck{}, fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if u.Host == "" {
		return entities.RobotsCheck{}, fmt.Errorf("url missing host")
	}
	u.Host = strings.ToLower(u.Host)

	testPath := u.EscapedPath()
	if testPath == "" {
		testPath = "/"
	}
	if u.RawQuery != "" {
		testPath = testPath + "?" + u.RawQuery
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	if err != nil {
		return entities.RobotsCheck{}, fmt.Errorf("create robots request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entities.RobotsCheck{}, fmt.Errorf("fetch robots.txt failed: %w", err)
	}
	defer resp.Body.Close()

	rc := entities.RobotsCheck{
		InputURL:   u.String(),
		RobotsURL:  robotsURL,
		StatusCode: resp.StatusCode,
		UserAgent:  userAgent,
		TestedPath: testPath,
		FetchedAt:  time.Now(),
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		rc.Allowed = true
		rc.Reason = "robots.txt not found (404) → allow all by convention"
		return rc, nil
	case http.StatusOK:
		// continue
	default:
		rc.Allowed = false
		rc.Reason = fmt.Sprintf("robots.txt returned non-200 (%d)", resp.StatusCode)
		return rc, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.RobotsCheck{}, fmt.Errorf("read robots.txt body failed: %w", err)
	}
	rc.RawRobots = string(body)

	robotsData, err := robotstxt.FromStatusAndBytes(resp.StatusCode, body)
	if err != nil {
		return entities.RobotsCheck{}, fmt.Errorf("parse robots.txt failed: %w", err)
	}

	group := robotsData.FindGroup(userAgent)
	if group.Test(testPath) {
		rc.Allowed = true
		rc.Reason = "allowed by robots.txt rules"
	} else {
		rc.Allowed = false
		rc.Reason = "disallowed by robots.txt rules"
	}

	// CrawlDelay เป็น time.Duration (ค่า 0 = ไม่กำหนด)
	if group.CrawlDelay > 0 {
		rc.CrawlDelay = group.CrawlDelay
	}

	if len(robotsData.Sitemaps) > 0 {
		rc.Sitemaps = append(rc.Sitemaps, robotsData.Sitemaps...)
	}

	return rc, nil
}
