package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func show(id int) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	sort.Slice(subs, func(i, j int) bool { return subs[i].Rating > subs[j].Rating })

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"title", "rating", "trust"})
	for _, s := range subs {
		table.Append([]string{
			s.Title,
			fmt.Sprintf("%0.2f", s.Rating),
			fmt.Sprintf("%0.2f", s.Trust),
		})
	}
	table.Render()
}
