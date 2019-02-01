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
)

type Topic struct {
	Name        string
	Submissions []*Submission
}

func topics(id int, input string) {
	subs := submissions(fmt.Sprintf("papercall.%d.json", id))

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

	t, err := template.New("topics").Parse(topicsT)
	check(err)
	err = t.Execute(os.Stdout, map[string]interface{}{
		"cfp":    1642, // cannot figure out how to get this from the API
		"Topics": topics,
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
  <h1>GopherCon 2019 topic breakdown</h1>
  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.6/umd/popper.min.js" integrity="sha384-wHAiFfRlMFy6i5SRaxvfOCifBUQy1xHdJ/yoi7FRNXMRBu5WHdZYu1hA6ZOblgut" crossorigin="anonymous"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/js/bootstrap.min.js" integrity="sha384-B0UglyR+jN6CkvvICOB2joaf5I4l3gm9GU6Hc1og6Ls7i6U/mkkaduKaBhlAXv9k" crossorigin="anonymous"></script>
<div class="container">
<base href="https://www.papercall.io/cfps/{{.cfp}}/submissions/">
{{ range .Topics }}
<h2>{{ .Name }}</h1>
{{ range .Submissions }}
<div class="row justify-content-md-left">
  <div class="col-6">
    <a href="{{.Id}}">{{ .Talk.Title }}</a>
  </div>
  <div class="col-2">
    {{ .Rating }}
  </div>  
  <div class="col-2">
    {{ .Trust }}
  </div>
</div>
{{ end }}
{{ end }}
</div>
</body>
</html>
`

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
