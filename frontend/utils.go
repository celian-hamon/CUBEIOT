package main

import (
	"bytes"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Admin    int    `json:"admin"`
	Password []byte `json:"password"`
	Salt     []byte `json:"salt"`
}

var client = &http.Client{}
var baseUri = "http://localhost:5000/"

func CheckCredentials(name string, password string) (bool, error) {
	userId := getIdFromName(name)
	req, _ := http.Post(baseUri+"user/"+strconv.Itoa(userId), "application/json", bytes.NewBuffer([]byte(`{"password":"`+password+`"}`)))
	if req.StatusCode != 200 {
		return false, nil
	} else {
		return true, nil
	}
}

type sonde struct {
	Id        int     `json:"id"`
	Active    int     `json:"active"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type sondeListe struct {
	Sondes []sonde `json:"sondes"`
}

func retreiveSondeList(name string) sondeListe {
	var sondeList sondeListe
	userId := getIdFromName(name)
	//make a request to the api to get the list of sondes of the user id
	user := strconv.Itoa(userId)
	response, err := http.Get(baseUri + "user/" + user + "/sonde")
	if err != nil {
		fmt.Println("erreur http :", err)
	}
	defer response.Body.Close()

	//Decode the data
	if err := json.NewDecoder(response.Body).Decode(&sondeList.Sondes); err != nil {
		fmt.Printf("ooopsss! an error occurred, please try again : ", err)
	}

	return sondeList
}
func createAccount(name string, password string, admin string) bool {
	//make a post request to the api to create the account
	var admin_value int
	if admin == "True" {
		admin_value = 1
	} else {
		admin_value = 0
	}
	response, _ := http.Post(baseUri+"user", "application/json", bytes.NewBuffer([]byte(`{"name":"`+name+`","password":"`+password+`","admin":`+strconv.Itoa(admin_value)+`}`)))

	if response.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

func updateAccount(name string, newName string) bool {
	//make a put request to the api to update the account
	userID := getIdFromName(name)
	user := getUser(userID)

	url := baseUri + "user/" + strconv.Itoa(user.Id)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(`{"name":"`+newName+`","admin":`+strconv.Itoa(user.Admin)+`}`)))
	req.Header.Add("profile", strconv.Itoa(user.Id))
	req.Header.Set("Content-type", "application/json")

	resp, _ := client.Do(req)

	if resp.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

func getUser(id int) User {
	var user User

	//make a request to the api to get the user
	req, err := http.NewRequest("GET", baseUri+"user/"+strconv.Itoa(id), nil)
	req.Header.Add("profile", strconv.Itoa(user.Id))
	response, _ := client.Do(req)
	if err != nil {
		fmt.Println("erreur http :", err)
	}
	defer response.Body.Close()
	//Decode the data
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		fmt.Printf("ooopsss! an error occurred, please try again : ", err)
	}

	return user
}

func retreiveSonde(name string, sondeId string) sonde {
	var sonde sonde

	userId := getIdFromName(name)

	req, _ := http.NewRequest("GET", baseUri+"sonde/"+sondeId, nil)
	req.Header.Add("profile", strconv.Itoa(userId))
	response, err := client.Do(req)

	if err != nil {
		fmt.Println("erreur http :", err)
	}
	defer response.Body.Close()

	//Decode the data
	if err := json.NewDecoder(response.Body).Decode(&sonde); err != nil {
		fmt.Printf("ooopsss! an error occurred, please try again : ", err)
	}

	return sonde
}

func getIdFromName(name string) int {

	resp, _ := http.Get(baseUri + "user/" + name + "/id")
	type user_id struct {
		Id int `json:"id"`
	}
	var userId user_id
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&userId); err != nil {
		fmt.Printf("ooopsss! an error occurred, please try again : ", err)
	}
	return userId.Id
}

func createSonde(lat string, long string, name string) bool {
	//make a post request to the api to create the sonde
	userId := getIdFromName(name)
	body := bytes.NewBuffer([]byte(`{"latitude":` + lat + `,"longitude":` + long + `}`))
	req, _ := http.NewRequest(http.MethodPost, baseUri+"sonde", body)
	req.Header.Set("Content-type", "application/json")
	req.Header.Add("profile", strconv.Itoa(userId))
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("erreur http :", err)
		return false
	}
	if res.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

func updateSonde(name string, sonde_id string, lat string, long string) bool {
	//make a put request to the api to update the prob
	userId := getIdFromName(name)
	body := bytes.NewBuffer([]byte(`{"latitude":` + lat + `,"longitude":` + long + `}`))
	req, _ := http.NewRequest(http.MethodPut, baseUri+"sonde/"+sonde_id, body)
	req.Header.Set("Content-type", "application/json")
	req.Header.Add("profile", strconv.Itoa(userId))
	resp, _ := client.Do(req)

	if resp.StatusCode == 201 {
		return true
	} else {
		return false
	}
}

type releve struct {
	Date        string  `json:"date"`
	Humidity    int     `json:"humidite"`
	Id          int     `json:"id"`
	SondeId     int     `json:"id_sonde"`
	Temperature float64 `json:"temperature"`
}

func retreiveReleveList(sondeId string) []releve {
	var releveList []releve
	//make a request to the api to get the list of releves of the sonde id
	response, err := http.Get("http://localhost:5000/releve" + "/sonde/" + sondeId)
	if err != nil {
		fmt.Println("erreur http :", err)
	}
	defer response.Body.Close()
	//Decode the data
	if err := json.NewDecoder(response.Body).Decode(&releveList); err != nil {
		fmt.Printf("ooopsss! an error occurred, please try again : ", err)
	}

	return releveList
}
