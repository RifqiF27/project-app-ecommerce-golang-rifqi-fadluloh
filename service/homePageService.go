package service

import (
	"ecommerce/model"
	"ecommerce/repository"
)

type HomePageService struct {
	Repo repository.HomePageRepository
}

func NewHomePageService(repo repository.HomePageRepository) HomePageService {
	return HomePageService{Repo: repo}
}

func (s *HomePageService) GetAllProductsService(name string, categoryID, limit, page int) ([]*model.Product, int, int, error) {
	return s.Repo.GetAllProducts(name, categoryID, limit, page)
}
func (s *HomePageService) GetByIdProductService(id int) (*model.ProductID, error) {
	return s.Repo.GetByIdProduct(id)
}
func (s *HomePageService) GetAllCategoriesService() ([]*model.Category, error) {
	return s.Repo.GetAllCategories()
}
func (s *HomePageService) GetAllBannersService() ([]*model.BannerWeeklyPromotionRecomment, error) {
	return s.Repo.GetAllBanners()
}
func (s *HomePageService) GetAllBestSellingProductsService(limit, page int) ([]*model.Product, int, int, error) {
	return s.Repo.GetAllBestSellingProducts(limit, page)
}
func (s *HomePageService) GetAllWeeklyPromotionProductsService() ([]*model.BannerWeeklyPromotionRecomment, error) {
	return s.Repo.GetAllWeeklyPromotionProducts()
}
func (s *HomePageService) GetAllRecommentsProductsService() ([]*model.BannerWeeklyPromotionRecomment, error) {
	return s.Repo.GetAllRecommentsProducts()
}
func (s *HomePageService) AddWishlistService(wishlist model.Wishlist)  error {
	return s.Repo.AddWishlist(wishlist)
}
func (s *HomePageService) DeleteWishlistService(id, userID int) error {
    return s.Repo.DeleteWishlist(id, userID)
}