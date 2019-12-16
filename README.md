# pc
A CLI for papercall.io conferences.

# Installation

`pc` is writen in Go. You will need Go 1.11 or later to build it.

# Papercall API access required

`pc` requires API access to download data.
API access is a paid feature, _`pc` will not work with the free papercall plan_.

To obtain your API token, go to the papercall website, click the events tab, choose your event, then choose the Organisers link on the right hand side. On the list of organisers for your event your API token will be listed under your name.

Once you have your API key export it in your shell

```
$ export PAPERCALL_API_TOKEN=ae91a85a4d25c005a91172d7b51ba9bfcfa3c95a
```

Addionally, you'll need your papercall event ID. This will be in the URL for your event. Once you have it, export it in your shell:

```
$ export PAPERCALL_EVENT_ID=12345
```

# Usage

`pc` operates on a cache of your event's data. To populate this cache run

    pc refresh

This will download all the submissions, and their ratings and cache them locally.

Once you have downloaded the data for your event, you can display it with

    pc show

# Example usage

Show all tutorials sorted by trust, in reverse order (lowest to highest)

    pc show -f tut -s trust -r

Show all the talks with the tag `testing`

    pc show -t testing

Show reviwers by completion count

    pc reviewers -s count

Show proposals that have been updated since the newest review

    pc todo

Show proposals that have been updated since the _oldest_ review

    pc todo -a
