package nikeParser

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func New(c *Config) *pr {
	return &pr{
		config:   c,
		filepath: "./static/shops/" + c.SaveAs + ".xml",

		//дополнительные данные
		limit:  1500,
		offset: 0,
	}
}

func (p *pr) Parse() {
	log.Printf("[%s] Started!\n", p.config.SaveAs)

	//get vase url
	baseURL, query := getBaseURL(p)

	var customOffersBuffer nodeCustomOffersList
	var urlCache map[string]bool = map[string]bool{}
	var processedOffers int

	for {
		//set query params
		query.Set("limit", strconv.Itoa(p.limit))
		query.Set("offset", strconv.Itoa(p.offset))
		baseURL.RawQuery = query.Encode()

		//request & read body
		resp, err := http.Get(baseURL.String())
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		//parse body from xml
		var catalog xmlYmlCatalog
		if err = xml.Unmarshal(b, &catalog); err != nil {
			log.Println(err)
			return
		}

		if len(catalog.Shop.OffersList.Offers) == 0 {
			break
		}

		//offers filter & customize
		for _, offer := range catalog.Shop.OffersList.Offers {
			if offer.CategoryID == p.config.CategoryID && !urlCache[offer.Url] {
				customOffersBuffer.Items = append(customOffersBuffer.Items, nodeCustomOffer{
					Title:        offer.Name,
					Url:          offer.Url,
					Firma:        "Nike",
					Price:        offer.Price,
					Color:        strings.TrimSpace(strings.Split(offer.Name, "-")[len(strings.Split(offer.Name, "-"))-1]),
					Article:      strings.Split(strings.Split(offer.Picture, "/")[len(strings.Split(offer.Picture, "/"))-1], "?")[0],
					Group:        1,
					Sex:          1,
					FreeShipping: 1,
				})

				urlCache[offer.Url] = true

				log.Printf("[%s] Add (#%d) %s\n", p.config.SaveAs, len(customOffersBuffer.Items), offer.Name)
			}
		}

		p.offset += p.limit
		processedOffers += len(catalog.Shop.OffersList.Offers)
	}

	// var customOffersBuffer nodeCustomOffersList
	// for _, offer := range offersBuffer.Offers {
	// 	customOffer := nodeCustomOffer{
	// 		Title: offer.Name,
	// 		Url:   offer.Url,
	// 		Firma: "Nike",
	// 		Price: offer.Price,

	// 		Color:   strings.TrimSpace(strings.Split(offer.Name, "-")[len(strings.Split(offer.Name, "-"))-1]),
	// 		Article: strings.Split(strings.Split(offer.Picture, "/")[len(strings.Split(offer.Picture, "/"))-1], "?")[0],

	// 		Group:        1,
	// 		Sex:          1,
	// 		FreeShipping: 1,
	// 	}
	// 	customOffersBuffer.Items = append(customOffersBuffer.Items, customOffer)
	// }

	//encode xml offers
	marshalCustomOffers, err := xml.Marshal(customOffersBuffer)
	if err != nil {
		log.Println(err)
		return
	}

	//open or create file with 777 right
	file, err := os.OpenFile(p.filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	//установка каретки в начала файла и запись
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(marshalCustomOffers)

	log.Printf("[%s] Saved; Processed %d offers\n", p.config.SaveAs, processedOffers)
}

// return base URL
func getBaseURL(p *pr) (*url.URL, url.Values) {
	baseURL, err := url.Parse(p.config.Address)
	if err != nil {
		panic(err)
	}

	query := baseURL.Query()
	query.Set("user", p.config.User)
	query.Set("code", p.config.Code)
	query.Set("feed_id", strconv.Itoa(p.config.FeedID))
	query.Set("format", "xml")

	return baseURL, query
}
