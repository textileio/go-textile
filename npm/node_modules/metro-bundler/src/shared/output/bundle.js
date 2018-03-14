/**
 * Copyright (c) 2015-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * 
 */

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};

const Server = require('../../Server');

const meta = require('./meta');
const relativizeSourceMap = require('../../lib/relativizeSourceMap');
const writeFile = require('./writeFile');





function buildBundle(packagerClient, requestOptions) {
  return packagerClient.buildBundle(_extends({},
  Server.DEFAULT_BUNDLE_OPTIONS,
  requestOptions, {
    isolateModuleIDs: true }));

}

function createCodeWithMap(
bundle,
dev,
sourceMapSourcesRoot)
{
  const map = bundle.getSourceMap({ dev });
  const sourceMap = relativizeSourceMap(
  typeof map === 'string' ? JSON.parse(map) : map,
  sourceMapSourcesRoot);
  return {
    code: bundle.getSource({ dev }),
    map: sourceMap };

}

function saveBundleAndMap(
bundle,
options,
log)



{const

  bundleOutput =




  options.bundleOutput,encoding = options.bundleEncoding,dev = options.dev,sourcemapOutput = options.sourcemapOutput,sourcemapSourcesRoot = options.sourcemapSourcesRoot;

  log('start');
  const origCodeWithMap = createCodeWithMap(bundle, !!dev, sourcemapSourcesRoot);
  const codeWithMap = bundle.postProcessBundleSourcemap(_extends({},
  origCodeWithMap, {
    outFileName: bundleOutput }));

  log('finish');

  log('Writing bundle output to:', bundleOutput);const

  code = codeWithMap.code;
  const writeBundle = writeFile(bundleOutput, code, encoding);
  const writeMetadata = writeFile(
  bundleOutput + '.meta',
  meta(code, encoding),
  'binary');
  Promise.all([writeBundle, writeMetadata]).
  then(() => log('Done writing bundle output'));

  if (sourcemapOutput) {
    log('Writing sourcemap output to:', sourcemapOutput);
    const map = typeof codeWithMap.map !== 'string' ?
    JSON.stringify(codeWithMap.map) :
    codeWithMap.map;
    const writeMap = writeFile(sourcemapOutput, map, null);
    writeMap.then(() => log('Done writing sourcemap output'));
    return Promise.all([writeBundle, writeMetadata, writeMap]);
  } else {
    return writeBundle;
  }
}

exports.build = buildBundle;
exports.save = saveBundleAndMap;
exports.formatName = 'bundle';