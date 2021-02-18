#!/usr/bin/env -S jx -J //github.com/yugui/jsonnetunit/raw/master
//
// Roll-up all the tests in this directory
//

local test = import 'jsonnetunit/test.libsonnet';
local jx = import 'jx.jsonnet';

local jx_test = { name: 'jx.jsonnet test' } + test.suite({
  testImport: {
    // Force evaluation of imports in `jx.jsonnet` by using std.prune over
    // all fields. This will cause a failure if a file cannot be imported.
    actual: std.prune(std.objectValuesAll(jx)),
    expect: [],  // empty because all fields are hidden and thus pruned.
  },
});

[
  jx_test,
  import 'array_test.jsonnet',
  import 'object_test.jsonnet',
  import 'op_test.jsonnet',
  import 'string_test.jsonnet',
  import 'value_test.jsonnet',
]
