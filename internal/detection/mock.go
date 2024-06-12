package detection

import (
	"bytes"
	"context"
	"errors"
)

type MockDetectionService[T string] struct {
	ctx    context.Context
	reader *bytes.Buffer
}

func (s *MockDetectionService[string]) NewService(ctx context.Context) DetectionService[string] {
	return &MockDetectionService[string]{
		ctx:    ctx,
		reader: nil,
	}
}

func (s *MockDetectionService[string]) SetBoardImageReader(reader *bytes.Buffer) {
	s.reader = reader
}

func (s *MockDetectionService[string]) GetBoardImageReader() (*bytes.Buffer, error) {
	return nil, errors.New("method not implemented")
}

func (s *MockDetectionService[string]) DetectImage() (string, error) {
	return "", errors.New("method not implemented")
}
