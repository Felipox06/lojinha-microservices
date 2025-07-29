package models

import "time"

type User struct{
	ID string `json:"id"`
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password, omnitempty"` //omnitempty = oculta
	Type UserType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserType string //tipo customizado para os valores

// Tipos de usuário, para funcionamento como enum
const (
	userTypeSeller UserType = "seller"
	UserTypeBuyer UserType = "buyer"
)


// Esse é o body que o sistema recebe do cliente
type CreateUserRequest struct {
	Name string `json:"name" binding:"required, min=3"`
	Email string `json:"email" binding:"required, email"`
	Password string `json:"password" binding: "required, min=6"`
	Type UserType `json:"type" binding:"required, oneof=seller buyer"`
}

// UserResponse é o que o sistema retorna ao cliente
type UserResponse struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Type UserType `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		Type: u.Type,
		CreatedAt: u.CreatedAt,
	}
}