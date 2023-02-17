package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

type SessionMiddleware struct {
	Abstract
	session *session.Session
}

func NewSessionMiddleware(appSession *session.Session) *SessionMiddleware {
	appSession.Init()
	return &SessionMiddleware{session: appSession}
}

func (sessionMiddleware SessionMiddleware) Process(ctx *gin.Context) {
	err := sessionMiddleware.session.Start(ctx)
	if err != nil {
		sessionMiddleware.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.Next()
}
