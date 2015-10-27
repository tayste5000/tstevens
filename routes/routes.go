package routes

import(
	//stdlib
	"net/http"

	//third party
	"github.com/zenazn/goji/web"

	//internal
	"github.com/tayste5000/tstevens/templates"
	"github.com/tayste5000/tstevens/routes/projects/param2drive"
)

func Add(mux *web.Mux) {
	
	/* Endpoint to handler config */
	mux.Get("/", home)
	mux.Get("/projects", projects)
	mux.Handle("/projects/p2drive/*", param2drive.AddRoutes("/projects/p2drive"))
	mux.Get("/contact", contact)
	mux.Get("/site-map", siteMap)
	mux.Get("/faq", faq)

}

/* Handlers */

func home(c web.C, w http.ResponseWriter, r *http.Request){

	if err := templates.Render(w, "home.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
	
func projects(c web.C, w http.ResponseWriter, r *http.Request){

	if err := templates.Render(w, "projects.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func contact(c web.C, w http.ResponseWriter, r *http.Request){

	if err := templates.Render(w, "contact.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func siteMap(c web.C, w http.ResponseWriter, r *http.Request){

	if err := templates.Render(w, "site-map.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func faq(c web.C, w http.ResponseWriter, r *http.Request){

	if err := templates.Render(w, "faq.html", nil); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}