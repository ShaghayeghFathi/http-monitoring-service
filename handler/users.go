package handler

import (
	"net/http"

	"github.com/ShaghayeghFathi/http-monitoring-service/auth"
	"github.com/ShaghayeghFathi/http-monitoring-service/common"
	"github.com/ShaghayeghFathi/http-monitoring-service/model"
	"github.com/labstack/echo"
)

type UserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewUserResponse(user *model.User) *UserResponse {
	token, _ := auth.GenerateJWT(user.ID)
	ur := &UserResponse{Username: user.Username, Token: token}
	return ur
}

type userAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func bindToAuthRequest(c echo.Context) (*userAuthRequest, error) {
	authRequest := &userAuthRequest{}
	if err := c.Bind(authRequest); err != nil {
		return nil, common.NewRequestError("error binding user request", err, http.StatusBadRequest)
	}
	return authRequest, nil
}

func (h *Handler) Login(c echo.Context) error {
	authRequest, err := bindToAuthRequest(c)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: authRequest.Username,
		Password: authRequest.Password,
	}

	// retrieving user from database
	u, err := h.dm.GetUserByUserName(user.Username)
	if err != nil || !u.ValidatePassword(user.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Username or Password", err)
	}
	return c.JSON(http.StatusOK, NewUserResponse(u))
}

func (h *Handler) SignUp(c echo.Context) error {
	authRequest, err := bindToAuthRequest(c)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: authRequest.Username,
		Password: authRequest.Password,
	}

	user.Password, _ = model.HashPassword(user.Password)
	err = h.dm.AddUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not add user to database", err)
	}

	return c.JSON(http.StatusCreated, NewUserResponse(user))
}

func (h *Handler) FetchAlerts(c echo.Context) error {
	userID := extractID(c)
	alerts, err := h.dm.FetchAlerts(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Could not get alerts from database", err)
	}
	serializedAlerts := make([]*model.SerializedUrl, 0, len(alerts))
	for _, url := range alerts {
		serializedAlerts = append(serializedAlerts, url.Serialize())
	}
	return c.JSON(http.StatusOK, serializedAlerts)
}
