#!/usr/bin/env -S jnx -J //github.com/yugui/jsonnetunit/raw/master

local array = import 'array.jsonnet';
local test = import 'jsonnetunit/test.libsonnet';

{ name: 'array.jsonnet test' } +
test.suite({
  // tests for array.make(any)
  testMakeNonArray: {
    actual: array.make(42),
    expect: [42],
  },
  testMakeArray: {
    actual: array.make([42]),
    expect: [42],
  },

  // tests for array.coalesce(arr)
  testCoalesceFirst: {
    actual: array.coalesce([1, null, 2]),
    expect: 1,
  },
  testCoalesceSecond: {
    actual: array.coalesce([null, 1, 2]),
    expect: 1,
  },
  testCoalesceEmpty: {
    actual: array.coalesce([]),
    expect: null,
  },
  testCoalesceNull: {
    actual: array.coalesce([null]),
    expect: null,
  },

  // tests for array.colapse(arr)
  testCollapseNulls: {
    actual: array.collapse([null, null]),
    expect: [],
  },
  testCollapseEmpty: {
    actual: array.collapse([]),
    expect: [],
  },
  testCollapseFull: {
    actual: array.collapse([1, 2, 3, 4]),
    expect: [1, 2, 3, 4],
  },
  testCollapsePartial: {
    actual: array.collapse([1, null, 2, null, null, 3, 4]),
    expect: [1, 2, 3, 4],
  },

  // tests for array.accumulate(arr, init)
  testAccumulate: {
    actual: array.accumulate([1, 1, 2, 3, 5, 8, 13]),
    expect: 33,
  },
  testAccumulateInit: {
    actual: array.accumulate([1, 1, 2, 3, 5, 8, 13], 21),
    expect: 54,
  },
  testAccumulateEmpty: {
    actual: array.accumulate([]),
    expect: null,
  },
  testAccumulateEmptyInit: {
    actual: array.accumulate([], 42),
    expect: 42,
  },
})
