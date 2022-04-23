package groupie

import (
	"fmt"
	"net/http"
	"strconv"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
		}
	}()
	if e := fn(w, r); e != nil {
		w.WriteHeader(e.Code)
		if err := ExecuteTmpl(w, "err", &e); err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", 500)
		}
		if e.Error != nil {
			fmt.Println(e.Error)
		}
	}
}

func Init() {
	http.Handle("/", appHandler(Home_page))
	http.Handle("/artists/", appHandler(Artist_page))
	http.Handle("/filter/", appHandler(Filter_page))
}

func Home_page(w http.ResponseWriter, r *http.Request) *appError {
	if r.URL.Path != "/" {
		return &appError{nil, "Page Not Found", 404}
	}
	if r.Method != http.MethodGet {
		return &appError{nil, "Method Not Allowed", 405}
	}
	if err := GetData(); err != nil {
		return &appError{err, "Internal Server Error", 500}
	}
	if err := ExecuteTmpl(w, "index", DATA); err != nil {
		return &appError{err, "Internal Server Error", 500}
	}
	return nil
}

func Artist_page(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		return &appError{nil, "Method Not Allowed", 405}
	}
	if err := GetData(); err != nil {
		return &appError{err, "Internal Server Error", 500}
	}
	id := r.URL.Path[9:]
	if id == "" {
		if err := ExecuteTmpl(w, "index", DATA); err != nil {
			return &appError{err, "Internal Server Error", 500}
		}
		return nil
	}
	artistId, _ := strconv.Atoi(id)
	if artistId < 1 || artistId > 52 {
		return &appError{nil, "Page Not Found", 404}
	}
	if err := ExecuteTmpl(w, "artist", DATA[artistId-1]); err != nil {
		return &appError{err, "Internal Server Error", 500}
	}
	return nil
}

func Filter_page(w http.ResponseWriter, r *http.Request) *appError {
	if r.Method != http.MethodGet {
		return &appError{nil, "Method Not Allowed", 405}
	}
	r.ParseForm()
	if r.Form.Has("search_input") {
		newData := ArtistFinder(r.Form.Get("search_input"))
		if err := ExecuteTmpl(w, "index", newData); err != nil {
			return &appError{err, "Internal Server Error", 500}
		}
		return nil
	}
	filter := &DataFilter{
		CreationDate: DataFilterParams{
			From: r.Form.Get("creation_date[from]"),
			To:   r.Form.Get("creation_date[to]"),
		},
		FirstAlbum: DataFilterParams{
			From: r.Form.Get("first_album_date[from]"),
			To:   r.Form.Get("first_album_date[to]"),
		},
		MembersNum: DataFilterParams{
			From: r.Form.Get("num_of_members[from]"),
			To:   r.Form.Get("num_of_members[to]"),
		},
		Location: r.Form.Get("location"),
	}
	filter.ArtistFilter()
	if err := ExecuteTmpl(w, "index", filter.Result); err != nil {
		return &appError{err, "Internal Server Error", 500}
	}
	return nil
}
