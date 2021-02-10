package pathparse

import (
	"bytes"
	"strings"
	"text/tabwriter"

	"github.com/gertd/go-pluralize"
)

// inflect pluralizes words.
var inflect = pluralize.NewClient()

// A Route is a Rails-style definition of a route.
type Route struct {
	// The prefix is the Rails-style key of the URL path, ie, `root_path`
	prefix string

	// The HTTP verb for this route.
	verb string

	// The pattern is the URL pattern of the route.
	pattern string
}

// prepend prepends the given string to the given slice.
func prepend(slice []string, s string) []string {
	return append([]string{s}, slice...)
}

// Routes returns a human readable string containing all routes.
func Routes(routes []*Route) string {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 1, 4, 1, ' ', 0)

	w.Write([]byte("Prefix\t"))
	w.Write([]byte("Verb\t"))
	w.Write([]byte("URI Pattern\n"))

	for _, r := range routes {
		w.Write([]byte(r.prefix))
		w.Write([]byte("\t"))
		w.Write([]byte(r.verb))
		w.Write([]byte("\t"))
		w.Write([]byte(r.pattern))
		w.Write([]byte("\n"))
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}

	return buf.String()
}

// ParseRoute generates the Rails-style route prefix from an HTTP verb and a
// URL path.
func ParseRoute(verb, pattern string) *Route {
	if pattern == "/" {
		return &Route{
			prefix:  "root",
			verb:    verb,
			pattern: pattern,
		}
	}

	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	paths := strings.Split(pattern, "/")
	words := make([]string, 0)

	last := len(paths) - 1

	for i, p := range paths {
		// If this is the last element, append it.
		if i == last {
			if inflect.IsSingular(p) {
				words = prepend(words, p)
				continue
			}
			words = append(words, p)
			continue
		}

		if p == ":id" {
			continue
		}

		if inflect.IsPlural(p) {
			words = append(words, inflect.Singular(p))
			continue
		}

		words = append(words, p)
	}

	return &Route{
		prefix:  strings.Join(words, "_"),
		verb:    verb,
		pattern: pattern,
	}
}
