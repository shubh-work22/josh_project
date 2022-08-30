package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type URLS struct {
	Name []string `json:"urls"`
}

var mp = make(map[string]string)

func Status(w http.ResponseWriter, r *http.Request) {
	go goStatus()
	io.WriteString(w, "Checking!\n")
}

func goStatus() {
	for {
		time.Sleep(60 * time.Second)
		var wg sync.WaitGroup

		for key := range mp {
			wg.Add(1)
			key := key

			go func() {
				defer wg.Done()
				resp, err := http.Get("https://" + key)

				if err != nil || resp.StatusCode != 200 {
					mp[key] = "DOWN"
					fmt.Println("Url : ", key, "status : DOWN")
					return
				} else {
					mp[key] = "UP"
					fmt.Println("Url : ", key, "status : 200 OK")
				}

			}()

		}
		wg.Wait()
	}
}

func PostWebsites(w http.ResponseWriter, r *http.Request) {
	webUrls := URLS{}
	err := json.NewDecoder(r.Body).Decode(&webUrls)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(webUrls)

	for _, url := range webUrls.Name {
		mp[url] = "DOWN"
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(webUrls)
	if err != nil {
		fmt.Println(err)
	}

	io.WriteString(w, "Welcome!\n")
}

func getWebsite(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("name") {

		w.WriteHeader(http.StatusOK)
		var temp = make(map[string]string)

		temp[r.URL.Query().Get("name")] = mp[r.URL.Query().Get("name")]

		err := json.NewEncoder(w).Encode(temp)
		if err != nil {
			fmt.Println(err)
		}

	} else {

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(mp)
		if err != nil {
			fmt.Println(err)
		}

	}

	io.WriteString(w, "Response Given!\n")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Website!")
}

func main() {
	fmt.Println("Server is Up and Running on port 8080!")

	http.HandleFunc("/", home)
	http.HandleFunc("/post", PostWebsites)
	http.HandleFunc("/status", Status)
	http.HandleFunc("/websites", getWebsite)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

//curl -X POST -H "Content-Type: application/json" -d '{"urls":["www.facebook.com","www.google.com","www.fakewebsite1.com","www.youtube.com"]}' http://localhost:8080/post
//
//curl http://localhost:8080/status
//
//curl "http://localhost:8080/websites?name=www.youtube.com"
