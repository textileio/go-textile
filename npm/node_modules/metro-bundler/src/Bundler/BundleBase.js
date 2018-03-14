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

const ModuleTransport = require('../lib/ModuleTransport');












class BundleBase {






  constructor() {
    this._finalized = false;
    this.__modules = [];
    this._assets = [];
    this._mainModuleId = undefined;
  }

  isEmpty() {
    return this.__modules.length === 0 && this._assets.length === 0;
  }

  getMainModuleId() {
    return this._mainModuleId;
  }

  setMainModuleId(moduleId) {
    this._mainModuleId = moduleId;
  }

  addModule(module) {
    if (!(module instanceof ModuleTransport)) {
      throw new Error('Expected a ModuleTransport object');
    }

    return this.__modules.push(module) - 1;
  }

  replaceModuleAt(index, module) {
    if (!(module instanceof ModuleTransport)) {
      throw new Error('Expeceted a ModuleTransport object');
    }

    this.__modules[index] = module;
  }

  getModules() {
    return this.__modules;
  }

  getAssets() {
    return this._assets;
  }

  addAsset(asset) {
    this._assets.push(asset);
  }

  finalize(options) {
    if (!options.allowUpdates) {
      Object.freeze(this.__modules);
      Object.freeze(this._assets);
    }

    this._finalized = true;
  }

  getSource(options) {
    this.assertFinalized();

    if (this._source) {
      return this._source;
    }

    this._source = this.__modules.map(module => module.code).join('\n');
    return this._source;
  }

  invalidateSource() {
    this._source = null;
  }

  assertFinalized(message) {
    if (!this._finalized) {
      throw new Error(
      message || 'Bundle needs to be finalized before getting any source');

    }
  }

  setRamGroups(ramGroups) {}}


module.exports = BundleBase;