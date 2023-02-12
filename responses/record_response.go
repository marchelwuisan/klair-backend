package responses

import "github.com/gofiber/fiber/v2"

type RecordResponse struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Data       *fiber.Map `json:"data"`
	RecordData *fiber.Map `json:"record_data"`
}