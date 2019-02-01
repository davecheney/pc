// pc is a command line tool to analyse PaperCall.io CFP results.
package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiKey  = kingpin.Flag("apikey", "PAPERCALL_API_TOKEN").Short('k').Default(os.Getenv("PAPERCALL_API_TOKEN")).String()
	eventid = kingpin.Flag("event", "event id.").Short('e').Default("1701").Int()
)

func main() {

	kingpin.Command("refresh", "Refresh event cache")
	showCmd := kingpin.Command("show", "show proposals").Default()
	format := showCmd.Flag("format", "presentation format.").Short('f').String()
	sort := showCmd.Flag("sort", "sort by which column").Short('s').Default("rating").String()
	reverse := showCmd.Flag("reverse", "reverse sort order").Short('r').Bool()
	tag := showCmd.Flag("tag", "filter only tag").Short('t').String()
	reviewersCmd := kingpin.Command("reviewers", "show reviewer completion.")
	revsort := reviewersCmd.Flag("sort", "sort by which column").Short('s').Default("count").String()
	todoCmd := kingpin.Command("todo", "show outstanding proposals.")
	reviewerID := todoCmd.Flag("reviewer", "filter for reviwer ID").Short('i').Default("0").Int()
	all := todoCmd.Flag("all", "show all outstanding todos (not just latest).").Short('a').Bool()
	kingpin.Command("speakers", "show speakers.")
	topicsCmd := kingpin.Command("topics", "build HTML page of proposals by topic")
	input := topicsCmd.Flag("input", "space separated list of topic name and proposals ids").Short('i').Required().String()

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Dave Cheney")
	kingpin.CommandLine.Help = "pc is a command line tool to analyse PaperCall.io CFP results"
	switch kingpin.Parse() {
	case "show":
		show(*eventid, *format, *sort, *reverse, *tag)
	case "refresh":
		refreshCache(*eventid)
	case "reviewers":
		reviewers(*eventid, *revsort)
	case "todo":
		todo(*eventid, *reviewerID, *all)
	case "speakers":
		speakers(*eventid)
	case "topics":
		topics(*eventid, *input)
	default:
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
