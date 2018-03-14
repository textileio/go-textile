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

const DeltaCalculator = require('./DeltaCalculator');

const createModuleIdFactory = require('../lib/createModuleIdFactory');var _require =

require('events');const EventEmitter = _require.EventEmitter;









































const globalCreateModuleId = createModuleIdFactory();

/**
                                                       * This class is in charge of creating the delta bundle with the actual
                                                       * transformed source code for each of the modified modules. For each modified
                                                       * module it returns a `DeltaModule` object that contains the basic information
                                                       * about that file. Modules that have been deleted contain a `null` module
                                                       * parameter.
                                                       *
                                                       * The actual return format is the following:
                                                       *
                                                       *   {
                                                       *     pre: [{id, module: {}}],   Scripts to be prepended before the actual
                                                       *                                modules.
                                                       *     post: [{id, module: {}}],  Scripts to be appended after all the modules
                                                       *                                (normally the initial require() calls).
                                                       *     delta: [{id, module: {}}], Actual bundle modules (dependencies).
                                                       *   }
                                                       */
class DeltaTransformer extends EventEmitter {









  constructor(
  bundler,
  resolver,
  deltaCalculator,
  options,
  bundleOptions)
  {
    super();this.













































































































































































































































































































































    _onFileChange = () => {
      this.emit('change');
    };this._bundler = bundler;this._resolver = resolver;this._deltaCalculator = deltaCalculator;this._getPolyfills = options.getPolyfills;this._polyfillModuleNames = options.polyfillModuleNames;this._bundleOptions = bundleOptions; // Only when isolateModuleIDs is true the Module IDs of this instance are
    // sandboxed from the rest.
    // Isolating them makes sense when we want to get consistent module IDs
    // between different builds of the same bundle (for example when building
    // production builds), while coupling them makes sense when we want
    // different bundles to share the same ids (on HMR, where we need to patch
    // the correct module).
    this._getModuleId = this._bundleOptions.isolateModuleIDs ? createModuleIdFactory() : globalCreateModuleId;this._deltaCalculator.on('change', this._onFileChange);}static create(bundler, options, bundleOptions) {return _asyncToGenerator(function* () {const resolver = yield bundler.getResolver();const deltaCalculator = new DeltaCalculator(bundler, resolver.getDependencyGraph(), bundleOptions);return new DeltaTransformer(bundler, resolver, deltaCalculator, options, bundleOptions);})();} /**
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             * Destroy the Delta Transformer and its calculator. This should be used to
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             * clean up memory and resources once this instance is not used anymore.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             */end() {this._deltaCalculator.removeListener('change', this._onFileChange);return this._deltaCalculator.end();} /**
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               * Main method to calculate the bundle delta. It returns a DeltaResult,
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               * which contain the source code of the modified and added modules and the
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               * list of removed modules.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               */getDelta() {var _this = this;return _asyncToGenerator(function* () {// If there is already a build in progress, wait until it finish to start
      // processing a new one (delta transformer doesn't support concurrent
      // builds).
      if (_this._currentBuildPromise) {yield _this._currentBuildPromise;}_this._currentBuildPromise = _this._getDelta();let result;try {result = yield _this._currentBuildPromise;} finally {_this._currentBuildPromise = null;}return result;})();}_getDelta() {var _this2 = this;return _asyncToGenerator(function* () {// Calculate the delta of modules.
      var _ref = yield _this2._deltaCalculator.getDelta();const modified = _ref.modified,deleted = _ref.deleted,reset = _ref.reset;const transformerOptions = yield _this2._deltaCalculator.getTransformerOptions();const dependencyEdges = _this2._deltaCalculator.getDependencyEdges(); // Get the transformed source code of each modified/added module.
      const modifiedDelta = yield _this2._transformModules(Array.from(modified.values()), transformerOptions, dependencyEdges);deleted.forEach(function (id) {modifiedDelta.set(_this2._getModuleId({ path: id }), null);}); // Return the source code that gets prepended to all the modules. This
      // contains polyfills and startup code (like the require() implementation).
      const prependSources = reset ? yield _this2._getPrepend(transformerOptions, dependencyEdges) : new Map(); // Return the source code that gets appended to all the modules. This
      // contains the require() calls to startup the execution of the modules.
      const appendSources = reset ? yield _this2._getAppend(dependencyEdges) : new Map(); // Inverse dependencies are needed for HMR.
      const inverseDependencies = _this2._getInverseDependencies(dependencyEdges);return { pre: prependSources, post: appendSources, delta: modifiedDelta, inverseDependencies, reset };})();}_getPrepend(transformOptions, dependencyEdges) {var _this3 = this;return _asyncToGenerator(function* () {// Get all the polyfills from the relevant option params (the
      // `getPolyfills()` method and the `polyfillModuleNames` variable).
      const polyfillModuleNames = _this3._getPolyfills({ platform: _this3._bundleOptions.platform }).concat(_this3._polyfillModuleNames); // The module system dependencies are scripts that need to be included at
      // the very beginning of the bundle (before any polyfill).
      const moduleSystemDeps = _this3._resolver.getModuleSystemDependencies({ dev: _this3._bundleOptions.dev });const modules = moduleSystemDeps.concat(polyfillModuleNames.map(function (polyfillModuleName, idx) {return _this3._resolver.getDependencyGraph().createPolyfill({ file: polyfillModuleName, id: polyfillModuleName, dependencies: [] });}));return yield _this3._transformModules(modules, transformOptions, dependencyEdges);})();}_getAppend(dependencyEdges) {var _this4 = this;return _asyncToGenerator(function* () {// Get the absolute path of the entry file, in order to be able to get the
      // actual correspondant module (and its moduleId) to be able to add the
      // correct require(); call at the very end of the bundle.
      const absPath = _this4._resolver.getDependencyGraph().getAbsolutePath(_this4._bundleOptions.entryFile);const entryPointModule = _this4._resolver.getModuleForPath(absPath); // First, get the modules correspondant to all the module names defined in
      // the `runBeforeMainModule` config variable. Then, append the entry point
      // module so the last thing that gets required is the entry point.
      const append = new Map(_this4._bundleOptions.runBeforeMainModule.map(function (path) {return _this4._resolver.getModuleForPath(path);}).concat(entryPointModule).filter(function (module) {return dependencyEdges.has(module.path);}).map(_this4._getModuleId).map(function (moduleId) {const code = `;require(${JSON.stringify(moduleId)})`;const name = 'require-' + String(moduleId);const path = name + '.js';return [moduleId, { code, map: null, name, source: code, path, type: 'require' }];}));if (_this4._bundleOptions.sourceMapUrl) {const code = '//# sourceMappingURL=' + _this4._bundleOptions.sourceMapUrl;append.set(_this4._getModuleId({ path: '/sourcemap.js' }), { code, map: null, name: 'sourcemap.js', path: '/sourcemap.js', source: code, type: 'comment' });}return append;})();} /**
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    * Converts the paths in the inverse dependendencies to module ids.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    */_getInverseDependencies(dependencyEdges) {const output = Object.create(null);for (const _ref2 of dependencyEdges.entries()) {var _ref3 = _slicedToArray(_ref2, 2);const path = _ref3[0];const inverseDependencies = _ref3[1].inverseDependencies;output[this._getModuleId({ path })] = Array.from(inverseDependencies).map(dep => this._getModuleId({ path: dep }));} /* $FlowFixMe(>=0.56.0 site=react_native_fb) This comment suppresses an
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             * error found when Flow v0.56 was deployed. To see the error delete this
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             * comment and run Flow. */return output;}_transformModules(modules, transformOptions, dependencyEdges) {var _this5 = this;return _asyncToGenerator(function* () {return new Map((yield Promise.all(modules.map(function (module) {return _this5._transformModule(module, transformOptions, dependencyEdges);}))));})();}_transformModule(module, transformOptions, dependencyEdges) {var _this6 = this;return _asyncToGenerator(function* () {const name = module.getName();const metadata = yield _this6._getMetadata(module, transformOptions);const edge = dependencyEdges.get(module.path);const dependencyPairs = edge ? edge.dependencies : new Map();const wrapped = _this6._bundleOptions.wrapModules ? _this6._resolver.wrapModule({ module, getModuleId: _this6._getModuleId, dependencyPairs, dependencyOffsets: metadata.dependencyOffsets || [], name, code: metadata.code, map: metadata.map, minify: _this6._bundleOptions.minify, dev: _this6._bundleOptions.dev }) : { code: _this6._resolver.resolveRequires(module, _this6._getModuleId, metadata.code, dependencyPairs, metadata.dependencyOffsets || []), map: metadata.map }; // Ignore the Source Maps if the output of the transformer is not our
      // custom rawMapping data structure, since the Delta bundler cannot process
      // them. This can potentially happen when the minifier is enabled (since
      // uglifyJS only returns standard Source Maps).
      const map = Array.isArray(wrapped.map) ? wrapped.map : undefined;return [_this6._getModuleId(module), { code: ';' + wrapped.code, map, name, source: metadata.source, path: module.path, type: _this6._getModuleType(module) }];})();}_getModuleType(module) {if (module.isAsset()) {return 'asset';}if (module.isPolyfill()) {return 'script';}return 'module';}_getMetadata(module, transformOptions) {var _this7 = this;return _asyncToGenerator(function* () {if (module.isAsset()) {const asset = yield _this7._bundler.generateAssetObjAndCode(module, _this7._bundleOptions.assetPlugins, _this7._bundleOptions.platform);return { code: asset.code, dependencyOffsets: asset.meta.dependencyOffsets, map: undefined, source: '' };}return yield module.read(transformOptions);})();}}module.exports = DeltaTransformer;