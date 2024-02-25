package delivery

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/ell1jah/bmstu_web/internal/pkg/httperror"
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

func (h *handler) GetPostComments(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	comments, err := h.commentService.GetPostComments(postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespCommentsFromComments(comments))
}

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
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusCreated, dto.RespCommentFromComment(comment))
}
