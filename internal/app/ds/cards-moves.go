package ds

type MoveCards struct {
	MoveID int   `gorm:"primaryKey" json:"move_id"`
	CardID int   `gorm:"primaryKey" json:"card_id"`
	Move   Moves `gorm:"primaryKey;foreignKey:MoveID"`
	Card   Cards `gorm:"primaryKey;foreignKey:CardID"`
	Food   int   `json:"food"`
}
