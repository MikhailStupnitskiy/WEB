package ds

type Cards struct {
	ID              int    `gorm:"primarykey" json:"id"`
	Multiplier      string `json:"multiplier"`
	TitleEn         string `json:"title_en"`
	TitleRu         string `json:"title_ru"`
	ImageUrl        string `json:"image_url"`
	Description     string `json:"description"`
	LongDescription string `json:"long_description"`
	Status          bool   `json:"status"`
}
