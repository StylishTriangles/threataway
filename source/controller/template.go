package controller

import (
	"fmt"
	"gowebapp/source/model/templates"
	"gowebapp/source/model/validate"
	"gowebapp/source/view"
	"net/http"
	"strings"
)

func templateGET(w http.ResponseWriter, r *http.Request) {
	view.RenderTemplate(w, "template", nil)
}

func templatePOST(w http.ResponseWriter, r *http.Request) {
	//TODO: security
	v := view.New("template")
	v.Vars["Name"] = strings.TrimSpace(r.FormValue("name"))
	v.Vars["Header"] = r.FormValue("header")
	v.Vars["Footer"] = r.FormValue("footer")
	if urlTemplate := r.FormValue("urlTemplate"); len(urlTemplate) > 0 {
		v.Vars["UrlTemplate"] = urlTemplate
	} else {
		v.Vars["UrlTemplate"] = "{URL}"
	}

	// Validate name
	errCount := 0
	if ok, err := validate.Name(v.Vars["Name"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["NameModal"] = "This field is required"
		}
	}
	if ok, err := validate.UrlTemplate(v.Vars["UrlTemplate"].(string)); !ok {
		errCount++
		switch err {
		case validate.ErrEmpty:
			v.Vars["UrlTemplateModal"] = "This field is required"
		case validate.ErrInvalidFormat:
			v.Vars["UrlTemplateModal"] = "This field must contain {URL}"
		}
	}

	if errCount != 0 {
		v.AddFlash(fmt.Sprintf("Error: %d inputs are invalid", errCount), view.FlashError)
		v.Render(w)
	} else {
		// TODO: this is an ugly hack just to show the user how's his input looking like parsed
		v2 := view.New("template_added")
		v2.Vars = v.Vars
		if urls := strings.Split(v.Vars["UrlTemplate"].(string), "{URL}"); len(urls) == 2 {
			prefix := urls[0]
			suffix := urls[1]
			preview := ""
			for _, url := range []string{"onet.pl", "wykop.pl", "google.com", "solve.edu.pl", "31.185.104.19", "facebook.com"} {
				preview += prefix + url + suffix + "</br>"
			}
			v2.Vars["Urls"] = preview
			templates.CreateNewTemplate(v.Vars["Name"].(string), v.Vars["Header"].(string), v.Vars["Footer"].(string), v.Vars["UrlTemplate"].(string))
		}
		v2.Render(w)
	}
}
