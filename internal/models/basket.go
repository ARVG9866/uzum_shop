package models

type Basket struct {
	Id         int64 `json:"id"`
	User_id    int64 `json:"user_id"`
	Product_id int64 `json:"product_id"`
	Count      int64 `json:"count"`
}

type UpdateBasket struct {
	Product_id int64 `json:"product_id"`
	Count      int64 `json:"count"`
}
