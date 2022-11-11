// Package traefikgeoip2 a geoip2 traefik plugin.
package traefikgeoip2

import (
	"context"
	"log"
	"net"
	"net/http"
	"text/template"

	"github.com/oschwald/geoip2-golang"
)

// Config the plugin configuration.
type Config struct {
	FromHeader    string `json:"fromheader,omitempty"`
	GeoDBLocation string `json:"geodblocation,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		FromHeader:    "X-Forwarded-For",
		GeoDBLocation: "./GeoLite2-City.mmdb",
	}
}

// TraefikGeoip2 a geoip traefik plugin.
type TraefikGeoip2 struct {
	next          http.Handler
	fromheader    string
	geodblocation string
	name          string
	template      *template.Template
}

// New created a new TraefikGeoip2 plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TraefikGeoip2{
		fromheader:    config.FromHeader,
		geodblocation: config.GeoDBLocation,
		next:          next,
		name:          name,
		template:      template.New("traefikgeoip2").Delims("[[", "]]"),
	}, nil
}

func (a *TraefikGeoip2) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	db, err := geoip2.Open(a.geodblocation)
	if err != nil {
		log.Fatal(err)
	}
	reqheader := req.Header.Get(a.fromheader)
	ip := net.ParseIP(reqheader)
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Geoip2-Country", record.Country.IsoCode)

	a.next.ServeHTTP(rw, req)
}
