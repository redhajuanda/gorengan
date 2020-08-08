package user

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/redhajuanda/gorengan/internal/auth"
	"github.com/redhajuanda/gorengan/internal/httperror"
	"github.com/redhajuanda/gorengan/internal/httpsuccess"
	"github.com/redhajuanda/gorengan/pkg/log"
	"github.com/redhajuanda/gorengan/pkg/pagination"
)

// RegisterService registers a new user service
func RegisterService(r echo.Group, service Service, logger log.Logger) {
	handler := handler{service, logger}

	r.Use(auth.IsLoggedIn)
	r.Use(auth.IsAdmin)

	// the following endpoints require a valid JWT
	r.GET("/users/:id", handler.get)
	r.GET("/users", handler.query)
	r.POST("/users", handler.create)
	r.PUT("/users/:id", handler.update)
	r.DELETE("/users/:id", handler.delete)
}

type handler struct {
	service Service
	logger  log.Logger
}

func (h handler) get(c echo.Context) error {
	user, err := h.service.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return httpsuccess.ResponseWithJSON(c, "", http.StatusOK, user)
}

func (h handler) query(c echo.Context) error {
	ctx := c.Request().Context()
	count, err := h.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request(), count)
	users, err := h.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = users
	return httpsuccess.ResponseWithJSON(c, "", http.StatusOK, pages)
}

func (h handler) create(c echo.Context) error {
	var input CreateUserRequest
	if err := c.Bind(&input); err != nil {
		h.logger.With(c.Request().Context()).Info(err)
		return httperror.BadRequest("")
	}
	user, err := h.service.Create(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return httpsuccess.ResponseWithJSON(c, "user created", http.StatusCreated, user)
}

func (h handler) update(c echo.Context) error {
	var input UpdateUserRequest
	if err := c.Bind(&input); err != nil {
		h.logger.With(c.Request().Context()).Info(err)
		return httperror.BadRequest("")
	}

	user, err := h.service.Update(c.Request().Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h handler) delete(c echo.Context) error {
	user, err := h.service.Delete(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return httpsuccess.ResponseWithJSON(c, "user deleted", http.StatusOK, user)
}
