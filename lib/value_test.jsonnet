#!/usr/bin/env -S jnx -J //github.com/yugui/jsonnetunit/raw/master

local test = import 'jsonnetunit/test.libsonnet';
local value = import 'value.jsonnet';

{ name: 'value.jsonnet test' } +
test.suite({
  local getFixture = { a: 1, b: null, c: [1, 2, 3], d: { x: 'x', y: { z: 'z' } } },
  testGetNoPath: {
    actual: value.get(getFixture, []),
    expect: getFixture,
  },
  testGetStringPath: {
    actual: value.get(getFixture, 'a'),
    expect: 1,
  },
  testGetPathOneLevel: {
    actual: value.get(getFixture, ['a']),
    expect: 1,
  },
  testGetArray: {
    actual: value.get(getFixture.c, 2),
    expect: 3,
  },
  testGetObjDefault: {
    actual: value.get(getFixture, 'non-existent', 'default'),
    expect: 'default',
  },
  testGetArrayDefault: {
    actual: value.get(getFixture.c, 3, 42),
    expect: 42,
  },
  testGetDefaultNull: {
    actual: value.get(getFixture, 'non-existent'),
    expect: null,
  },
  testGetPath: {
    actual: value.get(getFixture, ['d', 'y', 'z']),
    expect: 'z',
  },
  testGetNull: {
    actual: value.get(getFixture, 'b', 'not-null'),
    expect: null,
  },

  testZeroBool: {
    actual: value.zero(true),
    expect: false,
  },
  testZeroNumber: {
    actual: value.zero(42),
    expect: 0,
  },
  testZeroString: {
    actual: value.zero('hello world'),
    expect: '',
  },
  testZeroArray: {
    actual: value.zero([1, 2, 3, 4]),
    expect: [],
  },
  testZeroObject: {
    actual: value.zero({ a: 1, b:: 2, c::: 3 }),
    expect: {},
  },
  testZeroFunction: {
    actual: value.zero(value.zero),
    expectThat: function(actual) actual() == null,
  },

  testID: {
    actual: value.identity(42),
    expect: 42,
  },
  testIDNull: {
    actual: value.identity(),
    expect: null,
  },
})
