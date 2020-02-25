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
	responseCreator response.AppResponseCreator) PurchaseRepository {
	return &purchaseRepository{collectionUsers, collectionSessions, collectionPurchase, responseCreator}
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
	accessTokenExist := pr.checkAccessTokenExist(request.AccessToken)
	userId, isValid := pr.validateAccessToken(request.AccessToken)
	_, isExist := pr.checkUserExist(userId)
	if accessTokenExist {
		if isValid {
			if isExist {
				if pr.checkTokenNonExist(request.PurchaseToken) {
					if pr.insertPurchase(userId, request) {
						return pr.responseCreator.CreateResponse(response.RegisterPurchaseResponse{}, userId)
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
	return pr.responseCreator.CreateResponse(response.RegisterPurchaseErrorResponse{}, userId)
}

func (pr *purchaseRepository) VerificationPurchase(request request.PurchaseVerificationRequest) response.AppResponse {
	accessTokenExist := pr.checkAccessTokenExist(request.AccessToken)
	userId, isValid := pr.validateAccessToken(request.AccessToken)
	localUser, userExist := pr.checkUserExist(userId)
	localPurchase, purchaseExist := pr.checkPurchaseExist(localUser)
	if accessTokenExist {
		if isValid {
			if userExist {
				if purchaseExist {
					if pr.checkPaymentData(localPurchase) {

					}
				}
			}
		}
	}

	return pr.responseCreator.CreateResponse(response.RegisterPurchaseResponse{}, "userId")
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

func (pr *purchaseRepository) checkUserExist(userId string) (*domain.UserEntity, bool) {
	id, errHex := primitive.ObjectIDFromHex(userId)
	if errHex != nil {
		return nil, false
	}

	var localUser domain.UserEntity
	err := pr.collectionUsers.FindOne(utils.GetContext(), bson.D{{"_id", id}}).Decode(&localUser)
	if err != nil {
		log.Error("Count: ", err)
		return nil, false
	}

	if (localUser.Id == userId) && (!localUser.IsPremiumUser) {
		return nil, true
	}

	return &localUser, false
}

func (pr *purchaseRepository) checkTokenNonExist(purchaseToken string) bool {
	count, err := pr.collectionPurchase.CountDocuments(utils.GetContext(), bson.D{{"purchase_token", purchaseToken}})
	if err != nil {
		return false
	}

	if count == 1 {
		return false
	}

	return true
}

func (pr *purchaseRepository) checkAccessTokenExist(accessToken string) bool {
	count, err := pr.collectionSessions.CountDocuments(utils.GetContext(), bson.D{{"access_token", accessToken}})
	if err != nil {
		log.Error("CheckAccessTokenExist error: ", err)
	}

	if count == 1 {
		return true
	}

	pr.responseCreator.CreateResponse(response.InvalidToken{}, "")
	return false
}

func (pr *purchaseRepository) insertPurchase(userId string, request request.PurchaseRegisterRequest) bool {
	_, err := pr.collectionPurchase.InsertOne(utils.GetContext(), bson.D{
		{"user_id", userId},
		{"is_premium_user", true},
		{"purchase_token", request.PurchaseToken},
		{"order_id", request.OrderId},
		{"purchase_time", request.PurchaseTime},
		{"sku", request.Sku},
		{"access_token", request.AccessToken},
	})

	if err != nil {
		log.Errorf("InsertOne: ", err)
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
	id, errHex := primitive.ObjectIDFromHex(user.Id)
	if errHex != nil {
		return nil, false
	}

	var localPurchase domain.Purchase
	errFind := pr.collectionUsers.FindOne(utils.GetContext(), bson.D{{"_id", id}}).Decode(&localPurchase)
	if errFind != nil {
		return nil, false
	}

	return &localPurchase, true
}

func (pr *purchaseRepository) checkPaymentData(purchase *domain.Purchase) bool {
	androidPublisherService, serviceErr := androidpublisher.NewService(context.Background(), option.WithCredentialsFile("E:\\Dev\\GoDev\\MoonWriterService\\credentials.json"))
	if serviceErr != nil {
		log.Infof("CheckPaymentData NewService Error: ", serviceErr)
		return false
	}

	r, getErr := androidPublisherService.Purchases.Products.Get(
		config.APP_PACKAGE,
		purchase.Sku,
		purchase.PurchaseToken).Do()
	if getErr != nil {
		log.Infof("CheckPaymentData Products Get: ", getErr)
		return false
	}

	if r.HTTPStatusCode == http.StatusOK {
		return true
	}

	return false
}
