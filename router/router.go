package router

import (
	"ecommerce/handler"
	middleware_auth "ecommerce/middleware"
	"ecommerce/service"
	"ecommerce/util"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(checkoutHandler *handler.CheckoutHandler, homePageHandler *handler.HomePageHandler, authHandler *handler.AuthHandler, authService service.AuthService, log *zap.Logger) (*chi.Mux, error) {

	r := chi.NewRouter()

	config := util.ReadConfiguration()

	log.Info("Loaded configuration", zap.String("AppName", config.AppName), zap.String("Port", config.Port))
	log.Debug("Loaded configuration", zap.String("AppName", config.AppName), zap.String("Port", config.Port))

	r.Use(middleware.Logger)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info("Request received", zap.String("method", r.Method), zap.String("url", r.URL.String()))

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			log.Info("Request processed", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.Duration("duration", duration))
		})
	})

	// fileServer := http.FileServer(http.Dir("./uploads"))
	// r.Handle("/uploads/*", http.StripPrefix("/uploads", fileServer))

	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.LoginHandler)
		r.Post("/register", authHandler.RegisterHandler)
		r.Post("/logout", authHandler.LogoutHandler)
	})

	authMiddleware := middleware_auth.NewAuthMiddleware(authService, log)

	r.Group(func(r chi.Router) {
		r.Route("/api/account", func(r chi.Router) {
			r.With(authMiddleware.Middleware).Get("/address", authHandler.GetAllAddressHandler)
			r.With(authMiddleware.Middleware).Get("/detail-user", authHandler.GetDetailUserHandler)
			r.With(authMiddleware.Middleware).Put("/update-user", authHandler.UpdateUserHandler)
			r.With(authMiddleware.Middleware).Put("/address", authHandler.UpdateAddressUserHandler)
			r.With(authMiddleware.Middleware).Delete("/address", authHandler.DeleteAddressHandler)
			r.With(authMiddleware.Middleware).Post("/address-default", authHandler.SetDefaultAddressUserHandler)
			r.With(authMiddleware.Middleware).Post("/address", authHandler.CreateAddressHandler)

		})
	})

	r.Group(func(r chi.Router) {
		r.Route("/api/products", func(r chi.Router) {
			r.Get("/", homePageHandler.GetAllProductsHandler)
			r.Get("/best-selling", homePageHandler.GetAllBestSellingProductsHandler)
			r.Get("/weekly-promotion", homePageHandler.GetAllWeeklyPromotionProductsHandler)
			r.Get("/recomments", homePageHandler.GetAllRecommentsProductsHandler)
			r.Get("/{id}", homePageHandler.GetByIdProductHandler)

			r.With(authMiddleware.Middleware).Get("/carts", checkoutHandler.GetAllCartHandler)
			r.With(authMiddleware.Middleware).Post("/carts", checkoutHandler.AddCartHandler)
			r.With(authMiddleware.Middleware).Post("/orders", checkoutHandler.CreateOrderHandler)
			r.With(authMiddleware.Middleware).Put("/carts/{id}", checkoutHandler.UpdateCartHandler)
			r.With(authMiddleware.Middleware).Delete("/carts/{id}", checkoutHandler.DeleteCartHandler)
			r.With(authMiddleware.Middleware).Get("/total-carts", checkoutHandler.GetTotalCartHandler)
			r.With(authMiddleware.Middleware).Post("/wishlist", homePageHandler.AddWishlistHandler)
			r.With(authMiddleware.Middleware).Delete("/wishlist/{id}", homePageHandler.DeleteWishlistHandler)

		})
		r.Route("/api/categories", func(r chi.Router) {
			r.Get("/", homePageHandler.GetAllCategoriesHandler)
		})
		r.Route("/api/banners", func(r chi.Router) {
			r.Get("/", homePageHandler.GetAllBannersHandler)
		})
	})

	return r, nil
}
