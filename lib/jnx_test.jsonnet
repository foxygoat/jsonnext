#!/usr/bin/env -S jnx -J //github.com/yugui/jsonnetunit/raw/master
//
// Roll-up all the tests in this directory
//

local jnx = import 'jnx.jsonnet';
local test = import 'jsonnetunit/test.libsonnet';

local jnx_test = { name: 'jnx.jsonnet test' } + test.suite({
  testImport: {
    // Force evaluation of imports in `jnx.jsonnet` by using std.prune over
    // all fields. This will cause a failure if a file cannot be imported.
    actual: std.prune(std.objectValuesAll(jnx)),
    expect: [],  // empty because all fields are hidden and thus pruned.
  },
});

[
  jnx_test,
  import 'array_test.jsonnet',
  import 'object_test.jsonnet',
  import 'op_test.jsonnet',
  import 'string_test.jsonnet',
  import 'value_test.jsonnet',
]
