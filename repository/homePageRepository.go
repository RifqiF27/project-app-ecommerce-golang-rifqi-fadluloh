package repository

import (
	"database/sql"
	"ecommerce/model"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type HomePageRepository interface {
	GetAllProducts(name string, categoryID, limit, page int) ([]*model.Product, int, int, error)
	GetByIdProduct(id int) (*model.ProductID, error)
	GetAllCategories() ([]*model.Category, error)
	GetAllBanners() ([]*model.BannerWeeklyPromotionRecomment, error)
	GetAllBestSellingProducts(limit, page int) ([]*model.Product, int, int, error)
	GetAllWeeklyPromotionProducts() ([]*model.BannerWeeklyPromotionRecomment, error)
	GetAllRecommentsProducts() ([]*model.BannerWeeklyPromotionRecomment, error)
	AddWishlist(wishlist model.Wishlist) error
	DeleteWishlist(id, userID int) error
}

type homePageRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewHomePageRepository(db *sql.DB, logger *zap.Logger) HomePageRepository {
	return &homePageRepository{db: db, log: logger}
}

func (r *homePageRepository) GetAllProducts(name string, categoryID, limit, page int) ([]*model.Product, int, int, error) {

	query := `
	SELECT p.id, p.name, p.images->>0 AS thumbnail_image, p.price, 
	COALESCE(CASE WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
	ELSE 0 
	END, 0) AS discount_percentage,
 		CASE 
    WHEN COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) = 0 THEN 0
    ELSE CAST(p.price AS FLOAT) * (1 - CAST(COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) AS FLOAT) / 100.0) END AS discount_price,  
	COALESCE(AVG(r.rating), 0) as average_rating,
	COUNT(DISTINCT oi.order_id) AS sold,
	CASE WHEN CURRENT_DATE - p.created_at <= INTERVAL '30 days' THEN TRUE ELSE FALSE
    END AS is_new
	FROM products p
	JOIN categories c ON p.category_id = c.id
	LEFT JOIN ratings r ON p.id = r.product_id
	LEFT JOIN order_items oi ON p.id = oi.product_id
	LEFT JOIN weekly_promotions wp ON p.id = wp.product_id
	WHERE 1=1
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE 1=1
	`

	var params []interface{}
	paramIndex := 1

	if name != "" {
		query += ` AND p.name ILIKE $` + fmt.Sprint(paramIndex)
		countQuery += ` AND p.name ILIKE $` + fmt.Sprint(paramIndex)
		params = append(params, "%"+name+"%")
		paramIndex++
	}

	if categoryID > 0 {
		query += ` AND p.category_id = $` + fmt.Sprint(paramIndex)
		countQuery += ` AND p.category_id = $` + fmt.Sprint(paramIndex)
		params = append(params, categoryID)
		paramIndex++
	}

	query += ` GROUP BY p.id, p.name, p.images, p.price, wp.discount_percentage, p.created_at`
	query += ` ORDER BY p.id ASC`

	var totalItems int
	err := r.db.QueryRow(countQuery, params...).Scan(&totalItems)
	if err != nil {
		r.log.Error("Repository: failed to execute count query", zap.Error(err))
		return nil, 0, 0, err
	}

	totalPages := (totalItems + limit - 1) / limit
	offset := (page - 1) * limit
	query += ` LIMIT $` + fmt.Sprint(paramIndex) + ` OFFSET $` + fmt.Sprint(paramIndex+1)
	params = append(params, limit, offset)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		r.log.Error("Repository: failed to execute query", zap.Error(err))
		return nil, 0, 0, err
	}
	defer rows.Close()

	r.log.Info("Repository: executed query", zap.String("query", query), zap.Any("params", params))

	var results []*model.Product
	for rows.Next() {
		var result model.Product
		if err := rows.Scan(&result.ID, &result.Name, &result.ThumbnailImage, &result.Price, &result.Discount, &result.DiscountPrice,
			&result.AverageRating, &result.Sold, &result.IsNEW); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, 0, 0, err
		}
		results = append(results, &result)
	}

	return results, totalItems, totalPages, nil
}

