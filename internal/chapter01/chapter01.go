package chapter01

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	_ "io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	blackIndex = 1
)

func Runner() {
	fmt.Println("*** Chapter 01 ***")
	ShowCommandCaller()
	ShowArgs()
	InspectExecutionTimes()
	//ChangeColorToGreen()
	//FetchResponseToStdout()
	//FetchURLWithPrefix()
	//FetchWithResponseStatus()
	SetColorsWithURLParams()
}

func ShowCommandCaller() {
	fmt.Println("*** Running ShowCommandCaller ***")
	fmt.Println(os.Args[0])
}

func ShowArgs() {
	fmt.Println("*** Running ShowArgs ***")

	for i, arg := range os.Args {
		fmt.Printf("[%d] %s\n", i, arg)
	}
}

func InspectExecutionTimes() {
	fmt.Println("*** Running InspectExecutionTimes ***")

	start1 := time.Now()
	var s1, sep1 string

	for i := 1; i < len(os.Args); i++ {
		s1 += sep1 + os.Args[i]
		sep1 = " "
	}
	fmt.Println(s1)
	fmt.Println("Start time 1: ", time.Since(start1).Microseconds(), "ms")

	start2 := time.Now()
	s2, sep2 := "", ""
	for _, arg := range os.Args[1:] {
		s2 += sep2 + arg
		sep2 = " "
	}
	fmt.Println(s2)
	fmt.Println("Start time 2: ", time.Since(start2).Microseconds(), "ms")

	start3 := time.Now()
	fmt.Println(strings.Join(os.Args[1:], " "))
	fmt.Println("Start time 3: ", time.Since(start3).Microseconds(), "ms")
}

func ChangeColorToGreen() {
	fmt.Println("*** Running ChangeColorToGreen ***")

	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	handler := func(w http.ResponseWriter, r *http.Request) {
		lissajousGreen(w)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func lissajousGreen(out io.Writer) {
	palette := []color.Color{color.RGBA{0, 187, 0, 255}, color.Black}
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	_ = gif.EncodeAll(out, &anim)
}

func FetchResponseToStdout() {
	fmt.Println("*** Running FetchResponseToStdout ***")

	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

func FetchURLWithPrefix() {
	fmt.Println("*** Running FetchURLWithPrefix ***")

	for _, url := range os.Args[1:] {
		resp, err := fetchURL(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}

func fetchURL(url string) (*http.Response, error) {
	if !strings.HasPrefix(url, "http") {
		fmt.Println("*** Adding http before url ***")
		url = fmt.Sprintf("http://%s", url)
	}
	resp, err := http.Get(url)
	return resp, err
}

func FetchWithResponseStatus() {
	fmt.Println("*** Running FetchURLWithPrefix ***")

	for _, url := range os.Args[1:] {
		resp, err := fetchURL(url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			fmt.Printf("*** Request Status: %v ***\n", resp.StatusCode)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			fmt.Printf("*** Request Status: %v ***\n", resp.StatusCode)
			os.Exit(1)
		}
		fmt.Printf("\n*** Request Status: %v ***\n", resp.StatusCode)
	}
}

func SetColorsWithURLParams() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		lissajous(w, r)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(w http.ResponseWriter, r *http.Request) {
	var palette = []color.Color{color.White, color.Black}
	var cycles float64 = 5 // number of complete x oscillator revolutions
	const (
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	urlPath, _ := url.Parse(r.URL.Path)
	params, _ := url.ParseQuery(urlPath.RawQuery)
	if params.Has("cycles") {
		value := params.Get("cycles")
		cycles, _ = strconv.ParseFloat(value, 64)
	}

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	_ = gif.EncodeAll(w, &anim)
}
