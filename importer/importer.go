// Package importer provides a jsonnet Importer that loads files from https
// URLs or from the filesystem. It maintains a slice of search paths on which
// it searches for imports, allowing absolute and relative imports.
//
// A search path element may start with a double-slash (`//`) which is
// interpreted as a HTTPS url without the `https:` scheme prefix. A search path
// element starting with a single-slash (`/`) is an absolute path in the
// processes filesystem. A search path without a leading slash is resolved
// relative to the current working directory of the process.
//
// Only HTTPS is supported for network paths.
package importer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
)

var noContent = jsonnet.Contents{} //nolint:gochecknoglobals

// A URLFetcher retreives a URL returning a http.Response or an error. It
// is defined such that http.Client implements it, but allows a different
// implementation or a custom-configured http.Client to be provided to
// the Importer.
type URLFetcher interface {
	Get(url string) (*http.Response, error)
}

// Importer implements the jsonnet.Importer interface, allowing jsonnet code to
// be imported via https in addition to local files. Filenames starting with a
// double-slash (`//`) are fetched via HTTPS. Otherwise the path is treated as
// a local filesystem path.
//
// Once an import path is successfully fetched, either with data or a definitive
// not found result, that result is cached for the lifetime of the Importer.
// This is a requirement of the jsonnet.Importer interface so it is not possible
// for the same import statement from different files to result in different
// content. If an Importer is shared across multiple jsonnet.VM instances,
// the the cache will be shared too. There is no cache expiry logic.
type Importer struct {
	// SearchPath is an ordered slice of paths (network or local filesystem)
	// that is prepended to the imported filename if the filename is not
	// found. Searching stops when it is found.
	SearchPath []string

	// Fetcher is the URLFetcher used to fetch paths. The default is a
	// http.Client that has a `file:` scheme handler.
	Fetcher URLFetcher

	cache map[string]jsonnet.Contents
}

// AppendSearchFromEnv appends a list of search paths specified in the given
// environment variable to the search path list. The elements of the path in
// the variable are separated by the filepath.SplitList() delimiter.
func (i *Importer) AppendSearchFromEnv(envvar string) {
	for _, p := range filepath.SplitList(os.Getenv(envvar)) {
		if p == "" {
			continue
		}

		i.SearchPath = append(i.SearchPath, p)
	}
}

// Import loads imp from a file or a network location. If imp is a relative
// path, search for it relative to the directory of source and the search path
// elements. If the import found, return its contents and the absolute location
// where it was found. If it was not found, or there was an error reading
// the content, return the error.
//
// Import will cache the result and return it the next time that path is
// requested.
//
// This method is defined in the jsonnet.Importer interface:
//   https://godoc.org/github.com/google/go-jsonnet#Importer
func (i *Importer) Import(source, imp string) (jsonnet.Contents, string, error) {
	dir, _ := path.Split(source)
	content, location, err := i.search(imp, dir)

	if err == nil && content == noContent {
		err = fmt.Errorf("couldn't import %#v: not found", imp)
	}

	return content, location, err
}

func (i *Importer) search(imp, dir string) (jsonnet.Contents, string, error) {
	if path.IsAbs(imp) {
		content, err := i.readViaCache(imp)
		return content, imp, err
	}

	// try to import imp relative to source first, then the search path
	for _, p := range append([]string{dir}, i.SearchPath...) {
		location := path.Join(p, imp)
		content, err := i.readViaCache(location)
		// content found, or an error. Stop searching - we're done
		// TODO(camh): Keep searching on hard errors. If a fetch results
		// in a hard error, then no amount of fetching that resource
		// will ever work with the path we have. We should treat this
		// as "not found" and keep searching. A soft error, where it
		// may succeed later (transient failures) should make the import
		// fail. If we keep searching and return a result, the next time
		// we import we may return a different result if the transient
		// error has been resolved.
		if content != noContent || err != nil {
			return content, location, err
		}
	}

	return noContent, "", nil // not found
}

func (i *Importer) readViaCache(imp string) (jsonnet.Contents, error) {
	if i.cache == nil {
		i.cache = make(map[string]jsonnet.Contents)
	}

	if content, ok := i.cache[imp]; ok {
		return content, nil
	}

	content, err := fetch(addScheme(imp), i.fetcher())
	if err == nil {
		i.cache[imp] = content
	}

	return content, err
}

func (i *Importer) fetcher() URLFetcher {
	if i.Fetcher == nil {
		// Create a http.Client with a transport handling file:// urls
		// https://golang.org/pkg/net/http/#NewFileTransport
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

		i.Fetcher = &http.Client{Transport: t}
	}

	return i.Fetcher
}

func addScheme(imp string) string {
	if strings.HasPrefix(imp, "//") {
		return "https:" + imp
	}

	return "file://" + imp
}

func fetch(u string, f URLFetcher) (jsonnet.Contents, error) {
	resp, err := f.Get(u)
	if err != nil {
		return noContent, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return noContent, nil
	} else if resp.StatusCode != http.StatusOK {
		return noContent, fmt.Errorf("could not fetch %#v: %s", u, resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return noContent, err
	}

	return jsonnet.MakeContents(string(b)), nil
}
