package main

/*
Due to the fact that credentials were being passed around, I decied to roll a small golang webserver using
the martini framework as a means to better familiarize myself with the language. There are three routes:

	/conversation
	/conversation/:UUID
	/wakeup/:UUID 

All of these correspond in fairly straightforward fashion with the PullString API's endpoints, the exception
being /wakeup, which simply transmits and empty JSON in order to elicit follow-up content.
*/

import ("github.com/go-martini/martini"
  		"github.com/martini-contrib/render"
		"net/http"
		"bytes"
		"io/ioutil"
		"encoding/json"
		"os"

)

/* structs to store JSON data, including array of structs to contain "outputs" array - golang is quite opinionated
in this regard! */

type Text struct {
	Text string
}

type PullstringResp struct {
	Outputs []Text
	Conversation string
	Timed_Response_Interval float64
	Is_Fallback bool
}

var args []string = os.Args[1:]

var projID string = args[0]
var apiKey string = args[1]

/* function to encode JSON data and POST it to PullString REST API */

func pullStringReq(key string, val string, UUID string) PullstringResp {
	var jsonStr []byte
	if (len(key) == 0 && len(val) == 0) {
		jsonStr = []byte("{}")
	} else {
		jsonStr = []byte("{\"" + key + "\": \"" + val + "\"}")
	}
	println(string(jsonStr))
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
	println("RESPONSE DATA:")
	println(string(body))
	return msg

}

/* main logic for routing and rendering responses */

func main() {
	router := martini.Classic()
	router.Use(render.Renderer())

	/* fallback/default page and unknown routes */

	static := martini.Static("assets", martini.StaticOptions{Fallback: "/index.html"})
	router.NotFound(static, http.NotFound)

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
		println(text)
		r.JSON(200, msg)
		})

	router.Get("/wakeup/:UUID", func (w http.ResponseWriter, params martini.Params, req *http.Request, r render.Render) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		UUID := "/" + params["UUID"]
		msg := pullStringReq("", "", UUID)
		r.JSON(200, msg)
		})

	router.RunOnAddr(":3000")
}