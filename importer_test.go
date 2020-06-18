package jsonnext

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

type errReader struct {
	err error
}

func (r errReader) Read(p []byte) (int, error) {
	return 0, r.err
}

type fetcherMock struct {
	urls       []string
	statusCode int
	content    string
	getErr     error
	readErr    error
}

func (f *fetcherMock) Get(url string) (*http.Response, error) {
	f.urls = append(f.urls, url)

	var body io.Reader
	if f.readErr != nil {
		body = errReader{err: f.readErr}
	} else {
		body = strings.NewReader(f.content)
	}

	return &http.Response{StatusCode: f.statusCode, Body: ioutil.NopCloser(body)}, f.getErr
}

func TestImport(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	i := Importer{Fetcher: f}

	contents, location, err := i.Import("/a/y", "x")

	require.NoError(t, err)
	require.Equal(t, "hello world", contents.String())
	require.Equal(t, "/a/x", location)
}

func TestImportNotFound(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound, content: ""}
	i := Importer{Fetcher: f}

	_, _, err := i.Import("/a/y", "x")

	require.Error(t, err)
}

func TestImportNetpathSource(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound, content: ""}
	i := Importer{Fetcher: f}

	_, _, err := i.Import("//example.com/y", "x")

	require.Error(t, err)
	require.Equal(t, []string{"https://example.com/x"}, f.urls)
}

func TestImportRootNetpathSource(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound, content: ""}
	i := Importer{Fetcher: f}

	// Should try to import `//example.com/x`, not `//x`
	_, _, err := i.Import("//example.com", "x")

	require.Error(t, err)
	require.Equal(t, []string{"https://example.com/x"}, f.urls)
}

func TestSearchAbsolute(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	i := Importer{Fetcher: f}

	_, location, err := i.search("/x/y/z", "/a")

	require.NoError(t, err)
	require.Equal(t, "/x/y/z", location)
	require.Equal(t, []string{"file:///x/y/z"}, f.urls)
}

func TestSearchRelativeToDir(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound, content: ""}
	i := Importer{Fetcher: f}

	_, location, err := i.search("x", "/a")

	require.NoError(t, err)
	require.Equal(t, "", location)
	require.Equal(t, []string{"file:///a/x"}, f.urls)
}

func TestSearchPathUntilFound(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	i := Importer{Fetcher: f, SearchPath: []string{"/1", "/2"}}

	contents, location, err := i.search("x", "/a")

	require.NoError(t, err)
	require.Equal(t, "/a/x", location)
	require.Equal(t, []string{"file:///a/x"}, f.urls)
	require.Equal(t, "hello world", contents.String())
}

func TestSearchNetpath(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	i := Importer{Fetcher: f, SearchPath: []string{"/1", "/2"}}

	contents, location, err := i.search("x", "//example.com")

	require.NoError(t, err)
	require.Equal(t, "//example.com/x", location)
	require.Equal(t, []string{"https://example.com/x"}, f.urls)
	require.Equal(t, "hello world", contents.String())
}

func TestSearchPathToEnd(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound, content: ""}
	i := Importer{Fetcher: f, SearchPath: []string{"/1", "/2"}}

	contents, location, err := i.search("x", "/a")

	require.NoError(t, err)
	require.Equal(t, "", location)
	require.Equal(t, []string{"file:///a/x", "file:///1/x", "file:///2/x"}, f.urls)
	require.Equal(t, noContent, contents)
}

func TestReadViaCacheEmpty(t *testing.T) {
	i := Importer{Fetcher: &fetcherMock{statusCode: http.StatusOK, content: "hello world"}}

	require.NotContains(t, i.cache, "file:///test/path")
	contents, err := i.readViaCache("file:///test/path")

	require.NoError(t, err)
	require.Equal(t, "hello world", contents.String())
	require.Contains(t, i.cache, "file:///test/path")
	require.Equal(t, "hello world", i.cache["file:///test/path"].String())
}

func TestReadViaCacheNotEmpty(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	i := Importer{Fetcher: f}

	// populate the cache
	_, err := i.readViaCache("file:///test/path")
	require.NoError(t, err)

	// cached contents will be "hello world". fetcher is now returning "goodbye world"
	f.content = "goodbye world"
	contents, err := i.readViaCache("file:///test/path")

	require.NoError(t, err)
	require.Equal(t, "hello world", contents.String())
}

func TestReadViaCacheDoesntCacheErrors(t *testing.T) {
	i := Importer{Fetcher: &fetcherMock{statusCode: http.StatusBadRequest, content: ""}}

	require.NotContains(t, i.cache, "file:///test/path")
	_, err := i.readViaCache("file:///test/path")

	require.Error(t, err)
	require.NotContains(t, i.cache, "file:///test/path")
}

func TestFetcherNew(t *testing.T) {
	i := Importer{}

	actual := i.fetcher()

	require.IsType(t, &http.Client{}, actual)
	c := actual.(*http.Client)
	require.IsType(t, &http.Transport{}, c.Transport)
	ts := c.Transport.(*http.Transport)
	// Test to ensure a "file" scheme handler is registered. By trying to
	// register one, it will panic if one is already registered.
	require.Panics(t, func() { ts.RegisterProtocol("file", http.NewFileTransport(http.Dir("/"))) })
}

func TestFetcherExisting(t *testing.T) {
	expected := &fetcherMock{}
	i := Importer{Fetcher: expected}

	actual := i.fetcher()

	require.Equal(t, expected, actual)
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

func TestAddScheme(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	tests := map[string]struct{ input, expected string }{
		"absoluteFile": {"/tmp/foo", "file:///tmp/foo"},
		"relativeFile": {"foo/bar", fmt.Sprintf("file://%s/foo/bar", cwd)},
		"networkFile":  {"//example.com/path/to/file", "https://example.com/path/to/file"},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) { //nolint:wsl
			actual, err := addScheme(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestFetchOK(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, content: "hello world"}
	u := "file:///foo"
	content, err := fetch(u, f)
	require.NoError(t, err)
	require.Equal(t, u, f.urls[0])
	require.Equal(t, f.content, content.String())
}

func TestFetchNotFound(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusNotFound}
	content, err := fetch("file:///foo", f)
	require.NoError(t, err)
	require.Equal(t, noContent, content)
}

func TestFetchFetcherErr(t *testing.T) {
	f := &fetcherMock{getErr: errors.New("get error")}
	_, err := fetch("file:///foo", f)
	require.EqualError(t, err, "get error")
}

func TestFetchStatusErr(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusInternalServerError}
	_, err := fetch("file:///foo", f)
	require.Error(t, err)
}

func TestFetchReadErr(t *testing.T) {
	f := &fetcherMock{statusCode: http.StatusOK, readErr: errors.New("read error")}
	_, err := fetch("file:///foo", f)
	require.EqualError(t, err, "read error")
}
