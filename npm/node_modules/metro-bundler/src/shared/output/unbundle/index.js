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

const Server = require('../../../Server');

const asAssets = require('./as-assets');
const asIndexedFile = require('./as-indexed-file').save;




function buildBundle(packagerClient, requestOptions) {
  return packagerClient.buildBundle(_extends({},
  Server.DEFAULT_BUNDLE_OPTIONS,
  requestOptions, {
    unbundle: true,
    isolateModuleIDs: true }));

}

function saveUnbundle(
bundle,
options,
log)
{
  // we fork here depending on the platform:
  // while android is pretty good at loading individual assets, ios has a large
  // overhead when reading hundreds pf assets from disk
  return options.platform === 'android' && !options.indexedUnbundle ?
  asAssets(bundle, options, log) :
  /* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an
                                    * error found when Flow v0.54 was deployed. To see the error delete this
                                    * comment and run Flow. */
  asIndexedFile(bundle, options, log);
}

exports.build = buildBundle;
exports.save = saveUnbundle;
exports.formatName = 'bundle';