#!/usr/bin/env -S jx -J //github.com/yugui/jsonnetunit/raw/master
//
// Roll-up all the tests in this directory
//

[
  import 'array_test.jsonnet',
  import 'object_test.jsonnet',
  import 'op_test.jsonnet',
  import 'string_test.jsonnet',
  import 'value_test.jsonnet',
]
