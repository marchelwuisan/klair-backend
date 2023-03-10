package responses

import "github.com/gofiber/fiber/v2"

type WalletResponse struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Data       *fiber.Map `json:"data"`
	WalletData *fiber.Map `json:"wallet_data"`
}

type WalletGetResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
