package controller

import (
	"fmt"
	"gowebapp/source/model/validate"
	"gowebapp/source/view"
	"net/http"
	"strings"
)

func templateGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "template", nil)
}

func templatePOST(w http.ResponseWriter, r *http.Request) {
	v := view.New("template")
	v.Vars["Templatename"] = strings.TrimSpace(r.FormValue("templatename"))
	v.Vars["Header"] = r.FormValue("header")
	v.Vars["Footer"] = r.FormValue("footer")
	if urlformat := r.FormValue("urlformat"); len(urlformat) > 0 {
		v.Vars["Urlformat"] = urlformat
	} else {
		v.Vars["Urlformat"] = "{URL}"
	}

	// Validate templatename
	errCount := 0
	if ok, err := validate.Templatename(v.Vars["Templatename"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["TemplatenameModal"] = "This field is required"
		}
	}
	if ok, err := validate.Urlformat(v.Vars["Urlformat"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["UrlformatModal"] = "This field is required"
		case validate.ErrInvalidFormat:
			v.Vars["UrlformatModal"] = "This field must contain {URL}"
		}
	}

	if errCount != 0 {
		v.AddFlash(fmt.Sprintf("Error: %d inputs are invalid", errCount), view.FlashError)
		v.Render(w)
	} else {
		v.Render(w)
	}
}
