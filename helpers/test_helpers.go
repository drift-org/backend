/*
Helper functions that are used in test files.
*/
package helpers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	"net/http"
	"net/http/httptest"
)

/*
Creates a test context for testing purposes.
- query are the query parameters of the request
- body is the verbatim json body of the request

Example:
	CreateTestContext('?foo=bar&page=10', `{
		"hello": "foo"
	}`)
*/
func CreateTestContext(query string, body string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Add the request to the context. Technically, not all our requests are POST requests, but
	// we shouldn't be verifying this in our controllers, so it doesn't actually matter for testing.
	// The only thing that the context should be used for is extracting the body and/or query parameters.
	c.Request, _ = http.NewRequest("POST", "/"+query, bytes.NewBufferString(body))

	// Add the json header.
	c.Request.Header.Add("Content-Type", "application/json")
	return c
}

/*
Logs output for tests.
*/
func LogTest(message string) {
	_, _ = fmt.Fprintf(ginkgo.GinkgoWriter, message)
}
