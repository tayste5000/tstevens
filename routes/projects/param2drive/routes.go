package param2drive

import(
	//stdlib
	"os"
  "path"
  "io/ioutil"
	"os/exec"
	"encoding/json"
	"log"
  "encoding/gob"
  "strings"
  "net/http"
  "strconv"
  "regexp"
  "fmt"

	//third party
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/drive/v2"
  "google.golang.org/api/plus/v1"
  "code.google.com/p/xsrftoken"
  "github.com/gorilla/sessions"
  "github.com/nu7hatch/gouuid"
  "github.com/zenazn/goji/web"
  "github.com/zenazn/goji/web/middleware"
  "golang.org/x/net/context"
  "github.com/fatih/structs"
  "github.com/gorilla/schema"

	//internal
	"github.com/tayste5000/tstevens/templates"
)

type inputData struct{
	Sequence		string
	Options			[]string
}

type outputData struct{
	Name					string
	Sequence			string
	Range					string
	MW						float32
	EC280					float32
	PI 						float32
	Instability		float32
}

func (od outputData) text() string{
  output := ""

  if od.Name != "" {output += "Name: " + od.Name}
  if od.Range != "" {output += "\n\nRange: " + od.Range}
  if od.Sequence != "" {output += "\n\nSequence: \n\n" + od.Sequence}
  if od.MW != 0 {output += "\n\nMW: " + fmt.Sprint(od.MW) + " Da"}
  if od.EC280 != 0 {output += "\n\nEC280: " + fmt.Sprint(od.EC280) + " 1/(M*cm)"}
  if od.PI != 0 {output += "\n\npI: " + fmt.Sprint(od.PI)}
  if od.Instability != 0 {output += "\n\nInstability Index: " + fmt.Sprint(od.Instability)}

  return output
}

var store *sessions.CookieStore
var secretKey string
var redirectUrl string

func init() {
  err := godotenv.Load()

  if err != nil {
    log.Fatal("Error loading .env file")
  }

  secretKey = os.Getenv("SECRET_KEY")
  redirectUrl = os.Getenv("P2D_REDIRECT")
  
  log.Print(redirectUrl)
  log.Print(redirectUrl)

  store = sessions.NewCookieStore([]byte(secretKey))

  gob.Register(oauth2.Token{})

}

