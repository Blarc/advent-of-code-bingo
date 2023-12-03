package repo

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
	"github.com/google/uuid"
)

type BingoCardRepo struct {
	*postgres.Postgres
}

func NewBingoCardRepo(pg *postgres.Postgres) *BingoCardRepo {
	return &BingoCardRepo{pg}
}

func (b *BingoCardRepo) GetBingoCard(id uuid.UUID) (entity.BingoCard, error) {
	var bingoCard entity.BingoCard

	result := b.DB.
		First(&bingoCard, id)

	if result.Error != nil {
		return entity.BingoCard{}, fmt.Errorf("BingoCardRepo - GetBingoCard - b.DB.First: %w", result.Error)
	}

	return bingoCard, nil
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

func (b *BingoCardRepo) AddBingoCardToUser(user *entity.User, bingoCard *entity.BingoCard) error {
	result := b.DB.
		Model(user).
		Association("BingoCards").
		Append(bingoCard)

	if result.Error != nil {
		return fmt.Errorf("BingoCardRepo - AddBingoCardToUser - b.DB.Model.Association.Append: %w", result.Error)
	}

	return nil
}

func (b *BingoCardRepo) RemoveBingoCardFromUser(user *entity.User, bingoCard *entity.BingoCard) error {
	result := b.DB.
		Model(user).
		Association("BingoCards").
		Delete(bingoCard)

	if result.Error != nil {
		return fmt.Errorf("BingoCardRepo - RemoveBingoCardFromUser - b.DB.Model.Association.Delete: %w", result.Error)
	}

	return nil
}
