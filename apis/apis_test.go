/*
	Copyright © 2017-2018 Harald Sitter <sitter@kde.org>

	This program is free software; you can redistribute it and/or
	modify it under the terms of the GNU General Public License as
	published by the Free Software Foundation; either version 3 of
	the License or any later version accepted by the membership of
	KDE e.V. (or its successor approved by the membership of KDE
	e.V.), which shall act as a proxy defined in Section 14 of
	version 3 of the license.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package apis

import (
	"bytes"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func init() {
	router = gin.Default()
}

// XML data can only get unmarshalled against a concrete type and not a generic
// map interface. So, to assert XML data we need custom handlers.
// This is passed into a test case and then called as part of the assertions.
type equalResponseHandler func(t *testing.T, test apiTestCase, res *httptest.ResponseRecorder)

type apiTestCase struct {
	tag      string
	method   string
	url      string
	body     string
	status   int
	response string
	handler  equalResponseHandler
}

func equalJSON(t *testing.T, test apiTestCase, res *httptest.ResponseRecorder) {
	assert.Contains(t, res.Header().Get("Content-Type"), "application/json", test.tag)
	assert.JSONEq(t, test.response, res.Body.String(), test.tag)
}

func testAPI(method, URL, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, URL, bytes.NewBufferString(body))
	// NB: we work against the live data as we have insufficient data/api
	// separation, so this may eventually fail to meet expectation.
	// Unfortunately the live data doesn't handle 192.0.2.0/24, which would be
	// an RFC defined address solely for testing.
	req.RemoteAddr = "91.189.93.5:1234"
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, tests []apiTestCase) {
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			res := testAPI(test.method, test.url, test.body)
			assert.Equal(t, test.status, res.Code, test.tag)
			if test.response != "" {
				test.handler(t, test, res)
			}
		})
	}
}

func TestApisClientIPFromRemote(t *testing.T) {
	req := httptest.NewRequest("GET", "/foo", bytes.NewBufferString(""))
	req.RemoteAddr = "91.189.93.5:1234"
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	assert.Equal(t, net.ParseIP("91.189.93.5"), clientIP(c))
}

func TestApisClientIPFromQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/foo?ip=8.8.8.8", bytes.NewBufferString(""))
	req.RemoteAddr = "91.189.93.5:1234"
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	assert.Equal(t, net.ParseIP("8.8.8.8"), clientIP(c))
}

func TestApisClientIPv6(t *testing.T) {
	// FTR geolite2 features ipv6 support, so make sure we correctly parse ipv6
	req := httptest.NewRequest("GET", "/foo", bytes.NewBufferString(""))
	req.RemoteAddr = "[2001:1af8:4100:a08c:22::10]:1234"
	c, e := gin.CreateTestContext(httptest.NewRecorder())
	e.ForwardedByClientIP = true
	c.Request = req
	assert.Equal(t, net.ParseIP("2001:1af8:4100:a08c:22::10"), clientIP(c))
}
