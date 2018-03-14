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

'use strict';var _slicedToArray = function () {function sliceIterator(arr, i) {var _arr = [];var _n = true;var _d = false;var _e = undefined;try {for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {_arr.push(_s.value);if (i && _arr.length === i) break;}} catch (err) {_d = true;_e = err;} finally {try {if (!_n && _i["return"]) _i["return"]();} finally {if (_d) throw _e;}}return _arr;}return function (arr, i) {if (Array.isArray(arr)) {return arr;} else if (Symbol.iterator in Object(arr)) {return sliceIterator(arr, i);} else {throw new TypeError("Invalid attempt to destructure non-iterable instance");}};}();function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}

const DependencyGraph = require('../node-haste/DependencyGraph');

const defaults = require('../defaults');
const pathJoin = require('path').join;













































class Resolver {






  constructor(opts, depGraph) {
    this._getPolyfills = opts.getPolyfills;
    this._minifyCode = opts.minifyCode;
    this._postMinifyProcess = opts.postMinifyProcess;
    this._polyfillModuleNames = opts.polyfillModuleNames || [];
    this._depGraph = depGraph;
  }

  static load(opts) {return _asyncToGenerator(function* () {
      const depGraph = yield DependencyGraph.load({
        assetDependencies: [opts.assetRegistryPath],
        assetExts: opts.assetExts,
        extraNodeModules: opts.extraNodeModules,
        forceNodeFilesystemAPI: false,
        getTransformCacheKey: opts.getTransformCacheKey,
        globalTransformCache: opts.globalTransformCache,
        ignorePattern: opts.blacklistRE || / ^/ /* matches nothing */,
        maxWorkers: opts.maxWorkers,
        moduleOptions: {
          hasteImpl: opts.hasteImpl,
          resetCache: opts.resetCache,
          transformCache: opts.transformCache },

        platforms: opts.platforms,
        preferNativePlatform: true,
        providesModuleNodeModules: opts.providesModuleNodeModules,
        reporter: opts.reporter,
        resetCache: opts.resetCache,
        roots: opts.projectRoots,
        sourceExts: opts.sourceExts,
        transformCode: opts.transformCode,
        useWatchman: true,
        watch: opts.watch });

      return new Resolver(opts, depGraph);})();
  }

  getShallowDependencies(
  entryFile,
  transformOptions)
  {
    return this._depGraph.getShallowDependencies(entryFile, transformOptions);
  }

  getModuleForPath(entryFile) {
    return this._depGraph.getModuleForPath(entryFile);
  }

  getDependencies(
  entryPath,
  options,




  bundlingOptions,
  onProgress,
  getModuleId)
  {var _this = this;return _asyncToGenerator(function* () {const
      platform = options.platform;var _options$recursive = options.recursive;const recursive = _options$recursive === undefined ? true : _options$recursive,prependPolyfills = options.prependPolyfills;

      const resolutionResponse =


      yield _this._depGraph.getDependencies({
        entryPath,
        platform,
        options: bundlingOptions,
        recursive,
        onProgress });


      if (prependPolyfills) {
        _this._getPolyfillDependencies(platform).
        reverse().
        forEach(function (polyfill) {return resolutionResponse.prependDependency(polyfill);});
      }

      /* $FlowFixMe: monkey patching */
      resolutionResponse.getModuleId = getModuleId;
      return resolutionResponse.finalize();})();
  }

  getModuleSystemDependencies(_ref) {var _ref$dev = _ref.dev;let dev = _ref$dev === undefined ? true : _ref$dev;
    const prelude = dev ?
    pathJoin(__dirname, 'polyfills/prelude_dev.js') :
    pathJoin(__dirname, 'polyfills/prelude.js');

    const moduleSystem = defaults.moduleSystem;

    return [prelude, moduleSystem].map(moduleName =>
    this._depGraph.createPolyfill({
      file: moduleName,
      id: moduleName,
      dependencies: [] }));


  }

