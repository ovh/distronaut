# 🚀 Distronaut

**Distronaut** is a tool that travel through the internet to find distribution ISO links and metadata, making it easier for you to monitor new releases, burn install medias or build your own ISO bank archive.

## ⌨️ CLI

### 🔨 Installation

Using a pre-built binary:
```bash
go install github.com/ovh/distronaut@latest
```

Building from source
```bash
make
```

### 🪛 Usage

Use `fetch` command to retrieve a JSON from configured sources:
```bash
distronaut fetch -c config/sources.yml -f 'debian'
```

Output is similar to below:
```json
[
  {
    "source": "Debian",
    "family": "Linux",
    "distribution": "Debian (formerly Debian GNU/Linux)",
    "website": "http://www.debian.org/",
    "documentation": "http://www.debian.org/doc/",
    "status": "Active",
    "logo": "https://distrowatch.com/images/yvzhuwbpy/debian.png",
    "logo64": "data:image/png;base64,...",
    "versions": [
      {
        "url": "https://cdimage.debian.org/cdimage/archive/12.6.0/amd64/iso-cd/debian-12.6.0-amd64-netinst.iso",
        "hash": "sha256:ade3a4acc465f59ca2496344aab72455945f3277a52afc5a2cae88cdc370fa12",
        "hashfile": "https://cdimage.debian.org/cdimage/archive/12.6.0/amd64/iso-cd/SHA256SUMS",
        "version": "12.6.0",
        "arch": "amd64",
        "meta": {
          "release": "2023-06-10"
        }
      }
    ]
  }
]
```

Additional metadata are scrapped from [distrowatch.com](https://distrowatch.com).

## 👨‍💻 Programmatic usage

This package can also be imported within another golang codebase:
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ovh/distronaut/pkg/distro"
)

func main() {
	src, _ := distro.FetchSources("config/sources.yml", "debian")
	j, _ := json.MarshalIndent(src, "", "  ")
	fmt.Println(string(j))
}
```

## 🌕 Configuring new sources

Open your configuration file and add a new entry, containing:
- `name` is a friendly name for your new source (can also be used for filtering)
- `url` is the source URL pattern
  - each *route parameter* is indicated by a colon (`:`)
- `patterns` is a map containing:
  - `:*` are regex patterns that are matched by defined *route parameter*
    - each *route parameter* can be back-referenced using `\k<:name>` syntax
  - `.hash.*` contains all hash related settings 
    - `.hash.file` is a regex pattern matching the file containing hashes 
    - `.hash.pattern` is a regex pattern capturing the hash from a given iso (that can be back-referenced with `\k<iso>`)
    - `.hash.algo` is the name of the hash algorithm
  - `.meta.*` contains all metadata related settings
    - `.meta.source` must be set to `distrowatch` (only metadata source supported for now)
    - `.meta.id` is the distro handle on distrowatch
    - `.meta.version` is a regex pattern matching the version as it is referenced on distrowatch

Example source:
```yml
- name: Debian
  url: https://cdimage.debian.org/debian-cd/:version/:arch/iso-cd/:iso
  patterns:
    :version: ^(\d+\.\d+(?:\.\d+)?)\/$
    :arch: ^(amd64|arm64)\/$
    :iso: ^debian-\k<version>-\k<arch>-netinst\.iso$
    .hash.file: SHA256SUMS
    .hash.algo: sha256
    .hash.pattern: (?m)^([0-9a-f]{64})\s+\k<iso>
    .meta.source: distrowatch
    .meta.id: debian
    .meta.version: (\d+)
```

# 💪 Contributing

Please read our contribution guidelines first ([CONTRIBUTING.md](https://github.com/ovh/distronaut/blob/master/CONTRIBUTING.md)).

## 🧪 Testing

Run tests using:
```bash
make test
```

A mocked server will temporary be spawned on port 3000 to avoid performing real network requests.

# 📜 License
 
```
Copyright 2021 OVH SAS
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
    http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
