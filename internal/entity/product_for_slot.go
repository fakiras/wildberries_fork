package entity

// ProductForSlot represents a product for slot
type ProductForSlot struct {
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Discount int32  `json:"discount"`
	Image    string `json:"image"`
}