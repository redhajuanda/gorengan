package httperror

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/redhajuanda/gorengan/pkg/validation"
)

// CustomHTTPErrorHandler sets error response for different type of errors and logs
func CustomHTTPErrorHandler(err error, c echo.Context) {
	log.Print(err)
	code := http.StatusInternalServerError

	if _, ok := err.(validation.ValidationErrors); ok {
		err = BadRequest(err.Error())
	}

	if errors.Is(err, sql.ErrNoRows) {
		err = NotFound("")
	}

	if resp, ok := err.(ErrorResponse); ok {
		code = resp.StatusCode()
	}

	if resp, ok := err.(*echo.HTTPError); ok {
		code = resp.Code
		err = ErrorResponse{
			Status:  code,
			Message: resp.Message.(string),
		}
	}
	c.JSON(code, err)
}
