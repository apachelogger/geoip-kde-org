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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	geoip2 "github.com/oschwald/geoip2-golang"
)

type debugResource struct {
	db *geoip2.Reader
}

// ServeDebugResource sets up the semi-internal data inspection resource.
// Its format is entirely undefined and absolutely not meant to for consumption.
func ServeDebugResource(rg *gin.RouterGroup, db *geoip2.Reader) {
	r := &debugResource{db}
	rg.GET("/debug", r.get)
}

func (r *debugResource) get(c *gin.Context) {
	// If you are using strings that may be invalid, check that ip is not nil
	record, err := r.db.City(clientIP(c))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", record)

	c.JSON(http.StatusOK, record)
}
