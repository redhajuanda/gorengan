package httpsuccess

import "github.com/labstack/echo"

func ResponseWithJSON(c echo.Context, msg string, code int, data interface{}) error {
	if msg == "" {
		msg = "success"
	}
	res := Response{
		Status:  code,
		Message: msg,
		Data:    data,
	}
	return c.JSON(code, res)
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
