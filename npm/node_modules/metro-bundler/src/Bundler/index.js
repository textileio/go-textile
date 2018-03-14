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

'use strict';var _slicedToArray = function () {function sliceIterator(arr, i) {var _arr = [];var _n = true;var _d = false;var _e = undefined;try {for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {_arr.push(_s.value);if (i && _arr.length === i) break;}} catch (err) {_d = true;_e = err;} finally {try {if (!_n && _i["return"]) _i["return"]();} finally {if (_d) throw _e;}}return _arr;}return function (arr, i) {if (Array.isArray(arr)) {return arr;} else if (Symbol.iterator in Object(arr)) {return sliceIterator(arr, i);} else {throw new TypeError("Invalid attempt to destructure non-iterable instance");}};}();var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};function _toArray(arr) {return Array.isArray(arr) ? arr : Array.from(arr);}function _toConsumableArray(arr) {if (Array.isArray(arr)) {for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) arr2[i] = arr[i];return arr2;} else {return Array.from(arr);}}function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}

const assert = require('assert');
const crypto = require('crypto');
const debug = require('debug')('Metro:Bundler');
const emptyFunction = require('fbjs/lib/emptyFunction');
const fs = require('fs');
const Transformer = require('../JSTransformer');
const Resolver = require('../Resolver');
const Bundle = require('./Bundle');
const HMRBundle = require('./HMRBundle');
const ModuleTransport = require('../lib/ModuleTransport');
const imageSize = require('image-size');
const path = require('path');
const denodeify = require('denodeify');
const defaults = require('../defaults');
const toLocalPath = require('../node-haste/lib/toLocalPath');
const createModuleIdFactory = require('../lib/createModuleIdFactory');var _require =

require('./util');const generateAssetTransformResult = _require.generateAssetTransformResult,isAssetTypeAnImage = _require.isAssetTypeAnImage;var _require2 =






require('path');const pathSeparator = _require2.sep,joinPath = _require2.join,pathDirname = _require2.dirname,extname = _require2.extname;

const VERSION = require('../../package.json').version;


















































const sizeOf = denodeify(imageSize);var _require3 =





require('../Logger');const createActionStartEntry = _require3.createActionStartEntry,createActionEndEntry = _require3.createActionEndEntry,log = _require3.log;const






















































hasOwnProperty = Object.hasOwnProperty;

class Bundler {








  constructor(opts) {
    this._opts = opts;

    opts.projectRoots.forEach(verifyRootExists);

    const transformModuleStr = fs.readFileSync(opts.transformModulePath);
    const transformModuleHash = crypto.
    createHash('sha1').
    update(transformModuleStr).
    digest('hex');

    const stableProjectRoots = opts.projectRoots.map(p => {
      return path.relative(path.join(__dirname, '../../../..'), p);
    });

    const cacheKeyParts = [
    'metro-bundler-cache',
    VERSION,
    opts.cacheVersion,
    stableProjectRoots.
    join(',').
    split(pathSeparator).
    join('-'),
    transformModuleHash];


    this._getModuleId = createModuleIdFactory();

    let getCacheKey = options => '';
    if (opts.transformModulePath) {
      /* $FlowFixMe: dynamic requires prevent static typing :'(  */
      const transformer = require(opts.transformModulePath);
      if (typeof transformer.getCacheKey !== 'undefined') {
        getCacheKey = transformer.getCacheKey;
      }
    }

    const transformCacheKey = crypto.
    createHash('sha1').
    update(cacheKeyParts.join('$')).
    digest('hex');

    debug(`Using transform cache key "${transformCacheKey}"`);
    this._transformer = new Transformer(
    opts.transformModulePath,
    opts.maxWorkers,
    {
      stdoutChunk: chunk =>
      opts.reporter.update({ type: 'worker_stdout_chunk', chunk }),
      stderrChunk: chunk =>
      opts.reporter.update({ type: 'worker_stderr_chunk', chunk }) },

    opts.workerPath);


    const getTransformCacheKey = options => {
      return transformCacheKey + getCacheKey(options);
    };

    this._resolverPromise = Resolver.load({
      assetExts: opts.assetExts,
      assetRegistryPath: opts.assetRegistryPath,
      blacklistRE: opts.blacklistRE,
      extraNodeModules: opts.extraNodeModules,
      getPolyfills: opts.getPolyfills,
      getTransformCacheKey,
      globalTransformCache: opts.globalTransformCache,
      hasteImpl: opts.hasteImpl,
      maxWorkers: opts.maxWorkers,
      minifyCode: this._transformer.minify,
      postMinifyProcess: this._opts.postMinifyProcess,
      platforms: new Set(opts.platforms),
      polyfillModuleNames: opts.polyfillModuleNames,
      projectRoots: opts.projectRoots,
      providesModuleNodeModules:
      opts.providesModuleNodeModules || defaults.providesModuleNodeModules,
      reporter: opts.reporter,
      resetCache: opts.resetCache,
      sourceExts: opts.sourceExts,
      transformCode: (module, code, transformCodeOptions) =>
      this._transformer.transformFile(
      module.path,
      module.localPath,
      code,
      transformCodeOptions),

      transformCache: opts.transformCache,
      watch: opts.watch });


    this._projectRoots = opts.projectRoots;
    this._assetServer = opts.assetServer;

    this._getTransformOptions = opts.getTransformOptions;
  }

