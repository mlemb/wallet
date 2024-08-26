package endpoint

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mlemb/wallet/backend/internal/db"
	"github.com/mlemb/wallet/backend/internal/model"
)

func NewTransferEndpoint() *TransferEndpoint {
	return &TransferEndpoint{}
}

type TransferEndpoint struct{}

func (e *TransferEndpoint) Routes(server *echo.Echo) {
	server.GET("/transfers", e.List)
	server.POST("/transfers", e.Create)
	server.PUT("/transfers/:id", e.Update)
	server.DELETE("/transfers/:id", e.Delete)
}

func (e *TransferEndpoint) GroupRoutes(group *echo.Group) {
	group.GET("/transfers", e.List)
	group.GET("/transfers/:id", e.Get)
	group.POST("/transfers", e.Create)
	group.PUT("/transfers/:id", e.Update)
	group.DELETE("/transfers/:id", e.Delete)
}

func (e *TransferEndpoint) List(c echo.Context) error {
	ctx := c.Request().Context()

	var request ListTransferRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	params := &db.TransfersParams{
		Type:       request.Type,
		From:       request.From,
		To:         request.To,
		AmountFrom: request.AmountFrom,
		AmountTo:   request.AmountTo,
		TimeFrom:   request.TimeFrom,
		TimeTo:     request.TimeTo,
		Page:       request.Page,
		PageSize:   request.PageSize,
	}

	transfers, err := db.QueryTransfers(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := db.CountTransfers(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := &ListTransferResponse{
		Total: count,
	}

	for _, transfer := range transfers {
		response.Data = append(response.Data, &TransferModel{
			ID:     transfer.ID,
			Type:   transfer.Type,
			From:   transfer.From,
			To:     transfer.To,
			Amount: transfer.Amount,
			Time:   transfer.Time,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (e *TransferEndpoint) Get(c echo.Context) error {
	ctx := c.Request().Context()

	var request GetTransferRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transfer, err := db.QueryTransferByID(ctx, request.ID)
	if errors.Is(err, db.ErrTransferNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := &GetTransferResponse{
		Data: &TransferModel{
			ID:     transfer.ID,
			Type:   transfer.Type,
			From:   transfer.From,
			To:     transfer.To,
			Amount: transfer.Amount,
			Time:   transfer.Time,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (e *TransferEndpoint) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var request CreateTransferRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transfer := &model.Transfer{
		Type:   request.Type,
		From:   request.From,
		To:     request.To,
		Amount: request.Amount,
		Time:   request.Time,
	}

	if err := db.CreateTransfer(ctx, transfer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: return transfer id

	return c.NoContent(http.StatusCreated)
}

func (e *TransferEndpoint) Update(c echo.Context) error {
	ctx := c.Request().Context()

	var request UpdateTransferRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transfer := &model.Transfer{
		ID:     request.ID,
		Type:   request.Type,
		From:   request.From,
		To:     request.To,
		Amount: request.Amount,
		Time:   request.Time,
	}

	if err := db.UpdateTransfer(ctx, transfer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// TODO: return transfer id

	return c.NoContent(http.StatusOK)
}

func (e *TransferEndpoint) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	var request DeleteTransferRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := db.DeleteTransferByID(ctx, request.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
