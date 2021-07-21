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
		limit:  1000,
		offset: 0,
	}
}

func (p *pr) Parse() {
	log.Printf("[%s] Started!\n", p.config.SaveAs)
	baseURL, query := getBaseURL(p)

	var offersBuffer nodeOffersList
	var processedOffers int

	for {
		query.Set("limit", strconv.Itoa(p.limit))
		query.Set("offset", strconv.Itoa(p.offset))
		baseURL.RawQuery = query.Encode()

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

		var catalog xmlYmlCatalog
		if err = xml.Unmarshal(b, &catalog); err != nil {
			log.Println(err)
			return
		}

		if len(catalog.Shop.OffersList.Offers) == 0 {
			break
		}

		for _, offer := range catalog.Shop.OffersList.Offers {
			if offer.CategoryID == p.config.CategoryID {
				offersBuffer.Offers = append(offersBuffer.Offers, offer)
				log.Printf("[%s] Add (#%d) %s\n", p.config.SaveAs, len(offersBuffer.Offers), offer.Name)
			}
		}

		p.offset += p.limit
		processedOffers += len(catalog.Shop.OffersList.Offers)
	}

	var customOffersBuffer nodeCustomOffersList

	for _, offer := range offersBuffer.Offers {

		customOffer := nodeCustomOffer{
			Title: offer.Name,
			Url:   offer.Url,
			Firma: "Nike",
			Price: offer.Price,

			Color:   strings.TrimSpace(strings.Split(offer.Name, "-")[len(strings.Split(offer.Name, "-"))-1]),
			Article: strings.Split(strings.Split(offer.Picture, "/")[len(strings.Split(offer.Picture, "/"))-1], "?")[0],

			Group:        1,
			Sex:          1,
			FreeShipping: 1,
		}

		customOffersBuffer.Items = append(customOffersBuffer.Items, customOffer)
	}

	encodeOffers, err := xml.Marshal(customOffersBuffer)
	if err != nil {
		log.Println(err)
		return
	}

	file, err := os.OpenFile(p.filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	file.Write(encodeOffers)
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
