package ipattr

import (
	"encoding/json"
	"io"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPResolver struct {
	DB               *geoip2.Reader
	SubdivisionsData map[string]interface{}
}

func NewGeoIPResolver(dbFilePath, subdivisionsPath string) (*GeoIPResolver, error) {
	db, err := geoip2.Open(dbFilePath)
	if err != nil {
		return nil, err
	}

	subdivisions, err := jsonToMap(err, subdivisionsPath)
	if err != nil {
		return nil, err
	}
	return &GeoIPResolver{
		DB:               db,
		SubdivisionsData: subdivisions,
	}, nil
}

func jsonToMap(err error, subdivisionsPath string) (map[string]interface{}, error) {
	file, err := os.Open(subdivisionsPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (r *GeoIPResolver) ResolveIP(ipString string) (string, error) {
	ip := net.ParseIP(ipString)
	record, err := r.DB.City(ip)
	if err != nil {
		return "", err
	}
	if zh, exists := record.Subdivisions[0].Names["zh-CN"]; exists {
		return zh, nil
	}
	if zh, exists := r.SubdivisionsData[record.Subdivisions[0].Names["en"]]; exists {
		return zh.(string), nil
	}
	if en, exists := record.Subdivisions[0].Names["en"]; exists {
		return en, nil
	}
	return "", nil
}
