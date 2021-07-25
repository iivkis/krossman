package pumaParser

import "encoding/xml"

type nodeOffer struct {
	AttrID     int      `xml:"id,attr"`
	CategoryID string   `xml:"categoryId"`
	Name       string   `xml:"name"`
	Oldprice   float32  `xml:"oldprice"`
	Price      float32  `xml:"price"`
	Picture    string   `xml:"picture"`
	Url        string   `xml:"url"`
	VendorCode string   `xml:"vendorCode"`
	Params     []string `xml:"param"`
}

type nodeOffersList struct {
	XMLName xml.Name    `xml:"offers"`
	Offers  []nodeOffer `xml:"offer"`
}

type nodeShop struct {
	XMLName    xml.Name       `xml:"shop"`
	OffersList nodeOffersList `xml:"offers"`
}

type xmlYmlCatalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Shop    nodeShop `xml:"shop"`
}

func (offer *nodeOffer) getPrice() float32 {
	if offer.Oldprice < offer.Price {
		return offer.Oldprice
	}
	return offer.Price
}

func (offer *nodeOffer) getOldprice() float32 {
	if offer.Oldprice > offer.Price {
		return offer.Oldprice
	}
	return offer.Price
}
