package service

import (
	"room-planner/app"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/repository"
	"room-planner/request"
	"room-planner/token"
	"room-planner/util"

	"gorm.io/gorm"
)

type UserService struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
	Db          *gorm.DB
	JWTMaker    *token.JWTMaker
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository, sessionRepo repository.SessionRepository, jwtMaker *token.JWTMaker) *UserService {
	return &UserService{
		Db:          db,
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
		JWTMaker:    jwtMaker,
	}
}

func (s *UserService) RegisterUser(regReq request.RegisterRequest) (*model.User, error) {
	existingUser, err := s.UserRepo.GetByEmail(regReq.Email)
	if err != nil {
		return nil, err
	}

	if existingUser.Id != 0 {
		return nil, app.ErrAlreadyExist
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

func (s *UserService) LoginUser(loginReq request.LoginRequest) (*string, *string, *token.UserClaims, *model.User, error) {
	user, err := s.UserRepo.GetByEmail(loginReq.Email)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if user.Id == 0 || !util.VerifyPassword(loginReq.Password, user.Password) {
		return nil, nil, nil, nil, app.ErrInvalidCredentials
	}

	accessToken, _, err := s.JWTMaker.GenerateToken(user, constant.ACCESS_TOKEN_TTL)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	refreshToken, refreshClaims, err := s.JWTMaker.GenerateToken(user, constant.REFRESH_TOKEN_TTL)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	session := &model.Session{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: *refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}

	_, err = s.SessionRepo.Create(session)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return accessToken, refreshToken, refreshClaims, user, err
}

func (s *UserService) RenewUserAccessToken(oldRefreshToken string) (*string, *string, *token.UserClaims, error) {
	refreshClaims, err := s.JWTMaker.VerifyToken(oldRefreshToken)
	if err != nil {
		return nil, nil, nil, err
	}

	oldSession, err := s.SessionRepo.GetById(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	if oldSession.IsRevoked {
		return nil, nil, nil, app.ErrSessionRevoked
	}

	if oldSession.UserEmail != refreshClaims.Email {
		return nil, nil, nil, app.ErrInvalidSession
	}

	tx := s.Db.Begin()
	userTx := s.UserRepo.WithTx(tx)
	sessionTx := s.SessionRepo.WithTx(tx)

	user, err := userTx.GetByEmail(refreshClaims.Email)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	err = sessionTx.Revoke(oldSession.Id)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	accessToken, _, err := s.JWTMaker.GenerateToken(user, constant.ACCESS_TOKEN_TTL)
	if err != nil {
		return nil, nil, nil, err
	}

	refreshToken, refreshClaims, err := s.JWTMaker.GenerateToken(user, constant.REFRESH_TOKEN_TTL)
	if err != nil {
		return nil, nil, nil, err
	}

	newSession := &model.Session{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: *refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	}
	_, err = sessionTx.Create(newSession)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, nil, nil, err
	}

	return accessToken, refreshToken, refreshClaims, nil
}

func (s *UserService) Logout(refreshToken string) error {
	refreshClaims, err := s.JWTMaker.VerifyToken(refreshToken)
	if err != nil {
		return err
	}

	session, err := s.SessionRepo.GetById(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return err
	}

	if session.IsRevoked {
		return app.ErrSessionRevoked
	}

	if session.UserEmail != refreshClaims.Email {
		return app.ErrInvalidSession
	}
	err = s.SessionRepo.Revoke(session.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserById(id int64) (*model.User, error) {
	return s.UserRepo.GetById(id)
}
