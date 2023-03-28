package distro

import (
	"encoding/json"
	"fmt"
	"github.com/ovh/distronaut/internal/utils"
	"github.com/ovh/distronaut/pkg/distro/meta/distrowatch"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Release structure
type Release struct {
	Source        string     `json:"source"`
	Family        string     `json:"family"`
	Distribution  string     `json:"distribution"`
	Website       string     `json:"website"`
	Documentation string     `json:"documentation"`
	Status        string     `json:"status"`
	Logo          string     `json:"logo"`
	Logo64        string     `json:"logo64"`
	Versions      []*Version `json:"versions"`
}

// Version structure
type Version struct {
	Url     string            `json:"url"`
	Hash    string            `json:"hash"`
	Version string            `json:"version"`
	Arch    string            `json:"arch"`
	Meta    map[string]string `json:"meta"`
}

// Source structure
type Source struct {
	Name     string            `yaml:"name"`
	Url      string            `yaml:"url"`
	Patterns map[string]string `yaml:"patterns"`
	Export   string            `yaml:"export"`
}

// Set log level
func SetLogLevel(lv log.Level) {
	log.SetLevel(lv)
}

// Fetch sources from configuration file with progress bar
func FetchSourcesWithProgress(path string, filter string) ([]*Release, error) {
	bar := progressbar.NewOptions(1,
		progressbar.OptionSetDescription("distronaut"),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(16),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(false),
	)
	return fetchSources(path, filter, bar)
}

// Fetch sources from configuration file
func FetchSources(path string, filter string) ([]*Release, error) {
	return fetchSources(path, filter, nil)
}

// Fetch sources from configuration file
func fetchSources(path string, filter string, bar *progressbar.ProgressBar) ([]*Release, error) {
	var rs []*Release
	src, err := ListSources(path, filter)
	if err != nil {
		log.Errorf("failed to parse sources (%s)", err)
		return rs, err
	}
	if bar != nil {
		bar.ChangeMax(bar.GetMax() + len(src))
	}
	for _, s := range src {
		if bar != nil {
			bar.Describe(fmt.Sprintf("fetching: [cyan]%s[reset]", s.Name))
		}
		r, err := fetch(s.Name, s.Url, s.Patterns, bar)
		if bar != nil {
			bar.Add(1)
		}
		if err != nil {
			log.Errorf("failed to fetch source: <%s> (%s)", s.Name, err)
			continue
		}
		rs = append(rs, r)

		if s.Export != "" {
			bar.Describe(fmt.Sprintf("exporting: [cyan]%s[reset] to [cyan]%s[reset] ", s.Name, s.Export))
			exp := []*Release{r}
			if err := export(s.Export, exp); err != nil {
				log.Errorf("failed to export source: <%s> (%s)", s.Name, err)
				continue
			}
		}
	}
	if len(rs) < len(src) {
		err = fmt.Errorf("%d/%d sources successfully fetched", len(rs), len(src))
	}
	if bar != nil {
		bar.Clear()
	}
	return rs, err
}

// List sources from configuration file
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
func fetch(source string, uri string, pats map[string]string, bar *progressbar.ProgressBar) (*Release, error) {
	// Fetch metadata
	if pats[".meta.source"] != "" && pats[".meta.source"] != "distrowatch" {
		log.Warnf("unsupported meta source: %s", pats[".meta.source"])
	}
	meta := distrowatch.About(pats[".meta.id"])
	r := &Release{Source: source, Family: meta["family"], Documentation: meta["documentation"], Website: meta["website"], Distribution: meta["distribution"], Status: meta["status"], Logo: meta["logo"], Logo64: meta["logo64"]}

	// Fetch links and metadata
	links, err := utils.Scrap(uri, pats, bar)
	if err != nil {
		return nil, err
	}
	if bar != nil {
		bar.ChangeMax(bar.GetMax() + len(links))
	}
	for _, link := range links {
		v := &Version{Url: link.Url, Hash: link.Hash, Version: link.Version, Arch: link.Arch}
		if bar != nil {
			bar.Describe(fmt.Sprintf("checking: [cyan]%s[reset]", v.Version))
		}
		v.Meta = distrowatch.AboutVersion(pats[".meta.id"], utils.RegexCapture(pats[".meta.version"], link.Version))
		if bar != nil {
			bar.Add(1)
		}
		r.Versions = append(r.Versions, v)
	}

	// Debug
	log.Debugf("found: %+v", r)
	j, _ := json.Marshal(r)
	log.Debugf(string(j))

	return r, nil
}

func export(file string, res []*Release) error {
	j, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(j); err != nil {
		return err
	}

	return nil
}
