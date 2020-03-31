package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func postWebHook(h *http.Client, url string, body interface{}) error {
	jsonStr, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Error", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := h.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Println("Body discard error", err)
	}

	return nil
}
