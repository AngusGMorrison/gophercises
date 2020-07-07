package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also implements http.Handler) that will
// attempt to map any paths (keys in the map) to their corresponding URL (values that each key
// in the map points to, in string format). If the path is not provided in the map, then the
// fallback http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will attempt to open and parse the provided YAML file and then return an
// http.HandlerFunc (which also implements http.Handler) that will attempt to map any paths to
// their corresponding URL. If the path is not provided in the YAML, then the fallback http.Handler
// will be called instead.
func YAMLHandler(redirectData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return fileHandler(redirectData, yaml.Unmarshal, fallback)
}

// JSONHandler will attempt to open and parse the provided YAML file and then return an
// http.HandlerFunc (which also implements http.Handler) that will attempt to map any paths to
// their corresponding URL. If the path is not provided in the YAML, then the fallback http.Handler
// will be called instead.
func JSONHandler(redirectData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return fileHandler(redirectData, json.Unmarshal, fallback)
}

type unmarshaller func(data []byte, v interface{}) error

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func fileHandler(data []byte, unmarshal unmarshaller, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseRedirectData(data, unmarshal)
	if err != nil {
		return nil, err
	}
	pathMap := buildPathMap(pathURLs)
	return MapHandler(pathMap, fallback), nil
}

func parseRedirectData(data []byte, unmarshal unmarshaller) ([]pathURL, error) {
	var pathURLs []pathURL
	err := unmarshal(data, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

func buildPathMap(pathURLs []pathURL) map[string]string {
	pathMap := make(map[string]string)
	for _, pu := range pathURLs {
		pathMap[pu.Path] = pu.URL
	}
	return pathMap
}
