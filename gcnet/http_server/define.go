package http_server

import "github.com/gofiber/fiber/v2"

type Ctx struct {
	*fiber.Ctx
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
type ResponseHandlerFn func(*Ctx) (any, error)

type OriginHandlerFn func(*Ctx) error
