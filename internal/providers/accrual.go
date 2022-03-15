package providers

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	Accrual "go-musthave-diploma-tpl/gen/accrual"
	"go.uber.org/ratelimit"
	"net/http"
	"strconv"
	"time"
)

var ErrWrongStatusCode = errors.New("wrong http status code")
var ErrUnknownAccrualStatus = errors.New("unknown accrual status")
var ErrFatalError = errors.New("fatal error")
var ErrTooManyRequests = errors.New("too many requests")

type TooManyRequestsError struct {
	err           error
	RetryAfterSec int
}

func (e TooManyRequestsError) Error() string {
	return e.err.Error()
}

func (e TooManyRequestsError) Unwrap() error {
	return e.err
}

func NewTooManyRequestsError(retryAfterSec int) *TooManyRequestsError {
	return &TooManyRequestsError{
		RetryAfterSec: retryAfterSec,
		err:           ErrTooManyRequests,
	}
}

type Provider interface {
	GetAccrual(ctx context.Context, orderID string) chan *GetAccrualResponse
}

type GetAccrualResponse struct {
	Status  string
	Accrual *float32
	Err     error
}

type Client struct {
	accrualAPIClient Accrual.ClientWithResponsesInterface
	rateLimiter      ratelimit.Limiter
}

func (c Client) GetAccrual(ctx context.Context, orderID string) chan *GetAccrualResponse {
	resultCh := make(chan *GetAccrualResponse)
	go func() {
		c.rateLimiter.Take()
		resp := c.getAccrual(ctx, orderID)
		resultCh <- resp
		close(resultCh)
	}()
	return resultCh
}

func (c Client) getAccrual(ctx context.Context, orderID string) *GetAccrualResponse {
	resp, err := c.accrualAPIClient.GetOrderAccrualWithResponse(ctx, orderID)
	if err != nil {
		return &GetAccrualResponse{Err: err}
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		retryAfter := resp.HTTPResponse.Header.Get("Retry-After")
		retryAfterSec, err := strconv.Atoi(retryAfter)
		if err != nil {
			retryAfterSec = 5
		}
		return &GetAccrualResponse{Err: NewTooManyRequestsError(retryAfterSec)}
	}

	if resp.StatusCode() != 200 {
		return &GetAccrualResponse{Err: ErrWrongStatusCode}
	}

	log.Info().
		Str("orderID", orderID).
		Str("accrualStatus", string(resp.JSON200.Status)).
		Float32("accrual", *resp.JSON200.Accrual).
		Msg("get accrual result")

	return &GetAccrualResponse{Accrual: resp.JSON200.Accrual, Status: resp.JSON200.Status}
}

func NewProvider(accrualAPIClient Accrual.ClientWithResponsesInterface) *Client {
	rl := ratelimit.New(1000, ratelimit.Per(1*time.Minute)) // per second
	return &Client{accrualAPIClient: accrualAPIClient, rateLimiter: rl}
}
