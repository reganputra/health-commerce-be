package models

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"column:order_id;not null" json:"order_id"`
	ProductID uint    `gorm:"column:product_id;not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"column:quantity;not null" json:"quantity"`
	Price     float64 `gorm:"column:price;not null" json:"price"`
}
