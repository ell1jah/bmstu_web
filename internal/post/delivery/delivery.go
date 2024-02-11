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

type PostLogic interface {
	GetPost(userId, postId uint64) (*model.Post, error)
	GetUsersPosts(askerId, ownerId uint64) ([]*model.Post, error)
	GetPostsWithParams(userId uint64, params model.PostParams) ([]*model.Post, error)
	CreatePost(post *model.Post) error
	DeletePost(userId, postId uint64) error
	LikePost(userId, postId uint64) error
	DislikePost(userId, postId uint64) error
	UnratePost(userId, postId uint64) error
}

type handler struct {
	postService PostLogic
}

func NewHandler(postService PostLogic) *handler {
	return &handler{
		postService: postService,
	}
}

func (h *handler) SetRoutes(e *echo.Echo, auth echo.MiddlewareFunc) {
	e.POST("/posts", h.CreatePost, auth)
	e.DELETE("/posts/:postID", h.DeletePost, auth)
	e.PUT("/posts/:postID/like", h.LikePost, auth)
	e.PUT("/posts/:postID/dislike", h.DislikePost, auth)
	e.DELETE("/posts/:postID/unrate", h.UnratePost, auth)

	e.GET("/posts/:postID", h.GetPost, auth)
	e.GET("/posts", h.GetPostsWithParams, auth)
	e.GET("/users/:userID/posts", h.GetUsersPosts, auth)
}

func (h *handler) GetPost(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	post, err := h.postService.GetPost(userClaims.User.ID, postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespPostFromPost(post))
}

func (h *handler) GetUsersPosts(c echo.Context) error {
	ownerId, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	posts, err := h.postService.GetUsersPosts(userClaims.User.ID, ownerId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespPostsFromPosts(posts))
}

func (h *handler) GetPostsWithParams(c echo.Context) error {
	var reqParams dto.ReqPostParams
	err := c.Bind(&reqParams)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqParams)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	params := reqParams.ToPostParams()

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	posts, err := h.postService.GetPostsWithParams(userClaims.User.ID, *params)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusOK, dto.RespPostsFromPosts(posts))
}

func (h *handler) CreatePost(c echo.Context) error {
	var reqPost dto.ReqPost
	err := c.Bind(&reqPost)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	_, err = govalidator.ValidateStruct(reqPost)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	post := reqPost.ToPost()

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	post.UserID = userClaims.User.ID

	err = h.postService.CreatePost(post)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusCreated, dto.RespPostFromPost(post))
}

func (h *handler) DeletePost(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	err = h.postService.DeletePost(userId, postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) LikePost(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	err = h.postService.LikePost(userId, postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) DislikePost(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	err = h.postService.DislikePost(userId, postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) UnratePost(c echo.Context) error {
	postId, err := strconv.ParseUint(c.Param("postID"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	userClaims, ok := c.Get("user").(*jwt.Token).Claims.(*jwtManager.Claims)
	userId := userClaims.User.ID
	if !ok {
		c.Logger().Error(model.ErrInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, model.ErrInternalServerError.Error())
	}

	err = h.postService.UnratePost(userId, postId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.NoContent(http.StatusOK)
}
