package flatscan

import (
	"appengine"
	"appengine/datastore"
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
	"time"
)

type FlatOffer struct {
	RentN       int64
	Zip         int64
	District    string
	Street      string
	Rooms       float64
	Size        string
	Url         string
	Description string
	TimeUpdated int64
	Valid       bool
	ID          string
}

func (f FlatOffer) Key() string {
	md5Writer := md5.New()
	io.WriteString(md5Writer, f.Url)
	return fmt.Sprintf("%x", md5Writer.Sum(nil))
}

func (f FlatOffer) Type() string {
	return "FlatOffer"
}

func (f FlatOffer) AEKey(con appengine.Context) *datastore.Key {
	return datastore.NewKey(con, "counter", f.Key(), 0, nil)
}

const (
	maxRent  int64   = 600
	minRooms float64 = 2
	maxRooms float64 = 3
)

var zipCodes map[int64]bool = map[int64]bool{
	10117: true,
}
var districts []string = []string{"mitte", "friedrichshain", "prenzlauer", "kreuzberg", "neukölln", "treptow"}

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
func getRent(rent string) (rentN int64, err error) {
	// remove "Kaltmiete: "
	rent = rent[11:]

	// only if there's a price
	if strings.Contains(rent, "EUR") {
		// remove " EUR VB" or " EUR"
		rentParts := strings.SplitN(rent, " ", 2)

		// remove digits "1.150" > "1150"
		integer := strings.Replace(rentParts[0], ".", "", -1)

		// convert to integer
		rentN, err = strconv.ParseInt(integer, 10, 64)
	}

	return
}

func getAttributes(c appengine.Context, sel *goquery.Selection, index int) (attribute string) {
	if len(sel.Nodes) < index {
		return strings.Trim(sel.Nodes[index].FirstChild.Data, "\n\t ")
	} else {
		c.Infof(fmt.Sprintf("%v; %d", sel.Nodes, index))
	}

	return
}

func getZIPCode(location string) (zip int64, district string) {
	locationParts := strings.Split(location, " ")

	// convert to integer
	zip, err := strconv.ParseInt(locationParts[0], 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	district = locationParts[len(locationParts)-1]

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

func CheckOffer(offer *FlatOffer) (isWanted bool) {
	if offer.RentN > maxRent {
		return false
	}

	if !zipCodes[offer.Zip] {
		isWanted = false
		for _, wantedDistrict := range districts {
			isWanted = isWanted || strings.Contains(strings.ToLower(offer.District), wantedDistrict)
		}

		if !isWanted {
			return false
		}
	}

	if offer.Rooms < minRooms || offer.Rooms > maxRooms {
		return false
	}

	return true
}

func GetOffer(doc *goquery.Document, c appengine.Context) (offer *FlatOffer, err error) {
	offer = &FlatOffer{Valid: true, TimeUpdated: time.Now().Unix()}

	rentS := fmt.Sprintf("%v", doc.Find("#viewad-price").Get(0).FirstChild.Data)
	offer.RentN, _ = getRent(rentS)

	sel := doc.Find(".c-attrlist > dd").Find("span")

	index := 1

	if len(sel.Nodes) == 5 {
		offer.Street = getAttributes(c, sel, index)
		index++
	}

	offer.Zip, offer.District = getZIPCode(getAttributes(c, sel, index))
	index++
	offer.Rooms, _ = getRooms(getAttributes(c, sel, index))

	index++
	offer.Size = getAttributes(c, sel, index)

	/*
	   offer.Description = doc.Find("#viewad-description-text").Text()
	   	if len(offer.Description) > 500 {
	   		offer.Description = offer.Description[:499]
	   	}
	*/

	return
}
