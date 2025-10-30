package models

type CartItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CartID    uint    `gorm:"column:cart_id;not null" json:"cart_id"`
	ProductID uint    `gorm:"column:product_id;not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"column:quantity;not null" json:"quantity"`
}
