package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/samarthasthan/luganodes-task/internal/store/database"
	"github.com/samarthasthan/luganodes-task/internal/store/database/mysql/sqlc"
)

// For Dependency Injection we are using database.Database interface instead of concrete type
type Controller struct {
	mysql database.Database
	redis database.Database
}

func NewController(mysql database.Database, redis database.Database) *Controller {
	return &Controller{mysql: mysql, redis: redis}
}

func (c *Controller) GetDeposit(ctx echo.Context) error {
	// Get paggination page and limit
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	// Convert page and limit to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ctx.JSON(400, "Invalid page")
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ctx.JSON(400, "Invalid limit")
	}
	// Type assertion for MySQL
	mysql, ok := c.mysql.(*database.MySQL)
	if !ok {
		return ctx.JSON(500, "MySQL is not initialized")
	}

	// Get Deposit from MySQL
	deposits, err := mysql.Queries.GetDeposites(ctx.Request().Context(), sqlc.GetDepositesParams{
		Limit:  int32(limitInt),
		Offset: int32((pageInt - 1) * limitInt),
	})
	if err != nil {
		return ctx.JSON(500, err.Error())
	}

	return ctx.JSON(200, deposits)

}
