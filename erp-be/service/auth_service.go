package service

import (
	"errors"

	"erp-kafeinart/config"
	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"
)

// LoginResponse adalah struktur kembalian dari login yang diharapkan oleh auth.module.js di frontend:
// FE: const { data } = await api.post('/auth/login', credentials)
// commit('SET_AUTH', { user: data.user, token: data.token })
type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

// AuthService mendefinisikan bisnis logika autentikasi.
type AuthService interface {
	Login(email, password string) (*LoginResponse, error)
	GetProfile(userID uint) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService menginisialisasi AuthService.
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(email, password string) (*LoginResponse, error) {
	if email == "" || password == "" {
		return nil, errors.New("email dan password wajib diisi")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("kredensial login tidak valid")
	}

	// Verifikasi hash password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("kredensial login tidak valid")
	}

	// Generate JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpirationHours)
	if err != nil {
		return nil, errors.New("gagal membuat sesi otentikasi")
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *authService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return user, nil
}
