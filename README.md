# 🚀 Distronaut

**Distronaut** is a tool that travel through the internet to find distribution ISO download links and metadata, making it easier for you to monitor new releases or burn an install media.

## ⌨️ CLI

Use `fetch` command to retrieve a JSON from configured sources:
```bash
go run main.go fetch -c config/sources.yml -f 'debian'
```

Output is similar to below:
```json
[
  {
    "Source": "Debian",
    "Family": "Linux",
    "Distribution": "Debian (formerly Debian GNU/Linux)",
    "Website": "http://www.debian.org/",
    "Documentation": "http://www.debian.org/doc/",
    "Status": "Active",
    "Versions": [
      {
        "Url": "https://cdimage.debian.org/debian-cd/11.4.0/amd64/iso-cd/debian-11.4.0-amd64-netinst.iso",
        "Hash": "sha256:d490a35d36030592839f24e468a5b818c919943967012037d6ab3d65d030ef7f",
        "Version": "11.4.0",
        "Arch": "amd64",
        "Meta": {
          "release": "2021-08-14"
        }
      }
    ]
  }
]
```

Additional metadata is scrapped from [distrowatch.com](https://distrowatch.com).

## 🌕 Configuring new sources

Open your configuration file and add a new entry, containing:
- `name` is a friendly name for your new source (can also be used for filtering)
- `url` is the source url pattern
  - each *route parameter* is indicated by a colon (`:`)
- `patterns` is a map containing:
  - `:*` regexs patterns that are matched by *route parameter*
    - each *route parameter* can be back-referenced using `\k<:name>` syntax
  - `.hash.*` contains all hash related settings 
  - `.meta.*` contains all metadata related settings

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