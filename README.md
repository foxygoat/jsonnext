# jsonnext

Extension libraries for [go-jsonnet](https://github.com/google/go-jsonnet).
The name `jsonnext` is a contraction of `jsonnet` and `ext`.

[Documentation for this module](https://pkg.go.dev/foxygo.at/jsonnext)
is on pkg.go.dev.

## Importer

[`foxygo.at/jsonnext.Importer`](https://pkg.go.dev/foxygo.at/jsonnext#Importer)
is a
[`jsonnet.Importer`](https://pkg.go.dev/github.com/google/go-jsonnet#Importer)
implementation that can import from the local filesystem and from
`https` network sources. Normal paths, absolute and relative, are
imported from the local filesystem. Paths starting with a double-slash
(`//`) are imported as `https` URLs with an implicit `https:` prefix.
For example, the path
`//github.com/grafana/grafonnet-lib/raw/master/grafonnet/grafana.libsonnet`
refers to the file `grafana.libsonnet` on the `master` branch of the
`grafana/grafonnet-lib` repository on GitHub (using the "raw" path to
retrieve the contents of the file instead of the GitHub web page for
the file).

The `Importer` has a `SearchPath` of type `[]string` that can be
populated directly or by using the `AppendSearchFromEnv()` method. The
method takes the name of an environment variable, splits it on the
OS-specific ListSeparator and appends the elements to the existing
search path.

The default `Fetcher` for the importer is the default `http.Client`. It
can be replaced with any type that implements the `Get` method of
`http.Client`. Most likely it will be overridden with a `http.Client`
that has been constructed with a non-default configuration.

The importer maintains a cache of results as is required by the
`jsonnet.Importer` interface description. Positive and negative results
are cached and returned on subsequent calls to import the same path.
Errors retrieving a path are not cached and are returned as an error
results from the `Import` method.
