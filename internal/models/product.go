package models

type Product struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Count       int64   `json:"count"`
}

type GetAllProduct struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
