package chapter05

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"math"
	"net/http"
	"os"
	"strings"
)

func Runner() {
	fmt.Println("*** Chapter 05 ***")
}

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func FindLinks() {
	fmt.Println("*** Running FindLinks ***")

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	if n.FirstChild != nil {
		links = visit(links, n.FirstChild)
	}

	if n.NextSibling != nil {
		links = visit(links, n.NextSibling)
	}

	return links
}

func ElementsMapping() {
	fmt.Println("*** Running ElementsMapping ***")

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mappingElements := make(map[string]int)
	mapElements(mappingElements, doc)

	for element, value := range mappingElements {
		fmt.Println(element, value)
	}
}

func mapElements(elementsMap map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		elementsMap[n.Data]++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		mapElements(elementsMap, c)
	}
}

func showElementContent(n *html.Node) {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return
	}
	if n.Type == html.ElementNode {
		fmt.Println(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		showElementContent(n)
	}
}

func CountWordsAndImages(url string) (int, int, error) {
	fmt.Println("*** Running CountWordsAndImages ***")

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return 0, 0, err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		err = fmt.Errorf("got error parsing HTML: %s", err)
		return 0, 0, err
	}
	words, images := countWordsAndImages(doc)
	return words, images, nil
}

func countWordsAndImages(n *html.Node) (int, int) {
	var words, images int

	if n.Type != html.TextNode || n.Data != "img" {
		return 0, 0
	}

	wordFreq := make(map[string]int)

	if n.Type == html.TextNode {
		input := bufio.NewScanner(strings.NewReader(n.Data))
		input.Split(bufio.ScanWords)

		for input.Scan() {
			word := input.Text()
			wordFreq[word]++
			words++
		}
	}

	if n.Data == "img" {
		images++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countWordsAndImages(n)
	}
	return words, images
}

func Surface() {
	fmt.Println("*** Running Surface ***")

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (sx, sy float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx = width/2 + (x-y)*cos30*xyscale
	sy = height/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

func f(x, y float64) (r float64) {
	r = math.Hypot(x, y) // distance from (0,0)
	r = math.Sin(r) / r
	return
}

func HTMLOutline(urls []string) error {
	fmt.Println("*** Running HTMLOutline ***")

	for _, url := range urls {
		err := outline(url)
		if err != nil {
			fmt.Printf("got error: %s", err.Error())
			return err
		}
	}
	return nil
}

func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	//!+call
	forEachNode(doc, startElement, endElement)
	//!-call

	return nil
}

// !+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

//!-forEachNode

// !+startend
var depth int

func startElement(n *html.Node) {
	var attributes string
	var child string

	for _, attribute := range n.Attr {
		attributes += fmt.Sprintf("%s=%s ", attribute.Key, attribute.Val)
	}

	if n.FirstChild == nil {
		child = " /"
	}

	if n.Type == html.ElementNode {
		if attributes == "" {
			fmt.Printf("%*s<%s>%s\n", depth*2, "", n.Data, child)
		} else {
			fmt.Printf("%*s<%s %s>%s\n", depth*2, "", n.Data, attributes, child)
		}
		depth++
	} else if n.Type == html.TextNode || n.Type == html.CommentNode {
		for _, line := range strings.Split(n.Data, "\n") {
			line = strings.TrimSpace(line)
			if line != "" && line != "\n" {
				fmt.Printf("%*s%s\n", depth*2, "", line)
			}
		}
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
