package groupie

import (
	"strconv"
	"strings"
	"time"
)

type DataFilter struct {
	CreationDate DataFilterParams
	FirstAlbum   DataFilterParams
	MembersNum   DataFilterParams
	Location     string
	Result       []Artist
}

type DataFilterParams struct {
	From string
	To   string
}

func (v DataFilterParams) checkValues() (a int, b int) {
	if v.From != "" {
		a, _ = strconv.Atoi(v.From)
	} else {
		a = 0
	}
	if v.To != "" {
		b, _ = strconv.Atoi(v.To)
	} else {
		t := time.Now()
		b = t.Year()
	}
	return
}

func (f *DataFilter) ArtistFilter() {
	from_crD, to_crD := f.CreationDate.checkValues()
	from_firstA, to_firstA := f.FirstAlbum.checkValues()
	from_mNum, to_mNum := f.MembersNum.checkValues()
	for i := range DATA {
		if DATA[i].CreationDate < from_crD || DATA[i].CreationDate > to_crD {
			continue
		}
		firstA, _ := strconv.Atoi(DATA[i].FirstAlbum[6:])
		if firstA < from_firstA || firstA > to_firstA {
			continue
		}
		mNum := len(DATA[i].Members)
		if mNum < from_mNum || mNum > to_mNum {
			continue
		}
		foundLocation := false
		for key := range DATA[i].DatesLocations {
			location := strings.ReplaceAll(f.Location, ", ", "-")
			location = strings.ReplaceAll(location, " ", "_")
			if strings.Contains(key, strings.ToLower(location)) {
				foundLocation = true
			}
		}
		if !foundLocation {
			continue
		}
		f.Result = append(f.Result, DATA[i])
	}
}