  end() {
    this._transformer.kill();
    /* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an
                               * error found when Flow v0.54 was deployed. To see the error delete this
                               * comment and run Flow. */
    return this._resolverPromise.then(resolver =>
    resolver.
    getDependencyGraph().
    getWatcher().
    end());

  }

  bundle(options)




  {const
    dev = options.dev,minify = options.minify,unbundle = options.unbundle;
    const postProcessBundleSourcemap = this._opts.postProcessBundleSourcemap;
    return this._resolverPromise.
    then(resolver => resolver.getModuleSystemDependencies({ dev, unbundle })).
    then(moduleSystemDeps =>
    this._bundle(_extends({},
    options, {
      bundle: new Bundle({
        dev,
        minify,
        sourceMapUrl: options.sourceMapUrl,
        postProcessBundleSourcemap }),

      moduleSystemDeps })));


  }

  _sourceHMRURL(platform, hmrpath) {
    return this._hmrURL('', platform, 'bundle', hmrpath);
  }

  _sourceMappingHMRURL(platform, hmrpath) {
    // Chrome expects `sourceURL` when eval'ing code
    return this._hmrURL('//# sourceURL=', platform, 'map', hmrpath);
  }

  _hmrURL(
  prefix,
  platform,
  extensionOverride,
  filePath)
  {
    const matchingRoot = this._projectRoots.find(root =>
    filePath.startsWith(root));


    if (!matchingRoot) {
      throw new Error('No matching project root for ' + filePath);
    }

    // Replaces '\' with '/' for Windows paths.
    if (pathSeparator === '\\') {
      filePath = filePath.replace(/\\/g, '/');
    }

    const extensionStart = filePath.lastIndexOf('.');
    const resource = filePath.substring(
    matchingRoot.length,
    extensionStart !== -1 ? extensionStart : undefined);


    return (
      prefix +
      resource +
      '.' +
      extensionOverride +
      '?' +
      'platform=' + (
      platform || '') +
      '&runModule=false&entryModuleOnly=true');

  }

  hmrBundle(
  options,
  host,
  port)
  {
    return this._bundle(_extends({},
    options, {
      bundle: new HMRBundle({
        sourceURLFn: this._sourceHMRURL.bind(this, options.platform),
        sourceMappingURLFn: this._sourceMappingHMRURL.bind(
        this,
        options.platform) }),


      hot: true,
      dev: true }));

  }

  _bundle(_ref)

































