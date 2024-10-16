package ds

type Cards struct {
	ID              int `gorm:"primarykey"`
	Multiplier      string
	TitleEn         string
	TitleRu         string
	ImageUrl        string
	Description     string
	LongDescription string
	Status          bool
}
