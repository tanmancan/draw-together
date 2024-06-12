package detection

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
)

var loggerGenImg = basiclogger.BasicLogger{Namespace: "internal.detection.detect"}

type GenerateImageService interface {
	NewService(ctx context.Context) GenerateImageService
	SetBoardImageDataLayers([]*model.ImageData)
	GetBoardImageDataLayers() ([]*model.ImageData, error)
	BoardDrawingsToImageBuffer() (*bytes.Buffer, error)
}

type generatePngImageService struct {
	ctx          context.Context
	imgDataLayer []*model.ImageData
}

func NewGenerateImageService(ctx context.Context) GenerateImageService {
	s := generatePngImageService{}
	return s.NewService(ctx)
}

func (s *generatePngImageService) NewService(ctx context.Context) GenerateImageService {
	return &generatePngImageService{
		ctx:          ctx,
		imgDataLayer: nil,
	}
}

func (s *generatePngImageService) SetBoardImageDataLayers(dl []*model.ImageData) {
	s.imgDataLayer = dl
}

func (s *generatePngImageService) GetBoardImageDataLayers() ([]*model.ImageData, error) {
	if s.imgDataLayer == nil {
		return nil, errors.New("image data layers not found")
	}

	return s.imgDataLayer, nil
}

func (s *generatePngImageService) BoardDrawingsToImageBuffer() (*bytes.Buffer, error) {
	reqId := helpers.GetReqIdFromContext(s.ctx)
	dl, err := s.GetBoardImageDataLayers()
	if err != nil {
		loggerGenImg.LogDebug(reqId, "error generating board image", "error", err)
		return nil, err
	}

	var boardImg *image.RGBA
	for idx, d := range dl {
		imgDataReader := bytes.NewReader(d.Data)
		imgLayer, err := png.Decode(imgDataReader)
		if err != nil {
			loggerGenImg.LogError(reqId, "error decoding image layer. skipping layer.", "error", err, "layerIdx", idx)
			continue
		}

		if idx == 0 {
			boardImg = image.NewRGBA(imgLayer.Bounds())
			draw.Draw(
				boardImg,
				boardImg.Bounds(),
				&image.Uniform{
					color.White,
				},
				image.Point{},
				draw.Src,
			)
		}

		draw.Draw(
			boardImg,
			boardImg.Bounds(),
			imgLayer,
			imgLayer.Bounds().Min,
			draw.Over,
		)
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, boardImg)
	if err != nil {
		loggerGenImg.LogDebug(reqId, "error encoding png", "error", err)
		return nil, err
	}

	return buf, nil
}
