package storage

import (
	"context"

	"github.com/ARVG9866/uzum_shop/internal/models"
)

type IStorage interface {
	GetProduct(ctx context.Context, product_id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.GetAllProduct, error)
	CreateBasket(ctx context.Context, basket *models.Basket) error
	DeleteFromBasket(ctx context.Context, product_id int64, user_id int64) error
	UpdateBasket(ctx context.Context, basket *models.UpdateBasket, user_id int64) error
	GetAllBasket(ctx context.Context, user_id int64) ([]*models.Basket, error)
	CreateOrder(ctx context.Context, order *models.Order, user_id int64) (int64, error)
	DeleteOrder(ctx context.Context, order_id int64) error
	// UpdateProduct(ctx context.Context, product_id int64, count int64) error
	UpdateBasketForOrder(ctx context.Context, basket []*models.Basket) ([]*models.OrderProduct, error)
	AddToOrder(ctx context.Context, products []*models.OrderProduct, order_id int64) error
	EmptyBasket(ctx context.Context, user_id int64) error
	GetUserCoordinate(ctx context.Context, user_id int64) (*models.Coordinate, error)
	UpdateUserCoordinate(ctx context.Context, user *models.Coordinate, user_id int64) error
}
