package delivery

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ell1jah/bmstu_web/internal/pkg/httperror"
	"github.com/ell1jah/bmstu_web/model"
	"github.com/ell1jah/bmstu_web/model/dto"
)

type ImageLogic interface {
	GetImage(imageId string) (io.Reader, error)
	CreateImage(io.Reader) (string, error)
}

type handler struct {
	imageService ImageLogic
}

func NewHandler(imageService ImageLogic) *handler {
	return &handler{
		imageService: imageService,
	}
}

func (h *handler) SetRoutes(e *echo.Echo, auth echo.MiddlewareFunc) {
	e.GET("/images/:imageID", h.GetImage, auth)
	e.POST("/images", h.CreateImage, auth)
}

func (h *handler) GetImage(c echo.Context) error {
	imageId := c.Param("imageID")
	if imageId == "" {
		c.Logger().Error("empty imageId")
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	}

	image, err := h.imageService.GetImage(imageId)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.Stream(http.StatusOK, "Image/png", image)
}

func (h *handler) CreateImage(c echo.Context) error {
	file, err := c.FormFile("Image")
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Not Attachment in form")
	}
	src, err := file.Open()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	defer src.Close()

	id, err := h.imageService.CreateImage(src)
	if err != nil {
		c.Logger().Error(err)
		return httperror.HandleError(err)
	}

	return c.JSON(http.StatusCreated, dto.RespImageFromID(id))
}
