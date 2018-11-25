package controller

import (
	"net/http"
	"regexp"

	"gowebapp/source/shared/database"

	"github.com/gorilla/mux"
)

func deployedListHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(w, r) {
		return
	}

	m := mux.Vars(r)

	stmt, err := database.DB.Prepare(`SELECT urls.domain FROM listlists LEFT JOIN deployments ON listlists.idList = deployments.listID LEFT JOIN urls ON listlists.idURL = urls.idUrl WHERE deployments.url = ?  `)
	if err != nil {
		return
	}
	defer stmt.Close()

	var urlList []string

	rows, err := stmt.Query(m["deploy_name"])

	if err != nil {
		return
	}

	for rows.Next() {
		url := ""
		rows.Scan(&url)

		urlList = append(urlList, url)
	}

	stmt, err = database.DB.Prepare(`SELECT templates.header, templates.footer, templates.urlTemplate FROM deployments LEFT JOIN templates ON deployments.templateID = templates.templateID WHERE deployments.url = ?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var header string
	var footer string
	var template string

	stmt.QueryRow(m["deploy_name"]).Scan(&header, &footer, &template)

	w.Write([]byte(header))
	for _, url := range urlList {
		var re = regexp.MustCompile(`({URL})`)
		s := re.ReplaceAllString(template, url)

		w.Write([]byte(s))
	}
	w.Write([]byte(footer))
}
