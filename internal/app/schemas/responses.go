package schemas

import (
	"Evolution/internal/app/ds"
)

type GetAllCardsResponse struct {
	ID    int        `json:"move_ID"`
	Count int        `json:"count"`
	Cards []ds.Cards `json:"cards"`
}

type GetCardResponse struct {
	Card ds.Cards `json:"card"`
}

type CreateCardResponse struct{}

type DeleteCardResponse struct{}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMovesWithParamsResponse struct {
	Moves []ds.Moves
}

type GetAllMovesResponse struct {
	Moves []ds.Moves `json:"moves"`
}

type GetMoveResponse struct {
	Move      ds.Moves   `json:"moves"`
	Count     int        `json:"count"`
	MoveCards []ds.Cards `json:"move_cards"`
}

type DeleteCardFromMoveResponse struct{}

type UpdateOrderMilkReqMealsResponse struct{}
