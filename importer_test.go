package jsonnext

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendSearchFromEnv(t *testing.T) {
	i := Importer{}
	err := os.Setenv("IMPORTER_TEST", "//example.com/path:/local/path")
	require.NoError(t, err)

	i.AppendSearchFromEnv("IMPORTER_TEST")

	require.Equal(t, []string{"//example.com/path", "/local/path"}, i.SearchPath)
}

func TestAppendSearchFromEnvNonExistent(t *testing.T) {
	i := Importer{}
	err := os.Unsetenv("IMPORTER_TEST")
	require.NoError(t, err)

	i.AppendSearchFromEnv("IMPORTER_TEST")

	require.Empty(t, i.SearchPath)
}

func TestAppendSearchFromEnvEmpty(t *testing.T) {
	i := Importer{}
	err := os.Setenv("IMPORTER_TEST", "://example.com/path:::/local/path:")
	require.NoError(t, err)

	i.AppendSearchFromEnv("IMPORTER_TEST")

	require.Equal(t, []string{"//example.com/path", "/local/path"}, i.SearchPath)
}

func TestImportLocal(t *testing.T) {
	i := Importer{}
	contents, foundAt, err := i.Import("", "testdata/importer/hello.txt")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", contents.String())
	require.Equal(t, "testdata/importer/hello.txt", foundAt)
}

func TestImportLocalNotFound(t *testing.T) {
	i := Importer{}
	_, _, err := i.Import("", "testdata/importer/notfound.txt")
	require.Error(t, err)
}

// Make reading file by trying to read a directory. That causes read(2) to
// return EDIR.
func TestImportLocalReadError(t *testing.T) {
	i := Importer{}
	_, _, err := i.Import("", "testdata/importer")
	require.Error(t, err)
}

func TestImportLocalSourceRelative(t *testing.T) {
	i := Importer{}
	contents, foundAt, err := i.Import("testdata/importer/hello.txt", "mellow.txt")
	require.NoError(t, err)
	require.Equal(t, "mellow world\n", contents.String())
	require.Equal(t, "testdata/importer/mellow.txt", foundAt)
}

func TestImportLocalRelativeSearch(t *testing.T) {
	i := Importer{}
	i.SearchPath = []string{"testdata"}
	contents, foundAt, err := i.Import("", "importer/hello.txt")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", contents.String())
	require.Equal(t, "testdata/importer/hello.txt", foundAt)
}

func TestImportNetpath(t *testing.T) {
	s := httptest.NewTLSServer(http.FileServer(http.Dir("testdata")))
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	contents, foundAt, err := i.Import("", np+"/importer/hello.txt")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", contents.String())
	require.Equal(t, np+"/importer/hello.txt", foundAt)
}

func TestImportNetpathNotFound(t *testing.T) {
	s := httptest.NewTLSServer(http.FileServer(http.Dir("testdata")))
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	_, _, err := i.Import("", np+"/importer/notfound.txt")
	require.Error(t, err)
}

func TestImportNetpathSourceRelative(t *testing.T) {
	s := httptest.NewTLSServer(http.FileServer(http.Dir("testdata")))
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	contents, foundAt, err := i.Import(np, "importer/mellow.txt")
	require.NoError(t, err)
	require.Equal(t, "mellow world\n", contents.String())
	require.Equal(t, np+"/importer/mellow.txt", foundAt)
}

// Make the URLFetcher.Get() method return an error by not using the proper
// TLS client for it to work.
func TestImportNetpathFetchError(t *testing.T) {
	s := httptest.NewTLSServer(http.FileServer(http.Dir("testdata")))
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{} // dont use s.Client(), so we get a TLS error
	_, _, err := i.Import("", np+"/importer/hello.txt")
	require.Error(t, err)
}

// Make the httptest server return a non-404 error. By making open(2) fail
// with ELOOP (with a symlink that links to itself), we get a 500 error.
func TestImportNetpathHTTPError(t *testing.T) {
	s := httptest.NewTLSServer(http.FileServer(http.Dir("testdata")))
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	_, _, err := i.Import("", np+"/importer/ELOOP")
	require.Error(t, err)
}

type requestRecorder struct {
	requests []*http.Request
	next     http.Handler
}

func (rr *requestRecorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.requests = append(rr.requests, r)
	rr.next.ServeHTTP(w, r)
}

// Test importing from the cache by importing the same thing twice. We record
// the requests to the http.Handler and count how many times it was called.
func TestImportFromCache(t *testing.T) {
	rr := &requestRecorder{next: http.FileServer(http.Dir("testdata"))}
	s := httptest.NewTLSServer(rr)
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	require.Equal(t, 0, len(rr.requests))

	contents, foundAt, err := i.Import("", np+"/importer/hello.txt")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", contents.String())
	require.Equal(t, np+"/importer/hello.txt", foundAt)
	require.Equal(t, 1, len(rr.requests))

	contents, foundAt, err = i.Import("", np+"/importer/hello.txt")
	require.NoError(t, err)
	require.Equal(t, "hello world\n", contents.String())
	require.Equal(t, np+"/importer/hello.txt", foundAt)
	require.Equal(t, 1, len(rr.requests)) // still 1
}

// Test that an import of an absolute path does not go through the search
// path, by importing non-existent file.
func TestImportAbsolute(t *testing.T) {
	rr := &requestRecorder{next: http.FileServer(http.Dir("testdata"))}
	s := httptest.NewTLSServer(rr)
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	i.SearchPath = []string{np + "/importer", np + "/importer/https"}
	require.Equal(t, 0, len(rr.requests))
	_, _, err := i.Import("", np+"/importer/notfound.txt")
	require.Error(t, err)
	require.Equal(t, 1, len(rr.requests))
}

func TestImportSearchAllPaths(t *testing.T) {
	rr := &requestRecorder{next: http.FileServer(http.Dir("testdata"))}
	s := httptest.NewTLSServer(rr)
	defer s.Close()
	np := strings.TrimPrefix(s.URL, "https:")

	i := Importer{Fetcher: s.Client()}
	i.SearchPath = []string{np + "/importer", np + "/importer/https"}
	require.Equal(t, 0, len(rr.requests))
	_, _, err := i.Import("", "notfound.txt")
	require.Error(t, err)
	require.Equal(t, 2, len(rr.requests))
}

func TestPreserveNetRoot(t *testing.T) {
	tests := map[string]struct{ inputOrig, inputClean, expected string }{
		"single":     {"/a/b/c", "/a/b", "/a/b"},
		"double":     {"//a/b/c", "/a/b", "//a/b"},
		"emptyOrig":  {"", "/a/b", "/a/b"},
		"emptyClean": {"/a/b/c", "", ""},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) { //nolint:wsl
			actual := preserveNetRoot(tt.inputOrig, tt.inputClean)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestFetcherNew(t *testing.T) {
	i := Importer{}
	require.IsType(t, &http.Client{}, i.fetcher())
}

func TestFetcherExisting(t *testing.T) {
	expected := &http.Client{}
	i := Importer{Fetcher: expected}
	require.Equal(t, expected, i.fetcher())
}
