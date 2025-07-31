// A interface UserService é um contrato que descreve *o que um serviço de usuários deve saber fazer*, sem se preocupar com *como* ele faz.
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Felipox06/lojinha-microservices/services/user-service/internal/models"
)

// Interface = conjunto de métodos que um tipo deve implementar
type UserService interface {
	CreateUser(req models.CreateUserRequest) (*models.User, error) //Cria um usuário a partir de uma requisição e retorna o usuário criado ou erro
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id string, req models.CreateUserRequest) (*models.User, error)
	DeleteUser(id string) error
	ListUsers() ([]*models.User, error)
}

// lowercase = private
type userService struct {
	// Por equanto, armazena em memória
	// TODO: substituir por banco de dados
	users map[string]*models.User // map[chave]valor
}

// Construtor(factoy pattern): cria instâncias sem precisar expor os detalhes de como elas são construídas
// Retorna interface (dependency inversion)
func NewUserService() UserService {
	return &userService{ //retorna um ponteiro para a instância
		users: make(map[string]*models.User), //make inicializa o map
	}
}

// CreateUser implementa o método da interface
func (s *userService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	for _, user := range s.users {
		if user.Email == req.Email {
			return nil, errors.New("email já existe")
		}
	}

	user := &models.User{
		ID:        generateID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashPassword(req.Password),
		Type:      req.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.users[user.ID] = user

	return user, nil
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user com o id %s não foi encontrado", id)
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user não encontrado")
}

func (s *userService) UpdateUser(id string, req models.CreateUserRequest) (*models.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user com id %s não foi encontrado", id)
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Password != "" {
		user.Password = hashPassword(req.Password)
	}
	user.Type = req.Type
	user.UpdatedAt = time.Now()

	return user, nil
}

func (s *userService) DeleteUser(id string) error {
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user com id %s não foi encontrado", id)
	}

	delete(s.users, id)
	return nil
}

func (s *userService) ListUsers() ([]*models.User, error) {
	users := make([]*models.User, 0, len(s.users))

	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

func generateID() string {
	// TODO: usar UUID real
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func hashPassword(password string) string {
	// TODO: usar bcrypt
	return fmt.Sprintf("hashed_%s", password)
}
