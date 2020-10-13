// Package value contains utility functions for working with arbitrary values.
{
  local value = self,

  // get indexes the data structure @any successively with the elements of
  // @path and returns the value accessed. If indexing by any element of the
  // path results in an unknown or invalid field/index, the @default value is
  // returned.
  // @path can be either an array of indicies or a single value of an index.
  get(any, path, default=null)::
    if path == [] then any
    else if !std.isArray(path) then value.get(any, [path], default)
    else
      local idx = path[0];
      local validIdxType =
        (std.isArray(any) && std.isNumber(idx)) ||
        (std.isObject(any) && std.isString(idx));
      local knownIdx =
        (std.isArray(any) && idx >= 0 && idx < std.length(any)) ||
        (std.isObject(any) && idx in any);
      if validIdxType && knownIdx then value.get(any[idx], path[1:], default)
      else default,

  // zero returns a zero value for the type of the given argument. The zero
  // value for a function is a nullary function that returns null.
  zero(any)::
    if std.isBoolean(any) then false
    else if std.isNumber(any) then 0
    else if std.isString(any) then ''
    else if std.isArray(any) then []
    else if std.isObject(any) then {}
    else if std.isFunction(any) then function() null
    else null,

  // identity returns its argument. If no argument is provided, null is returned.
  identity(any=null):: any,
}
