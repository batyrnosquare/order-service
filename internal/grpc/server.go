package grpcOrder

import (
	"context"
	"errors"
	orderv1 "github.com/batyrnosquare/protos/gen/go/order"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order/internal/models"
	"order/internal/storage"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) (primitive.ObjectID, error)
	OrderByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error)
	DeleteOrder(ctx context.Context, id primitive.ObjectID) error
}

type orderService struct {
	orderv1.UnimplementedOrderServiceServer
	order OrderService
}

func Register(gRPC *grpc.Server, order OrderService) {
	orderv1.RegisterOrderServiceServer(gRPC, &orderService{order: order})
}

func (os *orderService) CreateOrder(ctx context.Context, r *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	order := models.Order{
		UserID:  r.GetUserId(),
		PizzaID: r.GetMenuItemId(),
	}
	id, err := os.order.CreateOrder(ctx, &order)
	if err != nil {
		if errors.Is(err, storage.ErrorOrderAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "order already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create order")
	}
	return &orderv1.CreateOrderResponse{Id: id.Hex()}, nil
}

func (os *orderService) OrderByID(ctx context.Context, r *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	id, err := primitive.ObjectIDFromHex(r.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	order, err := os.order.OrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, "failed to get order")
	}

	var res orderv1.Order
	err = copier.Copy(&res, order)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &orderv1.GetOrderResponse{
		Id:         res.Id,
		Name:       res.Name,
		Status:     res.Status,
		UserId:     res.UserId,
		MenuItemId: res.MenuItemId,
		Quantity:   res.Quantity,
		CreatedAt:  res.CreatedAt}, nil
}

func (os *orderService) DeleteOrder(ctx context.Context, r *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {
	id, err := primitive.ObjectIDFromHex(r.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	err = os.order.DeleteOrder(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrorOrderNotFound) {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete order")
	}
	return &orderv1.DeleteOrderResponse{}, nil
}
