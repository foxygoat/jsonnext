package jsonnext

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
)

const stdin = ""

var noContent = jsonnet.Contents{} //nolint:gochecknoglobals

// A URLFetcher retrieves a URL returning a http.Response or an error. It
// is defined such that http.Client implements it, but allows a different
// implementation or a custom-configured http.Client to be provided to
// the Importer.
type URLFetcher interface {
	Get(url string) (*http.Response, error)
}

// Importer implements the jsonnet.Importer interface, allowing jsonnet code to
// be imported via https in addition to local files. Filenames starting with a
// double-slash (`//`) are fetched via HTTPS using the Fetcher of the Importer.
// Otherwise the path is treated as a local filesystem path. An empty path, a
// path of "-" and a path of "/dev/stdin" are treated as standard input.
// "/dev/stdin" is handled specially as any imports in the code read from stdin
// should not be searched relative to "/dev".
//
// Once an import path is successfully fetched, either with data or a
// definitive not found result, that result is cached for the lifetime of the
// Importer. This is a requirement of the jsonnet.Importer interface so it is
// not possible for the same import statement from different files to result in
// different content. If an Importer is shared across multiple jsonnet.VM
// instances, the the cache will be shared too. There is no cache expiry logic.
type Importer struct {
	// SearchPath is an ordered slice of paths (network or local filesystem)
	// that is prepended to the imported filename if the filename is not
	// found. Searching stops when it is found.
	SearchPath []string

	// Fetcher is the URLFetcher used to fetch paths. The default is
	// &http.Client{}
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
	imp = mapStdin(imp)
	dir := path.Dir(source)
	if dir = preserveNetRoot(source, dir); dir == "//" {
		// There's no such thing as a "root" netpath. Preserve the
		// host part if that's all there was.
		dir = source
	}
	content, location, err := i.search(imp, dir)

	if err == nil && content == noContent {
		err = fmt.Errorf("could not read %#v: not found", imp)
	}

	return content, location, err
}

func (i *Importer) search(imp, dir string) (jsonnet.Contents, string, error) {
	if path.IsAbs(imp) || imp == stdin {
		content, err := i.readViaCache(imp)
		return content, imp, err
	}

	// try to import imp relative to source first, then the search path
	for _, p := range append([]string{dir}, i.SearchPath...) {
		location := preserveNetRoot(p, path.Join(p, imp))
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

	content, err := i.fetch(imp)
	if err != nil {
		return noContent, err
	}

	i.cache[imp] = content
	return content, nil
}

func (i *Importer) fetch(imp string) (jsonnet.Contents, error) {
	if imp == stdin {
		imp = "/dev/stdin"
	}
	r, err := i.open(imp)
	if r == nil || err != nil {
		return noContent, err
	}

	defer r.Close() //nolint:errcheck
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return noContent, err
	}

	return jsonnet.MakeContents(string(b)), nil
}

func (i *Importer) open(imp string) (io.ReadCloser, error) {
	if !isNetpath(imp) {
		r, err := os.Open(imp) //nolint:gosec // We want to open user specified paths.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return r, err
	}

	resp, err := i.fetcher().Get("https:" + imp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		_ = resp.Body.Close()
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("could not fetch %#v: %s", imp, resp.Status)
	}

	return resp.Body, nil
}

func (i *Importer) fetcher() URLFetcher {
	// TODO(camh): Consider whether this needs to be concurrency-safe
	if i.Fetcher == nil {
		i.Fetcher = &http.Client{}
	}
	return i.Fetcher
}

// preserveRoot takes an original path and an alteration of that path
// and returns the altered version with a double slash at the start if
// the original path starts with a double slash. This is needed as
// certain functions of the "path" package clean the path they return
// which removes the leading double slash. As that has special meaning
// to this importer, we need to preserve that prefix.
// TODO(camh): Create a "netpath" package with these preservation semantics.
func preserveNetRoot(orig, cleaned string) string {
	if isNetpath(orig) {
		return "/" + cleaned
	}
	return cleaned
}

func isNetpath(path string) bool {
	return len(path) > 1 && path[0] == '/' && path[1] == '/'
}

func mapStdin(path string) string {
	if path == "/dev/stdin" || path == "-" {
		path = stdin
	}
	return path
}
