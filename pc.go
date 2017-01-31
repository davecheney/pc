// pc is a command line tool to analyse PaperCall.io CFP results.
package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiKey    = kingpin.Flag("apikey", "PAPERCALL_API_TOKEN").Short('k').Default(os.Getenv("PAPERCALL_API_TOKEN")).String()
	proposals = kingpin.Command("proposals", "Fetch all proposals.").Default()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("Dave Cheney")
	kingpin.CommandLine.Help = "pc is a command line tool to analyse PaperCall.io CFP results"
	switch kingpin.Parse() {
	case "proposals":

	default:
		os.Exit(1)
	}

}
