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
	"net/http"

	"github.com/apachelogger/geoip-kde-org/models"
	"github.com/gin-gonic/gin"
	geoip2 "github.com/oschwald/geoip2-golang"
)

// We are muddying the waters a bit by merging api+service+data.
type ubiquityResource struct {
	db *geoip2.Reader
}

// ServeUbiquityResource sets up the ubiquity resource routes.
func ServeUbiquityResource(rg *gin.RouterGroup, db *geoip2.Reader) {
	r := &ubiquityResource{db}
	rg.GET("/v1/ubiquity", r.get)
}

/**
 * @api {get} /ubiquity Ubiquity
 *
 * @apiVersion 1.0.0
 * @apiGroup GeoIP
 * @apiName ubiquity
 *
 * @apiDescription Ubuiqity-style XML geoip data. This is equivalent to calling
 *    geoip.ubuntu.com/lookup which is where the actual data format comes from.
 *
 * @apiSuccessExample {xml} Success-Response:
 *   <Response>
 *   <script/>
 *   <Ip>193.81.57.56</Ip>
 *   <Status>OK</Status>
 *   <CountryCode>AT</CountryCode>
 *   <CountryCode3/>
 *   <CountryName>Austria</CountryName>
 *   <RegionCode>4</RegionCode>
 *   <RegionName>Upper Austria</RegionName>
 *   <City>Gmunden</City>
 *   <ZipPostalCode>4810</ZipPostalCode>
 *   <Latitude>47.9022</Latitude>
 *   <Longitude>13.7642</Longitude>
 *   <AreaCode>0</AreaCode>
 *   <TimeZone>Europe/Vienna</TimeZone>
 *   </Response>
 */
func (r *ubiquityResource) get(c *gin.Context) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip := clientIP(c)
	record, err := r.db.City(ip)
	if err != nil {
		panic(err)
	}

	data := models.NewUbiquityGeoIPFromGeoIP2Record(ip.String(), record)
	c.XML(http.StatusOK, data)
}
