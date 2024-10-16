package ds

type MoveCards struct {
	MoveID int   `gorm:"primaryKey"`
	CardID int   `gorm:"primaryKey"`
	Move   Moves `gorm:"primaryKey;foreignKey:MoveID"`
	Card   Cards `gorm:"primaryKey;foreignKey:CardID"`
	Food   int
}
