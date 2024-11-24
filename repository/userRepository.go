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
	GetDetailUser(id int) (*model.User, error)
	UpdateUser(userID int, name, email, phone, password string, address []string) (*model.User, error)
	CreateAddress(userID int, newAddress string) (*model.User, error)
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
func (r *authRepository) GetDetailUser(id int) (*model.User, error) {
    var addressJSON []byte
    var user model.User

    err := r.DB.QueryRow(`SELECT name, email, phone, address FROM users WHERE id = $1`, id).
        Scan(&user.Name, &user.Email, &user.Phone, &addressJSON)
    if err != nil {
        if err == sql.ErrNoRows {
            r.Log.Warn("Repository: User not found", zap.Int("id", id))
            return nil, fmt.Errorf("user with id %d not found", id)
        }
        r.Log.Error("Repository: Failed to query user", zap.Error(err))
        return nil, err
    }
	r.Log.Info("Repository: Executing query", zap.Int("user_id", id))


    if len(addressJSON) > 0 {
        if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
            r.Log.Error("Repository: Failed to unmarshal address JSON", zap.String("address_json", string(addressJSON)), zap.Error(err))
            return nil, fmt.Errorf("invalid address format for user id %d: %w", id, err)
        }
    }

    r.Log.Debug("Repository: Retrieved user details", zap.Any("user", user))
    return &user, nil
}

func (r *authRepository) UpdateUser(userID int, name, email, phone, password string, address []string) (*model.User, error) {
    r.Log.Debug("Repository: Updating user details", zap.Int("userID", userID), zap.String("name", name), zap.String("email", email), zap.String("phone", phone), zap.Any("address", address))

    // Serialize the Address field into JSON
    addressJSON, err := json.Marshal(address)
    if err != nil {
        r.Log.Error("Repository: Failed to marshal address JSON", zap.Error(err))
        return nil, fmt.Errorf("failed to serialize address: %w", err)
    }

    // Update query
    query := `
        UPDATE users
        SET name = $1, email = $2, phone = $3, address = $4, password = $5, updated_at = NOW()
        WHERE id = $6
        RETURNING id, name, email, phone, address, password
    `

    // Execute the query and scan the updated user
    var updatedUser model.User
    err = r.DB.QueryRow(query, name, email, phone, addressJSON, password, userID).Scan(
        &updatedUser.ID,
        &updatedUser.Name,
        &updatedUser.Email,
        &updatedUser.Phone,
        &addressJSON,
		&updatedUser.Password,
    )
    if err != nil {
        r.Log.Error("Repository: Failed to update user", zap.Error(err))
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    // Deserialize the JSON address field back into the struct
    if err := json.Unmarshal(addressJSON, &updatedUser.Address); err != nil {
        r.Log.Error("Repository: Failed to unmarshal address JSON", zap.Error(err))
        return nil, fmt.Errorf("failed to deserialize address: %w", err)
    }

    r.Log.Debug("Repository: User details updated successfully", zap.Int("userID", updatedUser.ID))
    return &updatedUser, nil
}

func (r *authRepository) CreateAddress(userID int, newAddress string) (*model.User, error) {
    // Ambil alamat yang ada
    var addressJSON []byte
    err := r.DB.QueryRow(`SELECT address FROM users WHERE id = $1`, userID).Scan(&addressJSON)
    if err != nil {
        if err == sql.ErrNoRows {
            // Jika tidak ada data alamat, kita akan buat array kosong
            addressJSON = []byte("[]")
        } else {
            return nil, fmt.Errorf("failed to fetch user address: %w", err)
        }
    }

    // Unmarshal JSON ke dalam array alamat
    var address []string
    if len(addressJSON) > 0 {
        if err := json.Unmarshal(addressJSON, &address); err != nil {
            return nil, fmt.Errorf("failed to unmarshal address JSON: %w", err)
        }
    }

    // Tambahkan alamat baru ke array
    address = append(address, newAddress)

    // Serialize kembali address ke JSON
    addressJSON, err = json.Marshal(address)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal address: %w", err)
    }

    // Update alamat di database
    query := `
        UPDATE users
        SET address = $1
        WHERE id = $2
        RETURNING id, name, email, phone, address, updated_at
    `
    
    var user model.User
    err = r.DB.QueryRow(query, addressJSON, userID).Scan(
        &user.ID, &user.Name, &user.Email, &user.Phone, &addressJSON, &user.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to update user address: %w", err)
    }
	if err := json.Unmarshal(addressJSON, &user.Address); err != nil {
        r.Log.Error("Repository: Failed to unmarshal address JSON", zap.Error(err))
        return nil, fmt.Errorf("failed to deserialize address: %w", err)
    }

    return &user, nil
}







