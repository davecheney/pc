package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func reviewers(id int, ssort string) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	reviewers := make(map[string]int)
	for _, s := range subs {
		for _, r := range s.Ratings {
			reviewers[r.Name]++
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "count", "completion (%)"})

	var keys []string
	for k := range reviewers {
		keys = append(keys, k)
	}
	switch ssort {
	case "name":
		sort.Strings(keys)
	case "count":
		fallthrough
	default:
		sort.Slice(keys, func(i, j int) bool {
			return reviewers[keys[i]] > reviewers[keys[j]]
		})
	}
	for _, k := range keys {
		table.Append([]string{
			k,
			fmt.Sprint(reviewers[k]),
			fmt.Sprintf("%0.2f", float64(reviewers[k])/float64(len(subs))*100),
		})
	}
	table.Render()
}
