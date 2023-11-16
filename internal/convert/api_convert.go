package convert

import (
	"time"

	"github.com/ARVG9866/uzum_shop/internal/models"
	pb "github.com/ARVG9866/uzum_shop/pkg/shop_v1"
)

func ModelToPbProduct(product *models.Product) *pb.Product {
	return &pb.Product{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Count:       product.Count,
	}
}

func ModelToPbProductAll(product *models.GetAllProduct) *pb.ProductForGetAll {
	return &pb.ProductForGetAll{
		Id:    product.Id,
		Name:  product.Name,
		Price: product.Price,
	}
}

func PbToModelBasket(basket *pb.Basket) *models.Basket {
	return &models.Basket{
		Id:         basket.Id,
		User_id:    basket.UserId,
		Product_id: basket.ProductId,
		Count:      basket.Count,
	}
}

func PbToModelUpdateBasket(basket *pb.BasketForUpdate) *models.UpdateBasket {
	return &models.UpdateBasket{
		Product_id: basket.ProductId,
		Count:      basket.Count,
	}
}

func ModelToPbBasket(basket *models.Basket) *pb.Basket {
	return &pb.Basket{
		Id:        basket.Id,
		UserId:    basket.User_id,
		ProductId: basket.Product_id,
		Count:     basket.Count,
	}
}

func PbToModelOrder(order *pb.Order) *models.CreateOrder {
	return &models.CreateOrder{
		Address:              order.Address,
		Coordinate_address_x: order.CoordinateAddressX,
		Coordinate_address_y: order.CoordinateAddressY,
		Coordinate_point_x:   order.CoordinatePointX,
		Coordinate_point_y:   order.CoordinatePointY,
		Courier_id:           order.CourierId,
	}
}

func GetModelOrder(createOrder *models.CreateOrder, products []*models.OrderProduct) *models.Order {
	return &models.Order{
		Products:             products,
		Address:              createOrder.Address,
		Coordinate_address_x: createOrder.Coordinate_address_x,
		Coordinate_address_y: createOrder.Coordinate_address_y,
		Coordinate_point_x:   createOrder.Coordinate_point_x,
		Coordinate_point_y:   createOrder.Coordinate_point_y,
		Create_at:            time.Now(),
		Start_at:             time.Now(),
		Courier_id:           createOrder.Courier_id,
		Delivery_status:      "New",
	}
}
