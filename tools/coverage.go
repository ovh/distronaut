package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var path string
	var file string
	var threshold float64
	var code bool
	var verbose bool

	flag.StringVar(&path, "path", "./...", "go test fileglob")
	flag.StringVar(&file, "file", "coverage", "coverage file")
	flag.Float64Var(&threshold, "threshold", 80.0, "coverage threshold")
	flag.BoolVar(&code, "code", true, "print coverage")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.Parse()

	html := fmt.Sprintf("%s.html", file)
	unsatisfactory := map[string]float64{}

	// Perform coverage
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := run(pwd, fmt.Sprintf("go test %s -coverprofile %s", path, file), verbose); err != nil {
		panic(err)
	}
	if err := run(pwd, fmt.Sprintf("go tool cover -html %s -o %s", file, html), verbose); err != nil {
		panic(err)
	}

	// Read results
	doc, err := htmlquery.LoadDoc(html)
	if err != nil {
		panic(err)
	}
	nodes := htmlquery.Find(doc, `//select[@id="files"]/option`)
	for _, node := range nodes {
		id := htmlquery.SelectAttr(node, "value")
		head := strings.TrimSpace(htmlquery.InnerText(node))

		// Coverage threshold
		matches := regexp.MustCompile(`([\s\S]+)\s\((\d+(?:\.\d+))%\)$`).FindStringSubmatch(head)
		file := matches[1]
		covered, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			panic(err)
		}
		if covered >= threshold {
			color.New(color.FgWhite, color.Bold, color.BgGreen).Print(fmt.Sprintf("✓ | %.1f%% | %s", covered, file))
		} else {
			color.New(color.FgWhite, color.Bold, color.BgRed).Print(fmt.Sprintf("✗ | %.1f%% | %s", covered, file))
			unsatisfactory[file] = covered
		}
		color.New().Print("\n")

		// Coverage code
		if code {
			content := htmlquery.FindOne(doc, fmt.Sprintf(`//pre[@id="%s"]`, id))
			for child := content.FirstChild; child != nil; child = child.NextSibling {
				text := htmlquery.InnerText(child)
				switch htmlquery.SelectAttr(child, "class") {
				case "cov8":
					color.New(color.FgGreen).Print(text)
				case "cov0":
					color.New(color.FgRed).Print(text)
				default:
					color.New(color.FgHiBlack).Print(text)
				}
			}
		}
	}

	// Apply threshold
	if len(unsatisfactory) > 0 {
		color.New(color.FgRed).Println(fmt.Sprintf("\n%d file(s) don't meet required threshold of %.1f%%", len(unsatisfactory), threshold))
		os.Exit(2)
	}
}

func run(pwd string, cmd string, verbose bool) error {
	if verbose {
		color.New(color.FgMagenta, color.Bold).Println(fmt.Sprintf(">>> %s", cmd))
	}
	c := exec.Command("/bin/bash", "-c", cmd)
	if pwd != "" {
		c.Dir = pwd
	}
	stdio, err := c.CombinedOutput()
	if verbose || err != nil {
		color.New(color.FgHiBlack).Println(strings.TrimSpace(string(stdio)))
	}
	return err
}
