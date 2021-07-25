package nikeParser

import "encoding/xml"

type nodeOffer struct {
	AttrID        int     `xml:"id,attr"`
	AttrAvailable bool    `xml:"available,attr"`
	CategoryID    string  `xml:"categoryId"`
	Description   string  `xml:"description"`
	Name          string  `xml:"name"`
	Picture       string  `xml:"picture"`
	Price         float32 `xml:"price"`
	Url           string  `xml:"url"`
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
