/**
 * Copyright (c) 2015-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * 
 * @format
 */

'use strict';let getOrderedDependencyPaths = (() => {var _ref = _asyncToGenerator(









  function* (
  deltaBundler,
  assetsServer,
  projectRoots,
  options)
  {
    const modules = yield Serializers.getAllModules(deltaBundler, options);

    const dependencies = yield Promise.all(
    Array.from(modules.values()).map((() => {var _ref2 = _asyncToGenerator(function* (module) {
        if (module.type !== 'asset') {
          return [module.path];
        } else {
          const assetData = yield assetsServer.getAssetData(
          toLocalPath(projectRoots, module.path),
          options.platform);


          return assetData.files;
        }
      });return function (_x5) {return _ref2.apply(this, arguments);};})()));


    const output = [];
    for (const dependencyArray of dependencies) {
      output.push.apply(output, _toConsumableArray(dependencyArray));
    }
    return output;
  });return function getOrderedDependencyPaths(_x, _x2, _x3, _x4) {return _ref.apply(this, arguments);};})();function _toConsumableArray(arr) {if (Array.isArray(arr)) {for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) arr2[i] = arr[i];return arr2;} else {return Array.from(arr);}}function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}const Serializers = require('../DeltaBundler/Serializers');const toLocalPath = require('../node-haste/lib/toLocalPath');

module.exports = getOrderedDependencyPaths;