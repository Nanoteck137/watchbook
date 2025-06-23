package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
)

const (
	ErrTypeInvalidAuth        pyrin.ErrorType = "INVALID_AUTH"
	ErrTypeUserAlreadyExists  pyrin.ErrorType = "USER_ALREADY_EXISTS"
	ErrTypeUserNotFound       pyrin.ErrorType = "USER_NOT_FOUND"
	ErrTypeApiTokenNotFound   pyrin.ErrorType = "API_TOKEN_NOT_FOUND"
	ErrTypeInvalidCredentials pyrin.ErrorType = "INVALID_CREDENTIALS"

	ErrTypeAnimeNotFound pyrin.ErrorType = "ANIME_NOT_FOUND"
	ErrTypeEpisodeNotFound pyrin.ErrorType = "EPISODE_NOT_FOUND"
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

func EpisodeNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeEpisodeNotFound,
		Message: "Episode not found",
	}
}

func UserAlreadyExists() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeUserAlreadyExists,
		Message: "User already exists",
	}
}

func ApiTokenNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeApiTokenNotFound,
		Message: "Api Token not found",
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
