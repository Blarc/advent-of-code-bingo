package repo

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
)

type BingoCardRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *BingoCardRepo {
	return &BingoCardRepo{pg}
}

func (b *BingoCardRepo) GetBingoCards() ([]entity.BingoCard, error) {
	var bingoCards []entity.BingoCard

	// TODO: Does this work, without .Model(&entity.BingoCard{})?
	result := b.DB.
		Find(&bingoCards)

	if result.Error != nil {
		return nil, fmt.Errorf("BingoCardRepo - GetBingoCards - b.DB.Find: %w", result.Error)
	}

	return bingoCards, nil
}

func (b *BingoCardRepo) GetBingoCardsWithUserCount() ([]entity.BingoCardDto, error) {
	var bingoCardsDto []entity.BingoCardDto

	result := b.DB.
		Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_card.bingo_card_id) as user_count").
		Joins("left join user_bingo_card on user_bingo_card.bingo_card_id = bingo_cards.id").
		Group("bingo_cards.id").
		Where("bingo_cards.public = true").
		Scan(&bingoCardsDto)

	if result.Error != nil {
		return nil, fmt.Errorf("BingoCardRepo - GetBingoCards - b.DB: %w", result.Error)
	}

	return bingoCardsDto, nil
}
