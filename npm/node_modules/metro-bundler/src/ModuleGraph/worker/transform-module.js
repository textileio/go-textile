/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * 
 * @format
 */

'use strict';var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};

const AssetPaths = require('../../node-haste/lib/AssetPaths');
const JsFileWrapping = require('./JsFileWrapping');
const Platforms = require('./Platforms');

const collectDependencies = require('./collect-dependencies');
const defaults = require('../../defaults');
const docblock = require('jest-docblock');
const generate = require('./generate');
const getImageSize = require('image-size');
const invariant = require('fbjs/lib/invariant');
const path = require('path');var _require =

require('../../Bundler/util');const isAssetTypeAnImage = _require.isAssetTypeAnImage;var _require2 =
require('path');const basename = _require2.basename;




















const NODE_MODULES = path.sep + 'node_modules' + path.sep;
const defaultTransformOptions = {
  dev: false,
  generateSourceMaps: true,
  hot: false,
  inlineRequires: false,
  platform: '',
  projectRoot: '' };

const defaultVariants = { default: {} };

const ASSET_EXTENSIONS = new Set(defaults.assetExts);

function transformModule(
content,
options)
{
  if (ASSET_EXTENSIONS.has(path.extname(options.filename).substr(1))) {
    return transformAsset(content, options);
  }

  const code = content.toString('utf8');
  if (options.filename.endsWith('.json')) {
    return transformJSON(code, options);
  }const

  filename = options.filename,transformer = options.transformer,polyfill = options.polyfill;var _options$variants = options.variants;const variants = _options$variants === undefined ? defaultVariants : _options$variants;
  const transformed = {};

  for (const variantName of Object.keys(variants)) {var _transformer$transfor =
    transformer.transform({
      filename,
      localPath: filename,
      options: _extends({}, defaultTransformOptions, variants[variantName]),
      src: code });const ast = _transformer$transfor.ast;

    invariant(ast != null, 'ast required from the transform results');
    transformed[variantName] = makeResult(ast, filename, code, polyfill);
  }

  let hasteID = null;
  if (filename.indexOf(NODE_MODULES) === -1) {
    hasteID = docblock.parse(docblock.extract(code)).providesModule;
    if (options.hasteImpl) {
      if (options.hasteImpl.enforceHasteNameMatches) {
        options.hasteImpl.enforceHasteNameMatches(filename, hasteID);
      }
      hasteID = options.hasteImpl.getHasteName(filename);
    }
  }

  return {
    details: {
      assetContent: null,
      code,
      file: filename,
      hasteID: hasteID || null,
      transformed,
      type: options.polyfill ? 'script' : 'module' },

    type: 'code' };

}

function transformJSON(json, options) {
  const value = JSON.parse(json);const
  filename = options.filename;
  const code = `__d(function(${JsFileWrapping.MODULE_FACTORY_PARAMETERS.join(
  ', ')
  }) { module.exports = \n${json}\n});`;

  const moduleData = {
    code,
    map: null, // no source map for JSON files!
    dependencies: [] };

  const transformed = {};

  Object.keys(options.variants || defaultVariants).forEach(
  key => transformed[key] = moduleData);


  const result = {
    assetContent: null,
    code: json,
    file: filename,
    hasteID: value.name,
    transformed,
    type: 'module' };


  if (basename(filename) === 'package.json') {
    result.package = {
      name: value.name,
      main: value.main,
      browser: value.browser,
      'react-native': value['react-native'] };

  }
  return { type: 'code', details: result };
}

function transformAsset(
content,
options)
{
  const filePath = options.filename;
  const assetData = AssetPaths.parse(filePath, Platforms.VALID_PLATFORMS);
  const contentType = path.extname(filePath).slice(1);
  const details = {
    assetPath: assetData.assetName,
    contentBase64: content.toString('base64'),
    contentType,
    filePath,
    physicalSize: getAssetSize(contentType, content),
    platform: assetData.platform,
    scale: assetData.resolution };

  return { details, type: 'asset' };
}

function getAssetSize(type, content) {
  if (!isAssetTypeAnImage(type)) {
    return null;
  }var _getImageSize =
  getImageSize(content);const width = _getImageSize.width,height = _getImageSize.height;
  return { width, height };
}

function makeResult(ast, filename, sourceCode) {let isPolyfill = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : false;
  let dependencies, dependencyMapName, file;
  if (isPolyfill) {
    dependencies = [];
    file = JsFileWrapping.wrapPolyfill(ast);
  } else {var _collectDependencies =
    collectDependencies(ast);dependencies = _collectDependencies.dependencies;dependencyMapName = _collectDependencies.dependencyMapName;
    file = JsFileWrapping.wrapModule(ast, dependencyMapName);
  }

  const gen = generate(file, filename, sourceCode, false);
  return { code: gen.code, map: gen.map, dependencies, dependencyMapName };
}

module.exports = transformModule;