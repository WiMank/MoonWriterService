package response

import (
	"fmt"
	"log"
	"net/http"
)

type AppResponse interface {
	PrintLog()
	GetStatusCode() int
}

type AppResponseCreator interface {
	CreateResponse(i interface{}, userName string) AppResponse
}

type concreteAppResponseCreator struct {
}

func NewAppResponseCreator() AppResponseCreator {
	return &concreteAppResponseCreator{}
}

func (c *concreteAppResponseCreator) CreateResponse(i interface{}, userName string) AppResponse {
	var appResponse AppResponse
	switch t := i.(type) {
	case UnauthorizedResponse:
		appResponse = &UnauthorizedResponse{
			Message: fmt.Sprintf("User [%s] unauthorized", userName),
			Code:    http.StatusUnauthorized,
			Desc:    http.StatusText(http.StatusUnauthorized),
		}
	case UserCreatedResponse:
		appResponse = &UserCreatedResponse{
			Message: fmt.Sprintf("User [%s] registration success!", userName),
			Code:    http.StatusCreated,
			Desc:    http.StatusText(http.StatusCreated),
		}
	case UserExistResponse:
		appResponse = &UserExistResponse{
			Message: fmt.Sprintf("User with the name [%s] is already registered", userName),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	case UserInsertErrorResponse:
		appResponse = &UserInsertErrorResponse{
			Message: "Internal server error during user registration!",
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case UserFindResponse:
		appResponse = &UserFindResponse{
			Message: fmt.Sprintf("User [%s] not found", userName),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	case TokenResponse:
		appResponse = &TokenResponse{
			Message:      fmt.Sprintf("Tokens created for [%s]", userName),
			Code:         http.StatusOK,
			Desc:         http.StatusText(http.StatusOK),
			RefreshToken: t.RefreshToken,
			ExpiresInR:   t.ExpiresInR,
			AccessToken:  t.AccessToken,
			ExpiresInA:   t.ExpiresInA,
		}
	case TokenErrorResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Token creation error for user [%s]", userName),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case SessionUpdateFailedResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Error updating session for user [%s]", userName),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case SessionInsertFailedResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Session insert error for user [%s]", userName),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	default:
		log.Fatal("Unknown Response")
	}
	appResponse.PrintLog()
	return appResponse
}
