package models

import "time"

type CreateOrder struct {
	Address            string      `json:"address"`
	Coordinate_address *Coordinate `json:"coordinate_address"`
	Coordinate_point   *Coordinate `json:"coordinate_point"`
	Courier_id         int64       `json:"courier_id"`
}

type Order struct {
	User_id            int64       `json:"user_id"`
	Address            string      `json:"address"`
	Coordinate_address *Coordinate `json:"coordinate_address"`
	Coordinate_point   *Coordinate `json:"coordinate_point"`
	Create_at          time.Time   `json:"create_at"`
	Start_at           time.Time   `json:"start_at"`
	Delivery_at        time.Time   `json:"delivery_at"`
	Courier_id         int64       `json:"courier_id"`
	Delivery_status    string      `json:"delivery_status"`
}

type Coordinate struct {
	X float32
	Y float32
}

type OrderProduct struct {
	Product_id int64   `json:"products_id"`
	Count      int64   `json:"count"`
	Price      float64 `json:"price"`
}
