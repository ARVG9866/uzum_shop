package shop_v1

import (
	"context"
	"errors"

	"github.com/ARVG9866/uzum_shop/internal/convert"
	"github.com/ARVG9866/uzum_shop/internal/models"
	"github.com/ARVG9866/uzum_shop/internal/storage"
	pb_login "github.com/ARVG9866/uzum_shop/pkg/login_v1"
	"google.golang.org/grpc/metadata"
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
	Login(ctx context.Context, login string, password string) (*models.Token, error)
}

type service struct {
	storage     storage.IStorage
	loginClient pb_login.LoginV1Client
}

func NewShopService(storage storage.IStorage, loginClient pb_login.LoginV1Client) IShopService {
	return &service{
		storage:     storage,
		loginClient: loginClient,
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
	user_id, err := s.GetUserFromToken(ctx)
	if err != nil {
		return err
	}
	err = s.storage.DeleteFromBasket(ctx, product_id, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddToBasket(ctx context.Context, basket *models.Basket) error {
	// check if exist => update
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
	// check exist
	res, err := s.GetProduct(ctx, basket.Product_id)
	if err != nil {
		return err
	}

	if res.Count < basket.Count {
		return errors.New("There are not enough products")
	}

	user_id, err := s.GetUserFromToken(ctx)
	if err != nil {
		return nil
	}

	err = s.storage.UpdateBasket(ctx, basket, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetBasket(ctx context.Context) ([]*models.Basket, error) {
	user_id, err := s.GetUserFromToken(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.storage.GetAllBasket(ctx, user_id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *service) CreateOrder(ctx context.Context, order *models.CreateOrder) (int64, error) {
	user_id, err := s.GetUserFromToken(ctx)
	if err != nil {
		return 0, err
	}

	if order.Coordinate_address == nil || (order.Coordinate_address.X == 0 || order.Coordinate_address.Y == 0) {
		coord, err := s.getUserCoordinate(ctx, user_id)
		if err != nil {
			return 0, err
		}

		order.Coordinate_address = coord
	}

	s.rewriteCoordinates(ctx, order.Coordinate_address, user_id)

	products, err := s.UpdateProductsForOrder(ctx)
	if err != nil {
		return 0, err
	}

	modelOrder := convert.GetModelOrder(order)

	order_id, err := s.storage.CreateOrder(ctx, modelOrder, user_id)
	if err != nil {
		return order_id, err
	}

	err = s.storage.AddToOrder(ctx, products, order_id)
	if err != nil {
		return order_id, err
	}

	err = s.storage.EmptyBasket(ctx, user_id)
	if err != nil {
		return order_id, err
	}

	return order_id, nil
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

func (s *service) getUserCoordinate(ctx context.Context, user_id int64) (*models.Coordinate, error) {
	coordinate, err := s.storage.GetUserCoordinate(ctx, user_id)
	if err != nil {
		return nil, err
	}

	if coordinate != nil {
		return coordinate, nil
	}

	return nil, errors.New("User doesn't have coordinates")
}

func (s *service) rewriteCoordinates(ctx context.Context, coordinate *models.Coordinate, user_id int64) error {

	err := s.storage.UpdateUserCoordinate(ctx, coordinate, user_id)
	if err != nil {
		return err
	}

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

func (s *service) Login(ctx context.Context, login string, password string) (*models.Token, error) {
	req := &pb_login.Login_Request{Login: login, Password: password}
	auth, err := s.loginClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return convert.GetToken(auth), nil
}

func (s *service) GetUserFromToken(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("Can't get context")
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	check, err := s.loginClient.Check(ctx, &pb_login.Check_Request{EndpointAddress: ""})
	if err != nil {
		return 0, err
	}

	return check.UserId, nil
}
