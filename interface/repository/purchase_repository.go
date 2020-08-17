package repository

import (
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"net/http"
	"strconv"
)

type purchaseRepository struct {
	db              *sqlx.DB
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type PurchaseRepository interface {
	RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse
	VerificationPurchase(request request.PurchaseVerificationRequest) response.AppResponse
}

func NewPurchaseRepository(
	db *sqlx.DB,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate,
) PurchaseRepository {
	return &purchaseRepository{db, responseCreator, validator}
}

func (pr *purchaseRepository) RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse {
	if request.ValidateRequest(pr.validator) {
		accessTokenExist := pr.checkAccessTokenExistInSession(request.Purchase.AccessToken)
		userId, accessTokenValid := validateAccessToken(request.Purchase.AccessToken)
		localUser, userExist := pr.checkUserExist(userId)

		if accessTokenExist {
			if accessTokenValid {
				if (userExist) && (!localUser.IsPremiumUser) {
					if pr.checkPurchaseTokenNonExist(request.Purchase.PurchaseToken) {
						if pr.insertPurchase(userId, request) {
							return pr.responseCreator.CreateResponse(response.RegisterPurchaseResponse{}, userId)
						} else {
							return pr.responseCreator.CreateResponse(response.InsertPurchaseErrorResponse{}, userId)
						}
					} else {
						return pr.responseCreator.CreateResponse(response.PurchaseTokenExistResponse{}, userId)
					}
				} else {
					return pr.responseCreator.CreateResponse(response.PurchaseUserExistResponse{}, userId)
				}
			} else {
				return pr.responseCreator.CreateResponse(response.InvalidToken{}, userId)
			}
		} else {
			return pr.responseCreator.CreateResponse(response.InvalidToken{}, userId)
		}
	}

	return pr.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func (pr *purchaseRepository) VerificationPurchase(request request.PurchaseVerificationRequest) response.AppResponse {
	if request.ValidateRequest(pr.validator) {
		accessTokenExist := pr.checkAccessTokenExistInSession(request.Purchase.AccessToken)
		userId, accessTokenValid := validateAccessToken(request.Purchase.AccessToken)
		localUser, userExist := pr.checkUserExist(userId)
		localPurchase, purchaseExist := pr.checkPurchaseExist(localUser)
		paymentExist := checkPaymentData(localPurchase)

		if accessTokenExist {
			if accessTokenValid {
				if userExist {
					if purchaseExist {
						if paymentExist {
							return pr.responseCreator.CreateResponse(response.PurchaseValidResponse{}, userId)
						} else {
							return pr.responseCreator.CreateResponse(response.CheckPaymentDataErrorResponse{}, userId)
						}
					} else {
						return pr.responseCreator.CreateResponse(response.PurchaseNotFoundResponse{}, userId)
					}
				} else {
					return pr.responseCreator.CreateResponse(response.PurchaseUserExistResponse{}, userId)
				}
			} else {
				return pr.responseCreator.CreateResponse(response.InvalidToken{}, userId)
			}
		} else {
			return pr.responseCreator.CreateResponse(response.InvalidToken{}, userId)
		}
	}

	return pr.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func validateAccessToken(accessToken string) (string, bool) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return cast.ToString(claims["user"]), true
	}

	return "", false
}

func (pr *purchaseRepository) checkUserExist(userName string) (*domain.UserEntity, bool) {
	var localUser domain.UserEntity
	findUserQuery := "SELECT * FROM users WHERE user_name =$1"
	err := pr.db.QueryRowx(findUserQuery, userName).StructScan(&localUser)

	if err != nil {
		return nil, false
	}

	return &localUser, true
}

func (pr *purchaseRepository) checkPurchaseTokenNonExist(purchaseToken string) bool {
	var exist bool
	checkPurchaseToken := "SELECT EXISTS (SELECT purchase_token FROM purchases WHERE purchase_token=$1)::bool"
	err := pr.db.QueryRowx(checkPurchaseToken, purchaseToken).Scan(&exist)

	if err != nil {
		return false
	}

	return !exist
}

func (pr *purchaseRepository) checkAccessTokenExistInSession(accessToken string) bool {
	var exist bool
	checkAccessToken := "SELECT EXISTS (SELECT access_token FROM sessions WHERE access_token=$1)::bool"
	err := pr.db.QueryRowx(checkAccessToken, accessToken).Scan(&exist)

	if err != nil {
		log.Info("checkAccessTokenExistInSession: ", err)
		return false
	}

	return exist
}

func (pr *purchaseRepository) insertPurchase(userName string, request request.PurchaseRegisterRequest) bool {
	insertPurchase := "INSERT INTO purchases  (user_name, purchase_token, access_token, order_id, purchase_time, sku) " +
		"VALUES ($1, $2, $3, $4, $5, $6)"

	updateUserStatus := "UPDATE users SET premium = true WHERE user_name=$1"

	tx, errTransaction := pr.db.Begin()
	if errTransaction != nil {
		return false
	}

	_, errInsert := tx.Exec(insertPurchase,
		userName,
		request.Purchase.PurchaseToken,
		request.Purchase.AccessToken,
		request.Purchase.OrderId,
		strconv.FormatInt(request.Purchase.PurchaseTime, 10),
		request.Purchase.Sku,
	)
	_, errUpdate := tx.Exec(updateUserStatus, userName)

	if (errInsert != nil) || (errUpdate != nil) {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return false
		}
		return false
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return false
	}

	return true
}

func (pr *purchaseRepository) checkPurchaseExist(user *domain.UserEntity) (*domain.Purchase, bool) {
	var localPurchase domain.Purchase
	checkPurchaseExist := "SELECT user_name FROM purchases WHERE user_name=$1"
	err := pr.db.QueryRowx(checkPurchaseExist, user.UserName).StructScan(&localPurchase)

	if err != nil {
		return nil, false
	}

	return &localPurchase, false
}

func checkPaymentData(purchase *domain.Purchase) bool {
	if purchase != nil {
		androidPublisherService, serviceErr := androidpublisher.NewService(
			context.Background(),
			option.WithCredentialsFile("credentials.json"),
		)

		if serviceErr != nil {
			return false
		}

		r, getErr := androidPublisherService.Purchases.Products.Get(
			config.AppPackage,
			purchase.Sku,
			purchase.PurchaseToken).Do()

		if getErr != nil {
			return false
		}

		if r.HTTPStatusCode == http.StatusOK {
			return true
		}

		return false
	}

	return false
}
