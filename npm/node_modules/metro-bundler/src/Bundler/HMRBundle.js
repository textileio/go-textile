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

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};var _slicedToArray = function () {function sliceIterator(arr, i) {var _arr = [];var _n = true;var _d = false;var _e = undefined;try {for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {_arr.push(_s.value);if (i && _arr.length === i) break;}} catch (err) {_d = true;_e = err;} finally {try {if (!_n && _i["return"]) _i["return"]();} finally {if (_d) throw _e;}}return _arr;}return function (arr, i) {if (Array.isArray(arr)) {return arr;} else if (Symbol.iterator in Object(arr)) {return sliceIterator(arr, i);} else {throw new TypeError("Invalid attempt to destructure non-iterable instance");}};}();

const BundleBase = require('./BundleBase');
const ModuleTransport = require('../lib/ModuleTransport');





class HMRBundle extends BundleBase {





  constructor(_ref)





  {let sourceURLFn = _ref.sourceURLFn,sourceMappingURLFn = _ref.sourceMappingURLFn;
    super();
    this._sourceURLFn = sourceURLFn;
    this._sourceMappingURLFn = sourceMappingURLFn;
    this._sourceURLs = [];
    this._sourceMappingURLs = [];
  }

  addModule(
  /* $FlowFixMe: broken OOP design: function signature should be the same */
  resolver,
  /* $FlowFixMe: broken OOP design: function signature should be the same */
  response,
  /* $FlowFixMe: broken OOP design: function signature should be the same */
  module,
  /* $FlowFixMe: broken OOP design: function signature should be the same */
  moduleTransport)
  {
    const dependencyPairs = response.getResolvedDependencyPairs(module);

    const dependencyPairsMap = new Map();
    for (const _ref2 of dependencyPairs) {var _ref3 = _slicedToArray(_ref2, 2);const relativePath = _ref3[0];const module = _ref3[1];
      dependencyPairsMap.set(relativePath, module.path);
    }

    const code = resolver.resolveRequires(
    module,
    /* $FlowFixMe: `getModuleId` is monkey-patched so may not exist */
    response.getModuleId,
    moduleTransport.code,
    dependencyPairsMap,
    /* $FlowFixMe: may not exist */
    moduleTransport.meta.dependencyOffsets);


    super.addModule(new ModuleTransport(_extends({}, moduleTransport, { code })));
    this._sourceMappingURLs.push(
    this._sourceMappingURLFn(moduleTransport.sourcePath));

    this._sourceURLs.push(this._sourceURLFn(moduleTransport.sourcePath));
    // inconsistent with parent class return type
    return Promise.resolve();
  }

  getModulesIdsAndCode() {
    return this.__modules.map(module => {
      return {
        id: JSON.stringify(module.id),
        code: module.code };

    });
  }

  getSourceURLs() {
    return this._sourceURLs;
  }

  getSourceMappingURLs() {
    return this._sourceMappingURLs;
  }}


module.exports = HMRBundle;