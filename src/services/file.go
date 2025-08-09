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
)

// IFileService defines the interface for file service operations.
// type IFileService interface {
// 	GetDocx_DocData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error)
// 	GetPdfData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error)
// 	ReadFileData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error)
// 	RemoveJsonBlock(text string) string
// }

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

	switch fileType {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/msword":
		fmt.Println("Detected DOCX/DOC file.")
		return GetDocx_DocData(file, ctx)
	case "application/pdf":
		fmt.Println("Detected PDF file.")
		return GetPdfData(file, ctx)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}

}

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
	fmt.Println("allText: ", allText)

	return allText, nil
}
func RemoveJsonBlock(text string) string {
	markdownJsonContentRegex := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	matches := markdownJsonContentRegex.FindStringSubmatch(text)

	if len(matches) > 1 {
		return matches[1]
	}
	return text
}

func SaveFileToDisk(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error) {
	// Open the file
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	// Create a new file on disk
	downloadDir := "fileDocs"
	docPath := filepath.Join(downloadDir, file.Filename)
	diskFile, err := os.Create(docPath)
	if err != nil {
		return "", err
	}
	defer diskFile.Close()

	// Copy the content to the new file
	if _, err := io.Copy(diskFile, openedFile); err != nil {
		return "", err
	}

	return docPath, nil
}
