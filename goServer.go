package main

import ("github.com/go-martini/martini"
		"net/http"
		"bytes"
		"io/ioutil"
)

const projID string = "e50b56df-95b7-4fa1-9061-83a7a9bea372"
const apiKey string = "9fd2a189-3d57-4c02-8a55-5f0159bff2cf"

func pullStringReq(key string, val string, UUID string) {
	jsonStr := []byte("{\"" + key + "\": \"" + val + "\"}")
	req, err := http.NewRequest(
		"POST",
		"https://conversation.pullstring.ai/v1/conversation" + UUID,
		bytes.NewBuffer(jsonStr),
		)
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		println("You did it wrong.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	println(string(body))
}

func main() {
	router := martini.Classic()
	router.Get("/conversation", func() {

		pullStringReq("project", projID, "")

		})

	router.Get("/conversation/:UUID", func(params martini.Params, r *http.Request) {

		text := r.URL.Query().Get("text")
		UUID := "/" + params["UUID"]

		pullStringReq("text", text, UUID)

		})

	router.RunOnAddr(":3000")
}