package handlers

import (
	"github.com/PainQlers/backend/models" // Import package models ของเราเข้ามา
	"github.com/gofiber/fiber/v2"
)

// GetUserByID ทำหน้าที่ดึงข้อมูล User ตาม ID ที่ส่งมา
func GetUserByID(c *fiber.Ctx) error {
	// ดึง id จาก path parameter เช่น /v1/users/123
	userID := c.Params("id")

	// จำลองข้อมูลจาก Database (ในอนาคตตรงนี้จะไปดึงจาก DB จริง)
	mockUser := models.User{
		ID:    userID,
		Name:  "Punsakorn",
		Email: "punsakorn.a@example.com",
	}

	// ส่งกลับเป็น JSON พร้อม HTTP Status 200 OK
	return c.Status(fiber.StatusOK).JSON(mockUser)
}
