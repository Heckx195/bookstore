package models

type Book struct {
	ID         int     `gorm:"primaryKey" json:"id"`
	Title      string  `gorm:"not null" json:"title"`
	AuthorID   int     `gorm:"not null" json:"author_id"`
	CategoryID int     `gorm:"not null" json:"category_id"`
	Price      float64 `gorm:"not null" json:"price"`
}
