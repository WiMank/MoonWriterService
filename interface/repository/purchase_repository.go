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
	DecodeRequest(r *http.Request) request.PurchaseRegisterRequest
	RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse
	VerificationPurchase(request request.PurchaseRegisterRequest) response.AppResponse
}

func NewPurchaseRepository(
	collectionUsers *mongo.Collection,
	collectionSessions *mongo.Collection,
	collectionPurchase *mongo.Collection,
	responseCreator response.AppResponseCreator) PurchaseRepository {
	return &purchaseRepository{collectionUsers, collectionSessions, collectionPurchase, responseCreator}
}

func (pr *purchaseRepository) DecodeRequest(r *http.Request) request.PurchaseRegisterRequest {
	var refreshTokensRequest request.PurchaseRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokensRequest); err != nil {
		log.Errorf("Decode VerificationPurchase error:\n", err)
	}
	return refreshTokensRequest
}

func (pr *purchaseRepository) RegisterPurchase(request request.PurchaseRegisterRequest) response.AppResponse {
	accessTokenExist := pr.checkAccessTokenExist(request.AccessToken)
	userId, isValid := pr.validateAccessToken(request.AccessToken)
	if accessTokenExist {
		if isValid {
			if pr.CheckUserExist(userId) {
				if pr.insertPurchase(userId, request) {
					return pr.responseCreator.CreateResponse(response.RegisterPurchaseResponse{}, userId)
				}
			}
		}
	}
	return pr.responseCreator.CreateResponse(response.RegisterPurchaseErrorResponse{}, userId)
}

func (pr *purchaseRepository) VerificationPurchase(request request.PurchaseRegisterRequest) response.AppResponse {
	ctx := context.Background()
	str1 := "E:\\Dev\\GoDev\\MoonWriterService\\credentials.json"
	androidpublisherService, err := androidpublisher.NewService(ctx, option.WithCredentialsFile(str1))
	if err != nil {
		log.Infof("VerificationPurchase NewService: ", err)
	}

	r, errr := androidpublisherService.Purchases.Products.Get("com.mwriter.moonwriter",
		"lemonade35",
		"bgeddnkemanoelbedjokoocc.AO-J1Owsddv6PUW4Ct4TWhiPqs0HgiL0wIBIjZgoWWaKPF9_nbti33qJQcSMzZFcBhrM-Lu7WJORZZr4m3C6iZ_wLuGyFLFp6UDTWR9syP27IAGq0lNo5NHtgNgtIXXTRISSW2g275ig").Do()
	if errr != nil {
		log.Infof("VerificationPurchase Get: ", errr)
	}
	log.Infof("VerificationPurchase Result: ", r)

	return nil
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

func (pr *purchaseRepository) CheckUserExist(userId string) bool {
	id, errHex := primitive.ObjectIDFromHex(userId)
	if errHex != nil {
		return false
	}

	var localUser domain.UserEntity
	err := pr.collectionUsers.FindOne(utils.GetContext(), bson.D{{"_id", id}}).Decode(&localUser)
	if err != nil {
		log.Error("Count: ", err)
		return false
	}

	if (localUser.Id == userId) && (localUser.UserType != config.PREMIUM) {
		return true
	}

	return false
}

func (pr *purchaseRepository) checkAccessTokenExist(accessToken string) bool {
	count, err := pr.collectionSessions.CountDocuments(utils.GetContext(), bson.D{{"access_token", accessToken}})
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
		{"purchase_token", request.PurchaseToken},
		{"order_id", request.OrderId},
		{"purchase_time", request.PurchaseTime},
		{"Sku", request.Sku},
		{"access_token", request.AccessToken},
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
		bson.D{{"$set", bson.D{{"user_type", config.PREMIUM}}}})
	if errUpdate != nil {
		return false
	}

	return true
}
