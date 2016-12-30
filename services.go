package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func index(c echo.Context) error {

	return c.JSON(http.StatusOK, "index")

}

func create(c echo.Context) error {

	return c.JSON(http.StatusOK, "create")

}

func read(c echo.Context) error {

	return c.JSON(http.StatusOK, "read")

}

func update(c echo.Context) error {

	return c.JSON(http.StatusOK, "update")

}

func delete(c echo.Context) error {

	return c.JSON(http.StatusOK, "delete")

}

func routes(prefix string) {
	g := e.Group(prefix)
	g.GET("/", index)
	g.POST("/", create)
	g.GET("/:id", read)
	g.PUT("/:id", update)
	g.DELETE("/:id", delete)
}