func intro(c web.C, w http.ResponseWriter, r *http.Request){
  if err := templates.Render(w, "projects-param2drive.html", nil); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func form(c web.C, w http.ResponseWriter, r *http.Request){

  b, err := ioutil.ReadFile("client_secret.json")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config, err := google.ConfigFromJSON(b, drive.DriveScope, drive.DriveAppdataScope, drive.DriveAppsReadonlyScope, drive.DriveFileScope)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config.RedirectURL = redirectUrl

  ctx := context.Background()

  authToken := c.Env["access-token"].(oauth2.Token)

  plusClient := config.Client(ctx, &authToken)

  plusService, err := plus.New(plusClient);

  user, err := plusService.People.Get("me").Do()
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  templateData := make(map[string]interface{})
  templateData["Image"] = user.Image.Url
  templateData["Name"] = user.DisplayName

  session, err := store.Get(r, "p2drive")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return
  }   

  if flashes := session.Flashes("error"); len(flashes) > 0{
    templateData["ErrorFlash"] = flashes[0]
    log.Print(flashes)
    log.Print(len(flashes))
    for _,val := range flashes {
      log.Print(val)
    }
  }

  if flashes := session.Flashes("success"); len(flashes) > 0{
    templateData["SuccessFlash"] = flashes[0]
    log.Print(flashes)
    log.Print(len(flashes))
    for _,val := range flashes {
      log.Print(val)
    }
  }

  session.Save(r,w)
	
	if err := templates.Render(w, "projects-param2drive-form.html", templateData); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func compute(c web.C, w http.ResponseWriter, r *http.Request){

	err := r.ParseForm(); if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

  var inputError []string

  aa_pattern := "^([AaCcDdEeFfGgHhIiKkLlMmNnPpQqRrSsTtUuVvWwYy]+)$"
  range_pattern := "^([0-9]+)-([0-9]+)$"

  aa_match, err := regexp.MatchString(aa_pattern, r.Form.Get("Sequence"))
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  range_match, err := regexp.MatchString(range_pattern, r.Form.Get("Range"))
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  if r.Form.Get("Name") == ""{
    inputError = append(inputError, "A name is required")
  }

  if !aa_match {
    inputError = append(inputError, "The sequence must be a valid amino acid sequence")
  }

  if !range_match{
    inputError = append(inputError, "The amino acid range must be of the form #-# eg. 21-63")
  } else {
    rangeArray := strings.Split(r.Form.Get("Range"), "-")
    rangeStart, err := strconv.Atoi(rangeArray[0])
    rangeStop, err := strconv.Atoi(rangeArray[1])

    if err != nil{
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    rangeLen := rangeStop - rangeStart + 1

    if rangeLen != len(r.Form.Get("Sequence")){
      inputError = append(inputError, "The sequence length does not match the range length")
    }
  }

  if r.Form.Get("Features") == "" {
    inputError = append(inputError, "At least one feature needs to be selected")
  }

  if len(inputError) > 0 {

    session, err := store.Get(r, "p2drive")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } 

    session.AddFlash(inputError, "error")
    session.Save(r,w)

    http.Redirect(w, r, "/projects/p2drive/input", 302)
    return
  }

  p2dInput := inputData{
    Sequence: r.Form.Get("Sequence"),
    Options: r.Form["Features"],
  } 

	jsonInput, err := json.Marshal(p2dInput); if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonOutput, err := exec.Command("python", "scripts/compute_params.py", string(jsonInput)).Output(); if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	p2dOutput := outputData{
		Name: r.Form.Get("Name"),
		Sequence: r.Form.Get("Sequence"),
		Range: r.Form.Get("Range"),
	}

	err = json.Unmarshal(jsonOutput, &p2dOutput); if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	params := structs.Map(p2dOutput)

	if err := templates.Render(w, "projects-param2drive-results.html", params); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submit(c web.C, w http.ResponseWriter, r *http.Request){

	err := r.ParseForm(); if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	toSubmit := new(outputData)
	decoder := schema.NewDecoder()
	decoder.Decode(toSubmit,r.Form)

  b, err := ioutil.ReadFile("client_secret.json")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config, err := google.ConfigFromJSON(b, drive.DriveScope, drive.DriveAppdataScope, drive.DriveAppsReadonlyScope, drive.DriveFileScope)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config.RedirectURL = redirectUrl

  ctx := context.Background()

  authToken := c.Env["access-token"].(oauth2.Token)

  DriveClient := config.Client(ctx, &authToken)

  DriveService, err := drive.New(DriveClient);

  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var parentId string

  //check if folder "param2drive_files" exists
  driveList, err := DriveService.Files.List().Q("title=\"param2drive_files\"").Do();
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  //TODO: stop doing bootleg shit like this
  //TODO: if there is a deleted param2drive folder this fails
  if len(driveList.Items) > 0 {
    parentId = driveList.Items[0].Id
  } else {
    newFolder := &drive.File{
      MimeType: "application/vnd.google-apps.folder",
      Title: "param2drive_files",
    }
    newFolder, err = DriveService.Files.Insert(newFolder).Do()
    if err != nil{ log.Fatal(err) }
    parentId = newFolder.Id
  }

  var parentFolder = make([]*drive.ParentReference, 1)
  parentFolder[0] = &drive.ParentReference{Id: parentId}

  //title file protein name-range-date
  newFile := drive.File{
    MimeType: "text/plain",
    Title: "p2drive.txt",
    Parents: parentFolder,
  }

  contents := strings.NewReader(toSubmit.text())

  createdFile, err := DriveService.Files.Insert(&newFile).Media(contents).Do()
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  session, err := store.Get(r, "p2drive")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return
  } 

  log.Print(createdFile)
  log.Print(createdFile.DefaultOpenWithLink)
  log.Print(createdFile.AlternateLink)
  log.Print(createdFile.EmbedLink)
  log.Print(createdFile.ExportLinks)

  session.AddFlash(createdFile.AlternateLink, "success")
  session.Save(r,w)

  http.Redirect(w,r,"/projects/p2drive/input", 302)
}

func auth(c web.C, w http.ResponseWriter, r *http.Request){

  log.Print(r.URL.String())

	session, err := store.Get(r, "p2drive")
  if err != nil {
      http.Error(w, err.Error(), 500)
      return
  } 

  csrfResponse := r.URL.Query().Get("state")
  code := r.URL.Query().Get("code")
  error := r.URL.Query().Get("error")

  if error == "access_denied" {
    http.Redirect(w,r,"/projects/p2drive/intro", 302)
  }

  b, err := ioutil.ReadFile("client_secret.json")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config, err := google.ConfigFromJSON(b,
    drive.DriveScope,
    drive.DriveAppdataScope,
    drive.DriveAppsReadonlyScope,
    drive.DriveFileScope,
    plus.PlusMeScope,
    plus.UserinfoProfileScope)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  config.RedirectURL = redirectUrl

  if csrfResponse == "" || code == ""{
    csrfId, err := uuid.NewV4(); if err != nil{
      http.Error(w, "couldn't generate uuid", http.StatusInternalServerError)
    }

    session.Values["id"] = csrfId.String()

    csrfToken := xsrftoken.Generate(secretKey, session.Values["id"].(string), "/projects/p2drive/auth")
    session.Values["csrf"] = csrfToken
    session.Save(r,w)

    authURL := config.AuthCodeURL(csrfToken, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

    http.Redirect(w,r,authURL, 302)
    return
  }

  if session.Values["csrf"] != csrfResponse {
    http.Redirect(w, r, "/projects", 400)
  }

  if !xsrftoken.Valid(session.Values["csrf"].(string), secretKey, session.Values["id"].(string), "/projects/p2drive/auth") {
    http.Redirect(w, r, "/projects", 400)
  }

  tok, err := config.Exchange(oauth2.NoContext, code)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  session.Values["access-token"] = tok
  session.Save(r,w)
  http.Redirect(w,r,"/projects/p2drive/input", 302)
}

func logout(c web.C, w http.ResponseWriter, r *http.Request){
  session, err := store.Get(r, "p2drive")
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  } 

  session.Options.MaxAge = -1
  session.Save(r,w)

  http.Redirect(w,r,"/projects/p2drive/intro", 302)

}

func checkAuth(c *web.C, h http.Handler) http.Handler{

  fn := func (w http.ResponseWriter, r *http.Request) {
    if endpoint := path.Base(r.URL.Path);
    endpoint == "intro" ||
    endpoint == "auth" ||
    endpoint == "logout" {
      h.ServeHTTP(w,r)
      return
    }

    session, err := store.Get(r, "p2drive")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    } 

    //authenticate gDrive if that hasn't been done already
    tok, ok := session.Values["access-token"]; if !ok {
      http.Redirect(w,r,"/projects/p2drive/auth", 302)
      return
    } 

    c.Env["access-token"] = tok

    h.ServeHTTP(w,r)
  }

  return http.HandlerFunc(fn)
}

func AddRoutes(p string) *web.Mux{
  mux := web.New()
  mux.Get("/intro", intro)
  mux.Get("/input", form)
  mux.Post("/input", compute)
  mux.Post("/submit", submit)
  mux.Get("/auth", auth)
  mux.Get("/auth/logout", logout)
  mux.Use(middleware.SubRouter)
  mux.Use(checkAuth)
  return mux
}