package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aclements/go-moremath/stats"
)

type Topic struct {
	Name        string
	Submissions []*Submission
}

func topics(id int, input string) {
	path := fmt.Sprintf("papercall.%d.json", id)
	subs := submissions(path)

	f, err := os.Open(input)
	check(err)
	sc := bufio.NewScanner(f)
	var topics []Topic
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		topic := Topic{
			Name: strings.ReplaceAll(fields[0], "-", " "),
		}
		for _, f := range fields[1:] {
			id, err := strconv.ParseInt(f, 10, 64)
			check(err)
			s, ok := subs[int(id)]
			if ok {
				topic.Submissions = append(topic.Submissions, s)
			}
		}
		sort.SliceStable(topic.Submissions, func(i, j int) bool {
			return topic.Submissions[i].Rating > topic.Submissions[j].Rating
		})
		topics = append(topics, topic)
	}

	sort.SliceStable(topics, func(i, j int) bool {
		return topics[i].Name < topics[j].Name
	})

	t, err := template.New("topics").Funcs(map[string]interface{}{
		"mean":   meanRating,
		"stddev": stddevRating,
		"diff":   diffRatings,
	}).Parse(topicsT)
	check(err)
	fi, err := os.Stat(path)
	check(err)
	err = t.Execute(os.Stdout, map[string]interface{}{
		"cfp":       1642, // cannot figure out how to get this from the API
		"Topics":    topics,
		"inputfile": fi,
	})
	check(err)
}

const topicsT = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">

  <title>GopherCon 2019 topic breakdown</title>
</head>
<body>
  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.6/umd/popper.min.js" integrity="sha384-wHAiFfRlMFy6i5SRaxvfOCifBUQy1xHdJ/yoi7FRNXMRBu5WHdZYu1hA6ZOblgut" crossorigin="anonymous"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/js/bootstrap.min.js" integrity="sha384-B0UglyR+jN6CkvvICOB2joaf5I4l3gm9GU6Hc1og6Ls7i6U/mkkaduKaBhlAXv9k" crossorigin="anonymous"></script>
<div class="container">
<h1>GopherCon 2019 topic breakdown</h1>
<p>Ratings last updated: {{.inputfile.ModTime}}</p>
<base href="https://www.papercall.io/cfps/{{.cfp}}/submissions/">
<ol>
{{ range .Topics }}
<li><h2>{{ .Name }}</h1>
{{ range .Submissions }}
<div class="row justify-content-md-left">
  <div class="col-5">
    <a href="{{.Id}}">{{ .Talk.Title }}</a>
  </div>
  <div class="col-3">{{.Talk.Format }}</div>
  <div class="col-2">{{ mean .Ratings | printf "%3.1f" }} ± {{ diff .Ratings | printf "%.0f%%" }}</div>
  <div class="col-2">{{.State }}</div>
</div>
{{ end }}
</li>
{{ end }}
</div>
</body>
</html>
`

func meanRating(ratings []Rating) float64 {
	values := make([]float64, 0, len(ratings))
	for _, r := range ratings {
		values = append(values, float64(r.Value))
	}
	return stats.Mean(values)
}

func stddevRating(ratings []Rating) float64 {
	return stats.StdDev(values(ratings))
}

func values(ratings []Rating) []float64 {
	values := make([]float64, 0, len(ratings))
	for _, r := range ratings {
		values = append(values, float64(r.Value))
	}
	return values
}

func computeStats(ratings []Rating) (min, max, mean float64) {
	// Discard outliers.
	values := stats.Sample{Xs: values(ratings)}
	q1, q3 := values.Quantile(0.25), values.Quantile(0.75)

	var rvalues []float64
	lo, hi := q1-1.5*(q3-q1), q3+1.5*(q3-q1)
	for _, value := range values.Xs {
		if lo <= value && value <= hi {
			rvalues = append(rvalues, value)
		}
	}

	// Compute statistics of remaining data.
	min, max = stats.Bounds(rvalues)
	mean = stats.Mean(rvalues)
	return
}

func diffRatings(ratings []Rating) float64 {
	min, max, mean := computeStats(ratings)
	diff := 1 - min/mean
	if d := max/mean - 1; d > diff {
		diff = d
	}
	return diff * 100
}

func submissions(path string) map[int]*Submission {
	f, err := os.Open(path)
	check(err)
	defer f.Close()

	var subs []Submission
	dec := json.NewDecoder(f)
	err = dec.Decode(&subs)
	check(err)

	result := make(map[int]*Submission)
	for i := range subs {
		s := &subs[i]
		result[s.Id] = s
	}
	return result
}
