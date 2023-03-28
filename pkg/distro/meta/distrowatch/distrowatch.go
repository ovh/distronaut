package distrowatch

import (
	"encoding/base64"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

// Metadata about distribution
func About(id string) map[string]string {
	meta := make(map[string]string)
	if id == "" {
		return meta
	}

	//Fetch distrowatch page
	page := fmt.Sprintf("https://distrowatch.com/table.php?distribution=%s", id)
	doc, err := htmlquery.LoadURL(page)
	if err != nil {
		log.Warnf("failed to open: <%s> (%s)", page, err)
		return meta
	}

	//OS family
	n := htmlquery.FindOne(doc, `//li/b[text() = 'OS Type:']/../a`)
	if n != nil {
		meta["family"] = htmlquery.InnerText(n)
	}

	//Status
	n = htmlquery.FindOne(doc, `//li/b[text() = 'Status:']/../font`)
	if n != nil {
		meta["status"] = htmlquery.InnerText(n)
	}

	//Distribution
	n = htmlquery.FindOne(doc, `//th[text() = 'Distribution']/../td`)
	if n != nil {
		meta["distribution"] = htmlquery.InnerText(n)
	}

	//Website
	n = htmlquery.FindOne(doc, `//th[text() = 'Home Page']/../td`)
	if n != nil {
		meta["website"] = htmlquery.InnerText(n)
	}

	//Documentation
	n = htmlquery.FindOne(doc, `//th[text() = 'Documentation']/../td`)
	if n != nil {
		lines := strings.Split(strings.TrimSpace(htmlquery.InnerText(n)), "\n")
		links := strings.Split(lines[0], "•")
		meta["documentation"] = strings.TrimSpace(links[0])
	}

	//Logo
	n = htmlquery.FindOne(doc, `//*[@class = 'TablesTitle']//img[@hspace = '32']`)
	if n != nil {
		meta["logo"] = fmt.Sprintf("https://distrowatch.com/%s", htmlquery.SelectAttr(n, "src"))
	}

	//Logo base64
	if meta["logo"] != "" {
		res, err := http.Get(meta["logo"])
		if res.StatusCode == 200 && err == nil {
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err == nil {
				b64 := base64.StdEncoding.EncodeToString(body)
				meta["logo64"] = fmt.Sprintf("data:image/png;base64,%s", b64)
			}
		}
	}

	log.Debugf("about %s: %+v", id, meta)
	return meta
}

// Metadata about version
func AboutVersion(id string, version string) map[string]string {
	meta := make(map[string]string)
	if id == "" || version == "" {
		return meta
	}
	//Fetch distrowatch page
	page := fmt.Sprintf("https://distrowatch.com/table.php?distribution=%s", id)
	doc, err := htmlquery.LoadURL(page)
	if err != nil {
		log.Warnf("failed to open: <%s> (%s)", page, err)
		return meta
	}

	//Search version column index
	expr, err := xpath.Compile(fmt.Sprintf(`count((//table//td[@class='TablesInvert'][contains(text(), '%s')])[1]/preceding-sibling::*)`, version))
	if err != nil {
		log.Warnf("failed to compile xpath query (%s)", err)
		return meta
	}
	i := int(expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64))
	if i == 0 {
		log.Debugf("no match for <%s>", version)
		return meta
	}
	log.Debugf("distrowatch column for %s: %d", id, i)

	//Release date
	n := htmlquery.FindOne(doc, fmt.Sprintf(`//table//th[text() = 'Release Date']/../td[%d]`, i))
	if n != nil {
		meta["release"] = strings.TrimSpace(htmlquery.InnerText(n))
	}

	//EOL date
	n = htmlquery.FindOne(doc, fmt.Sprintf(`//table//th[text() = 'End Of Life']/../td[%d]`, i))
	if n != nil {
		meta["eol"] = strings.TrimSpace(htmlquery.InnerText(n))
	}

	log.Debugf("about %s (%s): %+v", id, version, meta)
	return meta
}
