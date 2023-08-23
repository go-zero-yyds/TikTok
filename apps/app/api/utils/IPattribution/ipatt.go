package ipattr

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPResolver struct {
	DB       *geoip2.Reader
	JsonData map[string]interface{}
}

func NewGeoIPResolver(dbFilePath, jsonFilePath string) (*GeoIPResolver, error) {
	db, err := geoip2.Open(dbFilePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return &GeoIPResolver{
		DB:       db,
		JsonData: jsonData,
	}, nil
}

func (r *GeoIPResolver) Close() {
	r.DB.Close()
}

func (r *GeoIPResolver) ResolveIP(ipString string) (string, error) {
	ip := net.ParseIP(ipString)
	record, err := r.DB.City(ip)
	if err != nil {
		return "", err
	}

	if zh, exists := r.JsonData[record.City.Names["en"]]; exists {
		return fmt.Sprint(zh), nil
	}
	return record.City.Names["en"], nil
}
