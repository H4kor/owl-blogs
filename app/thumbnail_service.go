package app

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"strings"

	"golang.org/x/image/draw"
)

var MAX_WIDTH = 620

type ThumbnailService struct {
	repo repository.ThumbnailRepository
}

func NewThumbnailService(repo repository.ThumbnailRepository) *ThumbnailService {
	return &ThumbnailService{
		repo: repo,
	}
}

func (s *ThumbnailService) GetThumbnailForBinaryFileId(binaryFileId string) (*model.Thumbnail, error) {
	return s.repo.Get(binaryFileId)
}

func (s *ThumbnailService) CreateThumbnailForBinary(binary *model.BinaryFile) (*model.Thumbnail, error) {
	if strings.HasPrefix(binary.Mime(), "image") {
		img, _, err := image.Decode(bytes.NewReader(binary.Data))
		if err != nil {
			return nil, err
		}

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		var data []byte
		if imgWidth > MAX_WIDTH {
			h := MAX_WIDTH
			w := int(float32(imgHeight) * float32(MAX_WIDTH) / float32(imgWidth))
			if w == 0 {
				w = 1
			}
			// Set the expected size that you want:
			dst := image.NewRGBA(image.Rect(0, 0, w, h))
			// Resize:
			draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

			var b bytes.Buffer
			output := bufio.NewWriter(&b)
			jpeg.Encode(output, dst, nil)
			data = b.Bytes()
		} else {
			// no need for resizing, image is small enough
			data = binary.Data
		}
		return s.repo.Save(binary.Id, "image/jpeg", data)
	}
	return nil, ErrBinaryFileNotAnImage
}
