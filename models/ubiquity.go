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

package models

import (
	"encoding/xml"

	geoip2 "github.com/oschwald/geoip2-golang"
)

// UbiquityGeoIP is the data model for ubiquity-style output (compatible with geoip.ubuntu.com)
type UbiquityGeoIP struct {
	XMLName       xml.Name `xml:"Response"`
	IP            string   `xml:"Ip"`
	Status        string
	CountryCode   string
	CountryCode3  string
	CountryName   string
	RegionCode    string
	RegionName    string
	City          string
	ZipPostalCode string
	Latitude      float64
	Longitude     float64
	AreaCode      uint
	TimeZone      string
}

// NewUbiquityGeoIPFromGeoIP2Record creates a new ubiquity data entity from a geoip2 record
func NewUbiquityGeoIPFromGeoIP2Record(ip string, record *geoip2.City) UbiquityGeoIP {
	obj := UbiquityGeoIP{
		IP:            ip,
		Status:        "OK",
		CountryCode:   record.Country.IsoCode,
		CountryCode3:  "", // CountryCode3 is not part of
		CountryName:   record.Country.Names["en"],
		City:          record.City.Names["en"],
		ZipPostalCode: record.Postal.Code,
		Latitude:      record.Location.Latitude,
		Longitude:     record.Location.Longitude,
		AreaCode:      record.Location.MetroCode,
		TimeZone:      record.Location.TimeZone,
	}
	if len(record.Subdivisions) >= 1 {
		obj.RegionCode = record.Subdivisions[0].IsoCode
		obj.RegionName = record.Subdivisions[0].Names["en"]
	}
	return obj
}
