package repo

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
)

type BingoBoardRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BingoBoardRepo {
	return &BingoBoardRepo{pg}
}

func (b *BingoBoardRepo) GetBingoBoard(shortUuid string) (*entity.BingoBoard, error) {
	var bingoBoard entity.BingoBoard
	result := b.DB.
		Preload("Users").
		First(&bingoBoard, "substring(id::text, 1, 16) = ?", shortUuid)
	if result.Error != nil {
		return nil, fmt.Errorf("BingoBoardRepo - GetBingoBoard - b.DB: %w", result.Error)
	}
	return &bingoBoard, nil
}

// GetBingoCardsWithCount returns all bingo cards for a board and how many users have them
func (b *BingoBoardRepo) GetBingoCardsWithCount(shortUuid string) ([]entity.BingoCardDto, error) {
	var bingoCards []entity.BingoCardDto
	result := b.DB.
		Table("bingo_board_bingo_card").
		Select("bingo_cards.id, bingo_cards.description, count(distinct user_id) as user_count").
		Joins("left join bingo_cards ON bingo_cards.id = bingo_board_bingo_card.bingo_card_id").
		Joins(`left join user_bingo_card ON user_bingo_card.bingo_card_id = bingo_board_bingo_card.bingo_card_id AND user_bingo_card.user_id IN (
						SELECT user_id
						FROM user_bingo_board
						WHERE substring(bingo_board_id::text, 1, 16) = ?
					)`, shortUuid).
		Where("substring(bingo_board_bingo_card.bingo_board_id::text, 1, 16) = ?", shortUuid).
		Group("bingo_cards.id").
		Scan(&bingoCards)

	if result.Error != nil {
		return nil, fmt.Errorf("BingoBoardRepo - GetBingoCardsWithCount - b.DB: %w", result.Error)
	}

	return bingoCards, nil
}

func (b *BingoBoardRepo) CreatePersonalBingoBoard(bingoBoard *entity.BingoBoard) (*entity.BingoBoard, error) {
	result := b.DB.
		Model(bingoBoard.Users[0]).
		Association("PersonalBingoBoard").
		Append(&bingoBoard)

	if result.Error != nil {
		return nil, fmt.Errorf("BingoBoardRepo - CreateBingoBoard - b.DB: %s", result.Error())
	}

	return bingoBoard, nil
}

func (b *BingoBoardRepo) DeletePersonalBingoBoard(user *entity.User) error {
	// TODO: Check if this also deletes the associations
	result := b.DB.
		Unscoped().
		Delete(&user.PersonalBingoBoard)

	if result.Error != nil {
		return fmt.Errorf("BingoBoardRepo - DeletePersonalBingoBoard - b.DB: %w", result.Error)
	}

	return nil
}

func (b *BingoBoardRepo) AddUserToBingoBoard(user *entity.User, shortUuid string) error {
	result := b.DB.
		Find(entity.BingoBoard{}, "substring(id::text, 1, 16) = ?", shortUuid).
		Association("Users").
		Append(&user)

	if result.Error != nil {
		return fmt.Errorf("BingoBoardRepo - AddUserToBingoBoard - b.DB: %s", result.Error())
	}

	return nil
}

func (b *BingoBoardRepo) RemoveUserFromBingoBoard(user *entity.User, shortUuid string) error {
	result := b.DB.
		Find(entity.BingoBoard{}, "substring(id::text, 1, 16) = ?", shortUuid).
		Association("Users").
		Delete(&user)

	if result.Error != nil {
		return fmt.Errorf("BingoBoardRepo - RemoveUserFromBingoBoard - b.DB: %s", result.Error())
	}

	return nil
}
