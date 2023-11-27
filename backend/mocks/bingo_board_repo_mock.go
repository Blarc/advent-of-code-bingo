package mocks

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type BingoBoardRepoMock struct {
	mock.Mock
}

func (m *BingoBoardRepoMock) FindBingoBoard(shortUuid string) (*models.BingoBoard, error) {
	if err := m.Called(shortUuid).Error(0); err != nil {
		return nil, err
	}

	return &models.BingoBoard{
		ID:   uuid.UUID{},
		Name: "BingoBoardTest",
	}, nil
}
