package schemas

import (
	"Evolution/internal/app/ds"
)

type InfoForMove struct {
	Card ds.Cards `json:"card"`
	Food int      `json:"food"`
}

type GetAllCardsResponse struct {
	ID    int        `json:"move_ID"`
	Count int        `json:"count"`
	Cards []ds.Cards `json:"cards"`
}

type GetCardResponse struct {
	Card ds.Cards `json:"card"`
}

type GetCardByNameResponse struct {
	Cards []ds.Cards `json:"cards"`
}

type CreateCardResponse struct {
	ID              int
	MessageResponse string
}

type DeleteCardResponse struct {
	ID              int
	MessageResponse string
}

type UpdateCardResponse struct {
	ID              int
	MessageResponse string
}

type AddCardToMoveResponse struct {
	CardID          int
	MoveID          int
	MessageResponse string
}

type GetAllMovesWithParamsResponse struct {
	Moves []ds.Moves
}

type GetAllMovesResponse struct {
	Moves []ds.Moves `json:"moves"`
}

type GetMoveResponse struct {
	Move      map[string]interface{} `json:"moves" form:"moves"`
	Status    int                    `json:"status" form:"status"`
	MoveCards []InfoForMove          `json:"move_cards"`
}

type DeleteCardFromMoveResponse struct{}

type UpdateOrderMoveCardsResponse struct{}

type ResponseMessage struct {
}
