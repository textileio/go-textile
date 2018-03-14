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

const babel = require('babel-core');
const babelGenerate = require('babel-generator').default;
const babylon = require('babylon');









const assetPropertyBlacklist = new Set(['files', 'fileSystemLocation', 'path']);

function generateAssetCodeFileAst(
assetRegistryPath,
assetDescriptor)
{
  const properDescriptor = filterObject(
  assetDescriptor,
  assetPropertyBlacklist);

  const descriptorAst = babylon.parseExpression(
  JSON.stringify(properDescriptor));

  const t = babel.types;
  const moduleExports = t.memberExpression(
  t.identifier('module'),
  t.identifier('exports'));

  const requireCall = t.callExpression(t.identifier('require'), [
  t.stringLiteral(assetRegistryPath)]);

  const registerAssetFunction = t.memberExpression(
  requireCall,
  t.identifier('registerAsset'));

  const registerAssetCall = t.callExpression(registerAssetFunction, [
  descriptorAst]);

  return t.file(
  t.program([
  t.expressionStatement(
  t.assignmentExpression('=', moduleExports, registerAssetCall))]));



}

function generateAssetTransformResult(
assetRegistryPath,
assetDescriptor)




{var _babelGenerate =
  babelGenerate(
  generateAssetCodeFileAst(assetRegistryPath, assetDescriptor),
  { comments: false, compact: true });const code = _babelGenerate.code;

  const dependencies = [assetRegistryPath];
  const dependencyOffsets = [code.indexOf(assetRegistryPath) - 1];
  return { code, dependencies, dependencyOffsets };
}

// Test extension against all types supported by image-size module.
// If it's not one of these, we won't treat it as an image.
function isAssetTypeAnImage(type) {
  return (
    ['png', 'jpg', 'jpeg', 'bmp', 'gif', 'webp', 'psd', 'svg', 'tiff'].indexOf(
    type) !==
    -1);

}

function filterObject(object, blacklist) {
  const copied = Object.assign({}, object);
  for (const key of blacklist) {
    delete copied[key];
  }
  return copied;
}

function createRamBundleGroups(
ramGroups,
groupableModules,
subtree)
{
  // build two maps that allow to lookup module data
  // by path or (numeric) module id;
  const byPath = new Map();
  const byId = new Map();
  groupableModules.forEach(m => {
    byPath.set(m.sourcePath, m);
    byId.set(m.id, m.sourcePath);
  });

  // build a map of group root IDs to an array of module IDs in the group
  const result = new Map(
  ramGroups.map(modulePath => {
    const root = byPath.get(modulePath);
    if (root == null) {
      throw Error(`Group root ${modulePath} is not part of the bundle`);
    }
    return [
    root.id,
    // `subtree` yields the IDs of all transitive dependencies of a module
    new Set(subtree(root, byPath))];

  }));


  if (ramGroups.length > 1) {
    // build a map of all grouped module IDs to an array of group root IDs
    const all = new ArrayMap();
    for (const _ref of result) {var _ref2 = _slicedToArray(_ref, 2);const parent = _ref2[0];const children = _ref2[1];
      for (const module of children) {
        all.get(module).push(parent);
      }
    }

    // find all module IDs that are part of more than one group
    const doubles = filter(all, (_ref3) => {var _ref4 = _slicedToArray(_ref3, 2);let parents = _ref4[1];return parents.length > 1;});
    for (const _ref5 of doubles) {var _ref6 = _slicedToArray(_ref5, 2);const moduleId = _ref6[0];const parents = _ref6[1];
      const parentNames = parents.map(byId.get, byId);
      const lastName = parentNames.pop();
      throw new Error(
      `Module ${byId.get(moduleId) ||
      moduleId} belongs to groups ${parentNames.join(', ')}, and ${String(
      lastName)
      }. Ensure that each module is only part of one group.`);

    }
  }

  return result;
}

function* filter(iterator, predicate) {
  for (const value of iterator) {
    if (predicate(value)) {
      yield value;
    }
  }
}

/* $FlowFixMe(>=0.54.0 site=react_native_fb) This comment suppresses an error
   * found when Flow v0.54 was deployed. To see the error delete this comment and
   * run Flow. */
class ArrayMap extends Map {
  get(key) {
    let array = super.get(key);
    if (!array) {
      array = [];
      this.set(key, array);
    }
    return array;
  }}


module.exports = {
  createRamBundleGroups,
  generateAssetCodeFileAst,
  generateAssetTransformResult,
  isAssetTypeAnImage };