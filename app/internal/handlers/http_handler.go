package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/email-auth/app/internal/handlers/middlewares"
	"github.com/AleksandrVishniakov/email-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/email-auth/app/internal/services/user_service"
	"github.com/AleksandrVishniakov/email-auth/app/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	userService user_service.UserService
}

func NewHTTPHandler(userService user_service.UserService) *HTTPHandler {
	return &HTTPHandler{userService: userService}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSHeaders())

	router.GET("/auth", h.authUser)
	router.GET("/user/:email", h.getUserByEmail)
	router.GET("/verify/:email", h.verifyEmail)

	return router
}

func (h *HTTPHandler) authUser(c *gin.Context) {
	email := c.Request.URL.Query().Get("email")

	if email == "" {
		rErr := e.NewResponseError(http.StatusBadRequest, "empty email parameter")
		rErr.Abort(c)
		return
	}

	isUserAuthorized, err := h.userService.AuthUser(email)

	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, struct {
		IsUserAuthorized bool `json:"isUserAuthorized"`
	}{IsUserAuthorized: isUserAuthorized})
}

func (h *HTTPHandler) getUserByEmail(c *gin.Context) {
	email, exists := c.Params.Get("email")
	if !exists {
		rErr := e.NewResponseError(http.StatusBadRequest, "email was not provided")
		rErr.Abort(c)
		return
	}

	user, err := h.userService.GetUserByEmail(email)
	if errors.Is(err, user_repo.ErrUserNotFound) {
		rErr := e.NewResponseError(http.StatusNotFound, "email not found")
		rErr.Abort(c)
		return
	}

	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (h *HTTPHandler) verifyEmail(c *gin.Context) {
	email, exists := c.Params.Get("email")
	if !exists {
		rErr := e.NewResponseError(http.StatusBadRequest, "email was not provided")
		rErr.Abort(c)
		return
	}

	codeStr := c.Request.URL.Query().Get("code")
	if len(codeStr) != 6 {
		rErr := e.NewResponseError(http.StatusBadRequest, "invalid code length")
		rErr.Abort(c)
		return
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		rErr := e.NewResponseError(http.StatusBadRequest, err.Error())
		rErr.Abort(c)
		return
	}

	ok, err := h.userService.VerifyEmail(email, code)
	if errors.Is(err, user_repo.ErrUserNotFound) {
		rErr := e.NewResponseError(http.StatusNotFound, "email not found")
		rErr.Abort(c)
		return
	}

	if !ok {
		rErr := e.NewResponseError(http.StatusBadRequest, "email validation failed")
		rErr.Abort(c)
		return
	}

	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}
}
