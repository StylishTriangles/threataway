package email

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"gopkg.in/gomail.v2"
)

var (
	// d is the default dialer
	d *gomail.Dialer
	// templates is a html template collection
	templates    = make(map[string]*template.Template)
	ignoredFiles = []string{}
	from         = "Admin <example@example.com>"
)

// SMTPInfo holds details of SMTP server
type SMTPInfo struct {
	Username string
	Password string
	Host     string
	Port     int
	From     string
}

// Email is an email
type Email struct {
	// Tmpl is template's name
	Tmpl string
	// Vars is used to send data to templates
	Vars map[string]interface{}

	message *gomail.Message
}

// New creates new email
func New() *Email {
	e := &Email{
		Tmpl:    "",
		Vars:    make(map[string]interface{}),
		message: gomail.NewMessage(),
	}
	e.message.SetHeader("From", from)
	return e
}

// Send sends the email
func (e *Email) Send() error {
	// Use template
	if e.Tmpl != "" {
		var bb bytes.Buffer
		if err := templates[e.Tmpl].Execute(&bb, e.Vars); err != nil {
			return err
		}
		// In case we want to additionally support plaintext messages in the future
		e.message.AddAlternative("text/html", bb.String())
	}
	return d.DialAndSend(e.message)
}

// SetRecipient updates the Recipient of outgoing email
func (e *Email) SetRecipient(s string) {
	e.message.SetHeader("To", s)
}

// SetSender updates email's sender
func (e *Email) SetSender(s string) {
	e.message.SetHeader("From", s)
}

// SetSubject updates email's subject
func (e *Email) SetSubject(s string) {
	e.message.SetHeader("Subject", s)
}

// SetTemplate updates which template will be used
func (e *Email) SetTemplate(name string) {
	e.Tmpl = name
}

// SetVar sets template variable
func (e *Email) SetVar(key string, val interface{}) {
	e.Vars[key] = val
}

// Configure adds the settings for the SMTP server
func Configure(c *SMTPInfo) {
	d = &gomail.Dialer{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
	}
	from = c.From
}

// LoadTemplates loads templates from specified path.
// Panics if any errors are encountered.
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
					template.Must(template.ParseFiles(path + file.Name()))
			}
		}
	}
}
