package utils

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// Link structure
type Link struct {
	Url     string `json:url`
	Hash    string `json:hash`
	Version string `json:version`
	Arch    string `json:arch`
}

// Version regex
const REG_VERSION = `(?P<version>\d+[-.]\d+(?:[-.]\d+)?)`

// Arch regex
const REG_ARCH = `(?P<arch>i386|amd(?:64)?|arm(?:64)?(?:el|hf)?|mips(?:64)?(?:el)?|ppc(?:64)?(?:el)?|s390x?|x86_64|(?:32|64)bits?)`

// Scrap a distribution mirror
func Scrap(uri string, pats map[string]string) []*Link {
	//Parse url
	u, _ := url.Parse(uri)

	//Extract domain
	hs := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())

	//Split path (wihtout leading empty one)
	ps := strings.Split(u.Path, "/")
	_, ps = ps[0], ps[1:]

	//Init context
	vars := make(map[string]string)

	//Scrap links and hashes
	ls := unique(scrap(hs, ps, pats, vars))
	var links []*Link
	for _, l := range ls {
		links = append(links, &Link{Url: l})
	}
	if pats[".hash.file"] != "" && pats[".hash.algo"] != "" && pats[".hash.pattern"] != "" {
		scrapHashes(links, pats)
	} else {
		log.Debugf("missing hash extration settings, ignoring")
	}
	scrapMeta(links, pats)
	for _, link := range links {
		log.Debugf("found: %+v", link)
	}
	return links
}

// Scrap a distribution mirror (recursive function)
func scrap(curr string, ps []string, pats map[string]string, vars map[string]string) []string {
	var links []string
	for i, p := range ps {

		//Handle route param
		if _, ok := pats[p]; p[:1] == ":" && ok {
			log.Debugf("found route param: %s", p)

			//Evaluate route pattern
			pat := scrapPattern(pats[p], vars)
			log.Debugf("searching for: %s", pat)

			//Iterate over matching links
			as, _ := querySelector(curr, fmt.Sprintf(`//a[matches(@href, '%s')]`, pat))
			for _, a := range as {

				//Extract link href
				href := htmlquery.SelectAttr(a, "href")
				next, _ := urlJoinPath(curr, href)

				//Set route value (if a capturing group exists)
				nvars := copyMap(vars)
				re := regexp.MustCompile(pat)
				match := re.FindStringSubmatch(href)
				if len(match) > 1 {
					val := match[1]
					nvars[p[1:]] = val
					log.Debugf("found route param value: %s = %s", p, val)
				}

				//Resume link building
				links = append(links, scrap(next, ps[i+1:], pats, nvars)...)
			}
			return links
		}

		//Continue navigation
		curr, _ = urlJoinPath(curr, p)
	}
	links = append(links, curr)
	return links
}

// Template a scrap regex pattern with a context map string
func scrapPattern(pat string, mp map[string]string) string {
	for k, v := range mp {
		pat = strings.Replace(pat, fmt.Sprintf(`\k<%s>`, k), v, -1)
	}
	return pat
}

// Extract hash
func scrapHashes(links []*Link, pats map[string]string) {
	for _, link := range links {
		iso := path.Base(link.Url)
		vars := make(map[string]string)
		vars["iso"] = iso
		log.Debugf("searching hash for: %s", iso)

		//Read manifest
		page, _ := urlJoinPath(link.Url, "..")
		pat := strings.Replace(pats[".hash.file"], `\k<iso>`, iso, -1)

		doc, err := htmlquery.LoadURL(page)
		if err != nil {
			log.Warnf("failed to open: <%s> (%s)", page, err)
			continue
		}

		a := htmlquery.FindOne(doc, fmt.Sprintf(`//a[matches(@href, '%s')]`, pat))
		if a == nil {
			log.Warnf("no match for hash file: <%s>", pat)
			continue
		}
		hfile := htmlquery.SelectAttr(a, "href")
		log.Debugf("found hash file: %s", hfile)

		//Fetch manifest
		hlink, _ := urlJoinPath(link.Url, "..", hfile)
		res, err := http.Get(hlink)
		if err != nil {
			log.Warnf("failed to open <%s> (%s)", hlink, err)
			continue
		}
		log.Debugf("read <%s>", hlink)
		defer res.Body.Close()

		//Read manifest
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Warnf("failed to read <%s> (%s)", hlink, err)
			continue
		}

		//Extract hash
		hash := RegexCapture(scrapPattern(pats[".hash.pattern"], vars), string(body))
		if hash != "" {
			link.Hash = fmt.Sprintf("%s:%s", pats[".hash.algo"], hash)
		}
	}
}

// Extract metadata
func scrapMeta(links []*Link, pats map[string]string) {
	for _, link := range links {
		log.Debugf("parsing meta: %s", link.Url)

		//Parse version (remove arch to avoid collisions)
		ea := regexp.MustCompile(REG_ARCH)
		version := RegexCapture(REG_VERSION, ea.ReplaceAllString(link.Url, "arch"))
		if version != "" {
			link.Version = version
			log.Debugf("found meta version: %s", link.Version)
		}

		//Parse arch
		arch := RegexCapture(REG_ARCH, link.Url)
		if arch != "" {
			link.Arch = arch
			log.Debugf("found meta arch: %s", link.Arch)
		}
	}
}
