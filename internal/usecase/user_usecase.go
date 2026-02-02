package usecase

import (
	"errors"

	"github.com/jovan/mybanksoal-api/config"
	"github.com/jovan/mybanksoal-api/internal/entity"
	"github.com/jovan/mybanksoal-api/internal/repository"
	"github.com/jovan/mybanksoal-api/pkg/utils"
)

type UserUseCase interface {
	Register(username, password, role string) error
	Login(username, password string) (string, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewUserUseCase(userRepo repository.UserRepository, cfg *config.Config) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (u *userUseCase) Register(username, password, role string) error {
	// Check if user exists
	if _, err := u.userRepo.FindByUsername(username); err == nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
	}

	return u.userRepo.Create(user)
}

func (u *userUseCase) Login(username, password string) (string, error) {
	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, u.config)
	if err != nil {
		return "", err
	}

	return token, nil
}
