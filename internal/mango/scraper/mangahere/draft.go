package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"search/internal/manga"
)

func CrawlPage(link string, mChan chan<- *manga.Manga, eChan chan<- error) {
	response, err := http.Get(link)
	if err != nil {
		eChan <- err
		return
	}

	if response.StatusCode != 200 {
		eChan <- fmt.Errorf("status code not 200 while getting link: %s", link)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		eChan <- err
		return
	}
	defer response.Body.Close()

	document.Find(mangaSelector).Each(func(i int, s *goquery.Selection) {
		go func() {
			select {
			case mChan <- parseManga(s):
			default:
			}
		}()
	})

	return
}
