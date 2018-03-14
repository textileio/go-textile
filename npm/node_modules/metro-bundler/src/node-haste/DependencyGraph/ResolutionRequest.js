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

'use strict';var _slicedToArray = function () {function sliceIterator(arr, i) {var _arr = [];var _n = true;var _d = false;var _e = undefined;try {for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {_arr.push(_s.value);if (i && _arr.length === i) break;}} catch (err) {_d = true;_e = err;} finally {try {if (!_n && _i["return"]) _i["return"]();} finally {if (_d) throw _e;}}return _arr;}return function (arr, i) {if (Array.isArray(arr)) {return arr;} else if (Symbol.iterator in Object(arr)) {return sliceIterator(arr, i);} else {throw new TypeError("Invalid attempt to destructure non-iterable instance");}};}();

const AsyncTaskGroup = require('../lib/AsyncTaskGroup');
const MapWithDefaults = require('../lib/MapWithDefaults');
const ModuleResolution = require('./ModuleResolution');

const debug = require('debug')('Metro:DependencyGraph');
const isAbsolutePath = require('absolute-path');
const path = require('path');const

DuplicateHasteCandidatesError = require('jest-haste-map').ModuleMap.DuplicateHasteCandidatesError;const







UnableToResolveError = ModuleResolution.UnableToResolveError,isRelativeImport = ModuleResolution.isRelativeImport;






































class ResolutionRequest {




  constructor(options) {
    this._options = options;
    this._resetResolutionCache();
  }

  resolveDependency(fromModule, toModuleName) {
    const resHash = getResolutionCacheKey(fromModule.path, toModuleName);

    const immediateResolution = this._immediateResolutionCache[resHash];
    if (immediateResolution) {
      return immediateResolution;
    }

    const cacheResult = result => {
      this._immediateResolutionCache[resHash] = result;
      return result;
    };

    const resolver = this._options.moduleResolver;
    const platform = this._options.platform;

    if (
    !this._options.helpers.isNodeModulesDir(fromModule.path) &&
    !(isRelativeImport(toModuleName) || isAbsolutePath(toModuleName)))
    {
      const result = ModuleResolution.tryResolveSync(
      () => this._resolveHasteDependency(fromModule, toModuleName, platform),
      () =>
      resolver.resolveNodeDependency(fromModule, toModuleName, platform));

      return cacheResult(result);
    }

    return cacheResult(
    resolver.resolveNodeDependency(fromModule, toModuleName, platform));

  }

  _resolveHasteDependency(
  fromModule,
  toModuleName,
  platform)
  {
    const rs = this._options.moduleResolver;
    try {
      return rs.resolveHasteDependency(fromModule, toModuleName, platform);
    } catch (error) {
      if (error instanceof DuplicateHasteCandidatesError) {
        throw new AmbiguousModuleResolutionError(fromModule.path, error);
      }
      throw error;
    }
  }

  resolveModuleDependencies(
  module,
  dependencyNames)
  {
    const dependencies = dependencyNames.map(name =>
    this.resolveDependency(module, name));

    return [dependencyNames, dependencies];
  }

  getOrderedDependencies(_ref)









