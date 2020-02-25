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

func (c *concreteAppResponseCreator) CreateResponse(i interface{}, data string) AppResponse {
	var appResponse AppResponse
	switch t := i.(type) {
	case UnauthorizedResponse:
		appResponse = &UnauthorizedResponse{
			Message: fmt.Sprintf("User [%s] unauthorized", data),
			Code:    http.StatusUnauthorized,
			Desc:    http.StatusText(http.StatusUnauthorized),
		}
	case UserCreatedResponse:
		appResponse = &UserCreatedResponse{
			Message: fmt.Sprintf("User [%s] registration success!", data),
			Code:    http.StatusCreated,
			Desc:    http.StatusText(http.StatusCreated),
		}
	case UserExistResponse:
		appResponse = &UserExistResponse{
			Message: fmt.Sprintf("User with the name [%s] is already registered", data),
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
			Message: fmt.Sprintf("User [%s] not found", data),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	case TokenResponse:
		appResponse = &TokenResponse{
			Message:      t.Message,
			Code:         http.StatusOK,
			Desc:         http.StatusText(http.StatusOK),
			SessionId:    t.SessionId,
			AccessToken:  t.AccessToken,
			RefreshToken: t.RefreshToken,
		}
	case TokenErrorResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Token creation error for user [%s]", data),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case SessionUpdateFailedResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Error updating session for user [%s]", data),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case SessionInsertFailedResponse:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("Session insert error for user [%s]", data),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case InvalidSession:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("INVALID SESSION [%s]", data),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	case InvalidToken:
		appResponse = &TokenErrorResponse{
			Message: fmt.Sprintf("INVALID TOKEN [%s]", data),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	case RefreshSessionErrorResponse:
		appResponse = &TokenErrorResponse{
			Message: "Session Update Error",
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}
	case RegisterPurchaseResponse:
		appResponse = &RegisterPurchaseResponse{
			Message: "Upgrade to the pro version was successful!",
			Code:    http.StatusCreated,
			Desc:    http.StatusText(http.StatusCreated),
		}
	case RegisterPurchaseErrorResponse:
		appResponse = &RegisterPurchaseErrorResponse{
			Message: fmt.Sprintf("Error updating to pro version [%s]", data),
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}

	case PurchaseTokenExistResponse:
		appResponse = &PurchaseTokenExistResponse{
			Message: "A purchase has already been made in the past",
			Code:    http.StatusInternalServerError,
			Desc:    http.StatusText(http.StatusInternalServerError),
		}

	case PurchaseUserExistResponse:
		appResponse = &UserFindResponse{
			Message: fmt.Sprintf("User [%s] not found or already has premium status", data),
			Code:    http.StatusBadRequest,
			Desc:    http.StatusText(http.StatusBadRequest),
		}
	default:
		log.Fatal("Unknown Response")
	}
	appResponse.PrintLog()
	return appResponse
}
