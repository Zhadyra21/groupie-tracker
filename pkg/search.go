package groupie

import (
	"strconv"
	"strings"
)

type NoResultFound struct {
	Message string
}

func ArtistFinder(search string) []Artist {
	var data []Artist
	searching := strings.Split(search, " -> ")
	found := false
	for i := range DATA { //name, members
		if strings.Contains(DATA[i].Name, searching[0]) {
			data = append(data, DATA[i])
			found = true
			continue
		}
		for _, v := range DATA[i].Members {
			if strings.Contains(v, searching[0]) {
				data = append(data, DATA[i])
				found = true
				break
			}
		}
	}
	if found {
		return data
	}
	for i := range DATA { //first album, creation date
		if strings.Contains(DATA[i].FirstAlbum, searching[0]) {
			data = append(data, DATA[i])
			found = true
			continue
		}
		creationDate := strconv.Itoa(DATA[i].CreationDate)
		if creationDate == searching[0] {
			data = append(data, DATA[i])
			found = true
		}
	}
	if found {
		return data
	}
	for i := range DATA { //location
		for key := range DATA[i].DatesLocations {
			if strings.Contains(key, searching[0]) {
				data = append(data, DATA[i])
				found = true
				break
			}
		}
	}
	if found {
		return data
	}
	for i := range DATA { //concert dates
		for _, value := range DATA[i].DatesLocations {
			for _, dates := range value {
				if strings.Contains(dates, searching[0]) {
					data = append(data, DATA[i])
					goto next
				}
			}
		}
	next:
	}
	return nil
}