  {let response = _ref.response,transformOptions = _ref.transformOptions,onProgress = _ref.onProgress;var _ref$recursive = _ref.recursive;let recursive = _ref$recursive === undefined ? true : _ref$recursive;
    const entry = this._options.moduleCache.getModule(this._options.entryPath);

    response.pushDependency(entry);
    let totalModules = 1;
    let finishedModules = 0;

    let preprocessedModuleCount = 1;
    if (recursive) {
      this._preprocessPotentialDependencies(transformOptions, entry, count => {
        if (count + 1 <= preprocessedModuleCount) {
          return;
        }
        preprocessedModuleCount = count + 1;
        if (onProgress != null) {
          onProgress(finishedModules, preprocessedModuleCount);
        }
      });
    }

    const resolveDependencies = module =>
    Promise.resolve().then(() => {
      const cached = module.readCached(transformOptions);
      if (cached.result != null) {
        return this.resolveModuleDependencies(
        module,
        cached.result.dependencies);

      }
      return module.
      readFresh(transformOptions).
      then((_ref2) => {let dependencies = _ref2.dependencies;return (
          this.resolveModuleDependencies(module, dependencies));});

    });

    const collectedDependencies =


    new MapWithDefaults(module => collect(module));
    const crawlDependencies = (mod, _ref3) => {var _ref4 = _slicedToArray(_ref3, 2);let depNames = _ref4[0],dependencies = _ref4[1];
      const filteredPairs = [];

      dependencies.forEach((modDep, i) => {
        const name = depNames[i];
        if (modDep == null) {
          debug(
          'WARNING: Cannot find required module `%s` from module `%s`',
          name,
          mod.path);

          return false;
        }
        return filteredPairs.push([name, modDep]);
      });

      response.setResolvedDependencyPairs(mod, filteredPairs);

      const dependencyModules = filteredPairs.map((_ref5) => {var _ref6 = _slicedToArray(_ref5, 2);let m = _ref6[1];return m;});
      const newDependencies = dependencyModules.filter(
      m => !collectedDependencies.has(m));


      if (onProgress) {
        finishedModules += 1;
        totalModules += newDependencies.length;
        onProgress(
        finishedModules,
        Math.max(totalModules, preprocessedModuleCount));

      }

      if (recursive) {
        // doesn't block the return of this function invocation, but defers
        // the resulution of collectionsInProgress.done.then(...)
        dependencyModules.forEach(dependency =>
        collectedDependencies.get(dependency));

      }
      return dependencyModules;
    };

    const collectionsInProgress = new AsyncTaskGroup();
    function collect(module) {
      collectionsInProgress.start(module);
      const result = resolveDependencies(module).then(deps =>
      crawlDependencies(module, deps));

      const end = () => collectionsInProgress.end(module);
      result.then(end, end);
      return result;
    }

    function resolveKeyWithPromise(_ref7)

    {var _ref8 = _slicedToArray(_ref7, 2);let key = _ref8[0],promise = _ref8[1];
      return promise.then(value => [key, value]);
    }

    /* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an
       * error found when Flow v0.54 was deployed. To see the error delete this
       * comment and run Flow. */
    return Promise.all([
    // kicks off recursive dependency discovery, but doesn't block until it's
    // done
    collectedDependencies.get(entry),

    // resolves when there are no more modules resolving dependencies
    collectionsInProgress.done]).

    then((_ref9) => {var _ref10 = _slicedToArray(_ref9, 1);let rootDependencies = _ref10[0];
      return Promise.all(
      Array.from(collectedDependencies, resolveKeyWithPromise)).
      then(moduleToDependenciesPairs => [
      rootDependencies,
      new MapWithDefaults(() => [], moduleToDependenciesPairs)]);

    }).
    then((_ref11) => {var _ref12 = _slicedToArray(_ref11, 2);let rootDependencies = _ref12[0],moduleDependencies = _ref12[1];
      // serialize dependencies, and make sure that every single one is only
      // included once
      const seen = new Set([entry]);
      function traverse(dependencies) {
        dependencies.forEach(dependency => {
          if (seen.has(dependency)) {
            return;
          }

          seen.add(dependency);
          response.pushDependency(dependency);
          traverse(moduleDependencies.get(dependency));
        });
      }

      traverse(rootDependencies);
    });
  }

  /**
     * This synchronously look at all the specified modules and recursively kicks
     * off global cache fetching or transforming (via `readFresh`). This is a hack
     * that workaround the current structure, because we could do better. First
     * off, the algorithm that resolves dependencies recursively should be
     * synchronous itself until it cannot progress anymore (and needs to call
     * `readFresh`), so that this algo would be integrated into it.
     */
  _preprocessPotentialDependencies(
  transformOptions,
  module,
  onProgress)
  {
    const visitedModulePaths = new Set();
    const pendingBatches = [
    this.preprocessModule(transformOptions, module, visitedModulePaths)];

    onProgress(visitedModulePaths.size);
    while (pendingBatches.length > 0) {
      const dependencyModules = pendingBatches.pop();
      while (dependencyModules.length > 0) {
        const dependencyModule = dependencyModules.pop();
        const deps = this.preprocessModule(
        transformOptions,
        dependencyModule,
        visitedModulePaths);

        pendingBatches.push(deps);
        onProgress(visitedModulePaths.size);
      }
    }
  }

  preprocessModule(
  transformOptions,
  module,
  visitedModulePaths)
  {
    const cached = module.readCached(transformOptions);
    if (cached.result == null) {
      module.readFresh(transformOptions).catch(error => {
        /* ignore errors, they'll be handled later if the dependency is actually
                                                          * not obsolete, and required from somewhere */
      });
    }
    const dependencies =
    cached.result != null ?
    cached.result.dependencies :
    cached.outdatedDependencies;
    return this.tryResolveModuleDependencies(
    module,
    dependencies,
    visitedModulePaths);

  }

  tryResolveModuleDependencies(
  module,
  dependencyNames,
  visitedModulePaths)
  {
    const result = [];
    for (let i = 0; i < dependencyNames.length; ++i) {
      try {
        const depModule = this.resolveDependency(module, dependencyNames[i]);
        if (!visitedModulePaths.has(depModule.path)) {
          visitedModulePaths.add(depModule.path);
          result.push(depModule);
        }
      } catch (error) {
        if (!(error instanceof UnableToResolveError)) {
          throw error;
        }
      }
    }
    return result;
  }

  _resetResolutionCache() {
    /* $FlowFixMe(>=0.56.0 site=react_native_fb) This comment suppresses an
                            * error found when Flow v0.56 was deployed. To see the error delete this
                            * comment and run Flow. */
    this._immediateResolutionCache = Object.create(null);
  }

  getResolutionCache() {
    return this._immediateResolutionCache;
  }}


function getResolutionCacheKey(modulePath, depName) {
  return `${path.resolve(modulePath)}:${depName}`;
}

class AmbiguousModuleResolutionError extends Error {



  constructor(
  fromModulePath,
  hasteError)
  {
    super(
    `Ambiguous module resolution from \`${fromModulePath}\`: ` +
    hasteError.message);

    this.fromModulePath = fromModulePath;
    this.hasteError = hasteError;
  }}


ResolutionRequest.AmbiguousModuleResolutionError = AmbiguousModuleResolutionError;

module.exports = ResolutionRequest;