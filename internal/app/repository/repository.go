package repository

import (
	"Evolution/internal/app/ds"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetAllCards() ([]ds.Cards, error) {
	var cards []ds.Cards
	err := r.db.Where("status=true").Order("id ASC").Find(&cards).Error
	log.Println()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *Repository) GetCurrMove() ([]ds.Moves, error) {
	var currmove []ds.Moves
	err := r.db.Where("status=0").Find(&currmove).Error
	if err != nil {
		return nil, err
	}
	return currmove, nil
}

func (r *Repository) GetCardByInfo(cardText string) ([]ds.Cards, error) {
	var cards []ds.Cards
	err := r.db.Where("title_ru LIKE ?", "%"+cardText+"%").First(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *Repository) GetCardsIDsByMoveID(moveID int) ([]ds.MoveCards, error) {
	var CardsIDs []ds.MoveCards

	err := r.db.
		Where("move_cards.move_id = ?", moveID).Order("food ASC").Find(&CardsIDs).Error

	if err != nil {
		return nil, err
	}

	return CardsIDs, nil
}

func (r *Repository) CreateMove() ([]ds.Moves, error) {
	newMove := ds.Moves{
		Status:      0,
		DateCreate:  time.Now(),
		DateUpdate:  time.Now(),
		DateFinish:  time.Now(),
		CreatorID:   1,
		ModeratorID: 2,
	}
	err := r.db.Create(&newMove).Error
	if err != nil {
		return nil, err
	}
	move, err := r.GetLastMove()
	return move, err
}

func (r *Repository) GetLastMove() ([]ds.Moves, error) {
	var move []ds.Moves
	err := r.db.Order("date_create DESC").Find(&move).Error
	if err != nil {
		return nil, err
	}
	return move, nil
}

func (r *Repository) AddToMove(move_ID int, card_ID int) error {
	query := "INSERT INTO move_cards (move_id, card_id, food) VALUES (?, ?, 2)"
	err := r.db.Exec(query, move_ID, card_ID)
	if err != nil {
		return fmt.Errorf("failed to add to move: %w", err)
	}
	return nil
}

func (r *Repository) GetCardByID(id int) (*ds.Cards, error) { // ?
	card := &ds.Cards{}
	err := r.db.Where("id = ?", id).First(card).Error
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *Repository) GetMoveStatusByID(id int) (int, error) {
	move := &ds.Moves{}
	err := r.db.Where("id = ?", id).First(move).Error
	if err != nil {
		return -1, err
	}
	return move.Status, nil
}

func (r *Repository) GetCardsByID(cardID int) ([]ds.Cards, error) {
	var card []ds.Cards
	err := r.db.Where("id = ?", cardID).Find(&card).Error
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (r *Repository) DeleteMove(id int) error {
	err := r.db.Exec("UPDATE moves SET status = ? WHERE id = ?", 3, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateCard(card ds.Cards) error {
	return r.db.Create(card).Error
}
