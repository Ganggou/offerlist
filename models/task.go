package models

import (
  "fmt"
  "strconv"
  "strings"
	"time"

  "github.com/gocolly/colly"
  "github.com/antchfx/htmlquery"
)

type Task struct {
	ID                string    `db:"id, primarykey" json:"id"`
	Price             float64   `db:"threshold" json:"threshold"`
	Image             string    `db:"image" json:"image"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at" pg:",null"`
}

func (t *Task) FetchPrice() {
  c := colly.NewCollector()

  c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"

  c.OnResponse(func(r *colly.Response) {
    doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
    if err != nil {
        fmt.Println(err)
    }
    nodes := htmlquery.Find(doc, `//div[@class="a-row a-spacing-mini olpOffer"]`)
    img := htmlquery.Find(doc, `//div[@id="olpProductImage"]/a/img/@src`)
    if len(img) > 0 {
      t.Image = strings.Trim(htmlquery.InnerText(img[0]), " ")
    }
    for i := 0; i < len(nodes); i++ {
      priceNode := htmlquery.Find(nodes[i], `//span[@class="a-size-large a-color-price olpOfferPrice a-text-bold"]`)
      priceString := strings.Trim(htmlquery.InnerText(priceNode[0]), " ")
      price, err := strconv.ParseFloat(priceString[1:], 10)
      if (err != nil) {
        fmt.Println(err)
        continue
      }
      t.Price = price
      t.UpdatedAt = time.Now()
      break
    }
  })

  c.Visit("https://www.amazon.com/gp/offer-listing/" + t.ID)
}
