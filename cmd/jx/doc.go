// jx evaluates a jsonnet file and outputs it as JSON.
//
//   Usage: jx [<filename>]
//
//   Arguments:
//     [<filename>]    File to evaluate. stdin is used if omitted or "-"
//
//   Flags:
//     -h, --help                      Show context-sensitive help.
//     -J, --jpath=dir                 Add a library search dir
//     -V, --ext-str=var[=str]         Add extVar string (from environment if <str> is omitted)
//         --ext-code=var[=code]       Add extVar code (from environment if <code> is omitted)
//         --ext-str-file=var=file     Add extVar string from a file
//         --ext-code-file=var=file    Add extVar code from a file
//     -A, --tla-str=var[=str]         Add top-level arg string (from environment if <str> is omitted)
//         --tla-code=var[=code]       Add top-level arg code (from environment if <code> is omitted)
//         --tla-str-file=var=file     Add top-level arg string from a file
//         --tla-code-file=var=file    Add top-level arg code from a file
package main
