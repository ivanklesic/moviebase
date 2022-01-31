package html

import (
	"fmt"
	"net/http"
	"text/template"
)

type TemplateData map[string]interface{}

func Render(w http.ResponseWriter, r *http.Request, name string, data TemplateData) {
	tpl, err := template.New("base.html").Funcs(template.FuncMap{
		"containsId": containsIdHelper,
		"getIndexOfArray": getIndexOfArrayHelper,
	}).ParseFiles("html/views/base.html", "html/views/"+name+".html")
	
	if err != nil {
		err := fmt.Errorf("couldn't parse template files: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

	if err := tpl.Execute(w, data); err != nil {
		return
	}
}
