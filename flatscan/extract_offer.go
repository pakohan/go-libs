package flatscan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
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
	Valid       bool
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
		url, valid := adTitle.First().Attr("href")
		if valid {
			urls = append(urls, url)
		} else {
			fmt.Errorf(url)
		}
	})

	return
}

// func getRent converts a string containing the rent of a flat to an integer
// for example:
//  "Kaltmiete 1.150 EUR VB" becomes '1150'
func getRent(rent string) (rentN int64) {
	// remove "Kaltmiete: "
	rent = rent[11:]

	// only if there's a price
	if strings.Contains(rent, "EUR") {
		// remove " EUR VB" or " EUR"
		rentParts := strings.SplitN(rent, " ", 2)

		// remove digits "1.150" > "1150"
		integer := strings.Replace(rentParts[0], ".", "", -1)

		// convert to integer
		rentN, _ = strconv.ParseInt(integer, 10, 64)
	}

	return
}

func getAttributes(sel *goquery.Selection, index int) (attribute string) {
	return strings.Trim(sel.Nodes[index].FirstChild.Data, "\n\t ")
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

func getRooms(rooms string) (roomsN float64) {
	if strings.Contains(rooms, ">") || strings.Contains(rooms, "<") {
		rooms = strings.SplitN(rooms, " ", 2)[1]
	}

	rooms = strings.Replace(rooms, ".", ",", -1)

	roomsN, _ = strconv.ParseFloat(rooms, 64)

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

func GetOffer(doc *goquery.Document) (offer *FlatOffer) {
	offer = &FlatOffer{Valid: true}

	rentS := fmt.Sprintf("%v", doc.Find("#viewad-price").Get(0).FirstChild.Data)
	offer.RentN = getRent(rentS)

	sel := doc.Find(".c-attrlist > dd").Find("span")

	index := 1

	if len(sel.Nodes) == 5 {
		offer.Street = getAttributes(sel, index)
		index++
	}

	offer.Zip, offer.District = getZIPCode(getAttributes(sel, index))
	index++
	offer.Rooms = getRooms(getAttributes(sel, index))
	index++
	offer.Size = getAttributes(sel, index)

	offer.Description = doc.Find("#viewad-description-text").Text()
	if len(offer.Description) > 500 {
		offer.Description = offer.Description[:499]
	}

	return
}
