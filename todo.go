package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func todo(id int) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"title", "format", "tags", "reason", "url"})
	var rows int
	sort.Slice(subs, func(i, j int) bool { return subs[j].Updated.After(subs[i].Updated) })
	for _, s := range subs {
		reason := "no rating"
		if len(s.Ratings) > 0 {
			sort.Slice(s.Ratings, func(i, j int) bool {
				// show only reviews where the newest rating is _older_ than the
				// proposal's updated date
				j, i = i, j // comment this line to show any proposal whose oldest rating is older than the propsal.
				return s.Ratings[i].Updated.After(s.Ratings[j].Updated)
			})
			if s.Updated.After(s.Ratings[0].Updated) {
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
			fmt.Sprintf("https://papercall.io/cfps/%d/submissions/%d", id, s.Id),
		})
		rows++
	}
	table.SetFooter([]string{"Count", fmt.Sprintf("%d", rows), "", "", ""})
	table.Render()
}
