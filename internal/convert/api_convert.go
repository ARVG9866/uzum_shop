package convert

import (
	"time"

	"github.com/ARVG9866/uzum_shop/internal/models"
	pb_login "github.com/ARVG9866/uzum_shop/pkg/login_v1"
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
		Address: order.Address,
		Coordinate_address: &models.Coordinate{
			X: order.CoordinateAddress.X,
			Y: order.CoordinateAddress.Y,
		},
		Coordinate_point: &models.Coordinate{
			X: order.CoordinatePoint.X,
			Y: order.CoordinatePoint.Y,
		},
		Courier_id: order.CourierId,
	}
}

func GetModelOrder(createOrder *models.CreateOrder) *models.Order {
	return &models.Order{
		Address: createOrder.Address,
		Coordinate_address: &models.Coordinate{
			X: createOrder.Coordinate_address.X,
			Y: createOrder.Coordinate_address.Y,
		},
		Coordinate_point: &models.Coordinate{
			X: createOrder.Coordinate_point.X,
			Y: createOrder.Coordinate_point.Y,
		},
		Create_at:       time.Now(),
		Start_at:        time.Now(),
		Courier_id:      createOrder.Courier_id,
		Delivery_status: "New",
	}
}

func GetToken(auth *pb_login.Login_Response) *models.Token {
	return &models.Token{
		Refresh: auth.RefreshToken,
		Access:  auth.AccessToken,
	}
}
