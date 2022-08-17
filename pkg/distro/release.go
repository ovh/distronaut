package distro

import (
	log "github.com/sirupsen/logrus"
	"github.com/ovh/distronaut/internal/utils"
	"encoding/json"
	"gopkg.in/yaml.v3"
    "io/ioutil"
	"regexp"
	"fmt"
	"github.com/ovh/distronaut/pkg/distro/meta/distrowatch"
)

// Release structure
type Release struct {
	Source       string `json:source`
	Family       string `json:family`
	Distribution string `json:distribution`
	Website string `json:website`
	Documentation string `json:documentation`
	Status string `json:status`
	Versions []*Version `json:versions`
}

// Version structure
type Version struct {
	Url          string `json:url`
	Hash         string `json:hash`
	Version      string `json:version`
	Arch  		 string `json:arch`
	Meta		 map[string]string `json:meta`
}

// Source structure
type Source struct {
	Name string `yaml:string`
	Url string `yaml:string`
	Patterns map[string]string
}

//Fetch sources from configuration file
func Fetch(path string, filter string) ([]*Release) {
	var rs []*Release
	src, err := ListSources(path, filter)
	if err != nil {
		log.Errorf("failed to parse sources (%s)", err)
		return rs
	}
	for _, s := range src {
		rs = append(rs, fetch(s.Name, s.Url, s.Patterns))
	}
	return rs
}

//List sources from configuration file
func ListSources(path string, filter string) ([]*Source, error) {
	var src []*Source	
	var parsed []*Source	

	//Read configuration file
	file, err := ioutil.ReadFile(path)
    if err != nil {
        return src, err
    }
	err = yaml.Unmarshal(file, &parsed)
    if err != nil {
        return src, err
    }

	//Filter sources
	filter = fmt.Sprintf("(?i)%s", filter)
	for _, v := range parsed {
		if match, _ := regexp.MatchString(filter, v.Name); match {
			log.Debugf("retained source: %s", v.Name)
			src = append(src, v)
		}
	}
	return src, nil
}

// Fetch release and versions
func fetch(source string, uri string, pats map[string]string) (*Release) {
	// Fetch metadata
	if pats[".meta.source"] != "distrowatch" {
		log.Warnf("unsupported meta source: %s", pats[".meta.source"])
	}
	meta := distrowatch.About(pats[".meta.id"])
	r := &Release{Source:source, Family:meta["family"], Documentation:meta["documentation"], Website:meta["website"], Distribution:meta["distribution"], Status:meta["status"]}

	// Fetch links and metadata
	links := utils.Scrap(uri, pats)
	for _, link := range links {
		v := &Version{Url:link.Url, Hash:link.Hash, Version:link.Version, Arch:link.Arch}
		v.Meta = distrowatch.AboutVersion(pats[".meta.id"], utils.RegexCapture(pats[".meta.version"], link.Version))
		r.Versions = append(r.Versions, v)
	}

	// Debug
	log.Debugf("found: %+v", r)
	j, _ := json.Marshal(r)
	log.Debugf(string(j)) 

	return r
}