package rest

import (
	"errors"
)

var (
	ErrConvID               = errors.New("converting id error")
	ErrParseQuery           = errors.New("parsing query params error")
	ErrUnmarshal            = errors.New("cannot unmarshal")
	ErrRequestBodyReading   = errors.New("cannot read request body")
	ErrUnsupportedMediaType = errors.New("only application/json is allowed")
	ErrBadRequest           = errors.New("invalid json message received")
	ErrEmptyBody            = errors.New("request must contain elements in array")
	ErrEmptyValue           = errors.New("request contains empty value")
	ErrMissingQuery         = errors.New("missing query params")
	ErrMissingParam         = errors.New("missing url param")
	ErrEmptyID              = errors.New("id parameter is required")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidAuthHeader    = errors.New("invalid authorization header format")
	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrForbidden            = errors.New("forbidden")
	ErrFileTooLarge         = errors.New("file size exceeds the allowed limit")
	ErrFileRequired         = errors.New("file is required")
	ErrParamRequired        = errors.New("parameter is required")
)
