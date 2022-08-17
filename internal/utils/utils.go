package utils

import (
	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"regexp"
)

// Query matching nodes from given selector
func querySelector(uri string, sel string) ([]*html.Node, error) {
	//Open url
	var ns []*html.Node
	doc, err := htmlquery.LoadURL(uri)
	if err != nil {
		return ns, err
	}
	log.Debugf("opened <%s>", uri)

	//Query matching nodes
	ns, err = htmlquery.QueryAll(doc, sel)
	if err != nil {
		log.Warnf("%s", err)
	}
	if len(ns) == 0 {
		log.Warnf("no match for <%s>", sel)
	}
	return ns, nil
}

// Execute regex and stores named groups into a string map
func ParseRegex(pattern string, str string) map[string]string {
	//Execute regex
	r := regexp.MustCompile(pattern)
	parsed := make(map[string]string)
	matches := r.FindStringSubmatch(str)
	if len(matches) == 0 {
		log.Warnf("no matches: <%s> <%s>", pattern, str)
		return parsed
	}

	//Convert named groups to string map
	parsed["*"] = matches[0]
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(matches) {
			parsed[name] = matches[i]
			log.Debugf("parsed group: <%s> = <%s>", name, parsed[name])
		}
	}
	return parsed
}

// Execute regex and return first captured group
func RegexCapture(pattern string, str string) (string) {
	//Execute regex
	r := regexp.MustCompile(pattern)
	matches := r.FindStringSubmatch(str)
	if len(matches) < 2 {
		log.Warnf("no matches: <%s> <%s>", pattern, str)
		return ""
	}
	if len(matches) > 2 {
		log.Warnf("multiple group matches: <%s> <%s> <%+v>", pattern, str, matches[1:])
	}
	return matches[1]
}

// Deep copy a map string
func copyMap(mp map[string]string) map[string]string {
	cp := make(map[string]string)
	for k, v := range mp {
		cp[k] = v
	}
	return cp
}

// Remove duplicates from array
func unique(arr []string) ([]string){
	var set []string
	has := make(map[string]bool)
	for _, v := range arr {
		if _, ok := has[v]; !ok {
			has[v] = true
			set = append(set, v)
		}
	} 
	return set
}