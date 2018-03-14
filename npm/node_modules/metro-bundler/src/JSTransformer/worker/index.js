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

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};

const asyncify = require('async/asyncify');
const constantFolding = require('./constant-folding');
const extractDependencies = require('./extract-dependencies');
const inline = require('./inline');
const invariant = require('fbjs/lib/invariant');
const minify = require('./minify');



































































const transformCode = asyncify(
(
transformer,
filename,
localPath,
sourceCode,
options) =>
{
  invariant(
  !options.minify || options.transform.generateSourceMaps,
  'Minifying source code requires the `generateSourceMaps` option to be `true`');


  const isJson = filename.endsWith('.json');
  if (isJson) {
    sourceCode = 'module.exports=' + sourceCode;
  }

  const transformFileStartLogEntry = {
    action_name: 'Transforming file',
    action_phase: 'start',
    file_name: filename,
    log_entry_label: 'Transforming file',
    start_timestamp: process.hrtime() };


  const transformed = transformer.transform({
    filename,
    localPath,
    options: options.transform,
    src: sourceCode });


  invariant(
  transformed != null,
  'Missing transform results despite having no error.');


  var code, map;
  if (options.minify) {var _constantFolding =
    constantFolding(
    filename,
    inline(filename, transformed, options));code = _constantFolding.code;map = _constantFolding.map;

    invariant(code != null, 'Missing code from constant-folding transform.');
  } else {
    code = transformed.code;map = transformed.map;
  }

  if (isJson) {
    code = code.replace(/^\w+\.exports=/, '');
  } else {
    // Remove shebang
    code = code.replace(/^#!.*/, '');
  }

  const depsResult = isJson ?
  { dependencies: [], dependencyOffsets: [] } :
  extractDependencies(code, filename);

  const timeDelta = process.hrtime(
  transformFileStartLogEntry.start_timestamp);

  const duration_ms = Math.round((timeDelta[0] * 1e9 + timeDelta[1]) / 1e6);
  const transformFileEndLogEntry = {
    action_name: 'Transforming file',
    action_phase: 'end',
    file_name: filename,
    duration_ms,
    log_entry_label: 'Transforming file' };


  return {
    result: _extends({}, depsResult, { code, map }),
    transformFileStartLogEntry,
    transformFileEndLogEntry };

});


exports.transformAndExtractDependencies = (
transform,
filename,
localPath,
sourceCode,
options,
callback) =>
{
  /* $FlowFixMe: impossible to type a dynamic require */
  const transformModule = require(transform);
  transformCode(
  transformModule,
  filename,
  localPath,
  sourceCode,
  options,
  callback);

};

exports.minify = asyncify(
(filename, code, sourceMap) => {
  var result;
  try {
    result = minify.withSourceMap(code, sourceMap, filename);
  } catch (error) {
    if (error.constructor.name === 'JS_Parse_Error') {
      throw new Error(
      error.message +
      ' in file "' +
      filename +
      '" at line ' +
      error.line +
      ':' +
      error.col);

    }
    throw error;
  }
  return result;
});


exports.transformCode = transformCode; // for easier testing