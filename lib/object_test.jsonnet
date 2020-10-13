#!/usr/bin/env -S jx -J //github.com/yugui/jsonnetunit/raw/master

local test = import 'jsonnetunit/test.libsonnet';
local object = import 'object.jsonnet';

{ name: 'object.jsonnet test' } +
test.suite({
  local o = { a: 1, b: 2, c:: 3, d:: 4 },

  testKeys: {
    actual: object.keys(o),
    expect: ['a', 'b'],
  },
  testKeysEmpty: {
    actual: object.keys({}),
    expect: [],
  },

  testKeysAll: {
    actual: object.keysAll(o),
    expect: ['a', 'b', 'c', 'd'],
  },
  testKeysAllEmpty: {
    actual: object.keysAll({}),
    expect: [],
  },

  testKeysHidden: {
    actual: object.keysHidden(o),
    expect: ['c', 'd'],
  },
  testKeysHiddenEmpty: {
    actual: object.keysHidden({}),
    expect: [],
  },

  testVals: {
    actual: object.vals(o),
    expect: [1, 2],
  },
  testValsEmpty: {
    actual: object.vals({}),
    expect: [],
  },

  testValsAll: {
    actual: object.valsAll(o),
    expect: [1, 2, 3, 4],
  },
  testValsAllEmpty: {
    actual: object.valsAll({}),
    expect: [],
  },

  testValsHidden: {
    actual: object.valsHidden(o),
    expect: [3, 4],
  },
  testValsHiddenEmpty: {
    actual: object.valsHidden({}),
    expect: [],
  },

  testKV: {
    actual: object.kv(o),
    expect: [['a', 1], ['b', 2]],
  },
  testKVEmpty: {
    actual: object.kv({}),
    expect: [],
  },

  testKVAll: {
    actual: object.kvAll(o),
    expect: [['a', 1], ['b', 2], ['c', 3], ['d', 4]],
  },
  testKVAllEmpty: {
    actual: object.kvAll({}),
    expect: [],
  },

  testKVHidden: {
    actual: object.kvHidden(o),
    expect: [['c', 3], ['d', 4]],
  },
  testKVHiddenEmpty: {
    actual: object.kvHidden({}),
    expect: [],
  },

  testMake: {
    actual: object.make([['a', 1], ['b', 2]]),
    expect: o,  // only the visible fields are compared.
  },
  testMakeEmpty: {
    actual: object.make([]),
    expect: {},
  },

  testMakeHidden: {
    // Hidden object fields are not compared when testing for equality,
    // so we just make sure we get our input back out using object.kvAll.
    local input = [['c', 3], ['d', 4]],
    actual: object.kvAll(object.makeHidden(input)),
    expect: input,
  },
  testMakeHiddenEmpty: {
    // Hidden object fields are not compared when testing for equality,
    // so we just make sure we get our input back out using object.kvAll.
    actual: object.kvAll(object.makeHidden([])),
    expect: [],
  },

  testTransformKey: {
    // we need to use object.kvAll to compare the hidden fields
    actual: object.make(object.kvAll(object.transform(function(k, v) [k + 'x', v], o))),
    expect: { ax: 1, bx: 2, cx: 3, dx: 4 },
  },
  testTransformKeyVisible: {
    actual: object.transform(function(k, v) [k + 'x', v], o),
    expect: { ax: 1, bx: 2 },
  },
  testTransformKeyHidden: {
    // we need to use object.kvHidden to compare the hidden fields
    actual: object.make(object.kvHidden(object.transform(function(k, v) [k + 'x', v], o))),
    expect: { cx: 3, dx: 4 },
  },
  testTransformValue: {
    actual: object.transform(function(k, v) [k, v + 1], o),
    expect: { a: 2, b: 3 },  // only visible fields are compared
  },
  testTransformRemove: {
    actual: object.transform(function(k, v) null, o),
    expect: {},
  },

  testFilterKeep: {
    actual: object.kvAll(object.filter(function(k, v) true, o)),
    expect: [['a', 1], ['b', 2], ['c', 3], ['d', 4]],
  },
  testFilterRemove: {
    actual: object.kvAll(object.filter(function(k, v) false, o)),
    expect: [],
  },

  testRemoveField: {
    actual: object.kvAll(object.removeField(o, 'a')),
    expect: [['b', 2], ['c', 3], ['d', 4]],
  },
  testRemoveFields: {
    actual: object.kvAll(object.removeFields(o, ['a', 'c'])),
    expect: [['b', 2], ['d', 4]],
  },

  testInvert: {
    actual: object.invert({ a: 'foo', b: 'bar', c:: 'bar', d:: 'baz' }),
    expect: { foo: ['a'], bar: ['b'] },
  },
  testInvertAll: {
    actual: object.invertAll({ a: 'foo', b: 'bar', c:: 'bar', d:: 'baz' }),
    expect: { foo: ['a'], bar: ['b', 'c'], baz: ['d'] },
  },
})
