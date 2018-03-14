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

const AssetPaths = require('./lib/AssetPaths');
const Module = require('./Module');



class AssetModule extends Module {





  constructor(
  args,
  platforms)
  {
    super(args);var _AssetPaths$parse =
    AssetPaths.parse(this.path, platforms);const resolution = _AssetPaths$parse.resolution,name = _AssetPaths$parse.name,type = _AssetPaths$parse.type;
    this.resolution = resolution;
    this._name = name;
    this._type = type;
    this._dependencies = args.dependencies || [];
  }

  isHaste() {
    return false;
  }

  readCached() {
    return {
      /** $FlowFixMe: improper OOP design. AssetModule, being different from a
              * normal Module, shouldn't inherit it in the first place. */
      result: { dependencies: this._dependencies },
      outdatedDependencies: [] };

  }

  /** $FlowFixMe: improper OOP design. */
  readFresh() {
    return Promise.resolve({ dependencies: this._dependencies });
  }

  getName() {
    return super.getName().replace(/\/[^\/]+$/, `/${this._name}.${this._type}`);
  }

  hash() {
    return `AssetModule : ${this.path}`;
  }

  isJSON() {
    return false;
  }

  isAsset() {
    return true;
  }}


module.exports = AssetModule;