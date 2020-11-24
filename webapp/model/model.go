//Package model : Contains backend processing of data
package model

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
)

const daysLimit = 14
const maxUint = ^uint(0)
const minUint = 0

//MinDate :
const MinDate = 200901

//MaxDate :
const MaxDate = 201230

/*
some assumption:
1 day has 3 slots, morning, afternoon and evening
book only 2 months in advance to reduce initialization
*/

// ReserveDT : reservation date time obj
// Use 2 bst to store available and unavailable datetime
// bst key : date with format YYMMDDT
// where T 1 = morning, 2 = Afternoon, 3 = Evening
// bst satellite value is booking ID
type ReserveDT struct {
	available   Tree
	unavailable Tree
	date        Tree
}

//Reserve : Reserve venue
//Mutex
func (rdt *ReserveDT) Reserve(datetime int, bookID int) {
	rdt.available.Delete(datetime)
	rdt.unavailable.AddA(datetime, bookID)
}

func (rdt *ReserveDT) delReserve(datetime int) {
	rdt.unavailable.Delete(datetime)
	rdt.available.Add(datetime)
}

//ReadAvailable : Get all available date
func (rdt *ReserveDT) ReadAvailable() []int {
	return rdt.available.Flatten()
}

func (rdt *ReserveDT) init(days int) {
	t := time.Now()
	until := t.AddDate(0, 0, days)
	for t.Before(until) {
		year, month, day := t.Date()
		//Date format: YYMMDDT
		//1 is morning, 2 is afternoon, 3 is evening
		for i := 1; i <= 3; i++ {
			formatDate := fmt.Sprintf("%d%02d%02d%d", year%100, int(month), day, i)
			iformatDate, _ := strconv.Atoi(formatDate)
			rdt.available.AddA(iformatDate, 0)
		}
		formatDateOnly := fmt.Sprintf("%d%02d%02d", year%100, int(month), day)
		iformatDateOnly, _ := strconv.Atoi(formatDateOnly)
		rdt.date.Add(iformatDateOnly)
		t = t.AddDate(0, 0, 1)
	}

}

//GetBookingList return booking list
func (rdt *ReserveDT) GetBookingList() Queue {
	//rdt.unavailable.FlattenA
	datetime, bookid := rdt.unavailable.FlattenA()
	bookingQ := Queue{}
	for i := range datetime {
		bookingQ.enqueue(datetime[i], bookid[i], 3)
	}
	bookingQ.printAllNodes()
	return bookingQ
}

//GetDate :return  map[date][time][availability]
func (rdt *ReserveDT) GetDate(min int, max int) (map[int]map[int]string, []int) {
	rdates := rdt.date.Flatten()
	dates := make([]int, 0, len(rdates))
	for _, i := range rdates {
		if i >= min && i <= max {
			dates = append(dates, i)
		}
	}
	ravailable := rdt.available.Flatten()
	available := make([]int, 0, len(ravailable))
	for _, i := range ravailable {
		if i >= min*10 && i <= max*10+3 {
			available = append(available, i)
		}
	}
	runavailable := rdt.unavailable.Flatten()
	unavailable := make([]int, 0, len(runavailable))
	for _, i := range runavailable {
		if i >= min*10 && i <= max*10+3 {
			unavailable = append(unavailable, i)
		}
	}

	bookingData := make(map[int]map[int]string)

	for _, v := range available {
		dateOnly := v / 10
		times := v % 10
		_, exists := bookingData[dateOnly]
		if exists {
			bookingData[dateOnly][times] = "AVAILABLE"
		} else {
			timeData := make(map[int]string)
			timeData[times] = "AVAILABLE"
			bookingData[dateOnly] = timeData
		}

	}
	for _, v := range unavailable {
		dateOnly := v / 10
		times := v % 10
		_, exists := bookingData[dateOnly]
		if exists {
			bookingData[dateOnly][times] = "UNAVAILABLE"
		} else {
			timeData := make(map[int]string)
			timeData[times] = "UNAVAILABLE"
			bookingData[dateOnly] = timeData
		}
	}
	return bookingData, dates
}

type date struct {
	year  int
	month int
	day   int
}

//Venue :
type Venue struct {
	Capacity int
	Kind     string
	Location string
	Name     string
	Desc     string
}

