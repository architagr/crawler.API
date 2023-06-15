package controller

import (
	"LoginAPI/models"
	"LoginAPI/service"

	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	CreateUser(c *fiber.Ctx) error
}

type authController struct {
	service service.IAuthService
}

var authControllerObj IAuthController

func InitAuthController(serviceObj service.IAuthService) IAuthController {
	if authControllerObj == nil {
		authControllerObj = &authController{
			service: serviceObj,
		}
	}
	return authControllerObj
}

func (ctlr *authController) CreateUser(c *fiber.Ctx) error {
	filter := new(models.LoginDetails)
	if err := c.BodyParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	token, err := ctlr.service.CreateCognitoUser(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.JSON(token)
}
