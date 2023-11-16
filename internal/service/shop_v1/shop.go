package shop_v1

import (
	"context"
	"errors"

	"github.com/ARVG9866/uzum_shop/internal/convert"
	"github.com/ARVG9866/uzum_shop/internal/models"
	"github.com/ARVG9866/uzum_shop/internal/storage"
)

type IShopService interface {
	GetProduct(ctx context.Context, product_id int64) (*models.Product, error)
	GetProducts(ctx context.Context) ([]*models.GetAllProduct, error)
	DeleteProduct(ctx context.Context, product_id int64) error
	AddToBasket(ctx context.Context, basket *models.Basket) error
	UpdateBasket(ctx context.Context, basket *models.UpdateBasket) error
	GetBasket(ctx context.Context) ([]*models.Basket, error)
	CreateOrder(ctx context.Context, order *models.CreateOrder) (int64, error)
	CancelOrder(ctx context.Context, order_id int64) error
}

type service struct {
	storage storage.IStorage
}

func NewShopService(storage storage.IStorage) IShopService {
	return &service{
		storage: storage,
	}
}

func (s *service) GetProduct(ctx context.Context, product_id int64) (*models.Product, error) {
	res, err := s.storage.GetProduct(ctx, product_id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) GetProducts(ctx context.Context) ([]*models.GetAllProduct, error) {
	res, err := s.storage.GetAllProducts(ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) DeleteProduct(ctx context.Context, product_id int64) error {
	err := s.storage.DeleteFromBasket(ctx, product_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddToBasket(ctx context.Context, basket *models.Basket) error {
	res, err := s.GetProduct(ctx, basket.Product_id)
	if err != nil {
		return err
	}

	if res.Count < basket.Count {
		return errors.New("There are not enough products")
	}

	err = s.storage.CreateBasket(ctx, basket)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateBasket(ctx context.Context, basket *models.UpdateBasket) error {
	res, err := s.GetProduct(ctx, basket.Product_id)
	if err != nil {
		return err
	}

	if res.Count < basket.Count {
		return errors.New("There are not enough products")
	}

	err = s.storage.UpdateBasket(ctx, basket)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetBasket(ctx context.Context) ([]*models.Basket, error) {
	res, err := s.storage.GetAllBasket(ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) CreateOrder(ctx context.Context, order *models.CreateOrder) (int64, error) {
	if order.Coordinate_address_x == 0 && order.Coordinate_address_y == 0 {
		coord, err := s.getUserCoordinate(ctx)
		if err != nil {
			return 0, err
		}

		order.Coordinate_address_x = coord.X
		order.Coordinate_address_y = coord.Y
	}

	s.rewriteCoordinates(order.Coordinate_address_x, order.Coordinate_address_y)

	products, err := s.UpdateProductsForOrder(ctx)
	if err != nil {
		return 0, err
	}

	modelOrder := convert.GetModelOrder(order, products)

	res, err := s.storage.CreateOrder(ctx, modelOrder)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) CancelOrder(ctx context.Context, order_id int64) error {
	err := s.storage.DeleteOrder(ctx, order_id)
	if err != nil {
		return err
	}

	return nil
}

// func (s *service) updateProduct(ctx context.Context, product_id int64, count int64) error {
// 	err := s.storage.UpdateProduct(ctx, product_id, count)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *service) getUserCoordinate(ctx context.Context) (*models.Coordinate, error) {
	//todo
	return nil, errors.New("User doesn't have coordinates")
}

func (s *service) rewriteCoordinates(x float32, y float32) error {
	//todo
	return nil
}

func (s *service) UpdateProductsForOrder(ctx context.Context) ([]*models.OrderProduct, error) {
	basket, err := s.GetBasket(ctx)
	if err != nil {
		return nil, err
	}

	products, err := s.storage.UpdateBasketForOrder(ctx, basket)
	if err != nil {
		return nil, err
	}

	return products, nil
}
