package detection

import (
	"bytes"
	"context"
	"errors"
	"os"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/openapi/azurecv"
)

const API_VERSION = "2024-02-01"

var loggerAzurecv = basiclogger.BasicLogger{Namespace: "internal.detection.azure"}

type azurecvDetectionService[T azurecv.ImageAnalysisResult] struct {
	ctx    context.Context
	client *azurecv.APIClient
	reader *bytes.Buffer
}

func (s *azurecvDetectionService[T]) NewService(ctx context.Context) DetectionService[*azurecv.ImageAnalysisResult] {
	cfg := azurecv.NewConfiguration()
	client := azurecv.NewAPIClient(cfg)
	cfg.Servers = azurecv.ServerConfigurations{
		azurecv.ServerConfiguration{
			URL: config.AppConfig.DetectionConfig.AzureCvEndpoint,
		},
	}
	ctx = context.WithValue(
		ctx,
		azurecv.ContextAPIKeys,
		map[string]azurecv.APIKey{
			"apiKeyHeader": {
				Key: config.AppConfig.DetectionConfig.AzureCvKey,
			},
		},
	)
	return &azurecvDetectionService[T]{
		ctx:    ctx,
		client: client,
		reader: nil,
	}
}

func (s *azurecvDetectionService[T]) SetBoardImageReader(reader *bytes.Buffer) {
	s.reader = reader
}

func (s *azurecvDetectionService[T]) GetBoardImageReader() (*bytes.Buffer, error) {
	if s.reader == nil {
		return nil, errors.New("image reader is not set")
	}

	return s.reader, nil
}

func (s *azurecvDetectionService[T]) DetectImage() (*azurecv.ImageAnalysisResult, error) {
	reqID := helpers.GetReqIdFromContext(s.ctx)
	analyzeReq := s.client.DefaultAPI.ImageAnalysisAnalyze(s.ctx)
	analyzeReq = analyzeReq.ApiVersion(API_VERSION)
	analyzeReq = analyzeReq.Features(
		[]string{
			"caption",
		},
	)

	ir, err := s.GetBoardImageReader()
	if err != nil {
		loggerAzurecv.LogError(reqID, "error fetching image", "error", err)
		return nil, err
	}

	f, err := os.CreateTemp("", "*.png")
	if err != nil {
		loggerAzurecv.LogError(reqID, "error creating temp image", "error", err)
		return nil, err
	}

	defer (func() {
		err := f.Close()
		if err != nil {
			loggerAzurecv.LogError(reqID, "error closing resource", "error", err)
		}
		os.Remove(f.Name())
		if err != nil {
			loggerAzurecv.LogError(reqID, "error removing temp resource", "error", err)
		}
	})()

	os.WriteFile(f.Name(), ir.Bytes(), 0644)
	analyzeReq = analyzeReq.Body(f)

	res, _, err := analyzeReq.Execute()
	if err != nil {
		loggerAzurecv.LogError(reqID, "error analyzing image", "error", err)
		return nil, err
	}

	return res, nil
}
