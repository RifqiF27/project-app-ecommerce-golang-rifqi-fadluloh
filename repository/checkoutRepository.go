package repository

import (
	"context"
	"database/sql"
	"ecommerce/model"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type CheckoutRepository interface {
	GetAllCart(userID int) ([]*model.Checkout, error)
	GetTotalCart(userID int) (*model.Checkout, error)
	AddCart(cart model.Checkout) error
	DeleteCart(id, userID int) error
	UpdateCart(userID, productID, quantity int) (*model.Checkout, error)
	CreateOrder(userID int, productID []int) (*model.OrderResponse, error)
}

type checkoutRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewCheckoutRepository(db *sql.DB, logger *zap.Logger) CheckoutRepository {
	return &checkoutRepository{db: db, log: logger}
}

func (r *checkoutRepository) GetAllCart(userID int) ([]*model.Checkout, error) {
	query := `
	SELECT oi.id, p."name", p.images->>0 as image, oi.price, oi.quantity, oi.total 
	FROM products p 
	JOIN order_items oi ON p.id = oi.product_id 
	WHERE oi.user_id = $1 AND order_id IS NULL
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Checkout
	for rows.Next() {
		var result model.Checkout
		if err := rows.Scan(&result.ID, &result.Name, &result.Image, &result.Price, &result.Quantity, &result.TotalPrice); err != nil {
			r.log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}
func (r *checkoutRepository) GetTotalCart(userID int) (*model.Checkout, error) {
	query := `
	SELECT COUNT(*) FROM order_items
	WHERE user_id = $1 AND order_id IS NULL
	`
	var result model.Checkout
	err := r.db.QueryRow(query, userID).Scan(&result.TotalCarts)
	if err != nil {
		r.log.Error("Repository: failed to count total items in cart", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (r *checkoutRepository) AddCart(cart model.Checkout) error {
	query := `
	WITH DiscountedPrice AS (SELECT p.id AS product_id, p.price, p.price * (1 - CAST(COALESCE(CASE 
    WHEN MAX(wp.start_date) <= CURRENT_DATE AND MAX(wp.end_date) >= CURRENT_DATE THEN wp.discount_percentage
    ELSE 0 END, 0) AS FLOAT) / 100) AS discount_price
    FROM products p
    LEFT JOIN weekly_promotions wp ON p.id = wp.product_id
    WHERE p.id = $2
    GROUP BY p.id, p.price, wp.discount_percentage)
	INSERT INTO order_items (user_id, product_id, quantity, price, total)
	VALUES ($1, $2, 1, (SELECT discount_price FROM DiscountedPrice), (SELECT discount_price FROM DiscountedPrice))
	ON CONFLICT (user_id, product_id) 
	DO UPDATE SET 
	quantity = order_items.quantity + 1,
	total = (order_items.quantity + 1) * order_items.price
	RETURNING id
	`
	err := r.db.QueryRow(query, cart.UserID, cart.ProductID).Scan(&cart.ID)
	if err != nil {
		r.log.Error("Repository: Error executing query", zap.Error(err))
		return err
	}

	r.log.Info("Repository: cart added successfully", zap.Int("id", cart.ID))
	return nil
}

func (r *checkoutRepository) DeleteCart(id, userID int) error {
    query := `DELETE FROM order_items WHERE id = $1 AND user_id = $2 AND order_id IS NULL`
    res, err := r.db.Exec(query, id, userID)
    if err != nil {
        r.log.Error("Repository: Error executing query", zap.Error(err))
        return err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        r.log.Warn("Repository: No cart found for the given userID and id")
        return fmt.Errorf("no cart found")
    }

    r.log.Info("Repository: Cart deleted successfully", zap.Int("id", id), zap.Int("userID", userID))
    return nil
}

func (r *checkoutRepository) UpdateCart(userID, productID, quantity int) (*model.Checkout, error) {
	if quantity == 0 {
		// Delete the cart item
		query := `
			DELETE FROM order_items
			WHERE user_id = $1 AND product_id = $2 AND order_id IS NULL
			RETURNING id;
		`

		var deletedID int
		err := r.db.QueryRow(query, userID, productID).Scan(&deletedID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				r.log.Warn("Repository: No cart found for the given userID and productID", zap.Int("userID", userID), zap.Int("productID", productID), zap.Int("qty", quantity))
				return nil, fmt.Errorf("cart not found")
			}
			r.log.Error("Repository: Error deleting cart item", zap.Error(err))
			return nil, err
		}

		r.log.Info("Repository: cart item deleted successfully", zap.Int("id", deletedID))
		return nil, nil 
	}

	// Update the cart item
	query := `
		UPDATE order_items
		SET 
			quantity = $3::INTEGER,
			total = $3::INTEGER * price
		WHERE user_id = $1 AND product_id = $2 AND order_id IS NULL
		RETURNING id, quantity, total;
	`

	var result model.Checkout
	err := r.db.QueryRow(query, userID, productID, quantity).Scan(&result.ID, &result.Quantity, &result.TotalPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warn("Repository: No cart found for the given userID and productID", zap.Int("userID", userID), zap.Int("productID", productID), zap.Int("qty", quantity))
			return nil, fmt.Errorf("cart not found")
		}
		r.log.Error("Repository: Error updating cart quantity", zap.Error(err))
		return nil, err
	}

	r.log.Info("Repository: cart quantity updated successfully", zap.Int("id", result.ID))
	return &result, nil
}

func (r *checkoutRepository) CreateOrder(userID int, productID []int) (*model.OrderResponse, error) {
    ctx := context.Background()
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }

    // Step 1: Insert into `orders` table and get new order ID
    var orderID int
    var shippingAddress string
    insertOrderQuery := `
        WITH selected_items AS (
            SELECT 
                oi.user_id,
                SUM(oi.quantity * oi.price) AS total_price
            FROM 
                order_items oi
            WHERE 
                oi.user_id = $1
                AND oi.product_id = ANY($2)
                AND oi.order_id IS NULL
            GROUP BY oi.user_id
        )
        INSERT INTO orders (user_id, total_amount, shipping_address)
        SELECT si.user_id, si.total_price, u.address->>0 AS shipping_address
        FROM selected_items si
        JOIN users u ON si.user_id = u.id
        RETURNING id, shipping_address;
    `
    err = tx.QueryRowContext(ctx, insertOrderQuery, userID, productID).Scan(&orderID, &shippingAddress)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create order: %w", err)
    }

    // Step 2: Update `order_items` to link the new order_id
    updateOrderItemsQuery := `
        UPDATE order_items
        SET order_id = $1
        WHERE user_id = $2 
          AND product_id = ANY($3) 
          AND order_id IS NULL;
    `
    _, err = tx.ExecContext(ctx, updateOrderItemsQuery, orderID, userID, productID)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to update order items: %w", err)
    }

    // Step 3: Retrieve order details
    selectOrderDetailsQuery := `
        WITH selected_items AS (
            SELECT 
                p.name AS product_name,
                p.images->>0 AS image,
                oi.quantity,
                (oi.quantity * oi.price) AS subtotal_price
            FROM 
                order_items oi
            JOIN 
                products p ON oi.product_id = p.id
            WHERE 
                oi.user_id = $1
                AND oi.product_id = ANY($2)
                AND oi.order_id = $3
        )
        SELECT 
            si.product_name,
            si.image,
            si.subtotal_price,
            o.shipping
        FROM 
            selected_items si
        JOIN 
            orders o ON o.id = $3;
    `
    rows, err := tx.QueryContext(ctx, selectOrderDetailsQuery, userID, productID, orderID)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to retrieve order details: %w", err)
    }
    defer rows.Close()

    var items []model.OrderItem
    var shipping string
    var totalAmount float64

    for rows.Next() {
        var item model.OrderItem
        if err := rows.Scan(&item.ProductName, &item.Image, &item.SubtotalPrice, &shipping); err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to scan order item: %w", err)
        }
        totalAmount += item.SubtotalPrice
        items = append(items, item)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    // Build response
    response := &model.OrderResponse{
        OrderID:        orderID,
        Items:          items,
        Shipping:       shipping,
        TotalAmount:    totalAmount,
        ShippingAddress: shippingAddress, // Tambahkan alamat ke response
    }
    return response, nil
}
