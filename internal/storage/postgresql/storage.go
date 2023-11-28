package postgresql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ARVG9866/uzum_shop/internal/models"
	repo "github.com/ARVG9866/uzum_shop/internal/storage"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const productTable = "product"
const basketTable = "basket"
const orderTable = `"order"`
const orderProductTable = "order_product"
const userTable = "user"

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) repo.IStorage {
	return &storage{db: db}
}

func (s *storage) GetProduct(ctx context.Context, product_id int64) (*models.Product, error) {
	var product models.Product

	query := squirrel.Select("id", "name", "description", "price", "count").
		From(productTable).
		Where(squirrel.Eq{"id": product_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	err := query.QueryRowContext(ctx).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Count)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *storage) GetAllProducts(ctx context.Context) ([]*models.GetAllProduct, error) {
	var products []*models.GetAllProduct
	query := squirrel.Select("id", "name", "price").
		From(productTable).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var product models.GetAllProduct

		err = rows.Scan(&product.Id, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (s *storage) CreateBasket(ctx context.Context, basket *models.Basket) error {
	query := squirrel.Insert(basketTable).
		Columns("user_id", "product_id", "count").
		Values(basket.User_id, basket.Product_id, basket.Count).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *storage) DeleteFromBasket(ctx context.Context, product_id int64, user_id int64) error {
	query := squirrel.Delete(basketTable).
		Where(squirrel.Eq{"user_id": user_id, "product_id": product_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *storage) EmptyBasket(ctx context.Context, user_id int64) error {
	query := squirrel.Delete(basketTable).
		Where(squirrel.Eq{"user_id": user_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *storage) UpdateBasket(ctx context.Context, basket *models.UpdateBasket, user_id int64) error {
	query := squirrel.Update(basketTable).
		SetMap(map[string]interface{}{
			"count": basket.Count,
		}).
		Where(squirrel.Eq{"product_id": basket.Product_id, "user_id": user_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *storage) GetAllBasket(ctx context.Context, user_id int64) ([]*models.Basket, error) {
	query := squirrel.Select("id", "user_id", "product_id", "count").
		From(basketTable).
		Where(squirrel.Eq{"user_id": user_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	var basket []*models.Basket

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row models.Basket

		err = rows.Scan(&row.Id, &row.User_id, &row.Product_id, &row.Count)
		if err != nil {
			return nil, err
		}

		basket = append(basket, &row)
	}

	return basket, nil
}

func (s *storage) CreateOrder(ctx context.Context, order *models.Order, user_id int64) (int64, error) {
	query := squirrel.Insert(orderTable).
		Columns("user_id", "address", "coordinate_address_x", "coordinate_address_y", "coordinate_point_x",
			"coordinate_point_y", "create_at", "start_at", "courier_id", "delivery_status").
		Values(user_id, order.Address, order.Coordinate_address.X, order.Coordinate_address.Y, order.Coordinate_point.X,
			order.Coordinate_point.Y, order.Create_at, order.Start_at, order.Courier_id, order.Delivery_status).
		Suffix("RETURNING \"id\"").
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	var order_id int64

	err := query.QueryRowContext(ctx).Scan(&order_id)
	if err != nil {
		return 0, err
	}

	return order_id, nil
}

func (s *storage) DeleteOrder(ctx context.Context, order_id int64) error {
	str_query := fmt.Sprintf("SELECT product_id, count FROM %s WHERE order_id = %d", orderProductTable, order_id)
	rows, err := s.db.QueryContext(ctx, str_query)
	if err != nil {
		return err
	}

	defer rows.Close()

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	for rows.Next() {
		var product models.OrderProduct

		err = rows.Scan(&product.Product_id, &product.Count)
		if err != nil {
			return err
		}
		str_query = fmt.Sprintf("UPDATE %s SET count = count + %d WHERE id = %d",
			productTable, product.Count, product.Product_id)
		_, err := tx.ExecContext(ctx, str_query)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println()
			}
			return err
		}
	}

	query := squirrel.Delete(orderTable).
		Where(squirrel.Eq{"id": order_id}).
		RunWith(tx).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			log.Println()
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// func (s *storage) GetOrder(ctx context.Context, order_id int64) (*models.Order, error) {
// 	var order *models.Order

// 	query := squirrel.Select("user_id", "address", "coordinate_address_x", "coordinate_address_y", "coordinates_point_x",
// 		"coordinates_point_y", "created_at", "start_at", "delivery_at", "courier_id", "delivery_status").
// 		From(orderTable).
// 		Where(squirrel.Eq{"id": order_id}).
// 		RunWith(s.db).
// 		PlaceholderFormat(squirrel.Dollar)

// 	err := query.QueryRowContext(ctx).Scan(&order.User_id, order.Address, order.Coordinate_address_x,
// 		order.Coordinate_address_y, order.Coordinate_point_x, order.Coordinate_point_y, order.Create_at,
// 		order.Start_at, order.Delivery_at, order.Courier_id, order.Delivery_status)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return order, nil

// }

func (s *storage) AddToOrder(ctx context.Context, products []*models.OrderProduct, order_id int64) error {
	query := squirrel.Insert(orderProductTable).Columns("order_id", "product_id", "count", "price")

	for _, product := range products {
		query = query.Values(order_id, product.Product_id, product.Count, product.Price)
	}

	query = query.RunWith(s.db).PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (s *storage) UpdateBasketForOrder(ctx context.Context, basket []*models.Basket) ([]*models.OrderProduct, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}

	products := make([]*models.OrderProduct, 0, len(basket))

	for _, product := range basket {
		res, err := s.GetProduct(ctx, product.Product_id)
		if err != nil {
			return nil, err
		}

		if res.Count < product.Count {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println("Couldn't rollback")
			}
			return nil, errors.New("There are not enough products")
		}

		el := &models.OrderProduct{
			Product_id: product.Product_id,
			Count:      product.Count,
			Price:      res.Price,
		}

		products = append(products, el)

		query := squirrel.Update(productTable).
			SetMap(map[string]interface{}{
				"count": res.Count - product.Count,
			}).
			Where(squirrel.Eq{"id": product.Product_id}).
			RunWith(tx).
			PlaceholderFormat(squirrel.Dollar)

		_, err = query.ExecContext(ctx)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println()
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *storage) GetUserCoordinate(ctx context.Context, user_id int64) (*models.Coordinate, error) {
	var coordinate models.Coordinate

	query := squirrel.Select("id", "coordinate_address_x", "coordinate_address_y").
		From(userTable).
		Where(squirrel.Eq{"id": user_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	err := query.QueryRowContext(ctx).Scan(&coordinate.X, &coordinate.Y)
	if err != nil {
		return nil, err
	}

	return &coordinate, nil

}

func (s *storage) UpdateUserCoordinate(ctx context.Context, coordinate *models.Coordinate, user_id int64) error {
	query := squirrel.Update(userTable).
		SetMap(map[string]interface{}{
			"coordinate_address_x": coordinate.X,
			"coordinate_address_y": coordinate.Y,
		}).
		Where(squirrel.Eq{"user_id": user_id}).
		RunWith(s.db).
		PlaceholderFormat(squirrel.Dollar)

	if _, err := query.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}