func (r *homePageRepository) GetByIdProduct(id int) (*model.ProductID, error) {
	var product model.ProductID
	var imagesJSON []byte
	var variantJSON []byte

	query := `
	SELECT p.id, p.name, p.images, c.name as category_name, p.price,
	CASE WHEN c.variant = '{}' THEN NULL ELSE c.variant END AS variant,
    COALESCE(CASE WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
	ELSE 0 
	END, 0) AS discount_percentage,
 		CASE 
    WHEN COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) = 0 THEN 0
    ELSE CAST(p.price AS FLOAT) * (1 - CAST(COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) AS FLOAT) / 100.0) END AS discount_price,
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(DISTINCT oi.order_id) AS sold,
	CASE WHEN CURRENT_DATE - p.created_at <= INTERVAL '30 days' THEN TRUE ELSE FALSE
    END AS is_new
	FROM products p
	JOIN categories c ON p.category_id = c.id
	LEFT JOIN ratings r ON p.id = r.product_id
	LEFT JOIN order_items oi ON p.id = oi.product_id
	LEFT JOIN weekly_promotions wp ON p.id = wp.product_id
	WHERE p.id = $1
	GROUP BY p.id, p.name, p.images, c."name", p.price, c.variant, wp.discount_percentage, p.created_at
	ORDER BY p.id ASC;
	`
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &imagesJSON, &product.Category, &product.Price, &variantJSON, &product.Discount, &product.DiscountPrice, &product.AverageRating, &product.Sold, &product.IsNEW)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Error("Repository: failed to scan row", zap.Error(err))
		return nil, err
	}

	if len(imagesJSON) > 0 {
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			r.log.Error("Repository: failed to unmarshal images JSON", zap.Error(err))
			return nil, err
		}
	}
	if len(variantJSON) > 0 {
		if err := json.Unmarshal(variantJSON, &product.Variant); err != nil {
			r.log.Error("Repository: failed to unmarshal variant JSON", zap.Error(err))
			return nil, err
		}
	}

	return &product, nil
}

