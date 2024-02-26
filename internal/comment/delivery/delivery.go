package delivery

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	jwtManager "github.com/ell1jah/bmstu_web/internal/pkg/jwt"
	"github.com/ell1jah/bmstu_web/model"
	"github.com/ell1jah/bmstu_web/model/dto"
)

type CommentLogic interface {
	GetPostComments(postId uint64) ([]*model.Comment, error)
	CreateComment(comment *model.Comment) error
}

type handler struct {
	commentService CommentLogic
}

func NewHandler(commentService CommentLogic) *handler {
	return &handler{
		commentService: commentService,
	}
}

func (h *handler) SetRoutes(e *echo.Echo, auth echo.MiddlewareFunc) {
	e.GET("/posts/:postID/comments", h.GetPostComments, auth)
	e.POST("/posts/:postID/comments", h.CreateComment, auth)
}

// GetPostComments godoc
// @Summary      Get post comments
// @Description  Get post comments
// @Tags     comments
// @Produce  application/json
// @Param postID path int true "comment ID"
// @Success  200 {object} []*dto.RespComment "success get community"
// @Failure 405 {object} echo.HTTPError "Method Not Allowed"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 404 {object} echo.HTTPError "item not found"
// @Failure 401 {object} echo.HTTPError "no auth"
// @Router   /posts/{postID}/comments [get]
func (h *handler) GetPostComments(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	comments, err := h.commentService.GetPostComments(postId)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespCommentsFromComments(comments))
}

// CreateComment godoc
// @Summary      Create a comment
// @Description  Create a comment
// @Tags     	 comments
// @Accept	 application/json
// @Produce  application/json
// @Param    comment body dto.ReqComment true "comment info"
// @Success  200 {object} *dto.RespComment "success create comment"
// @Failure 405 {object} echo.HTTPError "invalid http method"
// @Failure 400 {object} echo.HTTPError "bad request"
// @Failure 422 {object} echo.HTTPError "unprocessable entity"
// @Failure 500 {object} echo.HTTPError "internal server error"
// @Failure 401 {object} echo.HTTPError "no auth"
// @Router   /posts/{postID}/comments [post]
func (h *handler) CreateComment(c echo.Context) error {
	var reqComment dto.ReqComment
	err := c.Bind(&reqComment)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqComment)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	comment := reqComment.ToComment()

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	comment.UserID = userClaims.User.ID

	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	comment.PostID = postId

	err = h.commentService.CreateComment(comment)
	if err != nil {
		c.Logger().Error(err)
		return handleError(err)
	}

	return c.JSON(http.StatusCreated, dto.RespCommentFromComment(comment))
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
