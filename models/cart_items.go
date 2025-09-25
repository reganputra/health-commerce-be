package models

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint
	Cart      Cart `gorm:"foreignKey:CartID"`
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int
}