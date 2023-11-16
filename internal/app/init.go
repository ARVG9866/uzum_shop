package app

import (
	"context"
	"log"

	"github.com/ARVG9866/uzum_shop/docs"
	"github.com/ARVG9866/uzum_shop/internal/api"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	pb_shop "github.com/ARVG9866/uzum_shop/pkg/shop_v1"
)

func (a *App) initDB() {
	sqlConnectionString := a.getSqlConnectionString()

	var err error
	a.db, err = sqlx.Open("postgres", sqlConnectionString)
	if err != nil {
		log.Fatal("failed to opening connection to db: ", err.Error())
	}

	if err = a.db.Ping(); err != nil {
		log.Fatal("failed to connect to the database: ", err.Error())
	}
}

func (a *App) initReDoc() {
	a.reDoc = docs.Initialize()

}

func (a *App) initHTTPServer(ctx context.Context) error {
	a.muxShop = gateway_runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb_shop.RegisterShopV1HandlerFromEndpoint(ctx, a.muxShop, a.appConfig.App.PortGRPC, opts) // why not PortHTTP?
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initGRPCServer() {
	a.grpcShopServer = grpc.NewServer()
	pb_shop.RegisterShopV1Server(
		a.grpcShopServer,
		&api.Shop{
			ShopService: a.getService(),
		},
	)
}
