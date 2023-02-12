package controllers

import (
	"context"
	"golanglearn/configs"
	"golanglearn/models"
	"golanglearn/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var walletCollection *mongo.Collection = configs.GetCollection(configs.DB, "wallets")

func CreateWallet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var wallet models.Wallet
	userId := c.Params("userId")

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&wallet); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&wallet); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newWallet := models.Wallet{
		ID:        primitive.NewObjectID(),
		UserId:    userId,
		Name:      wallet.Name,
		Type:      wallet.Type,
		Currency:  wallet.Currency,
		Balance:   wallet.Balance,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	result, err := walletCollection.InsertOne(ctx, newWallet)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.WalletResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, WalletData: &fiber.Map{"walletData": newWallet}})
}

func GetWallet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	walletId := c.Params("walletId")
	var wallet models.Wallet
	var records []models.Record
	// var record models.Record

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(walletId)

	err := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&wallet)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.WalletGetResponse{Status: http.StatusInternalServerError, Message: "error getting wallet", Data: &fiber.Map{"data": err.Error()}})
	}

	results, err2 := recordCollection.Find(ctx, bson.M{"$or": bson.A{bson.M{"walletid": walletId}, bson.M{"walletfromid": walletId}, bson.M{"wallettoid": walletId}}})
	if err2 != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error getting records", Data: &fiber.Map{"data": err2.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleRecord models.Record
		if err = results.Decode(&singleRecord); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		records = append(records, singleRecord)
	}

	wallet.Records = records

	return c.Status(http.StatusOK).JSON(responses.WalletGetResponse{Status: http.StatusOK, Message: walletId, Data: &fiber.Map{"data": wallet}})
}

func GetAllWallet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")

	var wallet []models.Wallet
	var records []models.Record
	// var record models.Record

	defer cancel()

	results, err := walletCollection.Find(ctx, bson.M{"userid": userId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.WalletGetResponse{Status: http.StatusInternalServerError, Message: " error getting wallet", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleWallet models.Wallet
		if err = results.Decode(&singleWallet); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		results2, err2 := recordCollection.Find(ctx, bson.M{"walletid": singleWallet.ID})
		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error getting records", Data: &fiber.Map{"data": err2.Error()}})
		}

		defer results2.Close(ctx)
		for results2.Next(ctx) {
			var singleRecord models.Record
			if err = results2.Decode(&singleWallet); err != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
			}

			records = append(records, singleRecord)
		}

		singleWallet.Records = records

		wallet = append(wallet, singleWallet)
	}

	// results, err2 := recordCollection.Find(ctx, bson.D{{"walletid", walletId}})
	// if err2 != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error getting records", Data: &fiber.Map{"data": err2.Error()}})
	// }

	// defer results.Close(ctx)
	// for results.Next(ctx) {
	// 	var singleRecord models.Record
	// 	if err = results.Decode(&singleRecord); err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	// 	}

	// 	records = append(records, singleRecord)
	// }

	// wallet.Records = records

	return c.Status(http.StatusOK).JSON(responses.WalletGetResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": wallet}})
}
