// Package op provides function values for the built-in functions.
// These can be useful with functions such as std.foldl where you
// want to fold with a built-in function.
{
  local op = self,

  add(a, b):: a + b,
  sub(a, b):: a - b,
  mul(a, b):: a * b,
  div(a, b):: a / b,
  mod(a, b):: a % b,
  neg(a):: -a,

  eq(a, b):: a == b,
  ne(a, b):: a != b,
  lt(a, b):: a < b,
  le(a, b):: a <= b,
  gt(a, b):: a > b,
  ge(a, b):: a >= b,

  lsh(a, b):: a << b,
  rsh(a, b):: a >> b,
  band(a, b):: a & b,
  bor(a, b):: a | b,
  bxor(a, b):: a ^ b,
  bcpl(a):: ~a,

  and(a, b):: a && b,
  or(a, b):: a || b,
  not(a):: !a,

  inobj(a, b):: a in b,
}
