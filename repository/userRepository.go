package repository

import (
	"database/sql"
	"ecommerce/model"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type AuthRepository interface {
	Create(user *model.User) error
	GetUserLogin(user model.User) (*model.User, error)
	CreateSession(session *model.Session) error
	GetSessionByToken(token string) (*model.Session, error)
	DeleteSession(token string) error
	GetAllAddress(id int) ([]*model.User, error)
}

type authRepository struct {
	DB  *sql.DB
	Log *zap.Logger
}

func NewAuthRepository(db *sql.DB, logger *zap.Logger) AuthRepository {
	return &authRepository{DB: db, Log: logger}
}

func (r *authRepository) Create(user *model.User) error {
	var query string
	var err error

	if user.Email != "" {
		query = `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
		r.Log.Info("Repository: Executing query", zap.String("query", query), zap.Any("params", []interface{}{user.Name, user.Email, user.Password}))
		err = r.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
	} else if user.Phone != "" {
		query = `INSERT INTO users (name, phone, password) VALUES ($1, $2, $3) RETURNING id`
		r.Log.Info("Repository: Executing query", zap.String("query", query), zap.Any("params", []interface{}{user.Name, user.Phone, user.Password}))
		err = r.DB.QueryRow(query, user.Name, user.Phone, user.Password).Scan(&user.ID)
	} else {
		r.Log.Error("Repository: Validation failed", zap.String("reason", "either email or phone must be provided"))
		return fmt.Errorf("either email or phone must be provided")
	}

	if err != nil {
		r.Log.Error("Repository: Error executing query", zap.Error(err))
		return err
	}

	r.Log.Info("Repository: User created successfully", zap.Int("userID", user.ID))
	return nil

}

func (r *authRepository) GetUserLogin(user model.User) (*model.User, error) {
	query := `SELECT id, name, email, phone, password FROM users WHERE (email = $1 OR phone = $2) AND password = $3`
	r.Log.Info("Repository: Executing query", zap.String("query", query), zap.Any("params", []interface{}{user.Email, user.Phone, user.Password}))

	var userResponse model.User
	var email sql.NullString
	var phone sql.NullString
	err := r.DB.QueryRow(query, user.Email, user.Phone, user.Password).Scan(&userResponse.ID, &userResponse.Name, &email, &phone, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			r.Log.Warn("Repository: no user found with the provided credentials")
			return nil, fmt.Errorf("invalid username or password")
		}

		r.Log.Error("Repository: Error executing query", zap.Error(err))
		return nil, err
	}

	r.Log.Info("Repository: User login successful", zap.Int("userID", userResponse.ID))

	return &userResponse, nil
}

func (r *authRepository) CreateSession(session *model.Session) error {
	query := "INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)"
	_, err := r.DB.Exec(query, session.UserID, session.Token, session.ExpiresAt)
	if err != nil {
		return nil
	}
	return nil
}

func (r *authRepository) GetSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	query := "SELECT user_id, token, expires_at FROM sessions WHERE token=$1"
	err := r.DB.QueryRow(query, token).Scan(&session.UserID, &session.Token, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) DeleteSession(token string) error {
	query := "DELETE FROM sessions WHERE token=$1"
	res, err := r.DB.Exec(query, token)
	if err != nil {
		fmt.Println("Error executing delete:", err)
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println("Rows affected by delete:", rowsAffected)
	if rowsAffected == 0 {
		fmt.Println("No session found with this token.")
	}

	return nil
}

func (r *authRepository) GetAllAddress(id int) ([]*model.User, error) {
	var addressJSON []byte
	rows, err := r.DB.Query(`SELECT id, address FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.User
	for rows.Next() {
		var result model.User
		if err := rows.Scan(&result.ID, &addressJSON); err != nil {
			r.Log.Error("Repository: failed to scan row", zap.Error(err))
			return nil, err
		}
		if len(addressJSON) > 0 {
			if err := json.Unmarshal(addressJSON, &result.Address); err != nil {
				r.Log.Error("Repository: failed to unmarshal images JSON", zap.Error(err))
				return nil, err
			}
		}
		results = append(results, &result)
	}

	return results, nil
}
