package storage

import (
	"context"

	"github.com/ARVG9866/uzum_shop/internal/models"
)

type IStorage interface {
	GetProduct(ctx context.Context, product_id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.GetAllProduct, error)
	CreateBasket(ctx context.Context, basket *models.Basket) error
	DeleteFromBasket(ctx context.Context, product_id int64) error
	UpdateBasket(ctx context.Context, basket *models.UpdateBasket) error
	GetAllBasket(ctx context.Context) ([]*models.Basket, error)
	CreateOrder(ctx context.Context, order *models.Order) (int64, error)
	DeleteOrder(ctx context.Context, order_id int64) error
	// UpdateProduct(ctx context.Context, product_id int64, count int64) error
	UpdateBasketForOrder(ctx context.Context, basket []*models.Basket) ([]*models.OrderProduct, error)
}
