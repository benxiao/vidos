package handler

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/kirillrdy/vidos/db"
	"github.com/kirillrdy/vidos/path"
	"github.com/kirillrdy/vidos/view"
)

func Upload(response http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
		http.Redirect(response, request, path.Root, http.StatusFound)
		return
	}

	//TODO fix assumption on buffer size
	request.ParseMultipartForm(1024 * 1024)
	form := request.MultipartForm
	formFiles := form.File[view.FormParamName]

	for _, formFile := range formFiles {
		processVideoFormFile(formFile)
	}

	http.Redirect(response, request, path.Root, http.StatusFound)
}

func processVideoFormFile(formFile *multipart.FileHeader) {

	log.Printf("Received %#v", formFile.Filename)

	//TODO does this needs to be closed ?
	file, err := formFile.Open()

	//TODO don't Fatal
	if err != nil {
		log.Fatal(err)
	}
	processVideoFromFile(file, formFile.Filename)
}

func processVideoFromFile(file io.ReadCloser, filename string) {

	video := db.Video{Filename: filename}
	db.Postgres.Save(&video)
	video.Save(file)

	//TODO Stop doing this as part of request
	video.CalculateDuration()
	video.GenerateThumbnail()

	//This can block so do in goroutine
	//TODO potentially dangerous
	go func() {
		db.EncodeVideo <- video.Id
	}()
}
