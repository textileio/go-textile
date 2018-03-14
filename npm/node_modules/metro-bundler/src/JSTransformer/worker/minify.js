/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * 
 * @format
 */

'use strict';

const uglify = require('uglify-es');







function noSourceMap(code) {
  return minify(code).code;
}

function withSourceMap(
code,
sourceMap,
filename)
{
  const result = minify(code, sourceMap);

  const map = JSON.parse(result.map);
  map.sources = [filename];
  return { code: result.code, map };
}

function minify(inputCode, inputMap) {
  const result = uglify.minify(inputCode, {
    mangle: { toplevel: true },
    output: {
      ascii_only: true,
      quote_style: 3,
      wrap_iife: true },

    sourceMap: {
      content: inputMap,
      includeSources: false },

    toplevel: true });


  if (result.error) {
    throw result.error;
  }

  return {
    code: result.code,
    map: result.map };

}

module.exports = {
  noSourceMap,
  withSourceMap };