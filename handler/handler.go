package handler

import (
	"github.com/ShaghayeghFathi/http-monitoring-service/auth"
	"github.com/ShaghayeghFathi/http-monitoring-service/db_manager"
	"github.com/ShaghayeghFathi/http-monitoring-service/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Handler struct {
	dm  *db_manager.DbManager
	sch *service.Scheduler
	ech *echo.Echo
}

func NewHandler(dm *db_manager.DbManager, sch *service.Scheduler) *Handler {
	h := &Handler{dm: dm, sch: sch, ech: echo.New()}
	h.defineRoutes()
	return h
}

func (h *Handler) defineRoutes() {

	h.ech.Use(auth.JWT())

	auth.AddToWhiteList("/users/login", "POST")
	auth.AddToWhiteList("/users", "POST")

	h.ech.POST("/users", h.SignUp)
	h.ech.POST("/users/login", h.Login)

	h.ech.GET("/urls", h.FetchURLs)
	h.ech.POST("/urls", h.CreateURL)
	h.ech.GET("/urls/:urlID", h.GetURLStats)

	h.ech.GET("/alerts", h.FetchAlerts)
}

func (h *Handler) Start(address string) {
	h.ech.Logger.Fatal(h.ech.Start(":8080"))
}

func extractID(c echo.Context) uint {
	e := c.Get("user").(*jwt.Token)
	claims := e.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))
	return id
}
