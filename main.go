	package main



	import (
		"encoding/json"
		"fmt"
		"os"
		"strconv"
		"strings"
		"github.com/PuerkitoBio/goquery"
		"log"
		"net/http"
	)

	type User struct {
		Id         string `json:"UID"`
		Name       string `json:"Username"`
		Group      string `json:"Usergroup"`
		Messages   string `json:"Amount of messages"`
		Reactions  string `json:"Amount of reactions"`
	}


	// getDocument accepts a thread URL, then fetches the page and returns a goQuery document
	func getDocument(url string) *goquery.Document {

		client := &http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Print(err)
			return nil
		}

		req.Header.Set("User-Agent", "telegram:@omnipf")
		resp, err := client.Do(req)

		if err != nil {
			log.Print(err)
		} else {

			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			//log.Printf("GET %s successfully\n", url)
			return doc
		}

		return nil
	}

	func ParseUserInformation( DocumentForm *goquery.Document) *User {
		var TempUserInfo User

		TempUserInfo.Name = strings.TrimSpace(DocumentForm.Find("span.username").Text())


		TempUserInfo.Group = strings.TrimSpace(DocumentForm.Find("strong").Text())

		if(TempUserInfo.Group == ""){
			TempUserInfo.Group = "New User"
		}

		TempUserInfo.Messages = strings.TrimSpace(DocumentForm.Find("dl.pairs.pairs--rows.pairs--rows--centered.fauxBlockLink").Find("dd").Text())

		var Message string

		Message = strings.TrimSpace(DocumentForm.Find("dl.pairs.pairs--rows.pairs--rows--centered").Text())



		TempUserInfo.Reactions = strings.Split(Message, "Reaction score\n\n")[1]


		if(TempUserInfo.Reactions == ""){
			TempUserInfo.Reactions = "0"}


		return &TempUserInfo
	}

	func main(){


		var url string
		var UsersAmount int
		print("Xenforo User parser v 0.1 | By OmniRF\n")


		fmt.Println("Enter forum url: ")

		// var then variable name then variable type

		// Taking input from user
		fmt.Scanln(&url)

		fmt.Println("Enter a user amount: ")

		// var then variable name then variable type

		// Taking input from user
		fmt.Scanln(&UsersAmount)


		var UsersInf []*User
		var ID string
		for i := 1; i < UsersAmount; i++ {

			var BasicUser *User
			var DocumentForm *goquery.Document


			ID = strconv.Itoa(i)

			DocumentForm = getDocument(url + "/members/" + ID + "/" )

			DocumentForm.Find("div.memberHeader").Each(func(i int, s *goquery.Selection) {

				BasicUser = ParseUserInformation(DocumentForm)
				BasicUser.Id = ID
				UsersInf = append(UsersInf, BasicUser)

				jsonrer, _ := json.Marshal(UsersInf)

				file, _ := os.Create("extendedData.json")
				file.WriteString(string(jsonrer))



				print("User " + BasicUser.Name + " Parsed \n")
			})

		}
	}
