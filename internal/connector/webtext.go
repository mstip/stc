package connector

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	ParentType string
	Text       string
	Path       string
	URL        string
}

type RawScrapeDataRow struct {
	ID   int64  `db:"id"`
	URL  string `db:"url"`
	Text string `db:"text"`
	Path string `db:"path"`
	Type string `db:"type"`
}

func scanNode(data *[]Info, url string, baseUrl string, scrapedUrls *[]string) func(i int, s *goquery.Selection) {
	return func(i int, s *goquery.Selection) {
		ignore := []string{"script", "svg", "link", "style"}
		if slices.Contains(ignore, goquery.NodeName(s)) {
			return
		}

		if goquery.NodeName(s) == "a" {
			href, ok := s.Attr("href")
			if ok {
				if strings.HasPrefix(href, baseUrl) {
					scrapeUrl(data, href, baseUrl, scrapedUrls)
				} else if strings.HasPrefix(href, "/") {
					if strings.HasSuffix(baseUrl, "/") {
						href = href[1:]
					}
					scrapeUrl(data, baseUrl+href, baseUrl, scrapedUrls)
				} else if !strings.HasPrefix(href, "http") {
					scrapeUrl(data, baseUrl+href, baseUrl, scrapedUrls)
				}
			}
		}

		if goquery.NodeName(s) == "#text" {
			t := strings.TrimSpace(s.Text())
			if len(t) > 1 {

				path := []string{goquery.NodeName(s.Parent())}
				s.Parent().ParentsUntil("body").Each(func(i int, s *goquery.Selection) {
					path = append(path, goquery.NodeName(s))
				})
				slices.Reverse(path)

				*data = append(*data, Info{
					ParentType: goquery.NodeName(s.Parent()),
					Text:       t,
					Path:       strings.Join(path, ">"),
					URL:        url,
				})
			}
		}

		if s.Contents().Length() == 0 {
			return
		}
		s.Contents().Each(scanNode(data, url, baseUrl, scrapedUrls))
	}
}

func scrapeUrl(data *[]Info, url string, baseUrl string, scrapedUrls *[]string) error {
	// skip already scraped urls
	if slices.Contains(*scrapedUrls, url) {
		return nil
	}

	// ignore common endings
	ignoreEndings := []string{".png", ".jpg", ".pdf", ".jpeg", ".svg", "."}
	for _, e := range ignoreEndings {
		if strings.HasSuffix(url, e) {
			return nil
		}
	}
	log.Printf("scrapeUrl %s", url)
	*scrapedUrls = append(*scrapedUrls, url)
	res, err := http.Get(url)
	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		log.Printf("contenttype %s is not text/html skip", contentType)
		return nil
	}
	if err != nil {
		log.Print(err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %s %d %s", url, res.StatusCode, res.Status)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
		return err
	}

	doc.Find("body").First().Contents().Each(scanNode(data, url, baseUrl, scrapedUrls))
	return nil
}

func scrapePage(data *[]Info, baseUrl string) error {
	scrapedUrls := []string{}
	return scrapeUrl(data, baseUrl, baseUrl, &scrapedUrls)
}

func WebText(url string) ([]Info, error) {
	infos := []Info{}
	err := scrapePage(&infos, url)
	if err != nil {
		return nil, err
	}
	return infos, err
}
