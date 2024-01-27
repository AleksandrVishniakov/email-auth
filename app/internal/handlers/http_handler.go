package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/repositories/user_repo"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/internal/services/user_service"
	"github.com/AleksandrVishniakov/url-shortener-auth/app/pkg/e"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type HTTPHandler struct {
	userService user_service.UserService
}

func NewHTTPHandler(userService user_service.UserService) *HTTPHandler {
	return &HTTPHandler{userService: userService}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

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

	err := h.userService.NewUser(email, c.Request.Host)

	if errors.Is(err, user_service.ErrEmailIsExists) {
		rErr := e.NewResponseError(http.StatusConflict, err.Error())
		rErr.Abort(c)
		return
	}

	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}
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

	hash := c.Request.URL.Query().Get("h")
	if len(hash) != 64 {
		rErr := e.NewResponseError(http.StatusBadRequest, "invalid hash")
		rErr.Abort(c)
		return
	}

	err := h.userService.VerifyEmail(email, hash)
	if errors.Is(err, user_repo.ErrUserNotFound) {
		rErr := e.NewResponseError(http.StatusNotFound, "email not found")
		rErr.Abort(c)
		return
	}

	if errors.Is(err, user_service.ErrEmailValidation) {
		rErr := e.NewResponseError(http.StatusBadRequest, "email validation failed")
		rErr.Abort(c)
		return
	}

	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}

	htmlBlob, err := template.ParseFiles("web/email_verified_page.html")
	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}

	err = htmlBlob.Execute(c.Writer, nil)
	if err != nil {
		rErr := e.NewResponseError(http.StatusInternalServerError, err.Error())
		rErr.Abort(c)
		return
	}
}
