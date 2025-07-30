package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/gofiber/fiber/v2"
	"github.com/ledongthuc/pdf"
)

// IFileService defines the interface for file service operations.
type IFileService interface {
	GetDocx_DocData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error)
	GetPdfData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error)
}

func GetDocx_DocData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error) {

	fmt.Println("GetDocx_DocData func")
	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
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

	// err = f.ChapterServices.ChapterrizedText(ctx, alltext)
	// if err != nil {
	// 	return "", err
	// }

	return alltext, nil
}

func GetPdfData(file *multipart.FileHeader, ctx *fiber.Ctx) (string, error) {

	openedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	fmt.Println("GetPdfData func call with file: ", file.Filename)

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return "", err
	}

	// fmt.Println("fileBytes: ", fileBytes)

	reader := bytes.NewReader(fileBytes)

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return "", err
	}

	var allText string // Renamed from alltext for consistency
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		fmt.Println("page: ", i)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		allText += content
	}

	allText = strings.ReplaceAll(allText, "\n", "")
	fmt.Println("allText: ", allText)

	// Pass the Fiber context (fCtx) to ChapterrizedText
	// err = f.ChapterServices.ChapterrizedText(ctx, allText) // fCtx is passed here
	// if err != nil {
	// 	return "", err
	// }

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
