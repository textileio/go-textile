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














/**
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              * This module contains many serializers for the Delta Bundler. Each serializer
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              * returns a string representation for any specific type of bundle, which can
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              * be directly sent to the devices.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              */let deltaBundle = (() => {var _ref = _asyncToGenerator(

  function* (
  deltaBundler,
  options)
  {var _ref2 =
    yield _build(deltaBundler, _extends({},
    options, {
      wrapModules: true }));const id = _ref2.id,delta = _ref2.delta;


    function stringifyModule(_ref3) {var _ref4 = _slicedToArray(_ref3, 2);let id = _ref4[0],module = _ref4[1];
      return [id, module ? module.code : undefined];
    }

    const bundle = JSON.stringify({
      id,
      pre: Array.from(delta.pre).map(stringifyModule),
      post: Array.from(delta.post).map(stringifyModule),
      delta: Array.from(delta.delta).map(stringifyModule),
      reset: delta.reset });


    return {
      bundle,
      numModifiedFiles: delta.pre.size + delta.post.size + delta.delta.size };

  });return function deltaBundle(_x, _x2) {return _ref.apply(this, arguments);};})();let fullSourceMap = (() => {var _ref5 = _asyncToGenerator(

  function* (
  deltaBundler,
  options)
  {var _ref6 =
    yield _build(deltaBundler, _extends({},
    options, {
      wrapModules: true }));const id = _ref6.id,delta = _ref6.delta;


    const deltaPatcher = DeltaPatcher.get(id).applyDelta(delta);

    return fromRawMappings(deltaPatcher.getAllModules()).toString(undefined, {
      excludeSource: options.excludeSource });

  });return function fullSourceMap(_x3, _x4) {return _ref5.apply(this, arguments);};})();let fullSourceMapObject = (() => {var _ref7 = _asyncToGenerator(

  function* (
  deltaBundler,
  options)
  {var _ref8 =
    yield _build(deltaBundler, _extends({},
    options, {
      wrapModules: true }));const id = _ref8.id,delta = _ref8.delta;


    const deltaPatcher = DeltaPatcher.get(id).applyDelta(delta);

    return fromRawMappings(deltaPatcher.getAllModules()).toMap(undefined, {
      excludeSource: options.excludeSource });

  });return function fullSourceMapObject(_x5, _x6) {return _ref7.apply(this, arguments);};})();

/**
                                                                                                 * Returns the full JS bundle, which can be directly parsed by a JS interpreter
                                                                                                 */let fullBundle = (() => {var _ref9 = _asyncToGenerator(
  function* (
  deltaBundler,
  options)
  {var _ref10 =
    yield _build(deltaBundler, _extends({},
    options, {
      wrapModules: true }));const id = _ref10.id,delta = _ref10.delta;


    const deltaPatcher = DeltaPatcher.get(id).applyDelta(delta);
    const code = deltaPatcher.getAllModules().map(function (m) {return m.code;});

    return {
      bundle: code.join('\n'),
      lastModified: deltaPatcher.getLastModifiedDate(),
      numModifiedFiles: deltaPatcher.getLastNumModifiedFiles() };

  });return function fullBundle(_x7, _x8) {return _ref9.apply(this, arguments);};})();let getAllModules = (() => {var _ref11 = _asyncToGenerator(

  function* (
  deltaBundler,
  options)
  {var _ref12 =
    yield _build(deltaBundler, _extends({},
    options, {
      wrapModules: true }));const id = _ref12.id,delta = _ref12.delta;


    return DeltaPatcher.get(id).
    applyDelta(delta).
    getAllModules();
  });return function getAllModules(_x9, _x10) {return _ref11.apply(this, arguments);};})();let _build = (() => {var _ref13 = _asyncToGenerator(

  function* (
  deltaBundler,
  options)
  {var _ref14 =
    yield deltaBundler.getDeltaTransformer(
    options);const deltaTransformer = _ref14.deltaTransformer,id = _ref14.id;


    return {
      id,
      delta: yield deltaTransformer.getDelta() };

  });return function _build(_x11, _x12) {return _ref13.apply(this, arguments);};})();function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}const DeltaPatcher = require('./DeltaPatcher');var _require = require('../Bundler/source-map');const fromRawMappings = _require.fromRawMappings;

module.exports = {
  deltaBundle,
  fullBundle,
  fullSourceMap,
  fullSourceMapObject,
  getAllModules };