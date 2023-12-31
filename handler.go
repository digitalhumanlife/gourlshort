package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok{
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type YamlData []struct{
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

type JsonData [] struct {
	Path string `json:"path"`
	URL string `json:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYML,_ := parseYML(yml)
	pathMap := buildMap(parsedYML)
	return MapHandler(pathMap, fallback),nil
}

func JSONHandler(jsonByte []byte, fallback http.Handler) (http.HandlerFunc, error){
	var jsonData JsonData
	err := json.Unmarshal([]byte(jsonByte), &jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return nil, err
	}

	pathJson := make(map[string]string)
	for _, pd := range jsonData {
		pathJson[pd.Path] = pd.URL
	}
	return MapHandler(pathJson, fallback), nil
}



func parseYML(yml []byte) (YamlData, error){
	var ymlData YamlData
	err := yaml.Unmarshal(yml, &ymlData)
	if err != nil {
		return nil, err
	}
	return ymlData, nil
}

func buildMap(ym YamlData) map[string]string{
	pathMap := make(map[string]string)
	for _, pd := range ym {
		pathMap[pd.Path] = pd.URL 
	}
	return pathMap
}
