package entity

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

type User struct {
	ID       string
	Email    string
	Password string
	Name     string
	Roles    []string
}

func NewUser(email, plainPassword, name string) (*User, error) {
	if email == "" || plainPassword == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &User{
		ID:       email,
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		Roles:    []string{RoleOperator},
	}, nil
}

func (u *User) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

func (u *User) SetRole(role string) error {
	validRoles := map[string]bool{
		RoleAdmin:    true,
		RoleOperator: true,
		RoleViewer:   true,
	}

	if !validRoles[role] {
		return fmt.Errorf("invalid role: %s", role)
	}

	u.Roles = []string{role}
	return nil
}

func (u *User) GetRole() string {
	if len(u.Roles) > 0 {
		return u.Roles[0]
	}
	return RoleOperator
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
