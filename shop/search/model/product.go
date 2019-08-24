package model

// Product is the Data Model for a Product Item
type Product struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ImageURL   string `json:"imageUrl"`
	Brand      string `json:"brand"`
	Price      int    `json:"price"`
	Discounted bool   `json:"discounted"`
}

// ProductList contains a slice of Products
type ProductList struct {
	Products []Product `json:"docs"`
}
