package logic

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

const (
	imageDir = "images/"
	pngExt   = ".png"
)

type logic struct {
}

func NewLogic() *logic {
	return &logic{}
}

func (l *logic) GetImage(imageId string) (io.Reader, error) {
	f, err := os.Open(imageDir + imageId + pngExt)
	if err != nil {
		return nil, errors.Wrap(err, "os create error")
	}

	return f, nil
}

func (l *logic) CreateImage(file io.Reader) (string, error) {
	id := xid.New().String()

	dst, err := os.Create(imageDir + id + pngExt)
	if err != nil {
		return "", errors.Wrap(err, "os create error")

	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", errors.Wrap(err, "io copy error")
	}

	return id, nil
}
