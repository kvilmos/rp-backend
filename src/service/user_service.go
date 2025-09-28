package service

import (
	"errors"
	"room-planner/model"
	"room-planner/repository"
	"room-planner/request"
	"room-planner/util"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepository
	Db       *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		UserRepo: repo,
		Db:       db,
	}
}

func (s *UserService) RegisterUser(regReq request.RegisterRequest) (*model.User, error) {
	existingUser, err := s.UserRepo.GetByEmail(regReq.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email has already been taken")
	}

	hashedPassword, err := util.HashPassword(regReq.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Username: regReq.Username,
		Email:    regReq.Email,
		Password: hashedPassword,
	}

	err = s.UserRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) LoginUser(loginReq request.LoginRequest) (*string, *string, *model.User, error) {
	user, err := s.UserRepo.GetByEmail(loginReq.Email)
	if err != nil {
		return nil, nil, nil, err
	}

	if user == nil || !util.VerifyPassword(loginReq.Password, user.Password) {
		return nil, nil, nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := util.GenerateJWT(*user)
	if err != nil {
		return nil, nil, nil, err
	}

	return accessToken, refreshToken, user, err
}

func (s *UserService) GetUserById(id int64) (*model.User, error) {
	return s.UserRepo.GetById(id)
}
