// Package jsonnext provides some types and functions to make working with
// jsonnet (github.com/google/go-jsonnet) easier and more featureful.
//
// Importer
//
// An implementation of jsonnet.Importer that supports netpaths - paths that
// start with a double-slash - // - for network paths. Netpaths are retrieved
// via https using a URLFetcher. Other paths are read from the filesystem.
//
// This form of netpath makes it simpler to use in a PATH-style,
// colon-separated environment variable where the colon in a URL would need to
// be escaped.

package jsonnext
