// pc is a command line tool to analyse PaperCall.io CFP results.
package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiKey  = kingpin.Flag("apikey", "PAPERCALL_API_TOKEN").Short('k').Default(os.Getenv("PAPERCALL_API_TOKEN")).String()
	eventid = kingpin.Flag("event", "event id.").Short('e').Default("274").Int()
)

func main() {

	kingpin.Command("refresh", "Refresh event cache")
	showCmd := kingpin.Command("show", "show proposals").Default()
	format := showCmd.Flag("format", "presentation format.").Short('f').Default(".+").String()
	sort := showCmd.Flag("sort", "sort by which column").Short('s').Default("rating").String()
	reverse := showCmd.Flag("reverse", "reverse sort order").Short('r').Bool()

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Dave Cheney")
	kingpin.CommandLine.Help = "pc is a command line tool to analyse PaperCall.io CFP results"
	switch kingpin.Parse() {
	case "show":
		show(*eventid, *format, *sort, *reverse)
	case "refresh":
		refreshCache(*eventid)
	default:
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
