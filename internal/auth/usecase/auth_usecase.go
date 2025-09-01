package usecase

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zercle/gofiber-skeleton/internal/auth"
	"github.com/zercle/gofiber-skeleton/internal/config"
	"github.com/zercle/gofiber-skeleton/internal/core/domain"
	"github.com/zercle/gofiber-skeleton/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userUsecase user.UserUsecase
	cfg         *config.Config
}

func NewAuthUsecase(userUsecase user.UserUsecase, cfg *config.Config) auth.AuthUsecase {
	return &authUsecase{
		userUsecase: userUsecase,
		cfg:         cfg,
	}
}

func (u *authUsecase) Login(ctx *fiber.Ctx, email, password string) (*auth.LoginResponse, *domain.Error) {
	foundUser, err := u.userUsecase.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "invalid email or password")
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gofiber-skeleton",
			Subject:   foundUser.ID.String(),
			Audience:  []string{"gofiber-skeleton"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.cfg.JWT.ExpiresIn)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
		UserID: foundUser.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(u.cfg.JWT.Secret))
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "failed to sign token")
	}

	return &auth.LoginResponse{
		AccessToken: ss,
	}, nil
}
