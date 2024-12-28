// Package object contains utility functions for working with objects.
{
  local object = self,
  local array = import 'array.jsonnet',

  // make returns an object from an array @arr of `[key, value]` pairs. The
  // fields in the returned object are all visible (non-hidden).
  make(arr):: { [kv[0]]: kv[1] for kv in arr },

  // makeHidden returns an object from an array @arr of `[key, value]` pairs.
  // The fields in the returned object are all hidden.
  makeHidden(arr):: array.accumulate([{ [kv[0]]:: kv[1] } for kv in arr], {}),

  // keys returns all the non-hidden field keys of @obj in an array. It is
  // identical to std.objectFields and exists here to fill out the 3x3 matrix
  // of functions for (keys, vals, kv) x (visble, all, hidden).
  keys(obj):: std.objectFields(obj),

  // keysAll returns all the field keys (hidden and non-hidden) of @obj in an
  // array. It is identical to std.objectFieldsAll and exists here to fill out
  // the 3x3 matrix of functions for (keys, vals, kv) x (visble, all, hidden).
  keysAll(obj):: std.objectFieldsAll(obj),

  // keysHidden returns all the hidden field keys of @obj in an array.
  keysHidden(obj):: std.setDiff(object.keysAll(obj), object.keys(obj)),

  // vals returns all the non-hidden field values of @obj in an array.
  vals(obj):: [obj[k] for k in object.keys(obj)],

  // valsAll returns all the field values (hidden and non-hidden) of @obj in an
  // array.
  valsAll(obj):: [obj[k] for k in object.keysAll(obj)],

  // valsHidden returns all the hidden field values of @obj in an array.
  valsHidden(obj):: [obj[k] for k in object.keysHidden(obj)],

  // kv returns all the non-hidden fields of @obj as an array of `[key, value]`
  // pairs.
  kv(obj):: [[k, obj[k]] for k in object.keys(obj)],

  // kvAll returns all the fields (hidden and non-hidden) of @obj as an array
  // of `[key, value]` pairs.
  kvAll(obj):: [[k, obj[k]] for k in object.keysAll(obj)],

  // kvHidden returns all the hidden fields of @obj as an array of
  // `[key, value]` pairs.
  kvHidden(obj):: [[k, obj[k]] for k in object.keysHidden(obj)],

  // map applies a function taking a key and value to each non-hidden field
  // of a the given object to return an array of values returned by each
  // invocation of the given function.
  map(fn, obj):: [fn(kv[0], kv[1]) for kv in object.kv(obj)],

  // transform calls @fn(key, value) on each field of @obj passing it the key
  // and the value of each field and using the result of @fn to build the
  // result of transform. It can remove fields, rename fields or change the
  // value of fields based on what the given function @fn returns.
  //
  // The signature of @fn is fn(key, value) => [newkey, newval] | null.
  // If @fn returns null, the field is removed from the result of transform.
  // Otherwise @fn should return a `[newkey, newvalue]` pair which is put into
  // the result of transform as a field. The visibility of the field in the
  // result object is the same as the visibility of the field passed to @fn.
  transform(fn, obj)::
    local kvToArgs(fn) = function(kv) fn(kv[0], kv[1]);
    local visible = array.collapse(std.map(kvToArgs(fn), object.kv(obj)));
    local hidden = array.collapse(std.map(kvToArgs(fn), object.kvHidden(obj)));
    object.make(visible) + object.makeHidden(hidden),

  // filter calls @fn on each field of @obj passing it the key and the value
  // of each field and returns an object with the fields for which @fn
  // returned true.
  filter(fn, obj)::
    local map_fn(k, v) = if fn(k, v) then [k, v] else null;
    object.transform(map_fn, obj),

  // removeField returns @obj with @field removed from it.
  removeField(obj, field):: object.removeFields(obj, [field]),

  // removeFields returns @obj with each field listed in @fields removed
  // from it.
  removeFields(obj, fields)::
    local remove = std.set(fields);
    local keep_fn(k, v) = !std.setMember(k, remove);
    object.filter(keep_fn, obj),

  // invert swaps keys and values in @obj. The values in the output object is a
  // list of keys of visible (non-hidden) fields from @obj. If any of the
  // values in @obj cannot be keys in an object (i.e. non-string values),
  // jsonnet will produce an evaluation error.
  // `{ a: 'foo', b: 'bar', c:: 'foo' }` -> `{ bar: ['b'], foo: ['a'] }`
  invert(obj):: __invert(obj, object.keys),

  // invertAll swaps keys and values in @obj. The values in the output object
  // is a list of all keys (visible and hidden) of fields from @obj. If any of
  // the values in @obj cannot be keys in an object (i.e. non-string values),
  // jsonnet will produce an evaluation error.
  // `{ a: 'foo', b: 'bar', c:: 'bar' }` -> `{ bar: ['b', 'c'], foo: ['a'] }`
  invertAll(obj):: __invert(obj, object.keysAll),

  local __invert(obj, selector) = array.accumulate([{ [obj[k]]+: [k] } for k in selector(obj)]),

  // asNamedArray transforms the object @obj into an array of the values in the
  // object, adding a name field of the value's key if it does not have a name
  // field. The default @nameField is 'name', but can be overridden.
  asNamedArray(obj, nameField='name', valueField=null)::
    local mkobj(v) = if valueField == null then v else { [valueField]: v };
    [{ [nameField]: kv[0] } + mkobj(kv[1]) for kv in object.kv(obj)],
}
