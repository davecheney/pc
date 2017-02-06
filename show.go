package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/aclements/go-moremath/stats"
	"github.com/olekukonko/tablewriter"
)

func show(id int, format string, ssort string, reverse bool, tag string) {
	f, err := os.Open(fmt.Sprintf("papercall.%d.json", id))
	check(err)

	var subs []*Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	rev := func(a, b int) (int, int) {
		if reverse {
			a, b = b, a
		}
		return a, b
	}

	switch strings.ToLower(ssort) {
	case "updated":
		sort.Slice(subs, func(i, j int) bool { i, j = rev(i, j); return subs[i].Updated.After(subs[j].Updated) })
	case "trust":
		sort.Slice(subs, func(i, j int) bool { i, j = rev(i, j); return subs[i].Trust > subs[j].Trust })
	case "stddev":
		sort.Slice(subs, func(i, j int) bool {
			i, j = rev(i, j)
			var s1, s2 stats.Sample
			for _, r := range subs[i].Ratings {
				s1.Xs = append(s1.Xs, float64(r.Value))
			}
			for _, r := range subs[j].Ratings {
				s2.Xs = append(s2.Xs, float64(r.Value))
			}
			return s1.StdDev() < s2.StdDev()
		})
	case "rating":
		fallthrough
	default:
		sort.Slice(subs, func(i, j int) bool { i, j = rev(i, j); return subs[i].Rating > subs[j].Rating })

	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"title", "format", "tags", "rating", "trust", "url"})
	format = strings.ToUpper(format)
	var rows int
	for _, s := range subs {
		f := strings.SplitN(strings.ToUpper(s.Format), " ", -1)[0]
		match, err := regexp.MatchString(format, f)
		check(err)
		if !match {
			continue
		}
		tags := strings.Join(s.Tags, " ")
		match, err = regexp.MatchString(tag, tags)
		check(err)
		if !match {
			continue
		}
		rows++
		var samp stats.Sample
		for _, r := range s.Ratings {
			samp.Xs = append(samp.Xs, float64(r.Value))
		}
		table.Append([]string{
			s.Title,
			f,
			tags,
			fmt.Sprintf("%0.2f (%0.2f)", samp.Mean(), samp.StdDev()),
			fmt.Sprintf("%0.2f", s.Trust),
			fmt.Sprintf("https://papercall.io/cfps/%d/submissions/%d", id, s.Id),
		})
	}
	table.SetFooter([]string{"Count", fmt.Sprintf("%d", rows), "", "", "", ""})
	table.Render()
}
