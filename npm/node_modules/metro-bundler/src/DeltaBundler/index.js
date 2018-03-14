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

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}

const DeltaTransformer = require('./DeltaTransformer');














/**
                                                         * `DeltaBundler` uses the `DeltaTransformer` to build bundle deltas. This
                                                         * module handles all the transformer instances so it can support multiple
                                                         * concurrent clients requesting their own deltas. This is done through the
                                                         * `deltaBundleId` options param (which maps a client to a specific delta
                                                         * transformer).
                                                         */
class DeltaBundler {





  constructor(bundler, options) {this._deltaTransformers = new Map();this._currentId = 0;
    this._bundler = bundler;
    this._options = options;
  }

  end() {
    this._deltaTransformers.forEach(DeltaTransformer => DeltaTransformer.end());
    this._deltaTransformers = new Map();
  }

  getDeltaTransformer(
  options)
  {var _this = this;return _asyncToGenerator(function* () {
      let bundleId = options.deltaBundleId;

      // If no bundle id is passed, generate a new one (which is going to be
      // returned as part of the bundle, so the client can later ask for an actual
      // delta).
      if (!bundleId) {
        bundleId = String(_this._currentId++);
      }

      let deltaTransformer = _this._deltaTransformers.get(bundleId);

      if (!deltaTransformer) {
        deltaTransformer = yield DeltaTransformer.create(
        _this._bundler,
        _this._options, _extends({},

        options, { // The Delta Bundler does not support minifying due to
          minify: false // issues generating the source maps (T21699790).
        }));


        _this._deltaTransformers.set(bundleId, deltaTransformer);
      }

      return {
        deltaTransformer,
        id: bundleId };})();

  }}


module.exports = DeltaBundler;