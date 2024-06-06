package grpcOrder_test

import (
	"context"
	orderv1 "github.com/batyrnosquare/protos/gen/go/order"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	grpcOrder "order/internal/grpc"
	"order/internal/storage"
	"testing"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*orderv1.CreateOrderResponse), args.Error(1)
}

func (m *MockOrderService) OrderByID(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*orderv1.GetOrderResponse), args.Error(1)

}

func (m *MockOrderService) DeleteOrder(ctx context.Context, req *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*orderv1.DeleteOrderResponse), args.Error(1)

}

func TestCreateOrder(t *testing.T) {
	mockOrderService := new(MockOrderService)
	orderService := grpcOrder.orderService{order: mockOrderService}
	ctx := context.Background()

	orderID := primitive.NewObjectID()
	mockOrderService.On("CreateOrder", ctx, mock.AnythingOfType("*models.Order")).Return(orderID, nil)

	req := &orderv1.CreateOrderRequest{
		UserId:     "user1",
		MenuItemId: "pizza1",
	}

	resp, err := orderService.CreateOrder(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, orderID.Hex(), resp.Id)

	mockOrderService.AssertExpectations(t)
}

func TestCreateOrder_AlreadyExists(t *testing.T) {
	mockOrderService := new(MockOrderService)
	orderService := grpcOrder.orderService{order: mockOrderService}
	ctx := context.Background()

	mockOrderService.On("CreateOrder", ctx, mock.AnythingOfType("*models.Order")).Return(primitive.NilObjectID, storage.ErrorOrderAlreadyExists)

	req := &orderv1.CreateOrderRequest{
		UserId:     "user1",
		MenuItemId: "pizza1",
	}

	resp, err := orderService.CreateOrder(ctx, req)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, status.Code(err), codes.AlreadyExists)

	mockOrderService.AssertExpectations(t)
}
