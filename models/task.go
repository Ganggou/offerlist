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
	case "CA":
		return amazonU(shortId, "www.amazon.cn")
	case "JA":
		return amazonU(shortId, "www.amazon.co.jp")
	case "UA":
		return amazonU(shortId, "www.amazon.com")
	case "AA":
		return amazonU(shortId, "www.amazon.com.au")
	case "TAO":
		return taobao(shortId)
	default:
	}
	return "err"
}

func amazonU(shortId, link string) string {
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
			priceLine := strings.Trim(htmlquery.InnerText(priceNode[0]), " ")
			parts := strings.Split(priceLine, " ")
			var priceString string
			if len(parts) == 0 {
				continue
			} else if len(parts) == 1 {
				priceString = priceLine[1:]
			} else {
				priceString = parts[len(parts)-1]
			}
			s := strings.ReplaceAll(priceString, ",", "")
			price, err = strconv.ParseFloat(s, 10)
			if err != nil {
				log.Println(err)
				continue
			}
			break
		}
		wg.Done()
	})

	c.Visit("https://" + link + "/gp/offer-listing/" + shortId)
	wg.Wait()
	result := int(price * 100)
	if result == 0 {
		result = -1
	}
	return strconv.Itoa(result)
}

func taobao(shortId string) string {
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