//Intersect : intersection of 2 int array
func Intersect(a []int, b []int) []int {
	m := make(map[int]bool)
	c := make([]int, 0)
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

//categorical data stored as map, numerical data stored as tree. this is to facilitate search
type venueDB struct {
	//Venues :
	Venues       map[int]Venue
	kindMap      map[string][]int
	locationMap  map[string][]int
	capacityTree *Tree
	counter      int
	VenueMap     map[int]string
}

// Query :
type Query struct {
	Location string
	CapMin   int
	CapMax   int
	Kind     string
	DateMin  int
	DateMax  int
}

func getMapKey(m map[int]Venue) []int {
	result := make([]int, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
func (vDB *venueDB) KindList() []string {
	keys := make([]string, 0, len(vDB.kindMap))
	for k := range vDB.kindMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (vDB *venueDB) LocationList() []string {
	keys := make([]string, 0, len(vDB.locationMap))
	for k := range vDB.locationMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (vDB *venueDB) Caps() (int, int) {
	caps := vDB.capacityTree.Flatten()
	minCap := caps[0]
	maxCap := caps[len(caps)-1]
	return minCap, maxCap
}

func (vDB *venueDB) Filter(q Query) (map[int]Venue, []int) {
	result := getMapKey(vDB.Venues)
	r1, _ := vDB.locationMap[q.Location]
	if q.Location != "Nil" {
		result = Intersect(result, r1)
	}
	r2, _ := vDB.kindMap[q.Kind]
	if q.Kind != "Nil" {
		result = Intersect(result, r2)
	}
	finalResult := make(map[int]Venue)
	finalOrder := make([]int, 0, len(result))
	for _, mapIndex := range result {
		v := vDB.Venues[mapIndex]
		if v.Capacity >= q.CapMin && v.Capacity <= q.CapMax {
			finalOrder = append(finalOrder, mapIndex)
			finalResult[mapIndex] = v
		}
	}
	return finalResult, finalOrder
}

func (vDB *venueDB) addMap(m map[string][]int, s string) {
	_, exists := m[s]
	if exists {
		m[s] = append(m[s], vDB.counter)
	} else {
		m[s] = make([]int, 0)
		m[s] = append(m[s], vDB.counter)
	}
}

func (vDB *venueDB) AddVenue(v Venue) error {
	_, exists := vDB.GetID(v.Name)
	if !exists {
		vDB.Venues[vDB.counter] = v
		vDB.addMap(vDB.kindMap, v.Kind)
		vDB.addMap(vDB.locationMap, v.Location)
		vDB.capacityTree.Add(v.Capacity)
		vDB.counter++
	} else {
		msg := fmt.Sprintf("Error, %s already exists!", v.Name)
		return errors.New(msg)
	}

	return nil
}

func (vDB *venueDB) GetID(name string) (int, bool) {
	for k, v := range vDB.Venues {
		if v.Name == name {
			return k, true
		}
	}
	return 0, false
}

//Booking :
type Booking struct {
	IDBook   int
	User     string
	VenueID  int
	Datetime int
}

// Bookings id start from 1
// venueReserve , k: venue id, v: DateTime
type bookingDB struct {
	Bookings     map[int]*Booking
	VenueReserve map[int]*ReserveDT
}

func (b *bookingDB) getBookingID(vid int, date int) int {
	return b.VenueReserve[vid].unavailable.root.find(date).bookID
}
func (b *bookingDB) getBookingDetails(vid int, date int) *Booking {
	bid := b.getBookingID(vid, date)
	return b.Bookings[bid]
}

// Reserve :
func (b *bookingDB) Reserve(venueID int, datetime int, user string) int {
	order := Booking{
		IDBook:   len(b.Bookings) + 1,
		User:     user,
		Datetime: datetime,
		VenueID:  venueID,
	}
	b.Bookings[order.IDBook] = &order
	b.VenueReserve[order.VenueID].Reserve(order.Datetime, order.IDBook)
	return order.IDBook
}

func (b *bookingDB) DelReserve(venueID int, datetime int) {
	b.VenueReserve[venueID].delReserve(datetime)
}

// Model : consolidate all neede obj
type Model struct {
	VenueDB   *venueDB
	BookingDB *bookingDB
}

//InitModel : creates all needed obj for the app
func InitModel() Model {
	venues := make(map[int]Venue)
	kindMap := make(map[string][]int)
	locationMap := make(map[string][]int)
	VenueMap := make(map[int]string)
	capacityTree := Tree{}
	venueDB := venueDB{
		Venues:       venues,
		kindMap:      kindMap,
		locationMap:  locationMap,
		capacityTree: &capacityTree,
		counter:      1,
		VenueMap:     VenueMap,
	}
	bookings := make(map[int]*Booking)
	reDT := make(map[int]*ReserveDT)
	bookingDB := bookingDB{
		Bookings:     bookings,
		VenueReserve: reDT,
	}
	model := Model{
		VenueDB:   &venueDB,
		BookingDB: &bookingDB,
	}
	return model
}

//AddVenue :
func (m *Model) AddVenue(v Venue) error {
	venueID := m.VenueDB.counter
	e := m.VenueDB.AddVenue(v)
	if e == nil {
		rdt := ReserveDT{}
		rdt.init(daysLimit)
		m.BookingDB.VenueReserve[venueID] = &rdt
		m.VenueDB.VenueMap[venueID] = v.Name
		return nil
	}
	return e
}

//data structure 1
//book a venue lead to a booking queue, vip and normal , add delay before booking is completed

//data structure 2
//construct a bst for available time, another bst for unavailable time

//error handling needed to be smoked in, i.e. different user cannot edit booking by other user
//cannot book past date
