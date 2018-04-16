/*
	Copyright © 2018 Harald Sitter <sitter@kde.org>

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
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apachelogger/geoip-kde-org/models"
	geoip2 "github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
)

func equalUbiquity(t *testing.T, test apiTestCase, res *httptest.ResponseRecorder) {
	assert.Contains(t, res.Header().Get("Content-Type"), "application/xml", test.tag)
	var expectedObj, actualObj models.UbiquityGeoIP
	if err := xml.Unmarshal([]byte(test.response), &expectedObj); err != nil {
		assert.FailNow(t, fmt.Sprintf("Input ('%s') needs to be valid xml.\nXML parsing error: '%s'", test.response, err.Error()), test.tag)
	}
	if err := xml.Unmarshal(res.Body.Bytes(), &actualObj); err != nil {
		assert.FailNow(t, fmt.Sprintf("Expected value ('%s') is not valid xml.\nXML parsing error: '%s'", res.Body.String(), err.Error()), test.tag)
	}
	assert.Equal(t, expectedObj, actualObj, test.tag)
}

func TestUbiquityResource(t *testing.T) {
	db, err := geoip2.Open("../GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ServeUbiquityResource(router.Group("/"), db)

	kdeDotOrg := `
<Response>
<script/>
<Ip>91.189.93.5</Ip>
<Status>OK</Status>
<CountryCode>GB</CountryCode>
<CountryCode3/>
<CountryName>United Kingdom</CountryName>
<RegionCode>ENG</RegionCode>
<RegionName>England</RegionName>
<City>London</City>
<ZipPostalCode>EC2V</ZipPostalCode>
<Latitude>51.5142</Latitude>
<Longitude>-0.0931</Longitude>
<AreaCode>0</AreaCode>
<TimeZone>Europe/London</TimeZone>
</Response>`
	runAPITests(t, []apiTestCase{
		{"t1 - get", "GET", "/v1/ubiquity", "", http.StatusOK, kdeDotOrg, equalUbiquity},
	})
}
