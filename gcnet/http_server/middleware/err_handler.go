package middleware

import (
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/http_server"
)

func NewErrHandler() http_server.OriginHandlerFn {
	return func(c *http_server.Ctx) error {
		err := c.Next()
		if err != nil {
			if errCode, ok := err.(errcode.ErrCode); ok {
				http_server.SendErrCode(c.Ctx, errCode)
			} else {
				http_server.SendError(c.Ctx, err)
			}
		}
		return err
	}
}