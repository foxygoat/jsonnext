#!/usr/bin/env -S jnx -J //github.com/yugui/jsonnetunit/raw/master

local test = import 'jsonnetunit/test.libsonnet';
local op = import 'op.jsonnet';

{ name: 'op.jsonnet test' } +
test.suite({
  testAdd: { actual: op.add(20, 22), expect: 42 },
  testSub: { actual: op.sub(66, 24), expect: 42 },
  testMul: { actual: op.mul(6, 7), expect: 42 },
  testDiv: { actual: op.div(252, 6), expect: 42 },
  testMod: { actual: op.mod(109, 67), expect: 42 },
  testNeg: { actual: op.neg(-42), expect: 42 },

  testEq1: { actual: op.eq(42, 42), expect: true },
  testEq2: { actual: op.eq(42, 41), expect: false },
  testNe1: { actual: op.ne(42, 42), expect: false },
  testNe2: { actual: op.ne(42, 41), expect: true },
  testLt1: { actual: op.lt(41, 42), expect: true },
  testLt2: { actual: op.lt(42, 42), expect: false },
  testLe1: { actual: op.le(41, 42), expect: true },
  testLe2: { actual: op.le(42, 42), expect: true },
  testLe3: { actual: op.le(43, 42), expect: false },
  testGt1: { actual: op.gt(42, 41), expect: true },
  testGt2: { actual: op.gt(42, 42), expect: false },
  testGe1: { actual: op.ge(42, 41), expect: true },
  testGe2: { actual: op.ge(42, 42), expect: true },
  testGe3: { actual: op.ge(42, 43), expect: false },

  testLsh: { actual: op.lsh(21, 1), expect: 42 },
  testRsh: { actual: op.rsh(168, 2), expect: 42 },
  testBand: { actual: op.band(171, 110), expect: 42 },
  testBor: { actual: op.bor(40, 10), expect: 42 },
  testBxor: { actual: op.bxor(78, 100), expect: 42 },
  testBcpl: { actual: op.bcpl(-43), expect: 42 },

  testAnd1: { actual: op.and(true, true), expect: true },
  testAnd2: { actual: op.and(true, false), expect: false },
  testAnd3: { actual: op.and(false, true), expect: false },
  testAnd4: { actual: op.and(false, false), expect: false },
  testOr1: { actual: op.and(true, true), expect: true },
  testOr2: { actual: op.or(true, false), expect: true },
  testOr3: { actual: op.or(false, true), expect: true },
  testOr4: { actual: op.or(false, false), expect: false },
  testNot1: { actual: op.not(true), expect: false },
  testNot2: { actual: op.not(false), expect: true },

  testInObj1: { actual: op.inobj('a', { a: 42 }), expect: true },
  testInObj2: { actual: op.inobj('b', { a: 42 }), expect: false },
})
