package service

import (
	"ecommerce/model"
	"ecommerce/repository"
)

type CheckoutService struct {
	Repo repository.CheckoutRepository
}

func NewCheckoutService(repo repository.CheckoutRepository) CheckoutService {
	return CheckoutService{Repo: repo}
}

func (s *CheckoutService) GetAllCartService(userID int) ([]*model.Checkout, error) {
	return s.Repo.GetAllCart(userID)
}
func (s *CheckoutService) GetTotalCartService(userID int) (*model.Checkout, error) {
	return s.Repo.GetTotalCart(userID)
}
func (s *CheckoutService)  AddCartService(cart model.Checkout) error {
	return s.Repo.AddCart(cart)
}
func (s *CheckoutService)  DeleteCartService(id, userID int) error {
	return s.Repo.DeleteCart(id, userID)
}
func (s *CheckoutService) UpdateCartService(userID, productID, quantity int) (*model.Checkout, error) {
	return s.Repo.UpdateCart(userID, productID, quantity)
}
func (s *CheckoutService) CreateOrderService(userID int, productID []int) (*model.OrderResponse, error) {
	return s.Repo.CreateOrder(userID, productID)
}