package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func show(id int) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	// sort.Slice(subs, func(i, j int) bool { return subs[i].Rating > subs[j].Rating })
	sort.Slice(subs, func(i, j int) bool { return subs[i].Updated.After(subs[j].Updated) })

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"title", "kind", "rating", "trust", "url"})
	for _, s := range subs {
		table.Append([]string{
			s.Title,
			strings.SplitN(strings.ToUpper(s.Format), " ", -1)[0],
			fmt.Sprintf("%0.2f", s.Rating),
			fmt.Sprintf("%0.2f", s.Trust),
			fmt.Sprintf("https://papercall.io/cfps/%d/submissions/%d", id, s.Id),
		})
	}
	table.Render()
}
