package fileHttpTransfer

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"vcfConverter/src/services/httpResponse"
)

const uploadPath = "./tmp"
const maxUploadSize = 200 * 1024 // 200 KB

func Up(w http.ResponseWriter, r *http.Request) (bool, string) {

	fileName := strconv.FormatInt(time.Now().UnixNano(), 10)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		httpResponse.RenderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return false, ""
	}

	file, header, err := r.FormFile("uploadFile")
	// fileHeader := make([]byte, 512)
	if err != nil {
		httpResponse.RenderError(w, "INVALID_FILE", http.StatusBadRequest)
		return false, ""
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		httpResponse.RenderError(w, "INVALID_FILE", http.StatusBadRequest)
		return false, ""
	}

	fileType := header.Header["Content-Type"][0]
	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		httpResponse.RenderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return false, ""
	}

	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	newFile, err := os.Create(newPath)
	if err != nil {
		httpResponse.RenderError(w, "CANT_SAVE_FILE", http.StatusInternalServerError)
		return false, ""
	}

	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil {
		httpResponse.RenderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return false, ""
	}

	return true, newPath
}
