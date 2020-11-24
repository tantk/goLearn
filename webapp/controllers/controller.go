//Package controller : middle ware to link front and backend
package controller

import (
	"fmt"
	config "gia/config"
	model "gia/model"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/microcosm-cc/bluemonday"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func santizeString(s string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(s)
}

//Booking : it is at this moment i realised i should have done a data structure that can
//make Booking struct always sorted by datetime when filtered by user.
// My BST is sorted by datetime but i think it require high time complexity to filter by user.
type Booking struct {
	IDBook    int
	User      string
	VenueID   int
	VenueName string
	Date      int
	Time      string
}

func convertBooking(booking model.Booking, mapping map[int]string) Booking {

	date := booking.Datetime / 10
	time := booking.Datetime % 10
	ctime := "nil time"
	switch time {
	case 1:
		ctime = "Morning"
	case 2:
		ctime = "Afternoon"
	case 3:
		ctime = "Evening"
	}

	b := Booking{
		IDBook:    booking.IDBook,
		User:      booking.User,
		VenueID:   booking.VenueID,
		VenueName: mapping[booking.VenueID],
		Date:      date,
		Time:      ctime,
	}
	return b
}

// User : User object
type User struct {
	Username string
	Password []byte
	First    string
	Last     string
	Bookings []int
}

type wrongUserError struct {
	user1, user2 string //error message
}

func (p *wrongUserError) Error() string {
	format := "Invalid credential user %s against user %s"
	return fmt.Sprintf(format, p.user1, p.user2)
}

type mapUsers map[string]User
type mapSessions map[string]string

//Ctl : controller that holds all needed obj
type Ctl struct {
	Users    mapUsers
	Sessions mapSessions
	Template *template.Template
	Model    model.Model
	Logging  *config.Logging
}

//getUser :
func (a *Ctl) getUser(res http.ResponseWriter, req *http.Request) User {
	// get current session cookie
	myCookie, err := req.Cookie(config.NCOOKIE)
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  config.NCOOKIE,
			Value: id.String(),
		}
	}
	http.SetCookie(res, myCookie)
	// if the user exists already, get user
	var myUser User
	if username, ok := a.Sessions[myCookie.Value]; ok {
		myUser = a.Users[username]
	}
	return myUser
}

// Index : landing page
func (a *Ctl) Index(res http.ResponseWriter, req *http.Request) {
	type pageData struct {
		User User
	}
	d := pageData{
		User: a.getUser(res, req),
	}
	a.Template.ExecuteTemplate(res, "index.html", d)
}

// Profile : user profile page
func (a *Ctl) Profile(res http.ResponseWriter, req *http.Request) {
	type pageData struct {
		User User
	}
	d := pageData{
		User: a.getUser(res, req),
	}
	if !a.alreadyLoggedIn(req) {
		a.Logging.Warning.Println("Unauthorised access to Profile from ", req.UserAgent())
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "profile.html", d)
}

// EditProfile :
func (a *Ctl) EditProfile(res http.ResponseWriter, req *http.Request) {
	type pageData struct {
		User User
	}
	d := pageData{
		User: a.getUser(res, req),
	}
	if !a.alreadyLoggedIn(req) {
		a.Logging.Warning.Println("Unauthorised access to EditProfile from ", req.UserAgent())
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		// get form values
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		d.User.First = firstname
		d.User.Last = lastname
		a.Users[d.User.Username] = d.User
		// redirect to profile
		a.Logging.Info.Println("Profile edited from ", req.UserAgent())
		http.Redirect(res, req, "/profile", http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "editProfile.html", d)
}

// Signup :
func (a *Ctl) Signup(res http.ResponseWriter, req *http.Request) {
	if a.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	var myUser User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := santizeString(req.FormValue("username"))
		password := santizeString(req.FormValue("password"))
		firstname := santizeString(req.FormValue("firstname"))
		lastname := santizeString(req.FormValue("lastname"))
		if username != "" {
			// check if username exist/ taken
			if _, ok := a.Users[username]; ok {
				http.Error(res, "Username already taken", http.StatusForbidden)
				a.Logging.Info.Println("Signup with existing username from ", req.UserAgent())
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:    config.NCOOKIE,
				Value:   id.String(),
				Expires: time.Now().Add(2 * time.Hour),
			}
			http.SetCookie(res, myCookie)
			a.Sessions[myCookie.Value] = username
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				a.Logging.Error.Println("Error with password from ", req.UserAgent())
				return
			}
			myUser = User{
				Username: username,

				Password: bPassword,
				First:    firstname,
				Last:     lastname,
			}
			a.Users[username] = myUser
			a.Logging.Info.Println("New user sign up from ", req.UserAgent())
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "signup.html", myUser)
}

