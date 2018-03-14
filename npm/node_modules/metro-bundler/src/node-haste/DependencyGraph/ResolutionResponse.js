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

'use strict';

const NO_OPTIONS = {};

class ResolutionResponse {






  // This is monkey-patched from Resolver.






  constructor(options) {
    this.dependencies = [];
    this.mainModuleId = null;
    this.mocks = null;
    this.numPrependedDependencies = 0;
    this.options = options;
    /* $FlowFixMe(>=0.56.0 site=react_native_fb) This comment suppresses an
                             * error found when Flow v0.56 was deployed. To see the error delete this
                             * comment and run Flow. */
    this._mappings = Object.create(null);
    this._finalized = false;
  }

  copy(properties)



  {var _properties$dependenc =




    properties.dependencies;const dependencies = _properties$dependenc === undefined ? this.dependencies : _properties$dependenc;var _properties$mainModul = properties.mainModuleId;const mainModuleId = _properties$mainModul === undefined ? this.mainModuleId : _properties$mainModul;var _properties$mocks = properties.mocks;const mocks = _properties$mocks === undefined ? this.mocks : _properties$mocks;

    const numPrependedDependencies =
    dependencies === this.dependencies ? this.numPrependedDependencies : 0;

    /* $FlowFixMe: Flow doesn't like Object.assign on class-made objects. */
    return Object.assign(new this.constructor(this.options), this, {
      dependencies,
      mainModuleId,
      mocks,
      numPrependedDependencies });

  }

  _assertNotFinalized() {
    if (this._finalized) {
      throw new Error('Attempted to mutate finalized response.');
    }
  }

  _assertFinalized() {
    if (!this._finalized) {
      throw new Error('Attempted to access unfinalized response.');
    }
  }

  finalize() {
    return Promise.resolve().then(() => {
      /* $FlowFixMe: _mainModule is not initialized in the constructor. */
      this.mainModuleId = this._mainModule.getName();
      this._finalized = true;
      return this;
    });
  }

  pushDependency(module) {
    this._assertNotFinalized();
    if (this.dependencies.length === 0) {
      this._mainModule = module;
    }

    this.dependencies.push(module);
  }

  prependDependency(module) {
    this._assertNotFinalized();
    this.dependencies.unshift(module);
    this.numPrependedDependencies += 1;
  }

  setResolvedDependencyPairs(
  module,
  pairs)

  {let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : NO_OPTIONS;
    if (!options.ignoreFinalized) {
      this._assertNotFinalized();
    }
    const hash = module.hash();
    if (this._mappings[hash] == null) {
      this._mappings[hash] = pairs;
    }
  }

  setMocks(mocks) {
    this.mocks = mocks;
  }

  getResolvedDependencyPairs(
  module)
  {
    this._assertFinalized();
    return this._mappings[module.hash()] || [];
  }}


module.exports = ResolutionResponse;