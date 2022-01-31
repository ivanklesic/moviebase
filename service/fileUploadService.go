package service

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	uuid "github.com/nu7hatch/gouuid"
)

func SaveImage(img multipart.File, folderName string) (string, string, error) {
	fileHeader := make([]byte, 512)

	if _, err := img.Read(fileHeader); err != nil {
		return "", "", err
	}

	if _, err := img.Seek(0, 0); err != nil {
		return "", "", err
	}

	filename, err := uuid.NewV4()

	if err != nil {
		return "", "", err
	}

	filenameString := filename.String()

	var extension string

	switch http.DetectContentType(fileHeader) {
		case "image/png":
			extension = ".png"
		case "image/jpeg":
			extension = ".jpeg"
		default:
			return "", "Wrong file type, only .png and .jpg/.jpeg are accepted", nil
	}

	dst, err := os.Create("public/images/" + folderName + "/" + filenameString + extension)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, img); err != nil {
		return "", "", err
	}

	return filenameString + extension, "", nil
}
