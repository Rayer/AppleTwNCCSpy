package AppleProductMonitor

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
	DataAccess       DataService
	FetchTarget      string
	FetchIntervalSec int
}

func NewCrawler(fetchTarget string, fetchIntervalSec int) *Crawler {
	return &Crawler{
		DataAccess:       NewSimpleFileDataAccess(),
		FetchTarget:      fetchTarget,
		FetchIntervalSec: fetchIntervalSec,
	}
}

type Product struct {
	Group   string
	Product string
	Model   string
	NCC     string
}

type Event struct {
	Added   []Product
	Removed []Product
}

func (c *Crawler) parse(source io.ReadCloser) ([]Product, error) {

	doc, err := goquery.NewDocumentFromReader(source)
	if err != nil {
		return nil, err
	}

	main := doc.Find("main")
	sections := main.Find("div#sections").Children()

	log.Debugf("Parsed %d categories...", sections.Length())
	products := make([]Product, 0)

	for i := 0; i < sections.Length(); i++ {
		current := sections.Eq(i)
		group := current.Find("h2").Text()
		items := current.Find("td")
		l := items.Length()
		var item, model, ncc string
		for j := 0; j < l; j++ {
			text := strings.Trim(items.Eq(j).Text(), " \n")
			switch j % 3 {
			case 0:
				item = text
			case 1:
				model = text
			case 2:
				ncc = text
				products = append(products, Product{
					Group:   group,
					Product: item,
					Model:   model,
					NCC:     ncc,
				})
			}
		}
	}
	log.Debugf("Parsed : %v", products)
	return products, nil
}

func (c *Crawler) fetchAndCompare(e chan<- Event) {

	resp, err := http.Get(c.FetchTarget)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Warnf("Fail to fetch from %s !", c.FetchTarget)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	recent, err := c.parse(resp.Body)

	if err != nil || len(recent) == 0 {
		return
	}

	old := c.DataAccess.LoadData()
	c.DataAccess.SaveData(recent)

	added := make([]Product, 0)
	removed := make([]Product, 0)

	for _, r := range recent {
		found := false
		for _, o := range old {
			if r.Model == o.Model {
				found = true
				break
			}
		}

		if found == false {
			added = append(added, r)
		}
	}

	for _, o := range old {
		found := false
		for _, r := range recent {
			if r.Model == o.Model {
				found = true
				break
			}
		}

		if found == false {
			removed = append(removed, o)
		}
	}

	e <- Event{
		Added:   added,
		Removed: removed,
	}
}

func (c *Crawler) Run(ctx context.Context) (e chan Event) {

	e = make(chan Event)
	go func() {
		timer := time.NewTicker(time.Duration(c.FetchIntervalSec) * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				c.fetchAndCompare(e)
			}
		}
	}()

	return e
}
