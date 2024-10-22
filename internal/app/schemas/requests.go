package schemas

import (
	"Evolution/internal/app/ds"
	"time"
)

type GetAllCardsRequest struct{}

type GetCardRequest struct {
	ID string
}

type CreateCardRequest struct {
	ds.Cards
}

type DeleteCardRequest struct {
	ID string
}

type UpdateCardRequest struct {
	ID string
	ds.Cards
}

type AddCardToMoveRequest struct {
	ID string
}

type ChangePicRequest struct {
	ID       string `json:"id"`
	ImageUrl string `json:"image_link"`
}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMovesWithParamsRequest struct {
	HavingStatus bool      `json:"is_status"`
	Status       int       `json:"status"`
	FromDate     time.Time `json:"from_date"`
	ToDate       time.Time `json:"to_date"`
}

type GetMoveRequest struct {
	ID string
}

type UpdateOrderMoveRequest struct {
	ID     int `json:"move_id"`
	CardID int `json:"card_id"`
	OrderO int `json:"order_o"`
}

type UpdateFieldsMoveRequest struct {
	ID     string `uri:"move" json:"id"`
	Player string `json:"player"`
	Stage  string `json:"stage"`
}

type DeleteMoveRequest struct {
	ID string
}

type FormMoveRequest struct {
	ID string
}

type FinishMoveRequest struct {
	ID     string
	Status int `json:"status"`
}

type DeleteCardFromMoveRequest struct {
	ID     string
	CardID int `json:"card_id"`
}

type UpdateFoodMoveCardRequest struct {
	ID     string
	CardID int `json:"card_id"`
	Food   int `json:"food"`
}

type CreateUserRequest struct {
	ds.Users
}
