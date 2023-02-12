package controllers

import (
	"context"
	"golanglearn/configs"
	"golanglearn/models"
	"golanglearn/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var budgetCollection *mongo.Collection = configs.GetCollection(configs.DB, "budgets")

func CreateBudget(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var budget models.Budget
	userId := c.Params("userId")

	defer cancel()

	//validate the request body
	if err := c.BodyParser(&budget); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BudgetResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&budget); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.BudgetResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newBudget := models.Budget{
		ID:       primitive.NewObjectID(),
		UserId:   userId,
		WalletId: budget.WalletId,
		Name:     budget.Name,
		// PeriodStart: budget.PeriodStart,
		// PeriodEnd:   budget.PeriodEnd,
		PeriodStart: time.Now().Unix(),
		PeriodEnd:   time.Now().Unix(),
		Amount:      budget.Amount,
		CategoryId:  budget.CategoryId,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		DataStatus:  1,
	}

	result, err := budgetCollection.InsertOne(ctx, newBudget)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.BudgetResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.BudgetResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
