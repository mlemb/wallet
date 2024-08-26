package endpoint

import "time"

type TransferModel struct {
	ID     int64     `json:"id"`
	Type   string    `json:"type"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

type ListTransferRequest struct {
	Type       string    `query:"type"`
	From       string    `query:"from"`
	To         string    `query:"to"`
	AmountFrom float64   `query:"amount_from"`
	AmountTo   float64   `query:"amount_to"`
	TimeFrom   time.Time `query:"time_from"`
	TimeTo     time.Time `query:"time_to"`
	Page       int       `query:"page"`
	PageSize   int       `query:"page_size"`
}

type ListTransferResponse struct {
	Data  []*TransferModel `json:"data"`
	Total int64            `json:"total"`
}

type GetTransferRequest struct {
	ID string `param:"id"`
}

type GetTransferResponse struct {
	Data *TransferModel `json:"data"`
}

type CreateTransferRequest struct {
	Type   string    `json:"type"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

type CreateTransferResponse struct {
	Data *TransferModel `json:"data"`
}

type UpdateTransferRequest struct {
	ID     int64     `param:"id"`
	Type   string    `json:"type"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

type UpdateTransferResponse struct {
}

type DeleteTransferRequest struct {
	ID int64 `param:"id"`
}

type DeleteTransferResponse struct {
}
