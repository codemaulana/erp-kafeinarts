package service

import (
	"errors"
	"strings"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"
)

type CreateUserDTO struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UpdateUserDTO struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type UserService interface {
	GetUsers(search, role string) ([]model.User, error)
	GetUserByID(id uint) (*model.User, error)
	CreateUser(dto *CreateUserDTO) (*model.User, error)
	UpdateUser(id uint, dto *UpdateUserDTO) (*model.User, error)
	DeleteUser(id uint, currentUserID uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers(search, role string) ([]model.User, error) {
	return s.repo.FindAll(strings.TrimSpace(search), strings.TrimSpace(role))
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return user, nil
}

func (s *userService) CreateUser(dto *CreateUserDTO) (*model.User, error) {
	name := strings.TrimSpace(dto.Name)
	email := strings.ToLower(strings.TrimSpace(dto.Email))
	password := strings.TrimSpace(dto.Password)
	role := strings.ToLower(strings.TrimSpace(dto.Role))

	if name == "" {
		return nil, errors.New("nama pengguna wajib diisi")
	}
	if email == "" {
		return nil, errors.New("email wajib diisi")
	}
	if password == "" {
		return nil, errors.New("password wajib diisi")
	}
	if len(password) < 6 {
		return nil, errors.New("password minimal harus 6 karakter")
	}
	if role == "" {
		role = "staff"
	}

	// Cek keunikan email
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar dalam sistem")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	perms := dto.Permissions
	if len(perms) == 0 {
		if role == "superadmin" || role == "admin" {
			perms = []string{"*"}
		} else {
			perms = []string{role}
		}
	}

	user := &model.User{
		Name:        name,
		Email:       email,
		Password:    hashedPassword,
		Role:        role,
		Permissions: model.PermissionList(perms),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id uint, dto *UpdateUserDTO) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	name := strings.TrimSpace(dto.Name)
	email := strings.ToLower(strings.TrimSpace(dto.Email))
	role := strings.ToLower(strings.TrimSpace(dto.Role))

	if name != "" {
		user.Name = name
	}
	if email != "" && email != user.Email {
		existing, err := s.repo.FindByEmail(email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != user.ID {
			return nil, errors.New("email sudah digunakan oleh pengguna lain")
		}
		user.Email = email
	}
	if role != "" {
		user.Role = role
	}

	password := strings.TrimSpace(dto.Password)
	if password != "" {
		if len(password) < 6 {
			return nil, errors.New("password baru minimal harus 6 karakter")
		}
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password baru")
		}
		user.Password = hashedPassword
	}

	if len(dto.Permissions) > 0 {
		user.Permissions = model.PermissionList(dto.Permissions)
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint, currentUserID uint) error {
	if id == currentUserID {
		return errors.New("anda tidak dapat menghapus akun yang sedang aktif digunakan")
	}

	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("pengguna tidak ditemukan")
	}

	return s.repo.Delete(id)
}
