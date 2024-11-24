package service

import (
	"ecommerce/model"
	"ecommerce/repository"
	"errors"
	"time"
)

type AuthService struct {
	RepoUser repository.AuthRepository
}

func NewAuthService(repoUser repository.AuthRepository) AuthService {
	return AuthService{repoUser}
}

func (us *AuthService) LoginService(user model.User) (*model.User, error) {
	return us.RepoUser.GetUserLogin(user)
}

func (s *AuthService) RegisterService(user model.User) error {
	return s.RepoUser.Create(&user)
}
func (s *AuthService) CreateSessionService(user model.Session) error {
	return s.RepoUser.CreateSession(&user)
}


func (s *AuthService) Logout(token string) error {
	return s.RepoUser.DeleteSession(token)
}
func (s *AuthService) GetAllAddressService(id int) ([]*model.User, error) {
	return s.RepoUser.GetAllAddress(id)
}

func (s *AuthService) VerifyToken(token string) (int, error) {
	session, err := s.RepoUser.GetSessionByToken(token)
	if err != nil || session == nil || session.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("invalid or expired token")
	}
	return session.UserID, nil
}

