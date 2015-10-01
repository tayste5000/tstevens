package routes

import(
	//stdlib
	"net/http"

	//third party
	"github.com/zenazn/goji/web"

	//internal
	"github.com/tayste5000/go-seed/templates"
)

func Add(mux *web.Mux) {
	
	/* Endpoint to handler config */
	mux.Get("/", home)
	mux.Get("/home", about)
	mux.Use(mux.Router)

}

/* Handlers */

func home(c web.C, w http.ResponseWriter, r *http.Request){

	// ourname := make(map[string]interface{})
 
	// ourname["First"] = "Billy"
	// ourname["Last"] = "Bob"

	if err := templates.Render(w, "home.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
	
func about(c web.C, w http.ResponseWriter, r *http.Request){

	// ourname := make(map[string]interface{})

	// ourname["First"] = "Billy"
	// ourname["Last"] = "Bob"

	if err := templates.Render(w, "about.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}