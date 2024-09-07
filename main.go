package main

import (
	"flag"
	"julien/julien"
	"strconv"

	"julien/web"
)

func start(j *julien.Julien, site *julien.Site, endpoint string) {
	jweb := web.New(j, site)
	jweb.Start(endpoint)
}

func main() {
	sitepath := flag.String("site", "index.md", "site markdown file")
	configpath := flag.String("config", "julien.yaml", "julien config file")
	port := flag.Int("port", 1234, "webserver port")
	host := flag.String("host", "localhost", "webserver host")
	flag.Parse()
	j := julien.DefaultJulien()
	site := julien.DefaultSite()
	julien.LoadSite(*sitepath, &site)
	julien.LoadConfig(*configpath, &j)
	endpoint := (*host) + ":" + strconv.Itoa(*port)
	start(&j, &site, endpoint)
}
