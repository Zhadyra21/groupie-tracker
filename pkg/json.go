package groupie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

var (
	Artists   = "https://groupietrackers.herokuapp.com/api/artists"
	Locations = "https://groupietrackers.herokuapp.com/api/locations"
	Dates     = "https://groupietrackers.herokuapp.com/api/dates"
	Relations = "https://groupietrackers.herokuapp.com/api/relation"
)

var DATA []Artist

type Artist struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	DatesLocations map[string][]string
}

type Relation struct {
	Index []Concert `json:"index"`
}

type Concert struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

func ParseBody(url string, x interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("no response from API")
	}
	defer r.Body.Close()
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, x); err != nil {
			return fmt.Errorf("couldn't unmarshal the data")
		}
	}
	return nil
}

func GetData() error {
	if err := ParseBody(Artists, &DATA); err != nil {
		return err
	}
	rel := &Relation{}
	if err := ParseBody(Relations, &rel); err != nil {
		return err
	}
	for i, v := range rel.Index {
		DATA[i].DatesLocations = v.DatesLocations
	}
	return nil
}

var templates, _ = template.ParseGlob("templates/*.html")

func ExecuteTmpl(w http.ResponseWriter, tmpl string, data interface{}) error {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		return fmt.Errorf("an error occured while reaching the %q tmpl", tmpl)
	}
	return nil
}
