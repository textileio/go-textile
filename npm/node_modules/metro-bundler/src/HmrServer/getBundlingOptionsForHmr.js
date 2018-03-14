/**
 * Copyright (c) 2015-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * @format
 * 
 */

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};



/**
                                                                                                                                                                                                                                                                      * Module to easily create the needed configuration parameters needed for the
                                                                                                                                                                                                                                                                      * bundler for HMR (since a lot of params are not relevant in this use case).
                                                                                                                                                                                                                                                                      */
module.exports = function getBundlingOptionsForHmr(
entryFile,
platform)
{
  // These are the really meaningful bundling options. The others below are
  // not relevant for HMR.
  const mainOptions = {
    deltaBundleId: null,
    entryFile,
    hot: true,
    minify: false,
    platform,
    wrapModules: false };


  return _extends({},
  mainOptions, {
    assetPlugins: [],
    dev: true,
    entryModuleOnly: false,
    excludeSource: false,
    generateSourceMaps: false,
    inlineSourceMap: false,
    isolateModuleIDs: false,
    onProgress: null,
    resolutionResponse: null,
    runBeforeMainModule: [],
    runModule: false,
    sourceMapUrl: '',
    unbundle: false });

};