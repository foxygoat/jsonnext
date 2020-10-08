#!/usr/bin/env -S jx -J //github.com/yugui/jsonnetunit/raw/master
//
// Roll-up all the tests in this directory
//

[
  import 'string_test.jsonnet',
  import 'value_test.jsonnet',
]
