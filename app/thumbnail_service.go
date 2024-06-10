package app

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"log/slog"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"strings"

	"golang.org/x/image/draw"
)

var MAX_WIDTH = 620

type ThumbnailService struct {
	repo repository.ThumbnailRepository
}

func NewThumbnailService(repo repository.ThumbnailRepository, bus *EventBus) *ThumbnailService {
	s := &ThumbnailService{
		repo: repo,
	}
	bus.Subscribe(s)
	return s
}

func (s *ThumbnailService) GetThumbnailForBinaryFileId(binaryFileId string) (*model.Thumbnail, error) {
	return s.repo.Get(binaryFileId)
}

func (s *ThumbnailService) JpegEncode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: 90})
}

func (s *ThumbnailService) CreateThumbnailForBinary(binary *model.BinaryFile) (*model.Thumbnail, error) {
	if strings.HasPrefix(binary.Mime(), "image/") {
		// determine data format
		format, _ := strings.CutPrefix(binary.Mime(), "image/")

		var encoder func(w io.Writer, m image.Image) error
		switch format {
		case "png":
			encoder = png.Encode
		case "jpeg":
			encoder = s.JpegEncode
		case "jpg":
			encoder = s.JpegEncode
		default:
			return nil, ErrBinaryFileUnsupportedFormat
		}

		img, _, err := image.Decode(bytes.NewReader(binary.Data))
		if err != nil {
			return nil, err
		}

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		var data []byte
		if imgWidth > MAX_WIDTH {
			w := MAX_WIDTH
			h := int(float32(imgHeight) * float32(MAX_WIDTH) / float32(imgWidth))
			if w == 0 {
				w = 1
			}
			// Set the expected size that you want:
			dst := image.NewRGBA(image.Rect(0, 0, w, h))
			// Resize:
			draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

			var b bytes.Buffer
			output := bufio.NewWriter(&b)
			encoder(output, dst)
			data = b.Bytes()
		} else {
			// no need for resizing, image is small enough
			data = binary.Data
		}
		return s.repo.Save(binary.Id, binary.Mime(), data)
	}
	return nil, ErrBinaryFileNotAnImage
}

func (s *ThumbnailService) NotifyBinaryFileDeleted(binaryFile model.BinaryFile) {
	slog.Info("Deleting Thumbnail for binary file", "id", binaryFile.Id)
	s.repo.Delete(binaryFile.Id)
}
