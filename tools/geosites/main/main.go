package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/ext/tools/geosites"
)

func main() {
	geoSiteList := new(router.GeoSiteList)
	geoSiteList.Entry = append(geoSiteList.Entry, &router.GeoSite{
		CountryCode: "CN",
		Domain:      geosites.GetGeoSiteCN(),
	})
	geoSiteList.Entry = append(geoSiteList.Entry, &router.GeoSite{
		CountryCode: "SPEEDTEST",
		Domain:      geosites.GetGeoSiteSpeedTest(),
	})

	geoSiteListBytes, err := proto.Marshal(geoSiteList)
	if err != nil {
		fmt.Println("failed to marshal geosites:", err)
		return
	}
	if err := ioutil.WriteFile("geosite.dat", geoSiteListBytes, 0777); err != nil {
		fmt.Println("failed to write geosite.dat.", err)
	}
}
