package flatscan

import (
	//"appengine"
	//"appengine/datastore"
	"code.google.com/p/go.net/html"
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
	"time"
)

type FlatOffer struct {
	RentN       float64
	Zip         int64
	District    string
	Street      string
	Rooms       float64
	Size        float64
	Url         string
	Description string
	TimeUpdated int64
	Valid       bool
	ID          string
	Title       string
}

type Context interface {
	Infof(format string, args ...interface{})
}

func (f FlatOffer) Key() string {
	md5Writer := md5.New()
	io.WriteString(md5Writer, f.Url)
	return fmt.Sprintf("%x", md5Writer.Sum(nil))
}

func (f FlatOffer) Type() string {
	return "FlatOffer"
}

const (
	maxRent  float64 = 600.0
	minRooms float64 = 2
	maxRooms float64 = 3
)

func ExtractLinks(doc *goquery.Document) (urls []string) {
	doc.Find(".ad-title").Each(func(i int, adTitle *goquery.Selection) {
		url, exists := adTitle.First().Attr("href")
		if exists {
			urls = append(urls, url)
		}
	})

	return
}

// func getRent converts a string containing the rent of a flat to an integer
// for example:
//  "Kaltmiete 1.150 EUR VB" becomes '1150'
func getRent(rent string) (rentN float64, err error) {
	/* remove "Kaltmiete: "
	rent = rent[11:]

	// only if there's a price
	if strings.Contains(rent, "EUR") {
		// remove " EUR VB" or " EUR"
		rentParts := strings.SplitN(rent, " ", 2)

		// remove digits "1.150" > "1150"
		integer := strings.Replace(rentParts[0], ".", "", -1)

		// convert to integer*/
	rentN, err = strconv.ParseFloat(rent, 64)

	return
}

func getAttributes(c Context, sel *goquery.Selection, index int) (attribute string) {
	if len(sel.Nodes) > index {
		return strings.Trim(sel.Nodes[index].FirstChild.Data, "\n\t ")
	} else {
		//c.Infof(fmt.Sprintf("%v; %d", sel.Nodes, index))
	}

	return
}

func getZIPCode(location string) (zip int64, district string) {
	locationParts := strings.Split(location, " ")

	fmt.Println(location)

	// convert to integer
	zip, err := strconv.ParseInt(locationParts[0], 10, 64)
	if err != nil {
		fmt.Println("getZIPCode", err)
	}

	after := strings.Split(location, "-")
	district = strings.Trim(after[len(after)-1], " ")

	return
}

func getSize(sizeString string) (size float64) {
	// convert to integer
	size, err := strconv.ParseFloat(sizeString, 64)
	if err != nil {
		fmt.Println("getSize", err)
	}

	return
}

func getRooms(rooms string) (roomsN float64, err error) {
	if strings.Contains(rooms, ">") || strings.Contains(rooms, "<") {
		rooms = strings.SplitN(rooms, " ", 2)[1]
	}

	rooms = strings.Replace(rooms, ".", ",", -1)

	roomsN, err = strconv.ParseFloat(rooms, 64)

	return
}

func GetOffer(doc *goquery.Document, c Context) (offer *FlatOffer, err error) {
	offer = &FlatOffer{Valid: true, TimeUpdated: time.Now().Unix()}

	// _ = fmt.Sprintf("%v", getAttributes(c, doc.Find("#viewad-price"), 0))
	nodes := doc.Find("meta[itemprop=price]").Nodes
	var attr []html.Attribute
	if len(nodes) > 0 {
		attr = nodes[0].Attr
	}

	if len(attr) >= 2 {
		offer.RentN, _ = getRent(attr[1].Val)
	} else {
		fmt.Println(attr)
	}

	offer.Street = getAttributes(c, doc.Find("#street-address"), 0)

	offer.Zip, offer.District = getZIPCode(getAttributes(c, doc.Find("#viewad-locality"), 0))

	sel := doc.Find("#viewad-attributes > dd > span")

	offer.Rooms, _ = getRooms(getAttributes(c, sel, 0))

	offer.Size = getSize(getAttributes(c, sel, 1))

	offer.Title = getAttributes(c, doc.Find("#viewad-title"), 0)

	/*
		offer.Description = doc.Find("#viewad-description-text").Text()
		if len(offer.Description) > 500 {
			offer.Description = offer.Description[:499]
		}
	*/
	return
}
