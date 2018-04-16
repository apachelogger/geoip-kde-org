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
type calamaresResource struct {
	db *geoip2.Reader
}

// ServeCalamaresResource sets up the calamares resource routes.
func ServeCalamaresResource(rg *gin.RouterGroup, db *geoip2.Reader) {
	r := &calamaresResource{db}
	rg.GET("/v1/calamares", r.get)
}

/**
 * @api {get} /calamares Calamares
 *
 * @apiVersion 1.0.0
 * @apiGroup GeoIP
 * @apiName calamares
 *
 * @apiDescription Calamares-style JSON geoip data. This endpont offers the
 *   JSON format defined by Calamares' locale module.
 *
 * @apiSuccessExample {json} Success-Response:
 *   {"time_zone":"Europe/Vienna"}
 */
func (r *calamaresResource) get(c *gin.Context) {
	// If you are using strings that may be invalid, check that ip is not nil
	record, err := r.db.City(clientIP(c))
	if err != nil {
		panic(err)
	}

	data := models.CalamaresGeoIP{TimeZone: record.Location.TimeZone}
	c.JSON(http.StatusOK, data)
}
