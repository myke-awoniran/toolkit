package toolkit

import (
	"crypto/rand"
	"errors"
	"net/http"
	"strings"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_+"

type Tools struct {
	MaxFileSize    int
	AllowFileTypes []string
}

func (tool *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}

// uploaded file is a struct used to save information about the uploaded file
type uploadedFile struct {
	OriginalFileName string
	NewFileName      string
	fileSize         float64
}

func (tool *Tools) uploadFile(req http.Request, uploadDir string, rename ...bool) ([]*uploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles *[]uploadedFile

	if tool.MaxFileSize == 0 {
		tool.MaxFileSize = 1024 * 1024 * 1024
	}

	err := req.ParseMultipartForm(int64(tool.MaxFileSize))

	if err != nil {
		return nil, errors.New("file is too large")
	}

	for _, fHeaders := range req.MultipartForm.File {
		for _, hdr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles *[]uploadedFile) (*[]uploadedFile, error) {
				var uploadedFile uploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)

				_, err = infile.Read(buff)

				if err != nil {
					return nil, err
				}

				allowed := false
				fileType := http.DetectContentType(buff)

				//allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}

				if len(tool.AllowFileTypes) > 0 {
					for _, x := range tool.AllowFileTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("upload file type not permitted")
				}

				_, err = infile.Seek(0, 0)

				if err != nil {
					return nil, err
				}
				return uploadedFiles, nil
			}(uploadedFiles)
		}
	}
	//return uploadedFiles, nil

}
