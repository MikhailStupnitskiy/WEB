package repository

import (
	"Evolution/internal/app/ds"
	"Evolution/internal/app/schemas"
	"Evolution/internal/token"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

type Repository struct {
	db *gorm.DB
	rc *redis.Client
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err = redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		rc: redisClient,
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

func (r *Repository) GetCardsIDsByMoveID(move_ID int) ([]int, error) {
	var CardsIds []int
	err := r.db.
		Model(&ds.MoveCards{}).
		Where("move_id = ?", move_ID).
		Pluck("card_id", &CardsIds).
		Error

	if err != nil {
		return nil, err
	}
	return CardsIds, nil
}

func (r *Repository) GetFoodByCardID(moveID int, cardID int) (int, error) {
	var CardsIDs ds.MoveCards

	err := r.db.Where("move_cards.move_id = ? and move_cards.card_id = ?", moveID, cardID).Find(&CardsIDs).Error

	if err != nil {
		return 0, err
	}

	return CardsIDs.Food, nil
}

func (r *Repository) CreateMove() (ds.Moves, error) {
	creator_id := 1
	newMove := ds.Moves{
		Status:     0,
		DateCreate: time.Now(),
		DateUpdate: time.Now(),
		CreatorID:  &creator_id,
		Stage:      "Кормление",
		Player:     "Игрок 1",
	}
	err := r.db.Create(&newMove).Error
	if err != nil {
		return ds.Moves{}, err
	}
	log.Println(1)
	move, err := r.GetLastMove()
	return move, err
}

func (r *Repository) GetLastMove() (ds.Moves, error) {
	var move ds.Moves
	err := r.db.Order("date_create DESC").Find(&move).Error
	if err != nil {
		return ds.Moves{}, err
	}
	return move, nil
}

func (r *Repository) AddToMove(move_ID int, card_ID int) error {
	moveCard := ds.MoveCards{
		MoveID: move_ID,
		CardID: card_ID,
	}
	err := r.db.Create(&moveCard).Error
	if err != nil {
		return fmt.Errorf("failed to add to move: %w", err)
	}

	return nil
}

func (r *Repository) GetCardByID(cardID string) (ds.Cards, error) {
	var card ds.Cards
	err := r.db.Where("id = ?", cardID).Find(&card).Error
	if err != nil {
		return ds.Cards{}, err
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
	var move ds.Moves
	if err := r.db.First(&move, "id = ?", id).Error; err != nil {
		return err
	}
	move.Status = 3 // Устанавливаем статус удаления
	if err := r.db.Save(&move).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateCard(card ds.Cards) error {
	if err := r.db.Create(&card).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteCardByID(id string) error {
	var card ds.Cards
	if err := r.db.First(&card, id).Error; err != nil {
		return err // Если запись не найдена или произошла ошибка, возвращаем её
	}
	card.ImageUrl = ""
	card.Status = false
	if err := r.db.Save(&card).Error; err != nil {
		return err // Возвращаем ошибку, если обновление не удалось
	}
	return nil // Возвращаем nil, если всё прошло успешно
}

func (r *Repository) UpdateCardByID(id string, card ds.Cards) error {
	var existingCard ds.Cards
	if err := r.db.First(&existingCard, "ID = ?", id).Error; err != nil {
		return err
	}

	existingCard.Multiplier = card.Multiplier
	existingCard.TitleEn = card.TitleEn
	existingCard.TitleRu = card.TitleRu
	existingCard.ImageUrl = card.ImageUrl
	existingCard.Description = card.Description
	existingCard.LongDescription = card.LongDescription
	existingCard.Status = card.Status

	err := r.db.Save(&existingCard).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangePicByID(id string, image string) error {
	// 1. Поиск записи по ID
	card := ds.Cards{}
	result := r.db.First(&card, "ID = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("запись с ID %s не найдена", id)
	}
	card.ImageUrl = image
	err := r.db.Save(&card).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteCardFromMove(id string, card_id int) error {
	var moveCard ds.MoveCards
	if err := r.db.Where("move_id = ? AND card_id = ?", id, card_id).First(&moveCard).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&moveCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateFoodMoveCard(id string, card_id int, food int) error {
	var moveCard ds.MoveCards
	if err := r.db.Where("move_id = ? AND card_id = ?", id, card_id).First(&moveCard).Error; err != nil {
		return err
	}
	moveCard.Food = food
	if err := r.db.Save(&moveCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAllMovesWithFilters(status int, having_status bool) ([]ds.Moves, error) {
	var moves []ds.Moves
	log.Println(status, having_status)
	db := r.db // Инициализируем db без фильтра по дате
	if having_status {
		db = db.Where("Status = ?", status) // Фильтр по статусу
	}
	err := db.Find(&moves).Error // Выборка записей из базы данных
	if err != nil {
		return nil, fmt.Errorf("failed to get move: %w", err)
	}
	return moves, nil
}

func (r *Repository) GetMoveByID(id int) (ds.Moves, error) { // ?
	move := ds.Moves{}
	err := r.db.Where("id = ?", id).First(&move).Error
	if err != nil {
		return ds.Moves{}, err
	}
	return move, nil
}

func (r *Repository) UpdateFieldsMove(request schemas.UpdateFieldsMoveRequest) error {
	var move ds.Moves
	// Загрузка записи из базы данных по ID
	if err := r.db.First(&move, "id = ?", request.ID).Error; err != nil {
		return err
	}
	if request.Player != "" {
		move.Player = request.Player
	}
	if request.Stage != "" {
		move.Stage = request.Stage
	}
	if err := r.db.Save(&move).Error; err != nil {
		return err
	}
	return nil // Возвращаем nil, если все прошло успешно
}

func (r *Repository) FormMove(id string) error {
	var move ds.Moves
	if err := r.db.First(&move, "id = ?", id).Error; err != nil {
		return err
	}
	if move.CreatorID == nil {
		err := fmt.Errorf("Unable to finish request. Probably some fields are empty")
		return err
	}
	move.Status = 1
	if err := r.db.Save(&move).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FinishMove(id string, status int) error {
	var move ds.Moves
	if err := r.db.First(&move, "id = ?", id).Error; err != nil {
		return err
	}
	if move.CreatorID == nil {
		err := fmt.Errorf("Unable to finish request. Probably some fields are empty")
		return err
	}
	mod_id := 2
	move.Status = status
	move.Cube = rand.IntN(12) + 1
	move.DateFinish = time.Now()
	move.ModeratorID = &mod_id
	if err := r.db.Save(&move).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateUser(user ds.Users) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) RegisterUser(user ds.Users) (error, int) {
	err := r.db.First(&user, "login = ?", user.Login).Error
	if err == nil {
		return err, 0 //пользователь с таким логином уже есть
	}
	if err = r.db.Create(&user).Error; err != nil {
		return err, 1 //пользователь создан
	}
	return nil, 2 // произошла ошибка с создание пользователя
}

func (r *Repository) LoginUser(user ds.Users) (error, string) {
	var db_user ds.Users
	err := r.db.First(&db_user, "login = ?", user.Login).Error
	if err != nil {
		return errors.New("Пользователь не существует"), ""
	}
	if user.Password != db_user.Password {
		return errors.New("Неверный пароль"), ""
	}

	currToken, err := token.GenerateJWTToken(db_user)
	if err != nil {
		return err, ""
	}
	err = r.SaveJWTToken(db_user.ID, currToken)
	if err != nil {
		return err, ""
	}
	return nil, currToken
}

func (r *Repository) SaveJWTToken(userID int, token string) error {
	exp := 1 * time.Hour
	userID_str := strconv.Itoa(userID)
	err := r.rc.Set(userID_str, token, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckActive(key string) (string, error) {
	result, err := r.rc.Get(key).Result()
	if err != nil {
		return "revoked", err
	}
	return result, err
}

func (r *Repository) LogoutUser(login string) error {
	var userInDB ds.Users
	err := r.db.First(&userInDB, "login = ?", login).Error
	if err != nil {
		return errors.New("Пользователь не существует")
	}

	userIDStr := strconv.Itoa(userInDB.ID)
	exp := 24 * time.Hour
	err = r.rc.Set(userIDStr, "revoked", exp).Err()
	if err != nil {
		return err
	}
	return nil
}
