package main

import ("github.com/go-martini/martini"
  		"github.com/martini-contrib/render"
		"net/http"
		"bytes"
		"io/ioutil"
		"encoding/json"

)
type Text struct {
	Text string
}

type PullstringResp struct {
	Outputs []Text
	Conversation string
	Timed_Response_Interval float64
}

const projID string = "e50b56df-95b7-4fa1-9061-83a7a9bea372"
const apiKey string = "9fd2a189-3d57-4c02-8a55-5f0159bff2cf"

func pullStringReq(key string, val string, UUID string) PullstringResp {
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

	var msg PullstringResp
	err = json.Unmarshal(body, &msg)
	if err != nil {
		println("You did it wrong.")
		return msg
	}
	println("JSON VALUES:")
	println(string(msg.Conversation))

	println(string(body))
	return msg

}

func main() {
	router := martini.Classic()
	router.Use(render.Renderer())

	router.Get("/conversation", func(w http.ResponseWriter, r render.Render) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		msg := pullStringReq("project", projID, "")
		r.JSON(200, msg)
		})

	router.Get("/conversation/:UUID", func(w http.ResponseWriter, params martini.Params, req *http.Request, r render.Render) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		text := req.URL.Query().Get("text")
		UUID := "/" + params["UUID"]
		msg := pullStringReq("text", text, UUID)

		r.JSON(200, msg)
		})

	router.RunOnAddr(":3000")
}