// Login :
func (a *Ctl) Login(res http.ResponseWriter, req *http.Request) {
	if a.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// process form submission
	if req.Method == http.MethodPost {
		username := santizeString(req.FormValue("username"))
		password := santizeString(req.FormValue("password"))
		// check if user exist with username
		myUser, ok := a.Users[username]
		if !ok {
			http.Error(res, "Username and/or password do not match", http.
				StatusForbidden)
			a.Logging.Info.Println("Unexisting username login from ", req.UserAgent())
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.
				StatusForbidden)
			a.Logging.Info.Println("Wrong password from ", req.UserAgent())
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:    config.NCOOKIE,
			Value:   id.String(),
			Expires: time.Now().Add(2 * time.Hour),
		}
		http.SetCookie(res, myCookie)
		a.Sessions[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusSeeOther)
		a.Logging.Info.Println("Successful login from ", req.UserAgent())
		return
	}
	a.Template.ExecuteTemplate(res, "login.html", nil)
}

// Logout :
func (a *Ctl) Logout(res http.ResponseWriter, req *http.Request) {
	if !a.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie(config.NCOOKIE)
	// delete the session
	delete(a.Sessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   config.NCOOKIE,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	a.Logging.Info.Println("User logout from ", req.UserAgent())
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func (a *Ctl) alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie(config.NCOOKIE)
	if err != nil {
		return false
	}
	username := a.Sessions[myCookie.Value]
	_, ok := a.Users[username]
	return ok
}

func prependStr(strs []string, str string) []string {
	strs = append(strs, str)
	copy(strs[1:], strs)
	strs[0] = str
	return strs
}

func reorderStr(strs []string, str string) []string {
	var i int
	for i = 0; i < len(strs); i++ {
		if strs[i] == str {
			break
		}
	}
	copy(strs[i:], strs[i+1:])
	strs[len(strs)-1] = ""
	strs = strs[:len(strs)-1]
	return prependStr(strs, str)
}

func mapNilAll(s string) string {
	if s == "All" {
		return "Nil"
	}
	return s
}

// Browse : venues
func (a *Ctl) Browse(res http.ResponseWriter, req *http.Request) {
	type pageData struct {
		User     User
		Venues   map[int]model.Venue
		Kind     []string
		Location []string
		MaxCap   int
		MinCap   int
		SMaxCap  int
		SMinCap  int
	}
	min, max := a.Model.VenueDB.Caps()
	data := pageData{
		User:     a.getUser(res, req),
		Venues:   a.Model.VenueDB.Venues,
		Kind:     prependStr(a.Model.VenueDB.KindList(), "All"),
		Location: prependStr(a.Model.VenueDB.LocationList(), "All"),
		MinCap:   min,
		MaxCap:   max,
	}
	if req.Method == http.MethodPost {
		venueKind := santizeString(req.FormValue("venueKind"))
		venueLocation := santizeString(req.FormValue("venueLocation"))
		venueMinCap, _ := strconv.Atoi(req.FormValue("venueMinCap"))
		venueMaxCap, _ := strconv.Atoi(req.FormValue("venueMaxCap"))
		q := model.Query{
			Location: mapNilAll(venueLocation),
			CapMin:   venueMinCap,
			CapMax:   venueMaxCap,
			Kind:     mapNilAll(venueKind),
		}
		venues, _ := a.Model.VenueDB.Filter(q)
		data.Venues = venues
		data.Kind = reorderStr(data.Kind, venueKind)
		data.Location = reorderStr(data.Location, venueLocation)
		data.SMaxCap = venueMaxCap
		data.SMinCap = venueMinCap
		a.Template.ExecuteTemplate(res, "browse.html", data)
		a.Logging.Trace.Println("Venue search from ", req.UserAgent())
		return
	}
	a.Template.ExecuteTemplate(res, "browse.html", data)
}

// Book :
func (a *Ctl) Book(res http.ResponseWriter, req *http.Request) {

	keys, ok := req.URL.Query()["venueId"]

	if !ok || len(keys[0]) < 1 {
		http.Redirect(res, req, "/browse", http.StatusSeeOther)
		return
	}
	vID, valid := strconv.Atoi(keys[0])
	if valid != nil {
		http.Redirect(res, req, "/browse", http.StatusSeeOther)
		return
	}
	venue := a.Model.VenueDB.Venues[vID]
	type data struct {
		Venue model.Venue
		Avail map[int]map[int]string
		Order []int
		Vid   int
		User  User
	}
	avail, order := a.Model.BookingDB.VenueReserve[vID].GetDate(model.MinDate, model.MaxDate)
	d := data{
		Venue: venue,
		Avail: avail,
		Order: order,
		Vid:   vID,
		User:  a.getUser(res, req),
	}
	a.Logging.Trace.Println("Booking attempt from ", req.UserAgent())
	a.Template.ExecuteTemplate(res, "book.html", &d)
}

// ConfirmBook :
func (a *Ctl) ConfirmBook(res http.ResponseWriter, req *http.Request) {

	vIDs, ok1 := req.URL.Query()["venueId"]
	dates, ok2 := req.URL.Query()["date"]
	times, ok3 := req.URL.Query()["time"]

	if !ok1 || len(vIDs[0]) < 1 || !ok2 || len(dates[0]) < 1 || !ok3 || len(times[0]) < 1 {
		a.Logging.Warning.Println("Incorrect booking parameter from ", req.UserAgent())
		http.Redirect(res, req, "/browse", http.StatusSeeOther)
		return
	}
	u := a.getUser(res, req)
	if !a.alreadyLoggedIn(req) {
		a.Logging.Warning.Println("Unauthorised booking attempt from ", req.UserAgent())
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	vID, valid1 := strconv.Atoi(vIDs[0])
	date, valid2 := strconv.Atoi(dates[0])
	time, valid3 := strconv.Atoi(times[0])
	if valid1 != nil || valid2 != nil || valid3 != nil {
		a.Logging.Error.Println("Error converting parameter to integer from ", req.UserAgent())
		http.Redirect(res, req, "/browse", http.StatusSeeOther)
		return
	}
	venue := a.Model.VenueDB.Venues[vID]
	type data struct {
		User  User
		Venue model.Venue
		Date  int
		Time  string
	}
	timing := "nil time"
	switch time {
	case 1:
		timing = "Morning"
	case 2:
		timing = "Afternoon"
	case 3:
		timing = "Evening"
	}
	d := data{
		User:  u,
		Venue: venue,
		Date:  date,
		Time:  timing,
	}
	if req.Method == http.MethodPost {
		username := u.Username
		datetime := date*10 + time
		// check if user exist with username
		bookingID := a.Model.BookingDB.Reserve(vID, datetime, username)
		u.Bookings = append(u.Bookings, bookingID)
		a.Logging.Info.Println("Booking confirmed from ", req.UserAgent())
		a.Users[u.Username] = u
		http.Redirect(res, req, "/book?venueId="+fmt.Sprint(vID), http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "confirmBook.html", &d)
}

// Find : check for existing string in string slice
func Find(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// ViewBook :
func (a *Ctl) ViewBook(res http.ResponseWriter, req *http.Request) {
	//a.Model.BookingDB.VenueReserve
	mapping := a.Model.VenueDB.VenueMap
	type pageData struct {
		User   User
		Venues map[string]model.Venue
		BkData map[string][]Booking
		Order  []string
	}
	data := pageData{
		User:   a.getUser(res, req),
		Venues: make(map[string]model.Venue),
		BkData: make(map[string][]Booking),
	}
	if !a.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	for _, i := range data.User.Bookings {
		mbooking := a.Model.BookingDB.Bookings[i]
		booking := convertBooking(*mbooking, mapping)
		venueName := booking.VenueName
		if Find(data.Order, venueName) {
			data.BkData[venueName] = append(data.BkData[venueName], booking)
		} else {
			data.Order = append(data.Order, venueName)
			data.Venues[venueName] = a.Model.VenueDB.Venues[booking.VenueID]
			data.BkData[venueName] = append(data.BkData[venueName], booking)
		}

	}
	sort.Strings(data.Order)

	a.Template.ExecuteTemplate(res, "viewBooking.html", data)
}

func removeInt(ints []int, in int) []int {
	var i int
	for i = 0; i < len(ints); i++ {
		if ints[i] == in {
			break
		}
	}
	copy(ints[i:], ints[i+1:])
	ints[len(ints)-1] = 0
	ints = ints[:len(ints)-1]
	return ints
}

// DeleteBook : cancellation
func (a *Ctl) DeleteBook(res http.ResponseWriter, req *http.Request) {

	u := a.getUser(res, req)
	if !a.alreadyLoggedIn(req) {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}
	bIDs, ok := req.URL.Query()["bID"]

	if !ok || len(bIDs[0]) < 1 {
		http.Redirect(res, req, "/viewBook", http.StatusSeeOther)
		return
	}

	bID, valid := strconv.Atoi(bIDs[0])
	booking := a.Model.BookingDB.Bookings[bID]

	if valid != nil || booking.User != u.Username || u.Username != "admin" {
		userE := wrongUserError{
			user1: u.Username,
			user2: booking.User}
		a.Logging.Error.Println(userE.Error(), " from ", req.UserAgent())
		http.Redirect(res, req, "/viewBook", http.StatusSeeOther)
		return
	}
	type pageData struct {
		User    User
		Booking Booking
	}

	d := pageData{
		User:    u,
		Booking: convertBooking(*booking, a.Model.VenueDB.VenueMap),
	}

	if req.Method == http.MethodPost {
		IDBook := req.FormValue("IDBook")
		fmt.Println(IDBook)
		bID, _ := strconv.Atoi(IDBook)
		booking := a.Model.BookingDB.Bookings[bID]
		if booking.User != u.Username {
			a.Logging.Warning.Println("Invalid credential POST booking deletion attempt from ",
				req.UserAgent())
			http.Redirect(res, req, "/viewBook", http.StatusSeeOther)
			return
		}

		// check if user exist with username
		a.Model.BookingDB.DelReserve(booking.VenueID, booking.Datetime)
		u.Bookings = removeInt(u.Bookings, bID)
		a.Users[u.Username] = u
		a.Logging.Info.Println("Booking cancelled from ", req.UserAgent())
		http.Redirect(res, req, "/viewBook", http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "deleteBooking.html", &d)
}

// AddVenue :
func (a *Ctl) AddVenue(res http.ResponseWriter, req *http.Request) {

	u := a.getUser(res, req)
	if u.Username != "admin" {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	type pageData struct {
		User User
	}
	d := pageData{
		User: u,
	}
	if req.Method == http.MethodPost {
		vName := santizeString(req.FormValue("name"))
		vKind := santizeString(req.FormValue("kind"))
		vLocation := santizeString(req.FormValue("location"))
		vDesc := santizeString(req.FormValue("desc"))
		vCap, _ := strconv.Atoi(req.FormValue("capacity"))
		// check if user exist with username
		for _, v := range a.Model.VenueDB.VenueMap {
			if v == vName {
				http.Redirect(res, req, "/addVenue", http.StatusSeeOther)
				return
			}
		}
		a.Model.AddVenue(model.Venue{
			Capacity: vCap,
			Kind:     vKind,
			Location: vLocation,
			Name:     vName,
			Desc:     vDesc,
		})
		a.Logging.Info.Println("Venue added from ", req.UserAgent())
		http.Redirect(res, req, "/browse", http.StatusSeeOther)
		return
	}
	a.Template.ExecuteTemplate(res, "addVenue.html", &d)
}
