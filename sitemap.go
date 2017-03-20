package sitemap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/klauspost/compress/gzip"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
	xmlns:xhtml="http://www.w3.org/1999/xhtml">`
	footer   = ` </urlset>`
	template = `
	 <url>
	   <loc>%s</loc>
       <xhtml:link rel="alternate" media="only screen and (max-width: 640px)" href="%s" />
	   <lastmod>%s</lastmod>
	   <changefreq>%s</changefreq>
	   <priority>%.1f</priority>
	 </url>`

	indexHeader = `<?xml version="1.0" encoding="UTF-8"?>
      <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
	indexFooter = `
</sitemapindex>
	`
	indexTemplate = `
    <sitemap>
       <loc>%s%s</loc>
       <lastmod>%s</lastmod>
    </sitemap>
	`
)

type Item struct {
	Loc        string
	LocMobile  string
	LastMod    time.Time
	Changefreq string
	Priority   float32
}

func (item *Item) String() string {
	//2012-08-30T01:23:57+08:00
	//Mon Jan 2 15:04:05 -0700 MST 2006
	return fmt.Sprintf(template, item.Loc, item.LocMobile, item.LastMod.Format("2006-01-02T15:04:05+08:00"), item.Changefreq, item.Priority)
}

func SiteMap(f string, items []*Item) error {
	var buffer bytes.Buffer
	buffer.WriteString(header)
	defer buffer.Reset()
	for _, item := range (items) {
		_, err := buffer.WriteString(item.String())
		if err != nil {
			return err
		}
	}
	fo, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fo.Close()
	defer func() {
		buffer.Reset()
	}()
	buffer.WriteString(footer)

	zip := gzip.NewWriter(fo)
	defer zip.Close()
	_, err = zip.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	return err
}

func SiteMapIndex(folder, indexFile, baseurl string) error {
	var buffer bytes.Buffer
	defer buffer.Reset()
	buffer.WriteString(indexHeader)
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") {
			fmt.Println(f.Name())
			s := fmt.Sprintf(indexTemplate, baseurl, f.Name(), time.Now().Format("2006-01-02T15:04:05+08:00"))
			buffer.WriteString(s)
		}
	}
	buffer.WriteString(indexFooter)
	err = ioutil.WriteFile(indexFile, buffer.Bytes(), 0755)
	return err
}
