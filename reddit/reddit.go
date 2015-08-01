package reddit
import (
	"net/http"
	"io/ioutil"
	"golang.org/x/net/html"
	"strings"
)

func GetPhoto()([]string) {

	resp, _ := http.Get("https://www.reddit.com/r/funny/")

	ret, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	doc, _ := html.Parse(strings.NewReader(string(ret)))
	var links []string
	getPhotoLinks(doc, &links)

	return links
}

func getPhotoLinks(n *html.Node, links *[]string) {

	if n.Type == html.ElementNode && n.Data == "a" && n.Parent.Data == "p" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.Contains(a.Val, "imgur.com"){
				*links = append(*links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getPhotoLinks(c, links)
	}
}