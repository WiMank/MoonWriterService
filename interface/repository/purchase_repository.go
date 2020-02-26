package repository

import (
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/interface/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
	"net/http"
)

type purchaseRepository struct {
	collectionUsers    *mongo.Collection
	collectionSessions *mongo.Collection
	collectionPurchase *mongo.Collection
	responseCreator    response.AppResponseCreator
	validator          *validator.Validate
}

type PurchaseRepository interface {
	DecodePurchaseRegisterRequest(r *http.Request) request.PurchaseRegisterRequest
	DecodeVerificationRequest(r *http.Request) request.PurchaseVerificationRequest
	RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse
	VerificationPurchase(request request.PurchaseVerificationRequest) response.AppResponse
}

func NewPurchaseRepository(
	collectionUsers *mongo.Collection,
	collectionSessions *mongo.Collection,
	collectionPurchase *mongo.Collection,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate,
) PurchaseRepository {
	return &purchaseRepository{
		collectionUsers,
		collectionSessions,
		collectionPurchase,
		responseCreator,
		validator,
	}
}

func (pr *purchaseRepository) DecodePurchaseRegisterRequest(r *http.Request) request.PurchaseRegisterRequest {
	var purchaseRegisterRequest request.PurchaseRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&purchaseRegisterRequest); err != nil {
		log.Errorf("DecodePurchaseRegisterRequest error:\n", err)
	}
	return purchaseRegisterRequest
}

func (pr *purchaseRepository) DecodeVerificationRequest(r *http.Request) request.PurchaseVerificationRequest {
	var purchaseVerificationRequest request.PurchaseVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&purchaseVerificationRequest); err != nil {
		log.Errorf("DecodeVerificationRequest error:\n", err)
	}
	return purchaseVerificationRequest
}

func (pr *purchaseRepository) RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse {
	if request.ValidateRequest(pr.validator) {
		accessTokenExist := pr.checkAccessTokenExist(request.Purchase.AccessToken)
		userId, accessTokenValid := pr.validateAccessToken(request.Purchase.AccessToken)
		_, userExist := pr.checkFreeUserExist(userId)
		if accessTokenExist {
			if accessTokenValid {
				if userExist {
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
		accessTokenExist := pr.checkAccessTokenExist(request.Purchase.AccessToken)
		userId, accessTokenValid := pr.validateAccessToken(request.Purchase.AccessToken)
		localUser, userExist := pr.checkUserExist(userId)
		localPurchase, purchaseExist := pr.checkPurchaseExist(localUser)
		if accessTokenExist {
			if accessTokenValid {
				if userExist {
					if purchaseExist {
						if pr.checkPaymentData(localPurchase) {
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

func (pr *purchaseRepository) validateAccessToken(accessToken string) (string, bool) {
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

func (pr *purchaseRepository) checkFreeUserExist(userId string) (*domain.UserEntity, bool) {
	id, errHex := primitive.ObjectIDFromHex(userId)
	if errHex != nil {
		return nil, false
	}

	var localUser domain.UserEntity
	err := pr.collectionUsers.FindOne(utils.GetContext(), bson.D{{"_id", id}}).Decode(&localUser)
	if err != nil {
		return nil, false
	}

	if (localUser.Id == userId) && (!localUser.IsPremiumUser) {
		return &localUser, true
	}
	return nil, false
}

func (pr *purchaseRepository) checkUserExist(userId string) (*domain.UserEntity, bool) {
	id, errHex := primitive.ObjectIDFromHex(userId)
	if errHex != nil {
		return nil, false
	}

	var localUser domain.UserEntity
	err := pr.collectionUsers.FindOne(utils.GetContext(), bson.D{{"_id", id}}).Decode(&localUser)
	if err != nil {
		return nil, false
	}
	return &localUser, true
}

func (pr *purchaseRepository) checkPurchaseTokenNonExist(purchaseToken string) bool {
	count, err := pr.collectionPurchase.CountDocuments(
		utils.GetContext(),
		bson.D{{"purchase_token", purchaseToken}})
	if err != nil {
		return false
	}
	if count == 1 {
		return false
	}
	return true
}

func (pr *purchaseRepository) checkAccessTokenExist(accessToken string) bool {
	count, err := pr.collectionSessions.CountDocuments(
		utils.GetContext(),
		bson.D{{"access_token", accessToken}})
	if err != nil {
		return false
	}
	if count == 1 {
		return true
	}
	return false
}

func (pr *purchaseRepository) insertPurchase(userId string, request request.PurchaseRegisterRequest) bool {
	_, err := pr.collectionPurchase.InsertOne(utils.GetContext(), bson.D{
		{"user_id", userId},
		{"is_premium_user", true},
		{"purchase_token", request.Purchase.PurchaseToken},
		{"order_id", request.Purchase.OrderId},
		{"purchase_time", request.Purchase.PurchaseTime},
		{"sku", request.Purchase.Sku},
		{"access_token", request.Purchase.AccessToken},
	})

	if err != nil {
		return false
	}

	id, errHex := primitive.ObjectIDFromHex(userId)
	if errHex != nil {
		return false
	}
	_, errUpdate := pr.collectionUsers.UpdateOne(
		utils.GetContext(),
		bson.D{{"_id", id}},
		bson.D{{"$set", bson.D{{"is_premium_user", true}}}})

	if errUpdate != nil {
		return false
	}
	return true
}

func (pr *purchaseRepository) checkPurchaseExist(user *domain.UserEntity) (*domain.Purchase, bool) {
	if user != nil {
		var localPurchase domain.Purchase
		errFind := pr.collectionPurchase.FindOne(
			utils.GetContext(),
			bson.D{{"user_id", user.Id}}).Decode(&localPurchase)

		if errFind != nil {
			return nil, false
		}
		return &localPurchase, true
	}
	return nil, false
}

func (pr *purchaseRepository) checkPaymentData(purchase *domain.Purchase) bool {
	if purchase != nil {
		androidPublisherService, serviceErr := androidpublisher.NewService(
			context.Background(),
			option.WithCredentialsFile("E:\\Dev\\GoDev\\MoonWriterService\\credentials.json"),
		)

		if serviceErr != nil {
			return false
		}

		r, getErr := androidPublisherService.Purchases.Products.Get(
			config.APP_PACKAGE,
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
