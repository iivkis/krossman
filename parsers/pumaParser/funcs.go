package pumaParser

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"krossman/parsers/customxml"
)

func New(c *Config) *pr {
	return &pr{
		config:   c,
		filepath: "./static/shops/" + c.SaveAs + ".xml",
		limit:    1500,
		offset:   0,
	}
}

func (p *pr) getContent() ([]byte, error) {
	baseURL, err := url.Parse(p.config.Address)
	if err != nil {
		return []byte{}, err
	}

	query := baseURL.Query()
	query.Set("user", p.config.User)
	query.Set("code", p.config.Code)
	query.Set("feed_id", strconv.Itoa(p.config.FeedID))
	query.Set("format", "xml")
	query.Set("limit", strconv.Itoa(p.limit))
	query.Set("offset", strconv.Itoa(p.offset))
	baseURL.RawQuery = query.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (p *pr) save(xmlData interface{}) error {
	data, err := xml.Marshal(xmlData)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(p.filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(data)

	return nil
}

func (p *pr) Parse() {
	log.Printf("[%s] Started!\n", p.config.SaveAs)

	var customOffersBuffer customxml.NodeCustomOffersList
	var articleCache = map[string]int{}
	var processedOffers int

	for {
		b, err := p.getContent()
		if err != nil {
			fmt.Println(err)
			return
		}

		var catalog xmlYmlCatalog
		if err = xml.Unmarshal(b, &catalog); err != nil {
			log.Println(err)
			return
		}

		if len(catalog.Shop.OffersList.Offers) == 0 {
			p.offset = 0
			break
		}

		for _, offer := range catalog.Shop.OffersList.Offers {

			if index, ok := articleCache[offer.VendorCode]; ok {
				customOffersBuffer.Items[index].SizesList.Sizes = append(customOffersBuffer.Items[index].SizesList.Sizes,
					customxml.NodeCustomOfferSize{Size: offer.Params[1]},
				)
				continue
			}

			if offer.CategoryID == p.config.CategoryID {
				customOffersBuffer.Items = append(customOffersBuffer.Items, customxml.NodeCustomOffer{
					Title:     offer.Name,
					Url:       offer.Url,
					Article:   offer.VendorCode,
					Firma:     "Puma",
					Color:     offer.Params[0],
					Price:     offer.getOldprice(),
					PriceFree: offer.getPrice(),
					PhotosList: customxml.NodeCustomOfferPhotosList{Photos: []customxml.NodeCustomOfferPhoto{
						{Main: 1, Url: offer.Picture},
					}},
					SizesList: customxml.NodeCustomOfferSizesList{
						Sizes: []customxml.NodeCustomOfferSize{
							{Size: offer.Params[1]},
						},
					},
					Group:        1,
					FreeShipping: 1,
					Sex:          1,
				})

				articleCache[offer.VendorCode] = len(customOffersBuffer.Items) - 1

				log.Printf("[%s] Add (#%d) %s\n", p.config.SaveAs, len(customOffersBuffer.Items), offer.Name)
			}
		}

		p.offset += p.limit
		processedOffers += len(catalog.Shop.OffersList.Offers)

	}

	if err := p.save(customOffersBuffer); err != nil {
		fmt.Println(err)
	}

	log.Printf("[%s] Saved; Processed %d offers\n", p.config.SaveAs, processedOffers)
}
