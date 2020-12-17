package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	GivenName  string
	Id         string
	Roles      string
	FamilyName string
}

type RespObject struct {
	Person User
}

func getName(index int) string {
	resp, err := http.Get(fmt.Sprintf("http://name-service.appspot.com/api/v1/names/%d.json", index))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data := RespObject{}
	user := User{}
	json.Unmarshal(body, &data)
	json.Unmarshal([]byte(data.Person), &user)
	fmt.Print(user)
	return data.Person.GivenName
}

func main() {
	fmt.Println(getName(1))
}
