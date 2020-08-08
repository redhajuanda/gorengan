package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/redhajuanda/gorengan/internal/httpsuccess"
	"github.com/redhajuanda/gorengan/pkg/log"
)

// RegisterService registers a new user service
func RegisterService(r echo.Group, service Service, logger log.Logger) {
	handler := handler{service, logger}

	r.POST("/login", handler.login)
}

type handler struct {
	service Service
	logger  log.Logger
}

func (h handler) login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	accessToken, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return httpsuccess.ResponseWithJSON(c, "access granted", http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
