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
'use strict';var _require =

require('./util');const combineSourceMaps = _require.combineSourceMaps,combineSourceMapsAddingOffsets = _require.combineSourceMapsAddingOffsets,joinModules = _require.joinModules;










module.exports = (_ref) => {let fixWrapperOffset = _ref.fixWrapperOffset,lazyModules = _ref.lazyModules,moduleGroups = _ref.moduleGroups,startupModules = _ref.startupModules;
  const options = fixWrapperOffset ? { fixWrapperOffset: true } : undefined;
  const startupModule = {
    code: joinModules(startupModules),
    id: Number.MIN_SAFE_INTEGER,
    map: combineSourceMaps(startupModules, undefined, options),
    sourcePath: '' };


  const map = combineSourceMapsAddingOffsets(
  [startupModule].concat(lazyModules),
  moduleGroups,
  options);

  delete map.x_facebook_offsets[Number.MIN_SAFE_INTEGER];
  return map;
};