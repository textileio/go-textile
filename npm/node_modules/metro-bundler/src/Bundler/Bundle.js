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

const _ = require('lodash');
const crypto = require('crypto');
const debug = require('debug')('Metro:Bundle');
const invariant = require('fbjs/lib/invariant');var _require =

require('./util');const createRamBundleGroups = _require.createRamBundleGroups;var _require2 =
require('./source-map');const fromRawMappings = _require2.fromRawMappings;var _require3 =
require('../lib/SourceMap');const isMappingsMap = _require3.isMappingsMap;














const SOURCEMAPPING_URL = '\n//# sourceMappingURL=';

class Bundle extends BundleBase {











  constructor()













  {var _ref = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};let sourceMapUrl = _ref.sourceMapUrl,dev = _ref.dev,minify = _ref.minify,ramGroups = _ref.ramGroups,postProcessBundleSourcemap = _ref.postProcessBundleSourcemap;
    super();
    this._sourceMap = null;
    this._sourceMapFormat = 'undetermined';
    this._sourceMapUrl = sourceMapUrl;
    this._numRequireCalls = 0;
    this._dev = dev;
    this._minify = minify;

    this._ramGroups = ramGroups;
    this._ramBundle = null; // cached RAM Bundle
    this.postProcessBundleSourcemap = postProcessBundleSourcemap;
  }

  addModule(
  /**
              * $FlowFixMe: this code is inherently incorrect, because it modifies the
              * signature of the base class function "addModule". That means callsites
              * using an instance typed as the base class would be broken. This must be
              * refactored.
              */
  resolver,






  resolutionResponse,
  module,
  /* $FlowFixMe: erroneous change of signature. */
  moduleTransport)

  {
    const index = super.addModule(moduleTransport);

    const dependencyPairs = resolutionResponse.getResolvedDependencyPairs(
    module);


    const dependencyPairsMap = new Map();
    for (const _ref2 of dependencyPairs) {var _ref3 = _slicedToArray(_ref2, 2);const relativePath = _ref3[0];const dependencyModule = _ref3[1];
      dependencyPairsMap.set(relativePath, dependencyModule.path);
    }

    return Promise.resolve(
    resolver.wrapModule({
      module,
      getModuleId: resolutionResponse.getModuleId,
      dependencyPairs: dependencyPairsMap,
      name: moduleTransport.name,
      code: moduleTransport.code,
      map: moduleTransport.map,
      dependencyOffsets: moduleTransport.meta ?
      moduleTransport.meta.dependencyOffsets :
      undefined,
      dev: this._dev })).


    then((_ref4) => {let code = _ref4.code,map = _ref4.map;
      return this._minify ?
      resolver.minifyModule({ code, map, path: module.path }) :
      { code, map };
    }).
    then((_ref5) => {let code = _ref5.code,map = _ref5.map;
      // If we get a map from the transformer we'll switch to a mode
      // were we're combining the source maps as opposed to
      if (map) {
        const usesRawMappings = isRawMappings(map);

        if (this._sourceMapFormat === 'undetermined') {
          this._sourceMapFormat = usesRawMappings ? 'flattened' : 'indexed';
        } else if (usesRawMappings && this._sourceMapFormat === 'indexed') {
          throw new Error(
          `Got at least one module with a full source map, but ${moduleTransport.sourcePath} has raw mappings`);

        } else if (
        !usesRawMappings &&
        this._sourceMapFormat === 'flattened')
        {
          throw new Error(
          `Got at least one module with raw mappings, but ${moduleTransport.sourcePath} has a full source map`);

        }
      }

      this.replaceModuleAt(
      index,
      new ModuleTransport(_extends({}, moduleTransport, { code, map })));

    });
  }

  finalize(options) {
    options = options || {};
    if (options.runModule) {
      /* $FlowFixMe: this is unsound, as nothing enforces runBeforeMainModule
                             * to be available if `runModule` is true. Refactor. */
      options.runBeforeMainModule.forEach(this._addRequireCall, this);
      /* $FlowFixMe: this is unsound, as nothing enforces the module ID to have
                                                                        * been set beforehand. */
      this._addRequireCall(this.getMainModuleId());
    }

    super.finalize(options);
  }

  _addRequireCall(moduleId) {
    const code = `;require(${JSON.stringify(moduleId)});`;
    const name = 'require-' + moduleId;
    super.addModule(
    new ModuleTransport({
      name,
      id: -this._numRequireCalls - 1,
      code,
      virtual: true,
      sourceCode: code,
      sourcePath: name + '.js',
      meta: { preloaded: true } }));


    this._numRequireCalls += 1;
  }

  _getInlineSourceMap(dev) {
    if (this._inlineSourceMap == null) {
      const sourceMap = this.getSourceMapString({ excludeSource: true, dev });
      /*eslint-env node*/
      const encoded = new Buffer(sourceMap).toString('base64');
      this._inlineSourceMap = 'data:application/json;base64,' + encoded;
    }
    return this._inlineSourceMap;
  }

  getSource(options) {
    this.assertFinalized();

    options = options || {};

    let source = super.getSource(options);

    if (options.inlineSourceMap) {
      source += SOURCEMAPPING_URL + this._getInlineSourceMap(options.dev);
    } else if (this._sourceMapUrl) {
      source += SOURCEMAPPING_URL + this._sourceMapUrl;
    }

    return source;
  }

  getUnbundle() {
    this.assertFinalized();
    if (!this._ramBundle) {
      const modules = this.getModules().slice();

      // separate modules we need to preload from the ones we don't
      var _partition = partition(modules, shouldPreload),_partition2 = _slicedToArray(_partition, 2);const startupModules = _partition2[0],lazyModules = _partition2[1];

      const ramGroups = this._ramGroups;
      let groups;
      this._ramBundle = {
        startupModules,
        lazyModules,
        get groups() {
          if (!groups) {
            groups = createRamBundleGroups(
            ramGroups || [],
            lazyModules,
            subtree);

          }
          return groups;
        } };

    }

    return this._ramBundle;
  }

  invalidateSource() {
    debug('invalidating bundle');
    super.invalidateSource();
    this._sourceMap = null;
  }

  /**
     * Combine each of the sourcemaps multiple modules have into a single big
     * one. This works well thanks to a neat trick defined on the sourcemap spec
     * that makes use of of the `sections` field to combine sourcemaps by adding
     * an offset. This is supported only by Chrome for now.
     */
  _getCombinedSourceMaps(options) {
    const result = {
      version: 3,
      file: this._getSourceMapFile(),
      sections: [] };


    let line = 0;
    this.getModules().forEach(module => {
      invariant(
      !Array.isArray(module.map),
      `Unexpected raw mappings for ${module.sourcePath}`);

      let map =
      module.map == null || module.virtual ?
      generateSourceMapForVirtualModule(module) :
      module.map;

      if (options.excludeSource && isMappingsMap(map)) {
        map = _extends({}, map, { sourcesContent: [] });
      }

      result.sections.push({
        offset: { line, column: 0 },
        map });

      line += module.code.split('\n').length;
    });

    return result;
  }

  getSourceMap(options) {
    this.assertFinalized();

    return this._sourceMapFormat === 'indexed' ?
    this._getCombinedSourceMaps(options) :
    this._fromRawMappings().toMap(undefined, options);
  }

  getSourceMapString(options) {
    if (this._sourceMapFormat === 'indexed') {
      return JSON.stringify(this.getSourceMap(options));
    }

    // The following code is an optimization specific to the development server:
    // 1. generator.toSource() is faster than JSON.stringify(generator.toMap()).
    // 2. caching the source map unless there are changes saves time in
    //    development settings.
    let map = this._sourceMap;
    if (map == null) {
      debug('Start building flat source map');
      map = this._sourceMap = this._fromRawMappings().toString(
      undefined,
      options);

      debug('End building flat source map');
    } else {
      debug('Returning cached source map');
    }
    return map;
  }

  getEtag() {
    var eTag = crypto.
    createHash('md5')
    /* $FlowFixMe: we must pass options, or rename the
                       * base `getSource` function, as it does not actually need options. */.
    update(this.getSource()).
    digest('hex');
    return eTag;
  }

  _getSourceMapFile() {
    return this._sourceMapUrl ?
    this._sourceMapUrl.replace('.map', '.bundle') :
    'bundle.js';
  }

  getJSModulePaths() {
    return (
      this.getModules()
      // Filter out non-js files. Like images etc.
      .filter(module => !module.virtual).
      map(module => module.sourcePath));

  }

  getDebugInfo() {
    return [
    /* $FlowFixMe: this is unsound as the module ID could be unset. */
    '<div><h3>Main Module:</h3> ' + this.getMainModuleId() + '</div>',
    '<style>',
    'pre.collapsed {',
    '  height: 10px;',
    '  width: 100px;',
    '  display: block;',
    '  text-overflow: ellipsis;',
    '  overflow: hidden;',
    '  cursor: pointer;',
    '}',
    '</style>',
    '<h3> Module paths and transformed code: </h3>',
    this.getModules().
    map(function (m) {
      return (
        '<div> <h4> Path: </h4>' +
        m.sourcePath +
        '<br/> <h4> Source: </h4>' +
        '<code><pre class="collapsed" onclick="this.classList.remove(\'collapsed\')">' +
        _.escape(m.code) +
        '</pre></code></div>');

    }).
    join('\n')].
    join('\n');
  }

  setRamGroups(ramGroups) {
    this._ramGroups = ramGroups;
  }

  _fromRawMappings() {
    return fromRawMappings(
    this.getModules().map(module => ({
      map: Array.isArray(module.map) ? module.map : undefined,
      path: module.sourcePath,
      source: module.sourceCode,
      code: module.code })));


  }}