  {let assetPlugins = _ref.assetPlugins,bundle = _ref.bundle,dev = _ref.dev,entryFile = _ref.entryFile,entryModuleOnly = _ref.entryModuleOnly,generateSourceMaps = _ref.generateSourceMaps,hot = _ref.hot,isolateModuleIDs = _ref.isolateModuleIDs,minify = _ref.minify;var _ref$moduleSystemDeps = _ref.moduleSystemDeps;let moduleSystemDeps = _ref$moduleSystemDeps === undefined ? [] : _ref$moduleSystemDeps,onProgress = _ref.onProgress,platform = _ref.platform,resolutionResponse = _ref.resolutionResponse,runBeforeMainModule = _ref.runBeforeMainModule,runModule = _ref.runModule,unbundle = _ref.unbundle;
    const onResolutionResponse =
    response =>
    {
      /* $FlowFixMe: looks like ResolutionResponse is monkey-patched
       * with `getModuleId`. */
      bundle.setMainModuleId(response.getModuleId(getMainModule(response)));
      if (entryModuleOnly && entryFile) {
        response.dependencies = response.dependencies.filter(module =>
        module.path.endsWith(entryFile || ''));

      } else {
        response.dependencies = moduleSystemDeps.concat(response.dependencies);
      }
    };
    const finalizeBundle = (_ref2) => {let
      finalBundle = _ref2.bundle,
      transformedModules = _ref2.transformedModules,
      response = _ref2.response,
      modulesByPath = _ref2.modulesByPath;return (






        this._resolverPromise.
        then(resolver =>
        Promise.all(
        transformedModules.map((_ref3) => {let module = _ref3.module,transformed = _ref3.transformed;return (
            finalBundle.addModule(resolver, response, module, transformed));}))).



        then(() => {
          return Promise.all(
          runBeforeMainModule ?
          runBeforeMainModule.map(path => this.getModuleForPath(path)) :
          []);

        }).
        then(runBeforeMainModules => {
          runBeforeMainModules = runBeforeMainModules.
          map(module => modulesByPath[module.path]).
          filter(Boolean);

          finalBundle.finalize({
            runModule,
            runBeforeMainModule: runBeforeMainModules.map(module =>
            /* $FlowFixMe: looks like ResolutionResponse is monkey-patched
                                                                     * with `getModuleId`. */
            response.getModuleId(module)),

            allowUpdates: this._opts.allowBundleUpdates });

          return finalBundle;
        }));};

    return this._buildBundle({
      entryFile,
      dev,
      minify,
      platform,
      bundle,
      hot,
      unbundle,
      resolutionResponse,
      onResolutionResponse,
      finalizeBundle,
      isolateModuleIDs,
      generateSourceMaps,
      assetPlugins,
      onProgress });

  }

  _buildBundle(_ref4)















  {let entryFile = _ref4.entryFile,dev = _ref4.dev,minify = _ref4.minify,platform = _ref4.platform,bundle = _ref4.bundle,hot = _ref4.hot,unbundle = _ref4.unbundle,resolutionResponse = _ref4.resolutionResponse,isolateModuleIDs = _ref4.isolateModuleIDs,generateSourceMaps = _ref4.generateSourceMaps,assetPlugins = _ref4.assetPlugins;var _ref4$onResolutionRes = _ref4.onResolutionResponse;let onResolutionResponse = _ref4$onResolutionRes === undefined ? emptyFunction : _ref4$onResolutionRes;var _ref4$onModuleTransfo = _ref4.onModuleTransformed;let onModuleTransformed = _ref4$onModuleTransfo === undefined ? emptyFunction : _ref4$onModuleTransfo;var _ref4$finalizeBundle = _ref4.finalizeBundle;let finalizeBundle = _ref4$finalizeBundle === undefined ? emptyFunction : _ref4$finalizeBundle;var _ref4$onProgress = _ref4.onProgress;let onProgress = _ref4$onProgress === undefined ? emptyFunction : _ref4$onProgress;
    const transformingFilesLogEntry = log(
    createActionStartEntry({
      action_name: 'Transforming files',
      entry_point: entryFile,
      environment: dev ? 'dev' : 'prod' }));



    const modulesByPath = Object.create(null);

    if (!resolutionResponse) {
      resolutionResponse = this.getDependencies({
        entryFile,
        rootEntryFile: entryFile,
        dev,
        platform,
        hot,
        onProgress,
        minify,
        isolateModuleIDs,
        generateSourceMaps: unbundle || minify || generateSourceMaps,
        prependPolyfills: true });

    }

    /* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an
       * error found when Flow v0.54 was deployed. To see the error delete this
       * comment and run Flow. */
    return Promise.all([
    this._resolverPromise,
    resolutionResponse]).
    then((_ref5) => {var _ref6 = _slicedToArray(_ref5, 2);let resolver = _ref6[0],response = _ref6[1];
      bundle.setRamGroups(response.options.ramGroups);

      log(createActionEndEntry(transformingFilesLogEntry));
      onResolutionResponse(response);

      // get entry file complete path (`entryFile` is a local path, i.e. relative to roots)
      let entryFilePath;
      if (response.dependencies.length > 1) {
        // skip HMR requests
        const numModuleSystemDependencies = resolver.getModuleSystemDependencies(
        { dev, unbundle }).
        length;

        const dependencyIndex =
        (response.numPrependedDependencies || 0) +
        numModuleSystemDependencies;

        if (dependencyIndex in response.dependencies) {
          entryFilePath = response.dependencies[dependencyIndex].path;
        }
      }

      const modulesByTransport = new Map();
      const toModuleTransport = module =>
      this._toModuleTransport({
        module,
        bundle,
        entryFilePath,
        assetPlugins,
        options: response.options,
        /* $FlowFixMe: `getModuleId` is monkey-patched */
        getModuleId: response.getModuleId,
        dependencyPairs: response.getResolvedDependencyPairs(module) }).
      then(transformed => {
        modulesByTransport.set(transformed, module);
        modulesByPath[module.path] = module;
        onModuleTransformed({
          module,
          response,
          bundle,
          transformed });

        return transformed;
      });

      const p = this._opts.postProcessModules;
      const postProcess = p ?
      modules => p(modules, entryFile, { dev, minify, platform }) :
      null;

      return Promise.all(response.dependencies.map(toModuleTransport)).
      then(postProcess).
      then(moduleTransports => {
        const transformedModules = moduleTransports.map(transformed => ({
          module: modulesByTransport.get(transformed),
          transformed }));

        return finalizeBundle({
          bundle,
          transformedModules,
          response,
          modulesByPath });

      }).
      then(() => bundle);
    });
  }

