package apis

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

// TODO(patrik): Remove
var logger = watchbook.DefaultLogger()

type UserCheckFunc func(user *database.User) error

func RequireAdmin(user *database.User) error {
	if user.Role != types.RoleSuperUser && user.Role != types.RoleAdmin {
		return InvalidAuth("user requires 'super_user' or 'admin' role")
	}

	return nil
}

func HasEditPrivilege(user *database.User) error {
	return RequireAdmin(user)
}

func User(app core.App, c pyrin.Context, checks ...UserCheckFunc) (*database.User, error) {
	user, err := getUser(app, c)
	if err != nil {
		return nil, err
	}

	for _, check := range checks {
		err := check(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func getUser(app core.App, c pyrin.Context) (*database.User, error) {
	apiTokenHeader := c.Request().Header.Get("X-Api-Token")
	if apiTokenHeader != "" {
		ctx := context.TODO()
		token, err := app.DB().GetApiTokenById(ctx, apiTokenHeader)
		if err != nil {
			if errors.Is(err, database.ErrItemNotFound) {
				return nil, InvalidAuth("invalid api token")
			}

			return nil, err
		}

		user, err := app.DB().GetUserById(c.Request().Context(), token.UserId)
		if err != nil {
			return nil, InvalidAuth("invalid api token")
		}

		return &user, nil
	}

	authHeader := c.Request().Header.Get("Authorization")
	tokenString := utils.ParseAuthHeader(authHeader)
	if tokenString == "" {
		return nil, InvalidAuth("invalid authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.Config().JwtSecret), nil
	})

	if err != nil {
		// TODO(patrik): Handle error better
		return nil, InvalidAuth("invalid authorization token")
	}

	jwtValidator := jwt.NewValidator(jwt.WithIssuedAt())

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if err := jwtValidator.Validate(token.Claims); err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		userId := claims["userId"].(string)
		user, err := app.DB().GetUserById(c.Request().Context(), userId)
		if err != nil {
			return nil, InvalidAuth("invalid authorization token")
		}

		return &user, nil
	}

	return nil, InvalidAuth("invalid authorization token")
}

func ConvertURL(c pyrin.Context, path string) string {
	host := c.Request().Host

	scheme := "http"

	h := c.Request().Header.Get("X-Forwarded-Proto")
	if h != "" {
		scheme = h
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}
