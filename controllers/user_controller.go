package controllers

import (
	"context"
	"golanglearn/configs"
	helper "golanglearn/helpers"
	"golanglearn/models"
	"golanglearn/responses"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()
var fireApp *auth.Client = configs.SetupFirebase()

// func HashPassword(password string) string {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return string(bytes)
// }

// VerifyPassword checks the input password while verifying it with the passward in the DB.
// func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
// 	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
// 	check := true
// 	msg := ""

// 	if err != nil {
// 		msg = fmt.Sprintf("login or password is incorrect")
// 		check = false
// 	}

// 	return check, msg
// }

func FirebaseGetUser(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	uid := c.Params("userFirebaseUid")
	var user models.User

	defer cancel()

	u, err := fireApp.GetUser(ctx, uid)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error getting user", Data: &fiber.Map{"data": err.Error()}})
	}

	err2 := userCollection.FindOne(ctx, bson.M{"firebaseuid": uid}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err2.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": user}, UserData: &fiber.Map{"userData": u}})
}

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Currency:  user.Currency,
		// Password:     password,
		Phone:         user.Phone,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		FirebaseUid:   user.FirebaseUid,
		PushToken:     user.PushToken,
		DataStatus:    1,
		IsFirstSignIn: 1,
		// User_id: user.ID.Hex(),
		// Token:        token,
		// RefreshToken: refreshToken,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}, UserData: &fiber.Map{"userData": newUser}})
}

func CreateFirebaseUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(false).
		PhoneNumber(user.Phone).
		Password(user.Password).
		DisplayName(user.Username).
		Disabled(false)
	u, err := fireApp.CreateUser(ctx, params)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error creating user in firebase", Data: &fiber.Map{"data": err.Error()}})
	}

	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Currency:  user.Currency,
		// Password:     password,
		Phone:         user.Phone,
		FirebaseUid:   u.UID,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		PushToken:     user.PushToken,
		DataStatus:    1,
		IsFirstSignIn: 1,
		// User_id: user.ID.Hex(),
		// Token:        token,
		// RefreshToken: refreshToken,
	}

	result, err2 := userCollection.InsertOne(ctx, newUser)
	if err2 != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error creating user in mongo", Data: &fiber.Map{"data": err2.Error()}})
	}
	log.Printf("Successfully created user: %v\n", result)

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": u}, UserData: &fiber.Map{"userData": newUser}})
}

func Login(c *fiber.Ctx) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var user models.Login
	var foundUser models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: user.Email, Data: &fiber.Map{"data": err.Error()}, UserData: &fiber.Map{"userData": user}})
	}

	// passwordIsValid, msg := VerifyPassword(user.Password, user.Password)
	// defer cancel()
	// if passwordIsValid != true {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: msg, Data: &fiber.Map{"data": err.Error()}})
	// }

	token, refreshToken, _ := helper.GenerateAllTokens(user.Email, foundUser.ID)

	helper.UpdateAllTokens(token, refreshToken, foundUser.ID.Hex())

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: foundUser.ID.Hex(), Data: &fiber.Map{"data": foundUser}})

}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}, UserData: &fiber.Map{"data": user}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"username": user.Username, "email": user.Email, "currency": user.Currency}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

func Test(c *fiber.Ctx) error {

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "test", Data: &fiber.Map{"data": "test"}},
	)
}
