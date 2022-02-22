package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

// internal page
func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		context := make(map[string]string)
		context["name"] = userName
		fp := path.Join("templates", "internal", "internal.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(response, context); err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

//region Session

//#region Session/Login_Logout
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	context := make(map[string]string)
	failed := request.FormValue("failed")
	if failed != "" {
		context["failed"] = "mdp"
	}
	fp := path.Join("templates", "login", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	login, _ := CheckCredentials(name, pass)
	redirectTarget := "/"
	status := http.StatusOK
	if name != "" && pass != "" {
		// .. check credentials ..
		if login {
			setSession(name, response)
			redirectTarget, status = "/internal", http.StatusFound
		} else {
			redirectTarget, status = "/?failed=mdp", http.StatusFound
		}
	}
	http.Redirect(response, request, redirectTarget, status)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

//#endregion Session/Login_Logout
//#region Session/Cookies
// cookie handling
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

//#endregion Session/Cokkies

//endregion Session

//region Profil

//#region Profil/Creation
func createAccountPageHandler(response http.ResponseWriter, request *http.Request) {
	context := make(map[string]string)
	fp := path.Join("templates", "create-account", "create-account.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func accountCreationHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	admin := request.FormValue("admin")
	redirectTarget := "/createAccount"
	status := http.StatusFound
	success := createAccount(name, pass, admin)
	if success {
		setSession(name, response)
		redirectTarget, status = "/internal", http.StatusFound
	} else {
		redirectTarget, status = "/createAccount", http.StatusFound
	}
	http.Redirect(response, request, redirectTarget, status)
}

//#endregion Profil/Creation
//#region Profil/Update
func profilUpdateHandler(response http.ResponseWriter, request *http.Request) {
	newName := request.FormValue("name")
	redirectTarget := "/profile"
	status := http.StatusFound
	success := updateAccount(getUserName(request), newName)
	if success {
		clearSession(response)
		setSession(newName, response)
		redirectTarget, status = "/internal", http.StatusFound
	} else {
		redirectTarget, status = "/profile/?failed=profil", http.StatusFound
	}
	http.Redirect(response, request, redirectTarget, status)
}

func profilPageHandler(response http.ResponseWriter, request *http.Request) {
	context := make(map[string]string)
	failed := request.FormValue("failed")
	succes := request.FormValue("succes")
	if failed != "" {
		context["failed"] = "profil"
	}
	if succes != "" {
		context["succes"] = "profil"
	}
	fp := path.Join("templates", "update_user", "updateUser.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

func profilePageHandler(response http.ResponseWriter, request *http.Request) {
	context := make(map[string]string)
	fp := path.Join("templates", "update_user", "updateUser.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

//#endregion Profil/Update

//endregion Profil

//region Sonde

//#region Sonde/Creation
func createSondeHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	context := make(map[string]string)
	failed := request.FormValue("failed")
	if failed != "" {
		context["failed"] = "sondecreation"
	}
	if userName != "" {
		fp := path.Join("templates", "create-sonde", "create-sonde.html")
		tmpl, _ := template.ParseFiles(fp)
		tmpl.Execute(response, context)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

func sondeCreationHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		latitude := request.FormValue("latitude")
		longitude := request.FormValue("longitude")
		redirectTarget := "/createSonde"
		status := http.StatusFound
		success := createSonde(latitude, longitude, userName)
		if success {
			redirectTarget, status = "/sondes", http.StatusFound
		} else {
			redirectTarget, status = "/createSonde/?failed=sondecreation", http.StatusFound
		}
		http.Redirect(response, request, redirectTarget, status)
	}
}

//#endregion Sonde/Creation
//#region Sonde/Update
func sondeUpdateHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		vars := mux.Vars(request)
		sondeId := vars["sondeId"]
		latitude := request.FormValue("latitude")
		longitude := request.FormValue("longitude")
		redirectTarget := "/updatesonde/{sondeId}"
		status := http.StatusFound
		success := updateSonde(userName, sondeId, latitude, longitude)
		if success {
			redirectTarget, status = "/updatesonde/{sondeId}/?succes=sonde", http.StatusFound
		} else {
			redirectTarget, status = "/updatesonde/{sondeId}/?failed=sonde", http.StatusFound
		}
		http.Redirect(response, request, redirectTarget, status)
	}
}

func updateSondePageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		vars := mux.Vars(request)
		sondeId := vars["sondeId"]
		sonde := retreiveSonde(userName, sondeId)
		context := make(map[string]string)
		context["Id"] = strconv.Itoa(sonde.Id)
		context["Active"] = strconv.Itoa(sonde.Active)
		context["Latitude"] = fmt.Sprintf("%.4f", sonde.Latitude)
		context["Longitude"] = fmt.Sprintf("%.4f", sonde.Longitude)
		fp := path.Join("templates", "update-sonde", "updateSonde.html")
		tmpl, _ := template.ParseFiles(fp)
		tmpl.Execute(response, context)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

func sondeUpdatePageHandler(response http.ResponseWriter, request *http.Request) {
	context := make(map[string]string)
	failed := request.FormValue("failed")
	succes := request.FormValue("succes")
	if failed != "" {
		context["failed"] = "sonde"
	}
	if succes != "" {
		context["succes"] = "sonde"
	}
	fp := path.Join("templates", "update-sonde", "updateSonde.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

//#endregion Sonde/Update
//#region Sonde/Delete

func sondeDeleteHandler(response http.ResponseWriter, request *http.Request) {
	userId := getIdFromName(getUserName(request))
	vars := mux.Vars(request)
	sondeId := vars["sondeId"]
	req, _ := http.NewRequest(http.MethodDelete, baseUri+"sonde/"+sondeId, nil)
	req.Header.Add("profile", strconv.Itoa(userId))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		http.Error(response, "Error", http.StatusInternalServerError)
	}
	http.Redirect(response, request, "/sondes", http.StatusFound)
}

//#endregion Sonde/Delete
//#region Sonde/Releves
func sondeReleveHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		vars := mux.Vars(request)
		sondeId := vars["sondeId"]
		releveList := retreiveReleveList(sondeId)
		context := make(map[string][]releve)
		context["releveList"] = releveList
		fp := path.Join("templates", "sonde_releve_list", "sonde-releve.html")
		tmpl, _ := template.ParseFiles(fp)
		tmpl.Execute(response, context)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

func relevesDeleteHandler(response http.ResponseWriter, request *http.Request) {
	userId := getIdFromName(getUserName(request))
	vars := mux.Vars(request)
	sondeId := vars["sondeId"]
	req, _ := http.NewRequest(http.MethodDelete, baseUri+"releves/"+sondeId, nil)
	req.Header.Add("profile", strconv.Itoa(userId))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		http.Error(response, "Error", http.StatusInternalServerError)
	}
	http.Redirect(response, request, "/sonde/{sondeId}", http.StatusFound)
}

//#endregion Sonde/Releves
//#region Sonde/State
func stateHandler(response http.ResponseWriter, request *http.Request) {
	userId := getIdFromName(getUserName(request))
	vars := mux.Vars(request)
	sondeId := vars["sondeId"]
	current := vars["current"]
	client := &http.Client{}
	if current == "activated" {
		req, _ := http.NewRequest(http.MethodPut, baseUri+"sonde/"+sondeId+"/desactivate", nil)
		req.Header.Add("profile", strconv.Itoa(userId))
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			http.Error(response, "Error", http.StatusInternalServerError)
		}
		http.Redirect(response, request, "/sonde/"+sondeId, http.StatusFound)
	} else if current == "desactivated" {
		req, _ := http.NewRequest(http.MethodPut, baseUri+"sonde/"+sondeId+"/activate", nil)
		req.Header.Add("profile", strconv.Itoa(userId))
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			http.Error(response, "Error", http.StatusInternalServerError)
		}
		http.Redirect(response, request, "/sonde/"+sondeId, http.StatusFound)
	}
	http.Redirect(response, request, "/sonde/"+sondeId, http.StatusFound)
}

//#endregion Sonde/State
//#region Sonde/Liste_Solo
func sondeHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		vars := mux.Vars(request)
		sondeId := vars["sondeId"]
		sonde := retreiveSonde(userName, sondeId)
		context := make(map[string]string)
		context["Id"] = strconv.Itoa(sonde.Id)
		context["Active"] = strconv.Itoa(sonde.Active)
		context["Latitude"] = fmt.Sprintf("%.4f", sonde.Latitude)
		context["Longitude"] = fmt.Sprintf("%.4f", sonde.Longitude)

		fp := path.Join("templates", "sonde", "sonde.html")
		tmpl, _ := template.ParseFiles(fp)
		tmpl.Execute(response, context)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

func sondeListHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		Sondes := retreiveSondeList(userName)
		var sondeList []map[string]string
		for _, sonde := range Sondes.Sondes {
			latitude := fmt.Sprintf("%.4f", sonde.Latitude)
			longitude := fmt.Sprintf("%.4f", sonde.Longitude)
			sondeAfter := map[string]string{
				"Id":        strconv.Itoa(sonde.Id),
				"Active":    strconv.Itoa(sonde.Active),
				"Latitude":  latitude,
				"Longitude": longitude,
			}
			sondeList = append(sondeList, sondeAfter)
		}
		context := make(map[string]interface{})
		context["sondes"] = sondeList
		fp := path.Join("templates", "sondes_liste", "sondes.html")
		tmpl, _ := template.ParseFiles(fp)
		tmpl.Execute(response, context)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}

//#endregion Sonde/Liste_Solo

//endregion Sonde

//region Custom Error Pages

//#region 404
func notFoundHandler(response http.ResponseWriter, request *http.Request) {
	redirectTarget := "/404"
	status := http.StatusFound
	http.Redirect(response, request, redirectTarget, status)
}

func notFoundCustomHandler(response http.ResponseWriter, request *http.Request) {
	fp := path.Join("templates", "404", "404.html")
	tmpl, err := template.ParseFiles(fp)
	context := make(map[string]string)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

//#endregion 404
//#region 403
func forbiddenCustomHandler(response http.ResponseWriter, request *http.Request) {
	fp := path.Join("templates", "403", "403.html")
	tmpl, err := template.ParseFiles(fp)
	context := make(map[string]string)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(response, context); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}

//#endregion 403

//endregion Custom Error Pages

// server main method

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/internal", internalPageHandler)
	router.HandleFunc("/createSonde", createSondeHandler)
	router.HandleFunc("/createSonde/", createSondeHandler).Queries("failed", "{failed}")
	router.HandleFunc("/updatesonde/{sondeId}", updateSondePageHandler)
	router.HandleFunc("/updatesonde/{sondeId}/", sondeUpdatePageHandler).Queries("failed", "{failed}")
	router.HandleFunc("/updatesonde/{sondeId}/", sondeUpdatePageHandler).Queries("succes", "{succes}")
	router.HandleFunc("/sondes", sondeListHandler)
	router.HandleFunc("/sonde/{sondeId}/releve", sondeReleveHandler)
	router.HandleFunc("/sonde/{sondeId}", sondeHandler).Methods("GET")
	router.HandleFunc("/createAccount", createAccountPageHandler)
	router.HandleFunc("/profile", profilePageHandler)
	router.HandleFunc("/profile/", profilPageHandler).Queries("failed", "{failed}")
	router.HandleFunc("/profile/", profilPageHandler).Queries("succes", "{succes}")
	router.HandleFunc("/404", notFoundCustomHandler)
	router.HandleFunc("/403", forbiddenCustomHandler)
	router.HandleFunc("/", indexPageHandler).Queries("failed", "{failed}")
	router.HandleFunc("/", indexPageHandler)

	router.HandleFunc("/account", accountCreationHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/sonde", sondeCreationHandler).Methods("POST")
	router.HandleFunc("/sonde/{sondeId}/", sondeUpdateHandler).Methods("POST")
	router.HandleFunc("/profil", profilUpdateHandler).Methods("POST")

	router.HandleFunc("/state/{sondeId}/{current}", stateHandler)
	router.HandleFunc("/sonde/{sondeId}", sondeDeleteHandler)

	router.HandleFunc("/releves/{sondeId}", relevesDeleteHandler)

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

//todo recup gps depuis ip
