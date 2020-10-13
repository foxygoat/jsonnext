// Package array contains utility functions for working with arrays.
{
  local array = self,
  local op = import 'op.jsonnet',
  local value = import 'value.jsonnet',

  // make returns its argument if it is an array, otherwise returns a single
  // element array containing the argument.
  make(any):: if std.isArray(any) then any else [any],

  // coalesce returns the first non-null value from the given array or null
  // if there are none.
  coalesce(arr)::
    if arr == [] then null
    else if arr[0] != null then arr[0]
    else array.coalesce(arr[1:]),

  // collapse returns the input array with any null values removed
  collapse(arr):: [e for e in arr if e != null],

  // accumulate adds all the elements in the array @arr together starting
  // with the value @init. If @init is not supplied or null, the zero value
  // of the first element of the array is used. If the array is empty, @init
  // is returned.
  //
  // The values in the array can be any types that can be added together with
  // the plus (+) operator. Numbers are summed, strings and arrays are
  // concatenated and objects are merged. Booleans and functions do not
  // support the plus operator. Strings can be added to any type resulting
  // in the other type being converted to a string before concatenating.
  accumulate(arr, init=null)::
    local zero = if std.length(arr) > 0 then value.zero(arr[0]) else null;
    std.foldl(op.add, arr, array.coalesce([init, zero])),
}
