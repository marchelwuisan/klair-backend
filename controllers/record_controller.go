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

var recordCollection *mongo.Collection = configs.GetCollection(configs.DB, "records")

func CreateRecord(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var record models.Record

	userId := c.Params("userId")

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&record); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&record); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// if record.Attachment != "" {
	// 	bucket, err := c..DefaultBucket()
	//     if err != nil {
	//         return err
	//     }

	//     object := bucket.Object(fileName)
	//     writer := object.NewWriter(ctx)
	// }

	newRecord := models.Record{
		ID:           primitive.NewObjectID(),
		UserId:       userId,
		WalletId:     record.WalletId,
		DebtorId:     record.DebtorId,
		WalletFromId: record.WalletFromId,
		WalletToId:   record.WalletToId,
		Amount:       record.Amount,
		Type:         record.Type,
		Category:     record.Category,
		Attachment:   record.Attachment,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	if record.Type == "income" {
		var wallet models.Wallet

		objId, _ := primitive.ObjectIDFromHex(record.WalletId)

		err2 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&wallet)
		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error finding wallet", Data: &fiber.Map{"data": err2.Error()}})
		}

		result, err := recordCollection.InsertOne(ctx, newRecord)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error inserting record", Data: &fiber.Map{"data": err.Error()}})
		}

		update := bson.M{"balance": wallet.Balance + record.Amount}
		result3, err3 := walletCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error updating wallet balance", Data: &fiber.Map{"data": err3.Error()}})
		}
		//get updated wallet details
		var updatedWallet models.Wallet
		if result3.MatchedCount == 1 {
			err3 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWallet)

			if err3 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
			}
		}
		return c.Status(http.StatusCreated).JSON(responses.RecordResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, RecordData: &fiber.Map{"recordData": newRecord}})

	} else if record.Type == "expense" {
		var wallet models.Wallet

		objId, _ := primitive.ObjectIDFromHex(record.WalletId)

		err2 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&wallet)
		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error getting wallet balance", Data: &fiber.Map{"data": err2.Error()}})
		}

		result, err := recordCollection.InsertOne(ctx, newRecord)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error inserting record", Data: &fiber.Map{"data": err.Error()}})
		}

		update := bson.M{"balance": wallet.Balance - record.Amount}
		result3, err3 := walletCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
		}
		//get updated wallet details
		var updatedWallet models.Wallet
		if result3.MatchedCount == 1 {
			err3 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWallet)

			if err3 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
			}
		}
		return c.Status(http.StatusCreated).JSON(responses.RecordResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, RecordData: &fiber.Map{"recordData": newRecord}})

	} else if record.Type == "transfer" {
		var wallet models.Wallet
		var wallet2 models.Wallet

		WalletFromObjId, _ := primitive.ObjectIDFromHex(record.WalletFromId)
		WalletToObjId, _ := primitive.ObjectIDFromHex(record.WalletToId)

		if record.WalletFromId == "" || record.WalletToId == "" {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "Wallets not found", Data: &fiber.Map{"data": record}})
		}

		err2 := walletCollection.FindOne(ctx, bson.M{"_id": WalletFromObjId}).Decode(&wallet)
		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error getting wallet balance", Data: &fiber.Map{"data": err2.Error()}})
		}

		err3 := walletCollection.FindOne(ctx, bson.M{"_id": WalletToObjId}).Decode(&wallet2)
		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error getting wallet balance", Data: &fiber.Map{"data": err2.Error()}})
		}

		result, err := recordCollection.InsertOne(ctx, newRecord)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error inserting record", Data: &fiber.Map{"data": err.Error()}})
		}

		update := bson.M{"balance": wallet.Balance - record.Amount}
		update2 := bson.M{"balance": wallet2.Balance + record.Amount}

		result2, err2 := walletCollection.UpdateOne(ctx, bson.M{"_id": WalletFromObjId}, bson.M{"$set": update})

		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err2.Error()}})
		}

		result3, err3 := walletCollection.UpdateOne(ctx, bson.M{"_id": WalletToObjId}, bson.M{"$set": update2})

		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
		}
		//get updated wallet details
		var updatedWallet models.Wallet
		if result2.MatchedCount == 1 {
			err2 := walletCollection.FindOne(ctx, bson.M{"_id": WalletFromObjId}).Decode(&updatedWallet)

			if err2 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err2.Error()}})
			}
		}
		var updatedWallet2 models.Wallet
		if result3.MatchedCount == 1 {
			err3 := walletCollection.FindOne(ctx, bson.M{"_id": WalletToObjId}).Decode(&updatedWallet2)

			if err3 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
			}
		}

		return c.Status(http.StatusCreated).JSON(responses.RecordResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, RecordData: &fiber.Map{"recordData": newRecord}})

	} else if record.Type == "recievable" {

		var wallet models.Wallet

		objId, _ := primitive.ObjectIDFromHex(record.WalletId)

		err2 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&wallet)
		if err2 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error finding wallet", Data: &fiber.Map{"data": err2.Error()}})
		}

		result, err := recordCollection.InsertOne(ctx, newRecord)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error inserting record", Data: &fiber.Map{"data": err.Error()}})
		}

		update := bson.M{"balance": wallet.Balance + record.Amount}
		result3, err3 := walletCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error updating wallet balance", Data: &fiber.Map{"data": err3.Error()}})
		}
		//get updated wallet details
		var updatedWallet models.Wallet
		if result3.MatchedCount == 1 {
			err3 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWallet)

			if err3 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
			}
		}
		return c.Status(http.StatusCreated).JSON(responses.RecordResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, RecordData: &fiber.Map{"recordData": newRecord}})

	} else if record.Type == "payable" {
		var wallet models.Wallet

		objId, _ := primitive.ObjectIDFromHex(record.WalletId)

		update := bson.M{"balance": wallet.Balance - record.Amount}
		result3, err3 := walletCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err3 != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
		}
		//get updated wallet details
		var updatedWallet models.Wallet
		if result3.MatchedCount == 1 {
			err3 := walletCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedWallet)

			if err3 != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err3.Error()}})
			}
		}
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.WalletResponse{Status: http.StatusInternalServerError, Message: "type not found", Data: &fiber.Map{"data": "error"}})
}

func GetRecordByType(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	recordType := c.Params("recordType")
	var records []models.Record
	// var record models.Record

	defer cancel()

	results, err := recordCollection.Find(ctx, bson.M{"$and": bson.A{bson.M{"userid": userId}, bson.M{"type": recordType}}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error getting records", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleRecord models.Record
		if err = results.Decode(&singleRecord); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		records = append(records, singleRecord)
	}

	return c.Status(http.StatusOK).JSON(responses.WalletGetResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": records}})
}

func GetRecordByCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	recordCategory := c.Params("recordCategory")
	var records []models.Record
	// var record models.Record

	defer cancel()

	results, err := recordCollection.Find(ctx, bson.M{"$and": bson.A{bson.M{"userid": userId}, bson.M{"category": recordCategory}}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.RecordResponse{Status: http.StatusInternalServerError, Message: "error getting records", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleRecord models.Record
		if err = results.Decode(&singleRecord); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		records = append(records, singleRecord)
	}

	return c.Status(http.StatusOK).JSON(responses.WalletGetResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": records}})
}
