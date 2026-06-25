package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/command"
)

// Middleware to check database connection status before running queries
func dbCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CORS preflight requests carry no credentials and must pass through
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		path := strings.Replace(c.Request.URL.Path, command.Opts.Prefix, "", -1)

		// Allow whitelisted paths
		if allowedPaths[path] {
			c.Next()
			return
		}

		// Check if session exists in single-session mode
		if !command.Opts.Sessions {
			if DbClient == nil {
				badRequest(c, errNotConnected)
				return
			}

			c.Next()
			return
		}

		// Determine session ID from the client request
		sid := getSessionId(c.Request)
		if sid == "" {
			badRequest(c, errSessionRequired)
			return
		}

		// Determine the database connection handle for the session
		conn := DbSessions.Get(sid)
		if conn == nil {
			badRequest(c, errNotConnected)
			return
		}

		c.Next()
	}
}

// Middleware to inject CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", command.Opts.CorsOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Session-Id, X-Token")
		c.Header("Access-Control-Expose-Headers", "*")
		c.Header("Access-Control-Max-Age", "3600")
	}
}

// tokenAuthMiddleware lets API clients authenticate with an access token
// instead of going through the browser session flow. When a token is provided
// (Authorization: Bearer <token>, X-Token header, or ?token=), the token is
// reused as the session id and a database connection is established for it on
// first use, so subsequent requests (e.g. /api/query) just work.
func tokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		url, err := TokenStore.Get(token)
		if err != nil {
			badRequest(c, err)
			c.Abort()
			return
		}

		// Reuse the token as the session id so repeated calls share a connection
		c.Request.Header.Set("x-session-id", token)

		if DbSessions.Get(token) == nil {
			cl, err := client.NewFromUrl(url, nil)
			if err != nil {
				badRequest(c, err)
				c.Abort()
				return
			}

			if _, err := cl.Info(); err != nil {
				cl.Close()
				badRequest(c, err)
				c.Abort()
				return
			}

			DbSessions.Add(token, cl)
		}

		c.Next()
	}
}

// extractToken pulls an access token from the request, preferring the
// Authorization bearer header, then the X-Token header, then the query string.
func extractToken(c *gin.Context) string {
	if auth := c.GetHeader("Authorization"); len(auth) > 7 && strings.EqualFold(auth[:7], "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	if tok := c.GetHeader("X-Token"); tok != "" {
		return tok
	}
	return c.Query("token")
}

func requireLocalQueries() gin.HandlerFunc {
	return func(c *gin.Context) {
		if QueryStore == nil {
			badRequest(c, "local queries are disabled")
			return
		}

		c.Next()
	}
}
