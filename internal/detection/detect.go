package detection

import (
	"bytes"
	"context"

	"github.com/tanmancan/openapi/azurecv"
)

type DetectionService[T interface{}] interface {
	NewService(ctx context.Context) DetectionService[T]
	SetBoardImageReader(reader *bytes.Buffer)
	GetBoardImageReader() (*bytes.Buffer, error)
	DetectImage() (T, error)
}

func NewDetectionService(ctx context.Context) DetectionService[*azurecv.ImageAnalysisResult] {
	s := azurecvDetectionService[azurecv.ImageAnalysisResult]{}
	return s.NewService(ctx)
}
