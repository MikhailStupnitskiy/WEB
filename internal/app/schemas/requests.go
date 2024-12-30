package schemas

import (
	"Evolution/internal/app/ds"
)

type GetAllCardsRequest struct{}

type GetCardRequest struct {
	ID string
}

type GetCardByNameRequest struct {
	Name string
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

type GetAllMovesWithParamsRequest struct {
	Status int `json:"status" form:"status"`
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
	ID    string `uri:"move" json:"id"`
	Stage string `json:"stage"`
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

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogoutUserRequest struct {
	Login string `json:"login"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
