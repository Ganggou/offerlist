package models

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

func FetchPrice(platformId, shortId string) string {
	switch platformId {
	case "UA":
		return amazonA(shortId)
	case "TAO":
		return Taobao(shortId)
	default:
	}
	return "err"
}

func amazonA(shortId string) string {
	c := colly.NewCollector()
	var wg sync.WaitGroup
	var price float64

	wg.Add(1)

	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"

	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Println(err)
		}
		nodes := htmlquery.Find(doc, `//div[@class="a-row a-spacing-mini olpOffer"]`)
		for i := 0; i < len(nodes); i++ {
			priceNode := htmlquery.Find(nodes[i], `//span[@class="a-size-large a-color-price olpOfferPrice a-text-bold"]`)
			priceString := strings.Trim(htmlquery.InnerText(priceNode[0]), " ")
			price, err = strconv.ParseFloat(priceString[1:], 10)
			if err != nil {
				log.Println(err)
				continue
			}
			break
		}
		wg.Done()
	})

	c.Visit("https://www.amazon.com/gp/offer-listing/" + shortId)
	wg.Wait()
	return strconv.Itoa(int(price * 100))
}

func Taobao(shortId string) string {
	c := colly.NewCollector()
	var wg sync.WaitGroup
	var price float64

	wg.Add(1)

	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"

	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Println(err)
		}
		priceNode := htmlquery.Find(doc, `//*[@class="tb-rmb-num"]`)
		priceString := strings.Trim(strings.Split(htmlquery.InnerText(priceNode[0]), "-")[0], " ")
		price, err = strconv.ParseFloat(priceString, 10)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	})

	c.Visit("https://item.taobao.com/item.htm?id=" + shortId)
	wg.Wait()
	return strconv.Itoa(int(price * 100))
}
