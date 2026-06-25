package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tuvistavie/securerandom"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/connection"
)

// CreateToken validates a database connection string, stores it and returns an
// opaque token that can later be exchanged for an authenticated session.
//
//	POST /api/token   url=postgres://user:password@host:port/db?sslmode=mode
//	-> { "token": "...", "url": "/token/..." }
func CreateToken(c *gin.Context) {
	url := c.Request.FormValue("url")
	if url == "" {
		badRequest(c, errURLRequired)
		return
	}

	// Validate and normalize the connection string before storing it.
	url, err := connection.FormatURL(command.Options{
		URL:      url,
		Passfile: command.Opts.Passfile,
	})
	if err != nil {
		badRequest(c, err)
		return
	}

	token, err := TokenStore.Create(url)
	if err != nil {
		badRequest(c, err)
		return
	}

	successResponse(c, gin.H{
		"token": token,
		// Ready-to-use link for the common (no-prefix) deployment. Behind a
		// custom prefix, prepend it to this path.
		"url": fmt.Sprintf("/token/%s", token),
	})
}

// ConnectWithToken exchanges a previously issued token for an authenticated
// session and redirects to the app with the session id, mirroring the
// backend-credential flow used by ConnectWithBackend.
func ConnectWithToken(c *gin.Context) {
	url, err := TokenStore.Get(c.Param("token"))
	if err != nil {
		badRequest(c, err)
		return
	}

	// Make a new session
	sid, err := securerandom.Uuid()
	if err != nil {
		badRequest(c, err)
		return
	}
	c.Request.Header.Add("x-session-id", sid)

	// Connect to the database
	cl, err := client.NewFromUrl(url, nil)
	if err != nil {
		badRequest(c, err)
		return
	}

	// Finalize session setup
	if _, err = cl.Info(); err == nil {
		err = setClient(c, cl)
	}
	if err != nil {
		cl.Close()
		badRequest(c, err)
		return
	}

	c.Redirect(302, fmt.Sprintf("/%s?session=%s", command.Opts.Prefix, sid))
}
