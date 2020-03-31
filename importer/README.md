# importer

[`foxygo.at/jsonnext/importer`](https://pkg.go.dev/foxygo.at/jsonnext/importer)
is a
[`jsonnet.Importer`](https://pkg.go.dev/github.com/google/go-jsonnet#Importer)
implementation that can import from the local filesystem and from
`https` network sources. Normal paths, absolute and relative, are
imported from the local filesystem. Paths starting with a double-slash
(`//`) are imported as `https` URLs with an implicit `https:` prefix.
For example, the path
`//github.com/grafana/grafonnet-lib/blob/master/grafonnet/grafana.libsonnet`
refers to the file `grafana.libsonnet` on the `master` branch of the
`grafana/grafonnet-lib` repository on GitHub.

The `Importer` has a `SearchPath` that can be populated directly or with
the `AppendSearchFromEnv` method which takes a name of an environment
variable, splits it on the OS-specific ListSeparator and appends the
elements to the existing search path.

The default `Fetcher` for the importer is the default `http.Client`. It
can be replaced with any type that implements the `Get` method of
`http.Client`. Most likely it will be overridden with a `http.Client`
that has been constructed with a non-default configuration.

The importer maintains a cache of results as is required by the
`jsonnet.Importer` interface description. Positive and negative results
are cached and returned on subsequent calls to import the same path.
Errors retrieving a path are not cached and are returned as an error
results from the `Import` method.

