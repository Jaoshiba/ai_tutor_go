package services

import (
	"bytes"
	"fmt"
	"strings"

	repo "go-fiber-template/domain/repositories"

	"io"
	"mime/multipart"

	"baliance.com/gooxml/document"
	"github.com/ledongthuc/pdf"
)

type FileService struct {
	FileRepository repo.IFileRepository
}

type IFileService interface {
	GetDocx_DocData(file *multipart.FileHeader) (string, error)
	GetPdfData(file *multipart.FileHeader) (string, error)
}

func NewFileService(fileRepository repo.IFileRepository) IFileService {
	return &FileService{
		FileRepository: fileRepository,
	}
}

func (f *FileService) GetDocx_DocData(file *multipart.FileHeader) (string, error) {

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
		return "", nil
	}

	var text string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text += run.Text() + "\n"
		}

	}

	chap, err := ChapterrizedText(text)
	if err != nil {
		return "", err
	}

	fmt.Println("doc: ", chap)

	return "", nil
}

func (f *FileService) GetPdfData(file *multipart.FileHeader) (string, error) {

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

	pdfReader, err := pdf.NewReader(reader, reader.Size())
	if err != nil {
		return "", err
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
			return "", err
		}

		alltext += content

	}

	alltext = strings.ReplaceAll(alltext, "\n", "")

	chaps, err := ChapterrizedText(alltext)
	if err != nil {
		return "", err
	}

	fmt.Println("doc: ", chaps)

	// fmt.Print("alltext: ", alltext)

	return alltext, nil
}
