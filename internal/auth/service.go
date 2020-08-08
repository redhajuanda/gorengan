package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redhajuanda/gorengan/internal/httperror"
	"github.com/redhajuanda/gorengan/pkg/log"
	"github.com/redhajuanda/gorengan/pkg/password"
	"github.com/redhajuanda/gorengan/pkg/validation"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, req LoginRequest) (string, error)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetUsername() string
	// GetRole returns user role
	GetRole() string
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	repo            Repository
	validation      *validation.CustomValidator
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, logger log.Logger, repo Repository) Service {
	return service{signingKey, tokenExpiration, logger, repo, validation.New()}
}

// LoginRequest holds request data for login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, req LoginRequest) (string, error) {
	err := s.validation.Validate(req)
	if err != nil {
		return "", err
	}
	if identity := s.authenticate(ctx, req.Email, req.Password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", httperror.Unauthorized("Cannot authenticate, invalid email or password")
}

// authenticate authenticates a user using email and password.
// If email and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, email, plainPwd string) Identity {
	logger := s.logger.With(ctx, "user", email)

	user, err := s.repo.Login(ctx, email)
	if err != nil {
		return nil
	}

	if email == user.Email && password.ComparePasswords(user.Password, []byte(plainPwd)) {
		logger.Infof("authentication successful")
		return user
	}

	logger.Infof("authentication failed")
	return nil
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       identity.GetID(),
		"username": identity.GetUsername(),
		"role":     identity.GetRole(),
		"exp":      time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
