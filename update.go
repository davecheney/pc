package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func refreshCache(id int) {
	req, err := http.NewRequest("GET", "https://www.papercall.io/api/v1/submissions", nil)
	check(err)
	req.Header.Add("Authorization", *apiKey)
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("expecting status 200, got %d", resp.Status)
	}
	var submissions []map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&submissions)
	check(err)

	f, err := os.Create(fmt.Sprintf("papercall.%d.json", id))
	check(err)
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(submissions)
	check(err)
}
