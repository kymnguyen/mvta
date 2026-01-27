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
	Username string
	Password string
	Roles    []string
}

func NewUser(username, plainPassword string) (*User, error) {
	if username == "" || plainPassword == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &User{
		ID:       username,
		Username: username,
		Password: string(hashedPassword),
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
