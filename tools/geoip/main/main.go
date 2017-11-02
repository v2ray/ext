// GeoIP generator
//
// Before running this file, the GeoIP database must be downloaded and present.
// To download GeoIP database: https://dev.maxmind.com/geoip/geoip2/geolite2/
// Inside you will find block files for IPv4 and IPv6 and country code mapping.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/ext/tools/conf"
)

var (
	countryCodeFile = flag.String("country", "", "Path to the country code file")
	ipv4File        = flag.String("ipv4", "", "Path to the IPv4 block file")
	ipv6File        = flag.String("ipv6", "", "Path to the IPv6 block file")
)

func getCountryCodeMap() (map[string]string, error) {
	countryCodeReader, err := os.Open(*countryCodeFile)
	if err != nil {
		return nil, err
	}
	defer countryCodeReader.Close()

	m := make(map[string]string)
	reader := csv.NewReader(countryCodeReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, line := range lines[1:] {
		id := line[0]
		countryCode := line[4]
		if len(countryCode) == 0 {
			continue
		}
		m[id] = strings.ToUpper(countryCode)
	}
	return m, nil
}

func getCidrPerCountry(file string, m map[string]string, list map[string][]*router.CIDR) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	reader := csv.NewReader(fileReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines[1:] {
		cidrStr := line[0]
		countryId := line[1]
		if countryCode, found := m[countryId]; found {
			cidr, err := conf.ParseIP(cidrStr)
			if err != nil {
				return err
			}
			cidrs := append(list[countryCode], cidr)
			list[countryCode] = cidrs
		}
	}
	return nil
}

func main() {
	flag.Parse()

	ccMap, err := getCountryCodeMap()
	if err != nil {
		fmt.Println("Error reading country code map:", err)
		return
	}

	cidrList := make(map[string][]*router.CIDR)
	if err := getCidrPerCountry(*ipv4File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv4 file:", err)
		return
	}
	if err := getCidrPerCountry(*ipv6File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv6 file:", err)
		return
	}

	geoIPList := new(router.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}

	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		fmt.Println("Error marshalling geoip list:", err)
	}

	if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0777); err != nil {
		fmt.Println("Error writing geoip to file:", err)
	}
}
