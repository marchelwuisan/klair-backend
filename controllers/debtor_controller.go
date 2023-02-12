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

var debtorCollection *mongo.Collection = configs.GetCollection(configs.DB, "debtors")

func CreateDebtor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var debtor models.Debtor
	userId := c.Params("userId")

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&debtor); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&debtor); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newDebtor := models.Debtor{
		ID:                primitive.NewObjectID(),
		UserId:            userId,
		Name:              debtor.Name,
		PayableBalance:    0,
		RecievableBalance: 0,
	}

	result, err := debtorCollection.InsertOne(ctx, newDebtor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.WalletResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, WalletData: &fiber.Map{"walletData": newDebtor}})
}

func GetDebtor(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	debtorId := c.Params("debtorId")
	var debtor models.Debtor
	var records []models.Record
	// var record models.Record

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(debtorId)

	err := debtorCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&debtor)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.WalletGetResponse{Status: http.StatusInternalServerError, Message: "error getting wallet", Data: &fiber.Map{"data": err.Error()}})
	}

	results, err2 := recordCollection.Find(ctx, bson.M{"debtorid": debtorId})
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

	debtor.Records = records

	return c.Status(http.StatusOK).JSON(responses.WalletGetResponse{Status: http.StatusOK, Message: debtorId, Data: &fiber.Map{"data": debtor}})
}
