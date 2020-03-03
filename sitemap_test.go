package sitemap


import (
	"testing"
)


func TestParseSitemapBivar(t *testing.T) {
	urls := ParseSitemap("https://bivar.com/sitemap.xml")
	if len(urls) == 0 {
		t.Errorf("urls should be greater than 0 but was %v", len(urls))
	}
}

func TestParseSitemapBpDotCom1(t *testing.T) {
	urls := ParseSitemap("https://bikepacking.com/sitemap_index.xml")
	if len(urls) == 0 {
		t.Errorf("urls should be greater than 0 but was %v", len(urls))
	}
}

func TestParseBadSitemap(t *testing.T) {
	urls := ParseSitemap("https://bikepacking.com/postsafdsfds-sitemap1.xml")
	if len(urls) > 0 {
		t.Errorf("urls should be 0 but was %v", len(urls))
	}
} 