package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func speakers(id int) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "email", "shirt size", "talks"})

	type p struct {
		Profile
		submissions []*Submission
	}

	profiles := make(map[string]p)
	for _, s := range subs {
		prof := profiles[s.Profile.Name]
		prof.Profile = *s.Profile
		switch s.State {
		case "submitted", "accepted":
			prof.submissions = append(prof.submissions, s)

		}
		profiles[s.Profile.Name] = prof
	}
	var keys []string
	for k := range profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var rows, n int
	for _, k := range keys {
		prof := profiles[k]
		table.Append([]string{
			prof.Name,
			prof.Email,
			prof.ShirtSize,
			summarise(prof.submissions),
		})
		rows++
		n += len(prof.submissions)
	}
	table.SetFooter([]string{"Count", fmt.Sprintf("%d", rows), "Proposals", fmt.Sprintf("%d", n)})
	table.Render()
}

func summarise(subs []*Submission) string {
	s := make(map[string]int)
	for _, sub := range subs {
		s[strings.SplitN(sub.Format, " ", -1)[0]]++
	}
	var r string
	if n, ok := s["Keynote"]; ok {
		r += fmt.Sprintf("Keynote (%d) ", n)
	}
	if n, ok := s["Tutorial"]; ok {
		r += fmt.Sprintf("Tutorial (%d) ", n)
	}
	if n, ok := s["Workshop"]; ok {
		r += fmt.Sprintf("Workshop (%d) ", n)
	}
	return strings.TrimRight(r, " ")
}
