package customxml

import "encoding/xml"

//sizes
type NodeCustomOfferSize struct {
	XMLName xml.Name `xml:"size"`
	Size    string   `xml:",chardata"`
}

type NodeCustomOfferSizesList struct {
	XMLName xml.Name              `xml:"sizes"`
	Sizes   []NodeCustomOfferSize `xml:"size"`
}

//photos
type NodeCustomOfferPhoto struct {
	Main int    `xml:"main,attr"`
	Url  string `xml:",chardata"`
}

type NodeCustomOfferPhotosList struct {
	XMLName xml.Name               `xml:"photos"`
	Photos  []NodeCustomOfferPhoto `xml:"photo"`
}

type NodeCustomOffer struct {
	XMLName      xml.Name                  `xml:"item"`
	Title        string                    `xml:"title"`
	Url          string                    `xml:"url"`
	Article      string                    `xml:"art"`
	Firma        string                    `xml:"firma"`
	Color        string                    `xml:"color"`
	Price        float32                   `xml:"price"`
	PriceFree    float32                   `xml:"price_free"`
	Group        int                       `xml:"group"`
	FreeShipping int                       `xml:"free_shipping"`
	Sex          int                       `xml:"sex"`
	PhotosList   NodeCustomOfferPhotosList `xml:"photos"`
	SizesList    NodeCustomOfferSizesList  `xml:"sizes"`
}

type NodeCustomOffersList struct {
	XMLName xml.Name          `xml:"items"`
	Items   []NodeCustomOffer `xml:"item"`
}