function generateSourceMapForVirtualModule(module) {
  // All lines map 1-to-1
  let mappings = 'AAAA;';

  for (let i = 1; i < module.code.split('\n').length; i++) {
    mappings += 'AACA;';
  }

  return {
    version: 3,
    sources: [module.sourcePath],
    names: [],
    mappings,
    file: module.sourcePath,
    sourcesContent: [module.sourceCode] };

}

function shouldPreload(_ref6) {let meta = _ref6.meta;
  return meta && meta.preloaded;
}

function partition(array, predicate) {
  const included = [];
  const excluded = [];
  array.forEach(item => (predicate(item) ? included : excluded).push(item));
  return [included, excluded];
}

function* subtree(
moduleTransport,
moduleTransportsByPath)

{let seen = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : new Set();
  seen.add(moduleTransport.id);const
  meta = moduleTransport.meta;
  invariant(
  meta != null,
  'Unexpected module transport without meta information: ' +
  moduleTransport.sourcePath);

  for (const _ref7 of meta.dependencyPairs || []) {var _ref8 = _slicedToArray(_ref7, 2);const path = _ref8[1].path;
    const dependency = moduleTransportsByPath.get(path);
    if (dependency && !seen.has(dependency.id)) {
      yield dependency.id;
      yield* subtree(dependency, moduleTransportsByPath, seen);
    }
  }
}

const isRawMappings = Array.isArray;

module.exports = Bundle;