package templates

import(
	//stdlib
	"log"
	"path/filepath"
	"html/template"
	"net/http"
	"fmt"

	//third party


	//internal
)

var Map map[string]*template.Template

func init() {
	if Map == nil{
		Map = make(map[string]*template.Template)
	}

	templatesDir := "templates/"

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil{
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "includes/*.html")
	if err != nil{
		log.Fatal(err)
	}	

	for _,layout := range layouts{
		files := append(includes, layout)
		Map[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

func Render(w http.ResponseWriter, name string, data map[string]interface{}) error {

	tmpl, ok := Map[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	tmpl.ExecuteTemplate(w, "base", data)

	return nil
}