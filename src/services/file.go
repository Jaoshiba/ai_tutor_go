package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/gofiber/fiber/v2"
	"github.com/ledongthuc/pdf"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

func ReadFileData(docPath string, ctx *fiber.Ctx) (string, error) {
	file, err := os.Open(docPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return "", fmt.Errorf("ไม่สามารถอ่านข้อมูลไฟล์: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("ไม่สามารถรีเซ็ตตัวชี้ไฟล์: %w", err)
	}

	fileType := http.DetectContentType(fileHeader)
	fmt.Println("File type: ", fileType)

	switch {
	case strings.HasPrefix(fileType, "application/vnd.openxmlformats-officedocument.wordprocessingml.document"),
		strings.HasPrefix(fileType, "application/msword"):
		fmt.Println("Detected DOCX/DOC file.")
		return GetDocx_DocData(file, ctx)

	case strings.HasPrefix(fileType, "application/pdf"):
		fmt.Println("Detected PDF file.")
		return GetPdfData(file, ctx)

	case strings.HasPrefix(fileType, "text/html"):
		fmt.Println("Detected HTML file.")
		return GetHtmlData(file, ctx)

	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// ---------- HTML → ดึงลิงก์ PDF ที่ฝังอยู่ แล้วอ่านต่อด้วย GetPdfData ----------
func GetHtmlData(file *os.File, ctx *fiber.Ctx) (string, error) {
	fmt.Println("GetHtmlData func called with file:", file.Name())

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading HTML file bytes:", err)
		return "", err
	}
	htmlContent := string(fileBytes)

	// 1) พยายามหา URL ของ PDF ที่ฝังอยู่ในหน้า
	if pdfURL := extractPDFURLFromHTML(htmlContent); pdfURL != "" {
		fmt.Println("Found embedded PDF URL:", pdfURL)

		// ตั้งชื่อไฟล์จากชื่อ HTML เดิม
		title := strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name()))
		if title == "" {
			title = "document"
		}
		title = sanitizeFilename(title)

		// ดาวน์โหลด PDF จริง
		if err := os.MkdirAll("fileDocs", 0o755); err != nil {
			return "", fmt.Errorf("ไม่สามารถสร้างโฟลเดอร์ปลายทาง: %w", err)
		}
		destPath := filepath.Join("fileDocs", title+".pdf")
		if err := downloadFile(pdfURL, destPath); err != nil {
			return "", fmt.Errorf("ดาวน์โหลด PDF ที่ฝังอยู่ไม่สำเร็จ: %w", err)
		}

		// เปิด PDF แล้วส่งต่อเข้า pipeline เดิม
		pf, err := os.Open(destPath)
		if err != nil {
			return "", fmt.Errorf("เปิดไฟล์ PDF ที่ดาวน์โหลดไม่สำเร็จ: %w", err)
		}
		defer pf.Close()

		return GetPdfData(pf, ctx)
	}

	// 2) ไม่เจอ PDF — คืน HTML ดิบ (พฤติกรรมเดิม)
	return htmlContent, nil
}

// ---------- Extractors: หา URL ของ PDF จาก HTML ----------
var (
	reMetaPDF = regexp.MustCompile(`(?i)<meta[^>]+name=["'](?:citation_pdf_url|bepress_citation_pdf_url)["'][^>]+content=["']([^"']+)["']`)
	reHrefPDF = regexp.MustCompile(`(?i)<a[^>]+href=["']([^"']+\.pdf(?:\?[^"']*)?)["']`)
	reSrcPDF  = regexp.MustCompile(`(?i)\b(?:data|src)=["']([^"']+\.pdf(?:\?[^"']*)?)["']`)
)

func extractPDFURLFromHTML(htmlStr string) string {
	// 0) meta citation_pdf_url / bepress_citation_pdf_url
	if m := reMetaPDF.FindStringSubmatch(htmlStr); len(m) > 1 {
		return htmlUnescape(m[1])
	}
	// 1) DOM: object/embed/iframe ชี้ไป .pdf
	if u := scanDOMForPDF(htmlStr); u != "" {
		return u
	}
	// 2) <a href="...pdf">
	if m := reHrefPDF.FindStringSubmatch(htmlStr); len(m) > 1 {
		return htmlUnescape(m[1])
	}
	// 3) fallback: data|src="...pdf"
	if m := reSrcPDF.FindStringSubmatch(htmlStr); len(m) > 1 {
		return htmlUnescape(m[1])
	}
	return ""
}

