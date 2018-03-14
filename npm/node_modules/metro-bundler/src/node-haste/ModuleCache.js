/**
 * Copyright (c) 2013-present, Facebook, Inc.
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

const AssetModule = require('./AssetModule');
const Module = require('./Module');
const Package = require('./Package');
const Polyfill = require('./Polyfill');

const toLocalPath = require('./lib/toLocalPath');









class ModuleCache {














  constructor(_ref,






















  platforms)
  {let assetDependencies = _ref.assetDependencies,depGraphHelpers = _ref.depGraphHelpers,extractRequires = _ref.extractRequires,getClosestPackage = _ref.getClosestPackage,getTransformCacheKey = _ref.getTransformCacheKey,globalTransformCache = _ref.globalTransformCache,moduleOptions = _ref.moduleOptions,roots = _ref.roots,reporter = _ref.reporter,transformCode = _ref.transformCode;
    this._assetDependencies = assetDependencies;
    this._getClosestPackage = getClosestPackage;
    this._getTransformCacheKey = getTransformCacheKey;
    this._globalTransformCache = globalTransformCache;
    this._depGraphHelpers = depGraphHelpers;
    /* $FlowFixMe(>=0.56.0 site=react_native_fb) This comment suppresses an
                                              * error found when Flow v0.56 was deployed. To see the error delete this
                                              * comment and run Flow. */
    this._moduleCache = Object.create(null);
    this._moduleOptions = moduleOptions;
    /* $FlowFixMe(>=0.56.0 site=react_native_fb) This comment suppresses an
                                          * error found when Flow v0.56 was deployed. To see the error delete this
                                          * comment and run Flow. */
    this._packageCache = Object.create(null);
    this._packageModuleMap = new WeakMap();
    this._platforms = platforms;
    this._transformCode = transformCode;
    this._reporter = reporter;
    this._roots = roots;
  }

  getModule(filePath) {
    if (!this._moduleCache[filePath]) {
      this._moduleCache[filePath] = new Module({
        depGraphHelpers: this._depGraphHelpers,
        file: filePath,
        getTransformCacheKey: this._getTransformCacheKey,
        globalTransformCache: this._globalTransformCache,
        localPath: toLocalPath(this._roots, filePath),
        moduleCache: this,
        options: this._moduleOptions,
        reporter: this._reporter,
        transformCode: this._transformCode });

    }
    return this._moduleCache[filePath];
  }

  getAllModules() {
    return this._moduleCache;
  }

  getAssetModule(filePath) {
    if (!this._moduleCache[filePath]) {
      /* FixMe: AssetModule does not need all these options. This is because
                                        * this is an incorrect OOP design in the first place: AssetModule, being
                                        * simpler than a normal Module, should not inherit the Module class.
                                        */
      this._moduleCache[filePath] = new AssetModule(
      {
        depGraphHelpers: this._depGraphHelpers,
        dependencies: this._assetDependencies,
        file: filePath,
        getTransformCacheKey: this._getTransformCacheKey,
        globalTransformCache: null,
        localPath: toLocalPath(this._roots, filePath),
        moduleCache: this,
        options: this._moduleOptions,
        reporter: this._reporter,
        transformCode: this._transformCode },

      this._platforms);

    }
    return this._moduleCache[filePath];
  }

  getPackage(filePath) {
    if (!this._packageCache[filePath]) {
      this._packageCache[filePath] = new Package({
        file: filePath });

    }
    return this._packageCache[filePath];
  }

  getPackageForModule(module) {
    if (this._packageModuleMap.has(module)) {
      const packagePath = this._packageModuleMap.get(module);
      // $FlowFixMe(>=0.37.0)
      if (this._packageCache[packagePath]) {
        return this._packageCache[packagePath];
      } else {
        this._packageModuleMap.delete(module);
      }
    }

    const packagePath = this._getClosestPackage(module.path);
    if (!packagePath) {
      return null;
    }

    this._packageModuleMap.set(module, packagePath);
    return this.getPackage(packagePath);
  }

  createPolyfill(_ref2) {let file = _ref2.file;
    /* $FlowFixMe: there are missing arguments. */
    return new Polyfill({
      depGraphHelpers: this._depGraphHelpers,
      file,
      getTransformCacheKey: this._getTransformCacheKey,
      localPath: toLocalPath(this._roots, file),
      moduleCache: this,
      options: this._moduleOptions,
      transformCode: this._transformCode });

  }

  processFileChange(type, filePath) {
    if (this._moduleCache[filePath]) {
      this._moduleCache[filePath].invalidate();
      delete this._moduleCache[filePath];
    }
    if (this._packageCache[filePath]) {
      this._packageCache[filePath].invalidate();
      delete this._packageCache[filePath];
    }
  }}


module.exports = ModuleCache;