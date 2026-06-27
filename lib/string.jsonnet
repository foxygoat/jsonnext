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

  find(str, char)::
    local _find(str, char, pos) =
      if std.length(str) == 0 then -1
      else if std.startsWith(str, char) then pos
      else _find(str[1:], char, pos + 1);
    if !std.isString(str) then
      error ('string.find first param must be a string, got ' + std.type(str))
    else if !std.isString(char) then
      error ('string.find secomd param must be a string, got ' + std.type(char))
    else if std.length(char) == 0 then -1  // empty string is not in string
    else _find(str, char, 0),

}
