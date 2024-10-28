package domain

import "time"

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Password  string     `json:"-"`
	IsActive  bool       `json:"isActive"`
	IsDeleted bool       `json:"isDeleted,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	UpdateActive(id string, isActive bool) error
	SoftDelete(id string) error
	HardDelete(id string) error
}

type UserUsecase interface {
	Register(user *User) error
	Login(email, password string) (string, error)
	VerifyOTP(email, otp string) error
	ResendOTP(email string) error
	GetUserByID(id string) (*User, error)
	DeleteUser(id string, permanent bool) error
}
