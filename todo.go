package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func todo(id, reviewerID int, all bool) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"title", "format", "tags", "reason", "updated", "url"})
	var rows int
	sort.Slice(subs, func(i, j int) bool { return subs[j].Updated.After(subs[i].Updated) })
	rev := func(a, b int) (int, int) {
		if all {
			a, b = b, a
		}
		return a, b
	}

	for _, s := range subs {
		reason := "no rating"
		if len(s.Ratings) > 0 {
			sort.Slice(s.Ratings, func(i, j int) bool {
				i, j = rev(i, j)
				return s.Ratings[i].Updated.After(s.Ratings[j].Updated)
			})
			reviewerIDX := 0
			for idx, reviewer := range s.Ratings {
				if reviewer.Id == reviewerID {
					reviewerIDX = idx
					break
				}
			}

			if s.Updated.After(s.Ratings[reviewerIDX].Updated) {
				reason = "proposal updated"
			} else {
				continue
			}
		}
		tags := strings.Join(s.Tags, " ")
		table.Append([]string{
			s.Title,
			strings.SplitN(strings.ToUpper(s.Format), " ", -1)[0],
			tags,
			reason,
			s.Updated.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("https://papercall.io/cfps/%d/submissions/%d", id, s.Id),
		})
		rows++
	}
	table.SetFooter([]string{"Count", fmt.Sprintf("%d", rows), "", "", "", ""})
	table.Render()
}
