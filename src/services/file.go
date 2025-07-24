package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"baliance.com/gooxml/document"
	"github.com/gofiber/fiber/v2" // Correctly imported and used
	"github.com/ledongthuc/pdf"
)

// FileService handles file processing and delegates to chapter services.
type FileService struct {
	ChapterServices IChapterService
}

// IFileService defines the interface for file service operations.
type IFileService interface {
	// Updated methods to accept *fiber.Ctx
	GetDocx_DocData(file *multipart.FileHeader, fCtx *fiber.Ctx) error
	GetPdfData(file *multipart.FileHeader, fCtx *fiber.Ctx) error
}

// NewFileService creates a new instance of FileService.
func NewFileService(chapterServices IChapterService) IFileService {
	return &FileService{
		ChapterServices: chapterServices,
	}
}

// GetDocx_DocData extracts text from .docx/.doc files and sends it for chapterization.
func (f *FileService) GetDocx_DocData(file *multipart.FileHeader, fCtx *fiber.Ctx) error { // Using fCtx for clarity

	fmt.Println("GetDocx_DocData func")
	openedFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open docx/doc file: %w", err) // Use %w for error wrapping
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return fmt.Errorf("failed to read docx/doc file content: %w", err)
	}

	reader := bytes.NewReader(fileBytes)

	// In Baliance gooxml/document, Read expects a seeker for the second argument (size is not required here).
	// If you just pass reader, it generally works for in-memory bytes.
	doc, err := document.Read(reader, reader.Size()) // reader.Size() is technically correct here for bytes.Reader
	if err != nil {
		return fmt.Errorf("failed to parse docx/doc document: %w", err)
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text() + "\n"
		}
	}

	// Pass the Fiber context (fCtx) to ChapterrizedText
	err = f.ChapterServices.ChapterrizedText(fCtx, text) // fCtx is passed here
	if err != nil {
		// If ChapterrizedText returns a fiber.Error, it should be propagated.
		// Otherwise, wrap it to provide more context.
		if _, ok := err.(*fiber.Error); ok {
			return err // Propagate Fiber errors directly
		}
		return fmt.Errorf("failed to chapterize docx/doc text: %w", err)
	}

	return nil
}

// GetPdfData extracts text from PDF files and sends it for chapterization.
func (f *FileService) GetPdfData(file *multipart.FileHeader, fCtx *fiber.Ctx) error { // Using fCtx for clarity

	openedFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return fmt.Errorf("failed to read PDF file content: %w", err)
	}

	reader := bytes.NewReader(fileBytes)

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return fmt.Errorf("failed to create PDF reader: %w", err)
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
			return fmt.Errorf("failed to get plain text from PDF page %d: %w", i, err)
		}
		allText += content
	}

	allText = strings.ReplaceAll(allText, "\n", "")

	// Pass the Fiber context (fCtx) to ChapterrizedText
	err = f.ChapterServices.ChapterrizedText(fCtx, allText) // fCtx is passed here
	if err != nil {
		// If ChapterrizedText returns a fiber.Error, it should be propagated.
		// Otherwise, wrap it to provide more context.
		if _, ok := err.(*fiber.Error); ok {
			return err // Propagate Fiber errors directly
		}
		return fmt.Errorf("failed to chapterize PDF text: %w", err)
	}

	return nil
}