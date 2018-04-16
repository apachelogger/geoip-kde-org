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

package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apachelogger/geoip-kde-org/apis"
	"github.com/coreos/go-systemd/activation"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

var db *geoip2.Reader

func downloadGeoLite2City() {
	// FIXME: we should symlink the current version to the fixed name, but
	//   store the actual files with a timestamp embedded. that way re-opening
	//   the fixed path loads always the latest file
	path := "GeoLite2-City.mmdb"
	name := path

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Writer the body to file
	// Reading from Body.Resp via Gzip and Bufio is substantially slower
	// than first downloading the entire body and reading from local. I am
	// not entirely sure why that is since bufio should make it fast :(
	tmpfile, err := ioutil.TempFile("", "geoip-geolite2-city")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer tmpfile.Close()
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		panic(err)
	}
	tmpfile.Seek(0, 0)

	gzip, err := gzip.NewReader(bufio.NewReader(tmpfile))
	if err != nil {
		panic(err)
	}
	defer gzip.Close()

	tarReader := tar.NewReader(gzip)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		info := header.FileInfo()
		if info.Name() != name {
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			panic(err)
		}

		break
	}
}

func downloadGeoLite2() bool {
	download := true
	if stat, err := os.Stat("GeoLite2-City.mmdb"); err == nil {
		if time.Since(stat.ModTime()).Hours() < 24*8 {
			download = false
		}
	}

	if download {
		downloadGeoLite2City()
	}

	return download
}

func main() {
	flag.Parse()

	downloadGeoLite2()

	var err error
	db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	log.Println("Ready to rumble...")
	router := gin.Default()

	rg := router.Group("/")
	{
		apis.ServeCalamaresResource(rg, db)
		apis.ServeUbiquityResource(rg, db)
		apis.ServeDebugResource(rg, db)
	}
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/doc")
	})
	router.StaticFS("/doc", http.Dir("doc"))
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	listeners, err := activation.Listeners(true)
	if err != nil {
		panic(err)
	}

	log.Println("starting servers")
	var servers []*http.Server
	for _, listener := range listeners {
		server := &http.Server{Handler: router}
		go server.Serve(listener)
		servers = append(servers, server)
	}

	if len(servers) == 0 {
		log.Println("servers empty. adding manual server")

		host := os.Getenv("HOST")
		port := os.Getenv("PORT")
		if len(port) <= 0 {
			port = "8080"
		}

		server := &http.Server{
			Addr:    host + ":" + port,
			Handler: router,
		}
		go server.ListenAndServe()
		servers = append(servers, server)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	quitTicker := time.NewTicker(14 * 24 * time.Hour)
	go func() {
		<-quitTicker.C
		quitTicker.Stop()
		close(quit)
		log.Println("Auto-terminating after 14 days to force a db update.")
	}()

	// Wait for some quit cause.
	// This could be INT, TERM, QUIT or the db update trigger.
	// We'll then do a zero downtime shutdown.
	// This relies on systemd managing the socket and us doing graceful listener
	// shutdown. Once we are no longer listening, the system starts backlogging
	// the socket until we get restarted and listen again.
	// Ideally this results in zero dropped connections.
	<-quit
	log.Println("servers are shutting down")
	quitTicker.Stop()

	for _, srv := range servers {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server Shutdown: %s", err)
		}
	}

	log.Println("Server exiting")
}