  getShallowDependencies(_ref7)

















  {var _this = this;let entryFile = _ref7.entryFile,rootEntryFile = _ref7.rootEntryFile,platform = _ref7.platform;var _ref7$dev = _ref7.dev;let dev = _ref7$dev === undefined ? true : _ref7$dev;var _ref7$minify = _ref7.minify;let minify = _ref7$minify === undefined ? !dev : _ref7$minify;var _ref7$hot = _ref7.hot;let hot = _ref7$hot === undefined ? false : _ref7$hot;var _ref7$generateSourceM = _ref7.generateSourceMaps;let generateSourceMaps = _ref7$generateSourceM === undefined ? false : _ref7$generateSourceM,transformerOptions = _ref7.transformerOptions;return _asyncToGenerator(function* () {
      if (!transformerOptions) {
        transformerOptions = (yield _this.getTransformOptions(rootEntryFile, {
          dev,
          generateSourceMaps,
          hot,
          minify,
          platform,
          prependPolyfills: false })).
        transformer;
      }

      const notNullOptions = transformerOptions;

      return _this._resolverPromise.then(function (resolver) {return (
          resolver.getShallowDependencies(entryFile, notNullOptions));});})();

  }

  getModuleForPath(entryFile) {
    return this._resolverPromise.then(resolver =>
    resolver.getModuleForPath(entryFile));

  }

  getDependencies(_ref8)























  {var _this2 = this;let entryFile = _ref8.entryFile,platform = _ref8.platform;var _ref8$dev = _ref8.dev;let dev = _ref8$dev === undefined ? true : _ref8$dev;var _ref8$minify = _ref8.minify;let minify = _ref8$minify === undefined ? !dev : _ref8$minify;var _ref8$hot = _ref8.hot;let hot = _ref8$hot === undefined ? false : _ref8$hot;var _ref8$recursive = _ref8.recursive;let recursive = _ref8$recursive === undefined ? true : _ref8$recursive;var _ref8$generateSourceM = _ref8.generateSourceMaps;let generateSourceMaps = _ref8$generateSourceM === undefined ? false : _ref8$generateSourceM;var _ref8$isolateModuleID = _ref8.isolateModuleIDs;let isolateModuleIDs = _ref8$isolateModuleID === undefined ? false : _ref8$isolateModuleID,rootEntryFile = _ref8.rootEntryFile,prependPolyfills = _ref8.prependPolyfills,onProgress = _ref8.onProgress;return _asyncToGenerator(function* () {
      const bundlingOptions = yield _this2.getTransformOptions(
      rootEntryFile,
      {
        dev,
        platform,
        hot,
        generateSourceMaps,
        minify,
        prependPolyfills });



      const resolver = yield _this2._resolverPromise;
      const response = yield resolver.getDependencies(
      entryFile,
      { dev, platform, recursive, prependPolyfills },
      bundlingOptions,
      onProgress,
      isolateModuleIDs ? createModuleIdFactory() : _this2._getModuleId);

      return response;})();
  }

