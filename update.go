package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Talk struct {
	Title       string `json:"title"`
	Abstract    string `json:"abstract"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	Level       string `json:"audience_level"`
	Format      string `json:"talk_format"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Submission struct {
	Id        int       `json:"id"`
	Confirmed bool      `json:"confirmed"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
	Rating    float64   `json:"rating"`
	Trust     float64   `json:"trust"`
	State     string    `json:"state"`
	Tags      []string  `json:"tag_list,omitempty"`
	Talk      `json:"talk"`
	Ratings   []Rating `json:"ratings"`
}

type Rating struct {
	Id       int       `json:"id"`
	Value    int       `json:"value"`
	Created  time.Time `json:"created_at"`
	Updated  time.Time `json:"updated_at"`
	Comments string    `json:"comments,omitempty"`
	User     `json:"user"`
}

func refreshCache(id int) {
	req, err := http.NewRequest("GET", "https://www.papercall.io/api/v1/submissions?per_page=9999", nil)
	check(err)
	req.Header.Add("Authorization", *apiKey)
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("expecting status 200, got %d", resp.Status)
	}
	var submissions []*Submission
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&submissions)
	check(err)

	log.Println("retrieved", len(submissions), "submissions")

	f, err := os.Create(fmt.Sprintf("papercall.%d.json", id))
	check(err)
	defer f.Close()

	var wg sync.WaitGroup
	wg.Add(len(submissions))
	sem := make(chan int, 10)
	for _, s := range submissions {
		go func(s *Submission) {
			defer wg.Done()
			sem <- s.Id
			defer func() {
				<-sem
			}()
			req, err := http.NewRequest("GET", fmt.Sprintf("https://www.papercall.io/api/v1/submissions/%d/ratings", s.Id), nil)
			check(err)
			req.Header.Add("Authorization", *apiKey)
			resp, err := http.DefaultClient.Do(req)
			check(err)
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Fatalf("expecting status 200, got %d", resp.Status)
			}
			dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&s.Ratings)
			check(err)
			log.Printf("submission %d: %d ratings", s.Id, len(s.Ratings))
		}(s)
	}
	wg.Wait()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(submissions)
	check(err)
}
