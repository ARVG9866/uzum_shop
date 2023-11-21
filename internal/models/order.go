package models

import "time"

type CreateOrder struct {
	Address              string  `json:"address"`
	Coordinate_address_x float32 `json:"coordinate_address_x"`
	Coordinate_address_y float32 `json:"coordinate_address_y"`
	Coordinate_point_x   float32 `json:"coordinate_point_x"`
	Coordinate_point_y   float32 `json:"coordinate_point_y"`
	Courier_id           int64   `json:"courier_id"`
}

type Order struct {
	User_id              int64     `json:"user_id"`
	Address              string    `json:"address"`
	Coordinate_address_x float32   `json:"coordinate_address_x"`
	Coordinate_address_y float32   `json:"coordinate_address_y"`
	Coordinate_point_x   float32   `json:"coordinate_point_x"`
	Coordinate_point_y   float32   `json:"coordinate_point_y"`
	Create_at            time.Time `json:"create_at"`
	Start_at             time.Time `json:"start_at"`
	Delivery_at          time.Time `json:"delivery_at"`
	Courier_id           int64     `json:"courier_id"`
	Delivery_status      string    `json:"delivery_status"`
}

type Coordinate struct {
	X float32
	Y float32
}

type OrderProduct struct {
	Product_id int64 `json:"products_id"`
	Count      int64 `json:"count"`
}
