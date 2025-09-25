package models

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	Order     Order `gorm:"foreignKey:OrderID"`
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int
	Price     float64
}