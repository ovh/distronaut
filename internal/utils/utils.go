package utils

import (
	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"net/url"
	"path"
	"regexp"
	"strings"
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
		log.Debugf("no match for <%s>", sel)
	}
	return ns, nil
}

// Execute regex and return first captured group
func RegexCapture(pattern string, str string) string {
	//Execute regex
	r := regexp.MustCompile(pattern)
	matches := r.FindStringSubmatch(str)
	if len(matches) < 2 {
		log.Debugf("no matches: <%s> <%s>", pattern, str)
		return ""
	}
	if len(matches) > 2 {
		log.Debugf("multiple group matches: <%s> <%s> <%+v>", pattern, str, matches[1:])
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
func unique(arr []string) []string {
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

// Polyfill for url.JoinPath added in golang 1.19
func urlJoinPath(base string, elem ...string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		log.Warnf("failed to parse url: %s (%s)", base, err)
		return "", err
	}
	if len(elem) > 0 {
		elem = append([]string{u.EscapedPath()}, elem...)
		p := path.Join(elem...)
		if strings.HasSuffix(elem[len(elem)-1], "/") && !strings.HasSuffix(p, "/") {
			p += "/"
		}
		u.Path = p
	}
	return u.String(), nil
}
