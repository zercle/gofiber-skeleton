package utils

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/valyala/fastjson"
)

var (
	// Session storage
	SessStore *session.Store
	// parser pool
	JsonParserPool *fastjson.ParserPool
)
