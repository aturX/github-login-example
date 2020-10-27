package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// type Conf struct {
// 	// 配置
// 	ClientID     string
// 	ClientSecret string
// 	RedirectURL  string
// }

// var conf = Conf{
// 	ClientID:     "e30a93b7d9adcae48369",
// 	ClientSecret: "d9dbef09174f36b06fd3f1acbdba86086352ed23",
// 	RedirectURL:  "http://localhost:8080/oauth/redirect",
// }

const clientID = "e30a93b7d9adcae48369"
const clientSecret = "d9dbef09174f36b06fd3f1acbdba86086352ed23"

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// Create a new redirect route route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth enpoint
		// to get our access token
		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		// Finally, send a response to redirect the user to the "welcome" page
		// with the access token
		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("监听失败，错误信息为:", err)
	}
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
