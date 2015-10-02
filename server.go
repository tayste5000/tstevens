package main

import(
	//stdlib
	"net/http"

	//third party
	"github.com/zenazn/goji"

	//internal
	"github.com/tayste5000/tstevens/routes"
)

func main() {

	/* Make all files in "public" folder available */
	fs := http.FileServer(http.Dir("public"))
	goji.Get("/public/*", http.StripPrefix("/public/", fs))

	/* attach imported routing info to mux */
	routes.Add(goji.DefaultMux)

	/* voila */
	goji.Serve()
}