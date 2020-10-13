// Package string provides utility functions for working with strings.
{
  local string = self,

  // contains returns true if @substr is in @str, false if it is not. @str and
  // @substr must be a string. If @substr is empty then true is returned.
  //
  // Note: The implementation uses `std.strReplace` as this is built-in in the
  // Go implementation and is likely the fastest way.
  contains(str, substr)::
    if !std.isString(str) then
      error ('string.contains first param must be a string, got ' + std.type(str))
    else if !std.isString(substr) then
      error ('string.contains second param must be a string, got ' + std.type(substr))
    else if std.length(substr) == 0 then
      true
    else
      std.strReplace(str, substr, '') != str,
}
