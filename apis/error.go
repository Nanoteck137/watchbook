package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
)

const (
	ErrTypeInvalidAuth        pyrin.ErrorType = "INVALID_AUTH"
	ErrTypeUserAlreadyExists  pyrin.ErrorType = "USER_ALREADY_EXISTS"
	ErrTypeUserNotFound       pyrin.ErrorType = "USER_NOT_FOUND"
	ErrTypeInvalidCredentials pyrin.ErrorType = "INVALID_CREDENTIALS"

	ErrTypeAnimeNotFound pyrin.ErrorType = "ANIME_NOT_FOUND"
)

func InvalidAuth(message string) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidAuth,
		Message: "Invalid auth: " + message,
	}
}

func AnimeNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeAnimeNotFound,
		Message: "Anime not found",
	}
}

func UserAlreadyExists() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeUserAlreadyExists,
		Message: "User already exists",
	}
}

func UserNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusUnauthorized,
		Type:    ErrTypeUserNotFound,
		Message: "User not found",
	}
}

func InvalidCredentials() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusUnauthorized,
		Type:    ErrTypeInvalidCredentials,
		Message: "Invalid Credentials",
	}
}
