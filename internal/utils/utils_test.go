package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScrap(t *testing.T) {
	//Scrap link and hash
	{
		a, _ := Scrap("http://0.0.0.0:3000/b/hash/:iso", map[string]string{
			":iso":          `^(.*\.iso)$`,
			".hash.file":    "checksum",
			".hash.algo":    "sha256",
			".hash.pattern": `(?m)^([0-9a-f]{64})\s+\k<iso>`,
		}, nil)
		assert.Equal(t, a, []*Link{
			&Link{
				Url:      "http://0.0.0.0:3000/b/hash/distro-1.0-amd64.iso",
				Hash:     "sha256:8fb3ccc2bce9c1a7ad9ee470c3aa53486ed42d57d3bf6f6346718e4f958954d2",
				Hashfile: "http://0.0.0.0:3000/b/hash/checksum",
				Version:  "1.0",
				Arch:     "amd64",
			}})
	}
	//Scrap link without hash settings
	{
		a, _ := Scrap("http://0.0.0.0:3000/b/hash/:iso", map[string]string{
			":iso": `^(.*\.iso)$`,
		}, nil)
		assert.Equal(t, a, []*Link{
			&Link{
				Url:     "http://0.0.0.0:3000/b/hash/distro-1.0-amd64.iso",
				Version: "1.0",
				Arch:    "amd64",
			}})
	}
}
