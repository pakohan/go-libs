package main

import (
	//"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pakohan/go-libs/flatscan"
	"net/http"
)

type C struct{}

func (c C) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

var urls = []string{
	// "/anzeigen/s-anzeige/wir-suchen-einen-nachmieter-fuer-unsere-2-1-2-zimmerwohnung-/205874631-203-3545?ref=search",
	// "/anzeigen/s-anzeige/stilvoll-moebliert-und-lichtdurchflutet-apartment-naehe-mauerpark/207461196-203-3491?ref=search",
	// "/anzeigen/s-anzeige/nachmieter-gesucht/206104331-203-3418?ref=search",
	// "/anzeigen/s-anzeige/dringend-nachmieter-gesucht-kueche-geschenkt-ab-01-07-o-01-08-!/207411536-203-3383?ref=search",
	// "/anzeigen/s-anzeige/bel-etage-in-charmanter-dahlemer-villa-mit-garten/207809356-203-3429?ref=search",
	// "/anzeigen/s-anzeige/zeigen-sie-groesse-mit-dachterrasse/207814751-203-3478?ref=search",
	// "/anzeigen/s-anzeige/penthouse-ueber-den-daechern-von-berlin/207816071-203-3521?ref=search",
	// "/anzeigen/s-anzeige/stilvoll-moebliert-und-lichtdurchflutet-apartment-naehe-mauerpark/207461196-203-3491?ref=search",
	// "/anzeigen/s-anzeige/suesse-1-zimmer-wohnung-ab-01-06-07-frei-besichtigung-20-05-12-00/207285575-203-3423?ref=search",
	"/anzeigen/s-anzeige/nachmieter-gesucht-2zkb-direkt-an-schoenhauser-allee-625-kalt/207946867-203-3490?ref=search",
}

func main() {
	/*
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return errors.New("")
			},
		}
	*/

	for _, url := range urls {
		str := fmt.Sprintf("http://kleinanzeigen.ebay.de%s", url)
		//fmt.Println(str)
		doc, err := loadDocument(str)

		if err != nil {
			panic(err)
		}

		c := C{}

		offer, err := flatscan.GetOffer(doc, c)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", *offer)
	}
}

func loadDocument(url string) (doc *goquery.Document, err error) {
	/*
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://kleinanzeigen.ebay.de/anzeigen/s-anzeige/wir-suchen-einen-nachmieter-fuer-unsere-2-1-2-zimmerwohnung-/205874631-203-3545?ref=search", nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("User-Agent", "Golang Spider Bot v. 3.0")
	*/

	resp, err := http.Get(url)
	if err != nil {
		doc = nil
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		doc = nil
	}

	return
}
