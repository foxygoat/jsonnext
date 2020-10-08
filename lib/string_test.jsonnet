#!/usr/bin/env -S jx -J //github.com/yugui/jsonnetunit/raw/master

local test = import 'jsonnetunit/test.libsonnet';
local string = import 'string.jsonnet';

{ name: 'string.jsonnet test' } +
test.suite({
  testStrContains: {
    actual: string.contains('hello', 'el'),
    expect: true,
  },

  testStrDoesntContain: {
    actual: string.contains('hello', 'xx'),
    expect: false,
  },

  testEmptyStrDoesntContain: {
    actual: string.contains('', 'x'),
    expect: false,
  },

  testStrContainsEmpty: {
    actual: string.contains('hello', ''),
    expect: true,
  },

  testEmptyContainsEmpty: {
    actual: string.contains('', ''),
    expect: true,
  },
})
