package rinha

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jellydator/validation"
)

type errorResponseCode int

const (
	internalServerErrorCode errorResponseCode = iota + 1
	malformedRequest
	requestValidationFailedCode
	clientNotFound
	insufficientBalance
)

var errorResponses = map[string]errorResponse{
	"internal_server_error": {
		Code:       internalServerErrorCode,
		Message:    "internal server error",
		StatusCode: http.StatusInternalServerError,
	},
	"malformed_request": {
		Code:       malformedRequest,
		Message:    "malformed request body",
		StatusCode: http.StatusBadRequest,
	},
	"request_validation_failed": {
		Code:       requestValidationFailedCode,
		Message:    "request validation failed",
		StatusCode: http.StatusUnprocessableEntity,
	},
	"client_not_found": {
		Code:       clientNotFound,
		Message:    "client not found",
		StatusCode: http.StatusNotFound,
	},
	"insufficient_balance": {
		Code:       insufficientBalance,
		Message:    "insufficient balance",
		StatusCode: http.StatusUnprocessableEntity,
	},
}

type errorResponse struct {
	Code       errorResponseCode `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	StatusCode int               `json:"-"`
}

type transactionBalance struct {
	AccountBalance int       `json:"total"`
	CreatedAt      time.Time `json:"data_extrato"`
	AccountLimit   int       `json:"limite"`
}

type listTransactionResponse struct {
	Balance      transactionBalance `json:"saldo"`
	Transactions []Transaction      `json:"ultimas_transacoes"`
}

func parseHTTPError(functionName string, err error) errorResponse {
	if _, ok := err.(validation.Errors); ok {
		er := errorResponses["request_validation_failed"]
		er.Details = err.Error()
		return er
	}

	switch err {
	case ErrClientNotFound:
		return errorResponses["client_not_found"]
	case ErrInsufficientBalance:
		return errorResponses["insufficient_balance"]
	default:
		slog.Error(functionName, "error", err)
		return errorResponses["internal_server_error"]
	}
}

func RunServer(cfg *Config, pool *pgxpool.Pool) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/clientes/:client_id/extrato", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("client_id"))
		if err != nil {
			er := errorResponses["malformed_request"]
			c.JSON(er.StatusCode, &er)
			return
		}

		client, err := GetClient(c.Request.Context(), pool, uint(id))
		if err != nil {
			er := parseHTTPError("GetClient", err)
			c.JSON(er.StatusCode, &er)
			return
		}

		response := listTransactionResponse{
			Balance: transactionBalance{
				AccountBalance: client.AccountBalance,
				CreatedAt:      time.Now().UTC(),
				AccountLimit:   -client.AccountLimit,
			},
			Transactions: client.Transactions,
		}

		c.JSON(200, response)
	})

	router.POST("/clientes/:client_id/transacoes", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("client_id"))
		if err != nil {
			er := errorResponses["malformed_request"]
			c.JSON(er.StatusCode, &er)
			return
		}
		transaction := Transaction{}

		if err := transaction.Validate(); err != nil {
			er := parseHTTPError("transaction.Validate", err)
			c.JSON(er.StatusCode, &er)
			return
		}

		if err := c.ShouldBindJSON(&transaction); err != nil {
			slog.Error("malformed request", "error", err.Error())
			er := errorResponses["malformed_request"]
			c.JSON(er.StatusCode, &er)
			return
		}

		balance, err := AddTransaction(c.Request.Context(), pool, uint(id), transaction)
		if err != nil {
			er := parseHTTPError("AddTransaction", err)
			c.JSON(er.StatusCode, &er)
			return
		}

		c.JSON(200, balance)
	})

	if err := router.Run(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)); err != nil {
		slog.Error("Run server error", "error", err.Error())
	}
}
