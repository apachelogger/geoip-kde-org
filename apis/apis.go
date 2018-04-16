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
	"net"

	"github.com/gin-gonic/gin"
)

func clientIP(c *gin.Context) net.IP {
	if ip := c.Query("ip"); len(ip) > 0 {
		return net.ParseIP(ip)
	} else if ip := c.ClientIP(); len(ip) > 0 {
		return net.ParseIP(ip)
	}
	panic("Couldn't resolve client IP")
}