func (r *homePageRepository) GetAllCategories() ([]*model.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *homePageRepository) GetAllBanners() ([]*model.BannerWeeklyPromotionRecomment, error) {
	rows, err := r.db.Query(`SELECT id, image, title, subtitle, path_page FROM banners`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.BannerWeeklyPromotionRecomment
	for rows.Next() {
		var result model.BannerWeeklyPromotionRecomment
		if err := rows.Scan(&result.ID, &result.Image, &result.Title, &result.Subtitle, &result.PathPage); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

func (r *homePageRepository) GetAllBestSellingProducts(limit, page int) ([]*model.Product, int, int, error) {

	query := `
	SELECT p.id, p.name, p.images->>0 AS thumbnail_image, p.price, 
    COALESCE(CASE WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
	ELSE 0 
	END, 0) AS discount_percentage,
 		CASE 
    WHEN COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) = 0 THEN 0
    ELSE CAST(p.price AS FLOAT) * (1 - CAST(COALESCE(
    CASE 
        WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
        ELSE 0
    END, 0) AS FLOAT) / 100.0) END AS discount_price,
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(DISTINCT oi.order_id) AS sold,
	CASE WHEN CURRENT_DATE - p.created_at <= INTERVAL '30 days' THEN TRUE ELSE FALSE
    END AS is_new 
	FROM products p
	JOIN categories c ON p.category_id = c.id
	LEFT JOIN ratings r ON p.id = r.product_id
	LEFT JOIN order_items oi ON p.id = oi.product_id
	LEFT JOIN orders o ON oi.order_id = o.id
	LEFT JOIN weekly_promotions wp ON p.id = wp.product_id
	WHERE DATE_TRUNC('month', o.created_at) = DATE_TRUNC('month', CURRENT_DATE)
	`

	countQuery := `
		SELECT COUNT(*)
        FROM (
            SELECT p.id
            FROM products p
            JOIN categories c ON p.category_id = c.id
            LEFT JOIN order_items oi ON p.id = oi.product_id
            LEFT JOIN orders o ON oi.order_id = o.id
            WHERE DATE_TRUNC('month', o.created_at) = DATE_TRUNC('month', CURRENT_DATE)
            GROUP BY p.id
        ) AS filtered_products
	`
	var params []interface{}
	paramIndex := 1

	query += ` GROUP BY p.id, p.name, p.images, p.price, wp.discount_percentage, p.created_at`
	query += ` ORDER BY sold DESC, p.id ASC`

	var totalItems int
	err := r.db.QueryRow(countQuery, params...).Scan(&totalItems)
	if err != nil {
		r.log.Error("Repository: failed to execute count query", zap.Error(err))
		return nil, 0, 0, err
	}

	totalPages := (totalItems + limit - 1) / limit
	offset := (page - 1) * limit
	query += ` LIMIT $` + fmt.Sprint(paramIndex) + ` OFFSET $` + fmt.Sprint(paramIndex+1)
	params = append(params, limit, offset)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		r.log.Error("Repository: failed to execute query", zap.Error(err))
		return nil, 0, 0, err
	}
	defer rows.Close()

	r.log.Info("Repository: executed query", zap.String("query", query), zap.Any("params", params))

	var results []*model.Product
	for rows.Next() {
		var result model.Product
		if err := rows.Scan(&result.ID, &result.Name, &result.ThumbnailImage, &result.Price, &result.Discount, &result.DiscountPrice,
			&result.AverageRating, &result.Sold, &result.IsNEW); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, 0, 0, err
		}
		results = append(results, &result)
	}

	return results, totalItems, totalPages, nil
}

func (r *homePageRepository) GetAllWeeklyPromotionProducts() ([]*model.BannerWeeklyPromotionRecomment, error) {
	query := `
	SELECT p.id, p.images->>0 AS thumbnail_image, p.title, p.subtitle 
	FROM products p
	JOIN weekly_promotions wp ON p.id = wp.product_id
	WHERE wp.start_date <= CURRENT_DATE AND wp.end_date >= CURRENT_DATE
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.BannerWeeklyPromotionRecomment
	for rows.Next() {
		var result model.BannerWeeklyPromotionRecomment
		if err := rows.Scan(&result.ID, &result.Image, &result.Title, &result.Subtitle); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}
func (r *homePageRepository) GetAllRecommentsProducts() ([]*model.BannerWeeklyPromotionRecomment, error) {
	query := `
	SELECT p.id, p.images->>0 AS thumbnail_image, p.title, p.subtitle
	FROM products p
	JOIN recomments r ON p.id = r.product_id
	ORDER BY p.id ASC;
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.BannerWeeklyPromotionRecomment
	for rows.Next() {
		var result model.BannerWeeklyPromotionRecomment
		if err := rows.Scan(&result.ID, &result.Image, &result.Title, &result.Subtitle); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

func (r *homePageRepository) AddWishlist(wishlist model.Wishlist) error {
	query := `INSERT INTO wishlists (user_id, product_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, wishlist.UserID, wishlist.ProductID).Scan(&wishlist.ID)
	if err != nil {
		r.log.Error("Repository: Error executing query", zap.Error(err))
		return fmt.Errorf("product not found or product already exists in wishlist")
	}

	r.log.Info("Repository: Wishlist added successfully", zap.Int("id", wishlist.ID))
	return nil
}

func (r *homePageRepository) DeleteWishlist(id, userID int) error {
	query := `DELETE FROM wishlists WHERE id = $1 AND user_id = $2`
	res, err := r.db.Exec(query, id, userID)
	if err != nil {
		r.log.Error("Repository: Error executing query", zap.Error(err))
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		r.log.Warn("Repository: No wishlist found for the given userID and id")
		return fmt.Errorf("no wishlist found")
	}

	r.log.Info("Repository: Wishlist deleted successfully", zap.Int("id", id), zap.Int("userID", userID))
	return nil
}
