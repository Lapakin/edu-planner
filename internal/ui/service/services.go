package service

import (
	"net/http/httputil"
	"net/url"

	"github.com/Lapakin/edu-planner/internal/logging"
)

func NewServices(l *logging.Logger, userManagementURL string, syllabusURL string) *Services {
	authTarget, _ := url.Parse(userManagementURL)
	syllabusTarget, _ := url.Parse(syllabusURL)
	return &Services{
		AuthProxy:     httputil.NewSingleHostReverseProxy(authTarget),
		SyllabusProxy: httputil.NewSingleHostReverseProxy(syllabusTarget),
		Logger:        l,
	}
}
