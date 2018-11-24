package view

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func init() {
	gob.Register(&Flash{})
}

// Flash Message
type Flash struct {
	Message string
	Class   string
}

// View attributes
type View struct {
	// Name is template's name
	Name string
	// Ext is template's file extension (defaults to .html)
	Ext string
	// Vars are view variables passed to template renderer
	Vars map[string]interface{}
}

var (
	// FlashError is a bootstrap class
	FlashError = "alert-danger"
	// FlashSuccess is a bootstrap class
	FlashSuccess = "alert-success"
	// FlashNotice is a bootstrap class
	FlashNotice = "alert-info"
	// FlashWarning is a bootstrap class
	FlashWarning = "alert-warning"

	// templates is a html template collection
	templates = make(map[string]*template.Template)

	rootTemplate = "base.html"
	ignoredFiles = []string{"blank.html", "base.html"}
	flashKey     = "_flashes"
	commonExt    = ".html"
)

// New returns a new View
func New(name string) *View {
	v := &View{}
	v.Name = name
	v.Ext = commonExt
	v.Vars = make(map[string]interface{})
	v.Vars[flashKey] = make([]Flash, 0)

	return v
}

// AddFlash adds a flash message to a view
func (v *View) AddFlash(message, class string) {
	v.Vars[flashKey] = append(v.Vars[flashKey].([]Flash), Flash{Message: message, Class: class})
}

// PeekFlashes returns flashes bound to current view
func (v *View) PeekFlashes() []Flash {
	return v.Vars[flashKey].([]Flash)
}

// Set is used to set a variable.
// Same can be accomplished by using v.Vars[key] = val
func (v *View) Set(key string, val interface{}) {
	v.Vars[key] = val
}

// Render renders template to the writer
func (v *View) Render(w http.ResponseWriter) {
	t, ok := templates[v.Name+v.Ext]
	if !ok {
		http.Error(w, "This page could not be loaded", http.StatusInternalServerError)
		log.Println("Template", v.Name, "does not exist")
		return
	}
	err := t.Execute(w, v.Vars)
	if err != nil {
		http.Error(w, "This page could not be loaded", http.StatusInternalServerError)
		log.Println("Template", v.Name, "Error:", err)
	}
}

// RenderTemplate renders a template (HTML)
func RenderTemplate(w http.ResponseWriter, tname string, i interface{}) {
	t, ok := templates[tname+commonExt]
	if !ok {
		http.Error(w, "This page could not be loaded", http.StatusInternalServerError)
		log.Println("Template", tname, "does not exist")
		return
	}
	err := t.Execute(w, i)
	if err != nil {
		http.Error(w, "This page could not be loaded", http.StatusInternalServerError)
		log.Println(err)
	}
}

// LoadTemplates loads templates from specified path.
// Panics if error is encountered.
func LoadTemplates(path string, root ...string) {
	// Add trailing "/" to path
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	// Initialize root path
	if len(root) == 0 {
		root = append(root, path)
	}
	// List directory files
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		ok := true
		for _, ignored := range ignoredFiles {
			if ignored == file.Name() {
				ok = false
				break
			}
		}
		if ok {
			if file.IsDir() {
				LoadTemplates(path+file.Name(), root[0])
			} else {
				// Cuts off root path from left side of the string
				tname := (path + file.Name())[len(root[0]):]
				templates[tname] =
					template.Must(template.ParseFiles(
						root[0]+rootTemplate, path+file.Name()))
			}
		}
	}
}
