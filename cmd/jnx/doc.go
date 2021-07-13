// jnx evaluates a jsonnet file and outputs it as JSON.
//
// Usage: jnx [<filename>]
//
// Arguments:
//   [<filename>]    File to evaluate. stdin is used if omitted or "-"
//
// Flags:
//   -h, --help                            Show context-sensitive help.
//   -J, --jpath=dir                       Add a library search dir
//       --max-stack=500                   Number of allowed stack frames of jsonnet VM
//       --max-trace=20                    Maximum number of stack frames output on error
//   -V, --ext-str=var[=str]               Set extVar string (str from env if omitted)
//       --ext-str-file=var[=filename]     Set extVar string from a file (filename from env if omitted)
//       --ext-code=var[=code]             Set extVar code (code from env if omitted)
//       --ext-code-file=var[=filename]    Set extVar code from a file (filename from env if omitted)
//   -A, --tla-str=var[=str]               Set top-level arg string (str from env if omitted)
//       --tla-str-file=var[=filename]     Set top-level arg string from a file (filename from env if omitted)
//       --tla-code=var[=code]             Set top-level arg code (code from env if omitted)
//       --tla-code-file=var[=filename]    Set top-level arg code from a file (filename from env if omitted)
package main
