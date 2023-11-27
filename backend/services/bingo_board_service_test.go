package services

import (
	"github.com/Blarc/advent-of-code-bingo/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFindBingoBoard(t *testing.T) {
	t.Run("test normal case controller find bingo board", func(t *testing.T) {
		bingoBoardRepoMock := new(mocks.BingoBoardRepoMock)
		bingoBoardRepoMock.On("FindBingoBoard", mock.AnythingOfType("string")).Return(nil)

		bingoBoardService := NewBingoBoardService(bingoBoardRepoMock)
		bingoBoard, _ := bingoBoardService.FindBingoBoard("test")

		t.Run("test normal case controller find bingo board", func(t *testing.T) {
			assert.Equal(t, "BingoBoardTest", bingoBoard.Name)
		})
	})
}