  getOrderedDependencyPaths(_ref9)











  {let entryFile = _ref9.entryFile,dev = _ref9.dev,platform = _ref9.platform,minify = _ref9.minify,generateSourceMaps = _ref9.generateSourceMaps;
    /* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an
                                                                                                                                                   * error found when Flow v0.54 was deployed. To see the error delete this
                                                                                                                                                   * comment and run Flow. */
    return this.getDependencies({
      entryFile,
      rootEntryFile: entryFile,
      dev,
      platform,
      minify,
      generateSourceMaps,
      prependPolyfills: true }).
    then((_ref10) => {let dependencies = _ref10.dependencies;
      const ret = [];
      const promises = [];
      const placeHolder = {};
      dependencies.forEach(dep => {
        if (dep.isAsset()) {
          const localPath = toLocalPath(this._projectRoots, dep.path);
          promises.push(this._assetServer.getAssetData(localPath, platform));
          ret.push(placeHolder);
        } else {
          ret.push(dep.path);
        }
      });

      return Promise.all(promises).then(assetsData => {
        assetsData.forEach((_ref11) => {let files = _ref11.files;
          const index = ret.indexOf(placeHolder);
          ret.splice.apply(ret, [index, 1].concat(_toConsumableArray(files)));
        });
        return ret;
      });
    });
  }

  _toModuleTransport(_ref12)















  {let module = _ref12.module,bundle = _ref12.bundle,entryFilePath = _ref12.entryFilePath,options = _ref12.options,getModuleId = _ref12.getModuleId,dependencyPairs = _ref12.dependencyPairs,assetPlugins = _ref12.assetPlugins;
    let moduleTransport;
    const moduleId = getModuleId(module);
    const transformOptions = options.transformer;

    if (module.isAsset()) {
      moduleTransport = this._generateAssetModule(
      bundle,
      module,
      moduleId,
      assetPlugins,
      transformOptions.platform);

    }

    if (moduleTransport) {
      return Promise.resolve(moduleTransport);
    }

    return module.
    read(transformOptions).
    then((_ref13) => {let code = _ref13.code,dependencies = _ref13.dependencies,dependencyOffsets = _ref13.dependencyOffsets,map = _ref13.map,source = _ref13.source;
      const name = module.getName();const

      preloadedModules = options.preloadedModules;
      const isPolyfill = module.isPolyfill();
      const preloaded =
      module.path === entryFilePath ||
      isPolyfill ||
      preloadedModules &&
      hasOwnProperty.call(preloadedModules, module.path);

      return new ModuleTransport({
        name,
        id: moduleId,
        code,
        map,
        meta: { dependencies, dependencyOffsets, preloaded, dependencyPairs },
        polyfill: isPolyfill,
        sourceCode: source,
        sourcePath: module.path });

    });
  }

  generateAssetObjAndCode(
  module,
  assetPlugins)

