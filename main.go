package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/redhajuanda/gorengan/config"
	"github.com/redhajuanda/gorengan/internal/auth"
	"github.com/redhajuanda/gorengan/internal/httperror"
	"github.com/redhajuanda/gorengan/internal/user"
	"github.com/redhajuanda/gorengan/pkg/log"

	_ "github.com/go-sql-driver/mysql"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

func main() {
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// Load config
	cfg := config.LoadDefault()

	// Connect DB
	connString := fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local&", cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		logger.Errorf("%v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
	}

	address := fmt.Sprintf(":%v", cfg.Server.PORT)
	server := http.Server{
		Addr:    address,
		Handler: buildHandlers(db, cfg, logger),
	}
	logger.Infof("server %v is running at %v", Version, address)

	if err := server.ListenAndServe(); err != nil {
		logger.Errorf("")
	}
}

func buildHandlers(db *sql.DB, cfg config.Config, logger log.Logger) http.Handler {
	r := echo.New()
	r.Pre(middleware.RemoveTrailingSlash())

	// Set custom HTTP error handler
	r.HTTPErrorHandler = httperror.CustomHTTPErrorHandler

	// Register user service
	user.RegisterService(
		*r.Group(""),
		user.NewService(user.NewRepository(db), logger),
		cfg,
		logger,
	)

	// Register auth service
	auth.RegisterService(
		*r.Group(""),
		auth.NewService(cfg.JWT.SigningKey, cfg.JWT.TokenExpiration, logger, auth.NewRepository(db)),
		logger,
	)

	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API Version: "+Version)
	})
	return r
}
