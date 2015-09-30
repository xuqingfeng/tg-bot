package google
import (
	"net/http"
	"heroku.com/tg-bot/util"
	"io/ioutil"
	"net/url"
	"strings"
	"golang.org/x/net/html"
)

func SearchInGoogle(query string) ([]string) {


	resp, err := http.Get("https://www.google.com/search?q="+url.QueryEscape(query))
	util.PanicIf(err)

	ret, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	util.PanicIf(err)

	doc, _ := html.Parse(strings.NewReader(string(ret)))

	var links []string
	getLinks(doc, &links)

	return links
}

func getLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" && n.Parent.Data != "li" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if strings.HasPrefix(a.Val, "/url?q=") {
					fullUrl := strings.Split(a.Val, "&sa")
					*links = append(*links, strings.TrimLeft(fullUrl[0], "/url?q="))
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getLinks(c, links)
	}
}