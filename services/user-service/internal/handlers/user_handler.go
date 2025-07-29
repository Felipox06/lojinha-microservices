package handlers

import(
	"encoding/json"
    "net/http"
    "log"
    
    "github.com/gorilla/mux"
    "github.com/Felipox06/lojinha-microservices/services/user-service/internal/models"
    "github.com/Felipox06/lojinha-microservices/services/user-service/internal/services"
)

type UserHandler struct{
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler{
	return &UserHandler{
		service: service,
	}
}

// CreateUser - POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request){
	// Primeira etapa: Decodificar JSON do body para struct
	var req models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		h.respondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}
	defer r.Body.Close()
    
	// Segunda etapa: Validar dados básicos
	if req.Name == "" || req.Email == "" || req.Password == "" {
		h.respondWithError(w, http.StatusBadRequest, "Campos requisitados faltando")
		return
	}

	// Terceira etap: Chamar service layer
	user, err := h.service.CreateUser(req)
	if err != nil{
		log.Printf("Erro criando user: %v", err)

		if err.Error() == "email já existe"{
			h.respondWithError(w, http.StatusConflict, "Email já em uso")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Falha em criar usuário")
		return
	}

	// Quarta etapa: Retornar resposta de sucesso
	h.respondWithJSON(w, http.StatusCreated, user.ToResponse())
}

// GetUser - GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r*http.Request){
	// Primeira etapa: Extrair ID da URL
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == ""{
		h.respondWithError(w, http.StatusBadRequest, "ID de user é necessário")
		return
	}

	// Segunda etapa: Buscar no service
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		log.Printf("Erro buscando user %s: %v", userID, err)
		h.respondWithError(w, http.StatusNotFound, "User não encontrado")
		return
	}

	// Terceira etapa: Retornar user(sem senha)
	h.respondWithJSON(w, http.StatusOK, user.ToResponse())
}

// ListUsers - Get /users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request){
	// Primeira etapa: Buscar users
	users, err := h.service.ListUsers()
	if err != nil {
		log.Printf("Erro lisitando users: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Falha ao buscar users")
	}

	// Segunda etapa: Converter para response(sem senhas)
	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users{
		responses = append(responses, user.ToResponse())
	}

	// Terceira etapa: Retornar lista
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"users": responses,
		"count": len(responses),
	})
}

// UpdateUser - PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request){
	// Primeira etapa: Extrair ID
	vais := mux.Vars(r)
	userID := vars["id"]

	// Segunda etapa: Decodificar body
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		h.respondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}
	defer r.Body.Close()

	// Terceira etapa: Atualizar via service
	user, err := h.service.UpdateUser(userID, req)
	if err != nil {
		log.Printf("Erro atualizando user %s: %v", userID, err)
		h.respondWithError(w, http.StatusNotFound, "User não encontrado")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user.ToResponse())
}

// DeleteUser - DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userID = vars["id"]

	if err := h.service.DeleteUser(userID); err != nil{
		log.Printf("Erro deletando user %s: %v", userID, err)
		h.respondWithError(w, http.StatusNotFound, "User não encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helpers para respostas padronizadas
func (h *UserHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        log.Printf("Erro codificando resposta: %v", err)
    }
}

func (h *UserHandler) respondWithError(w http.ResponseWriter, code int, message string) {
    h.respondWithJSON(w, code, map[string]string{
        "erro": message,
    })
}