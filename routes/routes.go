package routes

import (
	"golanglearn/controllers"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	// app.Get("/", controllers.Test)
	app.Post("/user/login", controllers.Login)
	// app.Post("/user", controllers.CreateUser)
	app.Post("/user", controllers.CreateFirebaseUser)
	app.Get("/user/:userId", controllers.GetAUser)
	app.Get("/user/firebase/:userFirebaseUid", controllers.FirebaseGetUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Post("/user/:userId/wallet/create", controllers.CreateWallet)
	app.Get("/user/:userId/wallet/:walletId", controllers.GetWallet)
	app.Get("/user/:userId/wallets/", controllers.GetAllWallet)
	app.Post("/user/:userId/record/create", controllers.CreateRecord)
	app.Get("/user/:userId/record/:recordType", controllers.GetRecordByType)
	app.Get("/user/:userId/record/:recordCategory", controllers.GetRecordByCategory)
	app.Post("/user/:userId/debtor/create", controllers.CreateDebtor)
	app.Post("/user/:userId/budget/create", controllers.CreateBudget)
}
