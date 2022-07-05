package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/budhip/example-user/model"

	userSrv "github.com/budhip/example-user/service"
)

// NewUserHandler will initialize the users/ resources endpoint
func NewUserHandler(app fiber.Router, userSrv userSrv.UserService) {
	app.Post("/", createUser(userSrv))
	app.Get("/:id", getUserByID(userSrv))
}

func createUser(userSrv userSrv.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("start createUser in handler")
		rpReq := &model.AddNewUserRequest{}

		if err := c.BodyParser(rpReq); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": err.Error(),
			})

		}

		result, err := userSrv.CreateUser(c.Context(), rpReq)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"status":       "error",
				"error_detail": err,
				"error":        err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"message": "Users Created",
			"data":    result,
		})
	}
}

func getUserByID(userSrv userSrv.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID, err := strconv.ParseUint(id, 0, 64)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		result, err := userSrv.GetUserByID(c.Context(), int64(userID))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"status":       "error",
				"error_detail": err,
				"error":        err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Users Found",
			"data":    result,
		})
	}
}