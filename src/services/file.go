package services

import (
	"bytes"
	"fmt"
	"strings"

	"io"
	"mime/multipart"

	"baliance.com/gooxml/document"
	"github.com/ledongthuc/pdf"
)

type FileService struct {
	ChapterServices IChapterService
}
type IFileService interface {
	GetDocx_DocData(file *multipart.FileHeader) error
	GetPdfData(file *multipart.FileHeader) error
}

func NewFileService(chapterServices IChapterService) IFileService {
	return &FileService{
		ChapterServices: chapterServices,
	}
}

func (f *FileService) GetDocx_DocData(file *multipart.FileHeader) error {

	fmt.Println("GetDocx_DocData func")
	openedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(fileBytes)

	doc, err := document.Read(reader, reader.Size())
	if err != nil {
		return err
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text() + "\n"
		}

	}

	err = f.ChapterServices.ChapterrizedText(text)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileService) GetPdfData(file *multipart.FileHeader) error {

	openedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(fileBytes)

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return err
	}

	var alltext string
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		fmt.Println("page: ", i)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			return err
		}

		alltext += content

	}

	alltext = strings.ReplaceAll(alltext, "\n", "")

	err = f.ChapterServices.ChapterrizedText(alltext)
	if err != nil {
		return err
	}

	return nil
}
