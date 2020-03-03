package sitemap

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type ILoc interface {
	Get() string
}

type Container []interface{}

func (x *Container) GetUrls() []string{
	var urls []string
	for _, v := range (*x) {
		urls = append(urls, v.(ILoc).Get())
	}
	return urls
}


type SitemapIndex struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc string `xml:"loc"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Urls []Url `xml:"url"`
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc string `xml:"loc"`
}

func ParseSitemap (sitemapUrl string) [] string {
	var urls []string

	sitemap := fetchSitemap(sitemapUrl)
	if sitemap == nil {
		return urls
	}

	// possibly the sitemap is of type index, pointing at
	// other sitemaps, or maybe it's just a url set (most common)
	newSitemapUrls := ParseSitemapIndex(sitemap)
	if newSitemapUrls != nil {
		for _, url := range newSitemapUrls {
			subUrls := ParseSitemap(url)
			urls = append(urls, subUrls...)
		}
	}

	// Assume the sitemap is just a url set, which seems quite common
	newUrls := ParseUrlSet(sitemap)
	if newUrls != nil {
		urls = append(urls, newUrls...)
		return urls
	}
	return urls
}

func fetchSitemap(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	sitemap, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return sitemap
}

func ParseSitemapIndex(sitemap []byte) []string {
	var sitemapIndex SitemapIndex

	err := xml.Unmarshal(sitemap, &sitemapIndex) 
	if err != nil {
		return nil
	}

	var urls []string
	for _, v := range sitemapIndex.Sitemaps {
		urls = append(urls, v.Loc)
	}
	return urls
}

func ParseUrlSet(sitemap []byte) []string {
	var urlSet UrlSet

	err := xml.Unmarshal(sitemap, &urlSet)
	if err != nil {
		// TODO: Might just be a text file
		// TODO: Might be atom
		// TODO: Might be RSS feed
		return nil
	}

	var urls []string
	for _, v := range urlSet.Urls {
		urls = append(urls, v.Loc)
	}
	return urls
}