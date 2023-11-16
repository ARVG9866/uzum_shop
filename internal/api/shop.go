package api

import (
	"context"

	"github.com/ARVG9866/uzum_shop/internal/convert"
	"github.com/ARVG9866/uzum_shop/internal/models"
	shop_v1 "github.com/ARVG9866/uzum_shop/internal/service/shop_v1"
	pb "github.com/ARVG9866/uzum_shop/pkg/shop_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Shop struct {
	pb.UnimplementedShopV1Server

	ShopService shop_v1.IShopService
}

func (s *Shop) GetProduct(ctx context.Context, req *pb.GetProduct_Request) (*pb.GetProduct_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	var product *models.Product
	product, err = s.ShopService.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	rtn := &pb.GetProduct_Response{
		Product: convert.ModelToPbProduct(product),
	}

	return rtn, nil
}

func (s *Shop) GetProducts(ctx context.Context, _ *emptypb.Empty) (*pb.GetProducts_Response, error) {
	res, err := s.ShopService.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*pb.ProductForGetAll, 0, len(res))

	for _, v := range res {
		products = append(products, convert.ModelToPbProductAll(v))
	}

	rtn := &pb.GetProducts_Response{
		Product: products,
	}

	return rtn, nil
}

func (s *Shop) AddProduct(ctx context.Context, req *pb.AddProduct_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = s.ShopService.AddToBasket(ctx, convert.PbToModelBasket(req.Basket))

	return &emptypb.Empty{}, err
}

func (s *Shop) UpdateBasket(ctx context.Context, req *pb.UpdateBasket_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = s.ShopService.UpdateBasket(ctx, convert.PbToModelUpdateBasket(req.Basket))

	return &emptypb.Empty{}, err
}

func (s *Shop) DeleteProduct(ctx context.Context, req *pb.DeleteProduct_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}
	err = s.ShopService.DeleteProduct(ctx, req.ProductId)

	return &emptypb.Empty{}, err
}

func (s *Shop) GetBasket(ctx context.Context, _ *emptypb.Empty) (*pb.GetBasket_Response, error) {
	res, err := s.ShopService.GetBasket(ctx)
	if err != nil {
		return nil, err
	}

	basket := make([]*pb.Basket, 0, len(res))

	for _, v := range res {
		basket = append(basket, convert.ModelToPbBasket(v))
	}

	rtn := &pb.GetBasket_Response{
		Basket: basket,
	}

	return rtn, nil
}

func (s *Shop) CreateOrder(ctx context.Context, req *pb.CreateOrder_Request) (*pb.CreateOrder_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	var res int64
	res, err = s.ShopService.CreateOrder(ctx, convert.PbToModelOrder(req.Order))

	return &pb.CreateOrder_Response{
		OrderId: res,
	}, err
}

func (s *Shop) CancelOrder(ctx context.Context, req *pb.CancelOrder_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = s.ShopService.CancelOrder(ctx, req.OrderId)

	return &emptypb.Empty{}, err
}

func (s *Shop) Healthz(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
