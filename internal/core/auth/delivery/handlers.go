package delivery

import (
	"fmt"
	"net/http"

	authUseCase "study_buddy/internal/core/auth/service"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"
	"study_buddy/pkg/utils"

	"github.com/gin-gonic/gin"
)

var _ Handlers = (*Handler)(nil)

type Handler struct {
	service authUseCase.Service
}

func NewAuthHandler(service authUseCase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type Handlers interface {
	LogIn(g *gin.Context)
	SignUp(g *gin.Context)
}

type LogInRequest struct {
	Username string `json:"username" form:"username"` //binding:"required,email,max=100"`
	Password string `json:"password" form:"password"` //binding:"required,min=6,max=100"`
}

func (h *Handler) LogIn(g *gin.Context) {
	g.Set(constants.HandlerName, "LogIn")
	var request LogInRequest
	if err := g.BindJSON(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	token, err := h.service.LogIn(g.Request.Context(), request.Username, request.Password)

	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.LogIn: %w", err))
		return
	}

	utils.SetAuthorizationToken(g, token)

	g.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type SignUpRequest struct {
	Username string `json:"username" form:"username"` //binding:"required,email,max=100"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"` //binding:"required,min=6,max=100"`
}

func (h *Handler) SignUp(g *gin.Context) {
	g.Set(constants.HandlerName, "SignUp")
	var request SignUpRequest

	if err := g.BindJSON(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	token, err := h.service.SignUp(g.Request.Context(), request.Username, request.Email, request.Password)

	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.SignUp: %w", err))
		return
	}

	utils.SetAuthorizationToken(g, token)

	g.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
