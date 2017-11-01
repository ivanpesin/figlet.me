package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satori/go.uuid"
)

const rID = 0

var figlet = os.Getenv("FIGLET_BIN")
var figletDir = os.Getenv("FIGLET_DIR")

var figletFonts = make(map[string]string)

func execFiglet(f, s string) string {
	cmd := exec.Command(figlet, "-f", figletFonts[f], s)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func getFidgetFonts() (r []string) {
	for k := range figletFonts {
		r = append(r, k)
	}
	return
}

func fidgetHandler(w http.ResponseWriter, r *http.Request) {

	id := uuid.NewV4().String()[:6]
	log.Printf("| %6s | %s %s ", id, r.Method, r.URL.Path)

	w.Header().Set("Content-type", "text/plain")

	switch {
	case r.URL.Path == "/status":
		io.WriteString(w, "All OK\n")
	case r.URL.Path == "/getFonts":
		io.WriteString(w, strings.Join(getFidgetFonts(), "\n"))
	case r.URL.Path == "/":
		var params url.Values
		if r.Method == "GET" {
			params = r.URL.Query()
		} else {
			r.ParseForm()
			params = r.Form
		}

		font := params.Get("font")
		if font == "" {
			font = "banner"
		}
		text := params.Get("text")
		if text == "" {
			text = "hello, world"
		}

		log.Printf("| %6s | --> font: %s", id, font)
		log.Printf("| %6s | --> text: %s", id, text)

		if _, ok := figletFonts[font]; !ok {
			log.Printf("| %6s | 404 | No such font", id)
			w.WriteHeader(http.StatusNotImplemented)
			io.WriteString(w, "404 No such font\n")
			return
		}

		io.WriteString(w, execFiglet(font, text))
	default:
		log.Printf("| %6s | 404 | Not found", id)
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "404 Not found\n")
		return
	}

	log.Printf("| %6s | 200 | Done", id)

}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.SetPrefix("[figlet.me] ")

	ip := flag.String("ip", "0.0.0.0", "IP address to serve on")
	port := flag.Int("port", 8080, "port to listen on")

	flag.Parse()

	// populate fonts
	cmd := exec.Command("/usr/bin/find", figletDir, "-name", "*.flf")
	out, _ := cmd.CombinedOutput()
	r := strings.Split(string(out), "\n")
	for _, file := range r {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			figletFonts[strings.TrimSuffix(filepath.Base(file), ".flf")] = file
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", fidgetHandler)

	log.Printf("Starting server on %s:%d", *ip, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *ip, *port), mux))
}
