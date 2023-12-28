package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercises/urlshort"
)

func main() {
	var isJson bool

	ymlfilePath := flag.String("yml", "", "Path to YAML file")
	jsonfilePath := flag.String("jso", "", "Path to json file")
	flag.Parse()
	var ymlcontent []byte
	var jsoncontent []byte


	if *ymlfilePath != "" {
		fmt.Println("FilePath: ", *ymlfilePath)
		fileContent,err := os.ReadFile(*ymlfilePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(fileContent))
		ymlcontent = fileContent
		isJson = false
	}

	if *jsonfilePath != "" {
		fmt.Println("FilePath: ", *jsonfilePath)
		fileContent,err := os.ReadFile(*jsonfilePath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(fileContent))
		jsoncontent = fileContent
		isJson = true
	}

	



	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
// 	yaml := `
// - path: /urlshort
//   url: https://github.com/gophercises/urlshort
// - path: /urlshort-final
//   url: https://github.com/gophercises/urlshort/tree/solution

	if isJson {
		jsonHandler, err := urlshort.JSONHandler([]byte(jsoncontent), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	}


	yamlHandler, err := urlshort.YAMLHandler([]byte(ymlcontent), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