func scanDOMForPDF(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	var find func(*html.Node) string
	find = func(n *html.Node) string {
		if n.Type == html.ElementNode {
			switch strings.ToLower(n.Data) {
			case "object", "embed", "iframe":
				if u := getAttr(n, "data"); strings.HasSuffix(strings.ToLower(u), ".pdf") {
					return htmlUnescape(u)
				}
				if u := getAttr(n, "src"); strings.HasSuffix(strings.ToLower(u), ".pdf") {
					return htmlUnescape(u)
				}
			case "meta":
				name := strings.ToLower(getAttr(n, "name"))
				if name == "citation_pdf_url" || name == "bepress_citation_pdf_url" {
					if u := getAttr(n, "content"); u != "" {
						return htmlUnescape(u)
					}
				}
			case "a":
				if u := getAttr(n, "href"); strings.HasSuffix(strings.ToLower(u), ".pdf") {
					return htmlUnescape(u)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if v := find(c); v != "" {
				return v
			}
		}
		return ""
	}
	return find(doc)
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

func htmlUnescape(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&#x2F;", "/",
		"&quot;", "\"",
		"&#39;", "'",
		"&lt;", "<",
		"&gt;", ">",
	)
	return replacer.Replace(s)
}

// ---------- Downloader ช่วยดาวน์โหลด PDF ----------
func downloadFile(fileURL, destPath string) error {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}
	// ใช้ UA ใกล้เคียงเบราว์เซอร์ ลดโอกาสโดนป้องกันบอท
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome Safari")
	req.Header.Set("Accept", "*/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 30x redirect ที่บางเว็บให้ .pdf ดาวน์โหลดผ่าน HTML ตัวกลาง
	if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		loc := resp.Header.Get("Location")
		if loc == "" {
			return fmt.Errorf("redirect แต่ไม่มี Location")
		}
		return downloadFile(loc, destPath)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ดาวน์โหลดล้มเหลว: status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

// ---------- PDF & DOCX ----------
func GetDocx_DocData(file *os.File, ctx *fiber.Ctx) (string, error) {
	fmt.Println("GetDocx_DocData func")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(fileBytes)

	doc, err := document.Read(reader, reader.Size())
	if err != nil {
		return "", err
	}

	var alltext string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			alltext += run.Text() + "\n"
		}
	}
	return alltext, nil
}

func GetPdfData(file *os.File, ctx *fiber.Ctx) (string, error) {
	fmt.Println("GetPdfData func call with file: ", file.Name())

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file bytes:", err)
		return "", err
	}
	reader := bytes.NewReader(fileBytes)

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		fmt.Println("Error creating PDF reader:", err)
		return "", err
	}

	var allText string
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Println("Error getting plain text from page:", err)
			return "", err
		}
		allText += content
	}
	allText = strings.ReplaceAll(allText, "\n", "")
	// fmt.Println("allText: ", allText)
	return allText, nil
}

// ---------- Utils ----------
func RemoveJsonBlock(text string) string {
	// 1) ตัดเคส ```json ... ```
	markdownJsonContentRegex := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	matches := markdownJsonContentRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}

	// 2) ถ้าไม่มีรั้วโค้ด ให้สแกนหา { ... } ก้อนแรกที่วงเล็บปิดครบ
	depth := 0
	start := -1
	for i, r := range text {
		if r == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if r == '}' {
			depth--
			if depth == 0 && start >= 0 {
				return text[start : i+1]
			}
		}
	}

	// 3) fallback: คืน text เดิม
	return text
}

func SaveFileToDisk(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error) {
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	downloadDir := "fileDocs"
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", err
	}
	docPath := filepath.Join(downloadDir, sanitizeFilename(file.Filename))
	diskFile, err := os.Create(docPath)
	if err != nil {
		return "", err
	}
	defer diskFile.Close()

	if _, err := io.Copy(diskFile, openedFile); err != nil {
		return "", err
	}
	return docPath, nil
}

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	illegal := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = illegal.ReplaceAllString(name, "-")
	if len(name) > 120 {
		name = name[:120]
	}
	if name == "" {
		name = "document"
	}
	return name
}

func extractContentsFromHTML(raw string) (string, error) {
	// ใช้ StrictPolicy เพื่อลบทุกอย่าง ยกเว้นข้อความธรรมดา
	p := bluemonday.StrictPolicy()
	cleanText := p.Sanitize(raw)

	// TrimSpace เพื่อลบช่องว่างที่เกินมา
	return strings.TrimSpace(cleanText), nil
}
