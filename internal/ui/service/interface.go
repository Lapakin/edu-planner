package service

import (
	"net/http/httputil"

	"github.com/Lapakin/edu-planner/internal/logging"
)

type Services struct {
	AuthProxy     *httputil.ReverseProxy
	SyllabusProxy *httputil.ReverseProxy
	Logger        *logging.Logger
}