  _getPolyfillDependencies(platform) {
    const polyfillModuleNames = this._getPolyfills({ platform }).concat(
    this._polyfillModuleNames);


    return polyfillModuleNames.map((polyfillModuleName, idx) =>
    this._depGraph.createPolyfill({
      file: polyfillModuleName,
      id: polyfillModuleName,
      dependencies: polyfillModuleNames.slice(0, idx) }));


  }

  resolveRequires(
  module,
  getModuleId,
  code,
  dependencyPairs)

  {let dependencyOffsets = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : [];
    const resolvedDeps = Object.create(null);

    // here, we build a map of all require strings (relative and absolute)
    // to the canonical ID of the module they reference
    for (const _ref2 of dependencyPairs) {var _ref3 = _slicedToArray(_ref2, 2);const name = _ref3[0];const path = _ref3[1];
      resolvedDeps[name] = getModuleId({ path });
    }

    // if we have a canonical ID for the module imported here,
    // we use it, so that require() is always called with the same
    // id for every module.
    // Example:
    // -- in a/b.js:
    //    require('./c') => require(3);
    // -- in b/index.js:
    //    require('../a/c') => require(3);
    return dependencyOffsets.
    reduceRight(
    (_ref4, offset) => {var _ref5 = _slicedToArray(_ref4, 2);let unhandled = _ref5[0],handled = _ref5[1];return [
      unhandled.slice(0, offset),
      replaceDependencyID(unhandled.slice(offset) + handled, resolvedDeps)];},

    [code, '']).

    join('');
  }

  wrapModule(_ref6)

















  {let module = _ref6.module,getModuleId = _ref6.getModuleId,dependencyPairs = _ref6.dependencyPairs,dependencyOffsets = _ref6.dependencyOffsets,name = _ref6.name,map = _ref6.map,code = _ref6.code;var _ref6$dev = _ref6.dev;let dev = _ref6$dev === undefined ? true : _ref6$dev;
    if (module.isJSON()) {
      code = `module.exports = ${code}`;
    }

    if (module.isPolyfill()) {
      code = definePolyfillCode(code);
    } else {
      const moduleId = getModuleId(module);

      code = this.resolveRequires(
      module,
      getModuleId,
      code,
      dependencyPairs,
      dependencyOffsets);

      code = defineModuleCode(moduleId, code, name, dev);
    }

    return { code, map };
  }

  minifyModule(_ref7)







  {var _this2 = this;let path = _ref7.path,code = _ref7.code,map = _ref7.map;return _asyncToGenerator(function* () {
      const minified = yield _this2._minifyCode(path, code, map);
      return yield _this2._postMinifyProcess(minified);})();
  }

  getDependencyGraph() {
    return this._depGraph;
  }}


function defineModuleCode(moduleName, code) {let verboseName = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : '';let dev = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : true;
  return [
  `__d(/* ${verboseName} */`,
  'function(global, require, module, exports) {', // module factory
  code,
  '\n}, ',
  `${JSON.stringify(moduleName)}`, // module id, null = id map. used in ModuleGraph
  dev ? `, null, ${JSON.stringify(verboseName)}` : '',
  ');'].
  join('');
}

function definePolyfillCode(code) {
  return [
  '(function(global) {',
  code,
  `\n})(typeof global !== 'undefined' ? global : typeof self !== 'undefined' ? self : this);`].
  join('');
}

const reDepencencyString = /^(['"])([^'"']*)\1/;
function replaceDependencyID(stringWithDependencyIDAtStart, resolvedDeps) {
  const match = reDepencencyString.exec(stringWithDependencyIDAtStart);
  const dependencyName = match && match[2];
  if (match != null && dependencyName in resolvedDeps) {const
    length = match[0].length;
    const id = String(resolvedDeps[dependencyName]);
    return (
      padRight(id, length) +
      stringWithDependencyIDAtStart.
      slice(length).
      replace(/$/m, ` // ${id} = ${dependencyName}`));

  } else {
    return stringWithDependencyIDAtStart;
  }
}

function padRight(string, length) {
  return string.length < length ?
  string + Array(length - string.length + 1).join(' ') :
  string;
}

module.exports = Resolver;