package delivery

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	jwtManager "github.com/ell1jah/bmstu_web/internal/pkg/jwt"
	"github.com/ell1jah/bmstu_web/model"
	"github.com/ell1jah/bmstu_web/model/dto"
)

type UserLogic interface {
	GetUserByID(id uint64) (*model.User, error)
	ChangePass(chpass *model.UserChangePass) error
	SignIn(user *model.User) (*model.User, error)
	SignUp(user *model.User) (*model.User, error)
}

type SessionManager interface {
	CreateSession(user *jwtManager.UserClaims) (string, error)
}

type handler struct {
	userService    UserLogic
	sessionManager SessionManager
}

func NewHandler(userService UserLogic, sessionManager SessionManager) *handler {
	return &handler{
		userService:    userService,
		sessionManager: sessionManager,
	}
}

func (h *handler) SetRoutes(e *echo.Echo, auth echo.MiddlewareFunc) {
	e.GET("/users/me", h.GetMe, auth)
	e.POST("/users/changepass", h.ChangePass, auth)

	e.POST("/users/signin", h.SignIn)
	e.POST("/users/signup", h.SignUp)
}

func (h *handler) GetMe(c echo.Context) error {
	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	user, err := h.userService.GetUserByID(userId)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespGetMeFromUser(user))
}

func (h *handler) ChangePass(c echo.Context) error {
	var reqPass dto.ReqСhangePass
	err := c.Bind(&reqPass)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqPass)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	chpass := reqPass.ToChangePass()

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	chpass.ID = userId
	err = h.userService.ChangePass(chpass)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrConflictPassword.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) SignIn(c echo.Context) error {
	var reqSign dto.ReqSign
	err := c.Bind(&reqSign)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqSign)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	sign := reqSign.ToUser()

	user, err := h.userService.SignIn(sign)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	token, err := h.sessionManager.CreateSession(jwtManager.FromModelUsertoUserClaims(user))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	return c.JSON(http.StatusOK, dto.RespTokenFromString(token))
}

func (h *handler) SignUp(c echo.Context) error {
	var reqSign dto.ReqSign
	err := c.Bind(&reqSign)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqSign)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	sign := reqSign.ToUser()

	user, err := h.userService.SignUp(sign)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	token, err := h.sessionManager.CreateSession(jwtManager.FromModelUsertoUserClaims(user))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	return c.JSON(http.StatusCreated, dto.RespTokenFromString(token))
}

func handleError(err error) *echo.HTTPError {
	causeErr := errors.Cause(err)
	switch {
	case errors.Is(causeErr, model.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, model.ErrNotFound.Error())
	case errors.Is(causeErr, model.ErrBadRequest):
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	case errors.Is(causeErr, model.ErrPermissionDenied):
		return echo.NewHTTPError(http.StatusForbidden, model.ErrPermissionDenied.Error())
	case errors.Is(causeErr, model.ErrInvalidPassword):
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrInvalidPassword.Error())
	case errors.Is(causeErr, model.ErrConflictPassword):
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrConflictPassword.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, causeErr.Error())
	}
}
