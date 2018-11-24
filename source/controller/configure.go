package controller

import (
	"fmt"
	"gowebapp/source/model/validate"
	"gowebapp/source/view"
	"net/http"
	"strings"
)

func configureGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "configure", nil)
}

func configurePOST(w http.ResponseWriter, r *http.Request) {
	v := view.New("configure")
	v.Vars["Configname"] = strings.TrimSpace(r.FormValue("configname"))
	v.Vars["Header"] = r.FormValue("header")
	v.Vars["Footer"] = r.FormValue("footer")
	if urlformat := r.FormValue("urlformat"); len(urlformat) > 0 {
		v.Vars["Urlformat"] = urlformat
	} else {
		v.Vars["Urlformat"] = "{url}"
	}

	// Validate configname
	errCount := 0
	if ok, err := validate.Configname(v.Vars["Configname"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["ConfignameModal"] = "This field is required"
		}
	}
	if ok, err := validate.Urlformat(v.Vars["Urlformat"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["UrlformatModal"] = "This field is required"
		case validate.ErrInvalidFormat:
			v.Vars["UrlformatModal"] = "This field must contain {url}"
		}
	}

	if errCount != 0 {
		v.AddFlash(fmt.Sprintf("Error: %d inputs are invalid", errCount), view.FlashError)
		v.Render(w)
	} else {
		v.Render(w)
	}
}

func configuredHandler(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "configure", nil)
}