  {let platform = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : null;
    const localPath = toLocalPath(this._projectRoots, module.path);
    var assetUrlPath = joinPath('/assets', pathDirname(localPath));

    // On Windows, change backslashes to slashes to get proper URL path from file path.
    if (pathSeparator === '\\') {
      assetUrlPath = assetUrlPath.replace(/\\/g, '/');
    }

    const isImage = isAssetTypeAnImage(extname(module.path).slice(1));

    return this._assetServer.
    getAssetData(localPath, platform).
    then(assetData => {
      return Promise.all([
      isImage ? sizeOf(assetData.files[0]) : null,
      assetData]);

    }).
    then(res => {
      const dimensions = res[0];
      const assetData = res[1];
      const scale = assetData.scales[0];
      const asset = {
        __packager_asset: true,
        fileSystemLocation: pathDirname(module.path),
        httpServerLocation: assetUrlPath,
        width: dimensions ? dimensions.width / scale : undefined,
        height: dimensions ? dimensions.height / scale : undefined,
        scales: assetData.scales,
        files: assetData.files,
        hash: assetData.hash,
        name: assetData.name,
        type: assetData.type };


      return this._applyAssetPlugins(assetPlugins, asset);
    }).
    then(asset => {var _generateAssetTransfo =




      generateAssetTransformResult(this._opts.assetRegistryPath, asset);const code = _generateAssetTransfo.code,dependencies = _generateAssetTransfo.dependencies,dependencyOffsets = _generateAssetTransfo.dependencyOffsets;
      return {
        asset,
        code,
        meta: { dependencies, dependencyOffsets, preloaded: null } };

    });
  }

  _applyAssetPlugins(
  assetPlugins,
  asset)
  {
    if (!assetPlugins.length) {
      return asset;
    }var _assetPlugins = _toArray(

    assetPlugins);const currentAssetPlugin = _assetPlugins[0],remainingAssetPlugins = _assetPlugins.slice(1);
    /* $FlowFixMe: dynamic requires prevent static typing :'(  */
    const assetPluginFunction = require(currentAssetPlugin);
    const result = assetPluginFunction(asset);

    // If the plugin was an async function, wait for it to fulfill before
    // applying the remaining plugins
    if (typeof result.then === 'function') {
      return result.then(resultAsset =>
      this._applyAssetPlugins(remainingAssetPlugins, resultAsset));

    } else {
      return this._applyAssetPlugins(remainingAssetPlugins, result);
    }
  }

  _generateAssetModule(
  bundle,
  module,
  moduleId)


  {let assetPlugins = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : [];let platform = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : null;
    return this.generateAssetObjAndCode(
    module,
    assetPlugins,
    platform).
    then((_ref14) => {let asset = _ref14.asset,code = _ref14.code,meta = _ref14.meta;
      bundle.addAsset(asset);
      return new ModuleTransport({
        name: module.getName(),
        id: moduleId,
        code,
        meta,
        sourceCode: code,
        sourcePath: module.path,
        virtual: true });

    });
  }

  getTransformOptions(
  mainModuleName,
  options)







  {var _this3 = this;return _asyncToGenerator(function* () {
      const getDependencies = function (entryFile) {return (
          _this3.getDependencies(_extends({},
          options, {
            enableBabelRCLookup: _this3._opts.enableBabelRCLookup,
            entryFile,
            projectRoots: _this3._projectRoots,
            rootEntryFile: entryFile,
            prependPolyfills: false })).
          then(function (r) {return r.dependencies.map(function (d) {return d.path;});}));};const

      dev = options.dev,hot = options.hot,platform = options.platform;
      const extraOptions = _this3._getTransformOptions ?
      yield _this3._getTransformOptions(
      [mainModuleName],
      { dev, hot, platform },
      getDependencies) :

      {};var _extraOptions$transfo =

      extraOptions.transform;const transform = _extraOptions$transfo === undefined ? {} : _extraOptions$transfo;

      return {
        transformer: {
          dev,
          minify: options.minify,
          platform,
          transform: {
            enableBabelRCLookup: _this3._opts.enableBabelRCLookup,
            dev,
            generateSourceMaps: options.generateSourceMaps,
            hot,
            inlineRequires: transform.inlineRequires || false,
            platform,
            projectRoot: _this3._projectRoots[0] } },


        preloadedModules: extraOptions.preloadedModules,
        ramGroups: extraOptions.ramGroups };})();

  }

  getResolver() {
    return this._resolverPromise;
  }}


function verifyRootExists(root) {
  // Verify that the root exists.
  assert(fs.statSync(root).isDirectory(), 'Root has to be a valid directory');
}

function getMainModule(_ref15) {let dependencies = _ref15.dependencies;var _ref15$numPrependedDe = _ref15.numPrependedDependencies;let numPrependedDependencies = _ref15$numPrependedDe === undefined ? 0 : _ref15$numPrependedDe;
  return dependencies[numPrependedDependencies];
}

module.exports = Bundler;