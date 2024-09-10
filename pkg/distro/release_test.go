package distro

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
	"testing"
)

var sources string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	sources = path.Join(filename, "../../..", "tests/sources.yml")
}

func TestSetLogLevel(t *testing.T) {
	//Change log level
	{
		p := log.GetLevel()
		SetLogLevel(log.DebugLevel)
		assert.Equal(t, log.GetLevel(), log.DebugLevel)
		SetLogLevel(p)
		assert.Equal(t, log.GetLevel(), p)
	}
}

func TestListSources(t *testing.T) {
	//List source
	{
		a, err := ListSources(sources, "Test source listing")
		assert.Nil(t, err)
		assert.Len(t, a, 1)
		assert.Equal(t, *a[0], Source{
			Name: "Test source listing",
			Url:  "http://0.0.0.0:3000",
			Patterns: map[string]string{
				".test": "test",
			},
		})
	}
	//List sources fail because of inexsting file
	{
		_, err := ListSources(path.Join(sources, "../sources..yml"), "Test source listing")
		assert.NotNil(t, err)
	}
	//List sources fail because of invalid source file
	{
		_, err := ListSources(path.Join(sources, "../sources.bad.yml"), "Test source listing")
		assert.NotNil(t, err)
	}
}

func TestFetchSources(t *testing.T) {
	//Fecth sources
	{
		a, err := FetchSources(sources, "Test source fetching")
		assert.Nil(t, err)
		assert.Len(t, a, 1)
		assert.Len(t, a[0].Versions, 2)
		assert.Equal(t, *a[0], Release{
			Source:        "Test source fetching",
			Family:        "",
			Distribution:  "",
			Website:       "",
			Documentation: "",
			Status:        "",
			Versions: []*Version{
				&Version{
					Url:      "http://0.0.0.0:3000/a/distro/server/releases/1.0.0/amd64/distro-server-1.0.0-amd64.iso",
					Hash:     "sha256:9c004c39885a275f63b9f3dff6e775d8a9de663dc71e53f32efb9fe214b4b3e1",
					Hashfile: "http://0.0.0.0:3000/a/distro/server/releases/1.0.0/amd64/checksum",
					Version:  "1.0.0",
					Arch:     "amd64",
					Meta:     map[string]string{},
				},
				&Version{
					Url:      "http://0.0.0.0:3000/a/distro/server/releases/2.0.0/amd64/distro-server-2.0.0-amd64.iso",
					Hash:     "sha256:90a507f511e42581a17a80d648915519f4b5bdad39e443f1ae30fc9a1b3c5c24",
					Hashfile: "http://0.0.0.0:3000/a/distro/server/releases/2.0.0/amd64/checksum",
					Version:  "2.0.0",
					Arch:     "amd64",
					Meta:     map[string]string{},
				},
			},
		})
	}

	//Fetch sources fail because of invalid source file
	{
		_, err := FetchSources(path.Join(sources, "../sources.bad.yml"), ".*")
		assert.NotNil(t, err)
	}
}
