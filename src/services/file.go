package services

import (
	"bytes"
	"fmt"
	"go-fiber-template/domain/entities"
	"strings"

	"io"
	"mime/multipart"

	"baliance.com/gooxml/document"
	"github.com/ledongthuc/pdf"
)

type IFileService interface {
	GetDocx_DocData(file *multipart.FileHeader) ([]entities.ChapterDataModel, error)
	GetPdfData(file *multipart.FileHeader) ([]entities.ChapterDataModel, error)
}

func GetDocx_DocData(file *multipart.FileHeader) ([]entities.ChapterDataModel, error) {

	var chapters []entities.ChapterDataModel

	fmt.Println("GetDocx_DocData func")
	openedFile, err := file.Open()
	if err != nil {
		return chapters, err
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return chapters, err
	}

	reader := bytes.NewReader(fileBytes)

	doc, err := document.Read(reader, reader.Size())
	if err != nil {
		return chapters, err
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text() + "\n"
		}

	}

	chapters, err = ChapterrizedText(text)
	if err != nil {
		return chapters, err
	}

	return chapters, nil
}

func GetPdfData(file *multipart.FileHeader) ([]entities.ChapterDataModel, error) {

	var chapters []entities.ChapterDataModel

	openedFile, err := file.Open()
	if err != nil {
		return chapters, err
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return chapters, err
	}

	reader := bytes.NewReader(fileBytes)

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return chapters, err
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
			return chapters, err
		}

		alltext += content

	}

	alltext = strings.ReplaceAll(alltext, "\n", "")

	chapters, err = ChapterrizedText(alltext)
	if err != nil {
		return chapters, err
	}

	return chapters, nil
}
