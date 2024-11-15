package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(m map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := m[r.URL.Path]; ok {
			http.Redirect(w, r, val, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// PathURL defines the structure of an object in a YAML array file.
type PathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// ParseYAML converts YAML data from bytes to an array of PathURL objects.
// If for any reason the data cannot be converted, an error is returned.
func ParseYAML(yml []byte) ([]PathURL, error) {

	var m []PathURL

	err := yaml.Unmarshal(yml, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// BuildMap converts an array of PathURL objects into a map.
func BuildMap(data []PathURL) map[string]string {

	m := make(map[string]string)

	for _, val := range data {
		m[val.Path] = val.URL
	}

	return m
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	parsedYAML, err := ParseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathMap := BuildMap(parsedYAML)

	return MapHandler(pathMap, fallback), nil
}
