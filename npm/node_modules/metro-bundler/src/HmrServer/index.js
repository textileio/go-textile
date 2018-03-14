/**
 * Copyright (c) 2015-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * @format
 * 
 */

'use strict';var _slicedToArray = function () {function sliceIterator(arr, i) {var _arr = [];var _n = true;var _d = false;var _e = undefined;try {for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {_arr.push(_s.value);if (i && _arr.length === i) break;}} catch (err) {_d = true;_e = err;} finally {try {if (!_n && _i["return"]) _i["return"]();} finally {if (_d) throw _e;}}return _arr;}return function (arr, i) {if (Array.isArray(arr)) {return arr;} else if (Symbol.iterator in Object(arr)) {return sliceIterator(arr, i);} else {throw new TypeError("Invalid attempt to destructure non-iterable instance");}};}();var _extends = Object.assign || function (target) {for (var i = 1; i < arguments.length; i++) {var source = arguments[i];for (var key in source) {if (Object.prototype.hasOwnProperty.call(source, key)) {target[key] = source[key];}}}return target;};function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}

const formatBundlingError = require('../lib/formatBundlingError');
const getBundlingOptionsForHmr = require('./getBundlingOptionsForHmr');
const querystring = require('querystring');
const url = require('url');var _require =





require('../Logger');const createActionStartEntry = _require.createActionStartEntry,createActionEndEntry = _require.createActionEndEntry,log = _require.log;










/**
                                                                                                                                                              * The HmrServer (Hot Module Reloading) implements a lightweight interface
                                                                                                                                                              * to communicate easily to the logic in the React Native repository (which
                                                                                                                                                              * is the one that handles the Web Socket connections).
                                                                                                                                                              *
                                                                                                                                                              * This interface allows the HmrServer to hook its own logic to WS clients
                                                                                                                                                              * getting connected, disconnected or having errors (through the
                                                                                                                                                              * `onClientConnect`, `onClientDisconnect` and `onClientError` methods).
                                                                                                                                                              */
class HmrServer {



  constructor(packagerServer, reporter) {
    this._packagerServer = packagerServer;
    this._reporter = reporter;
  }

  onClientConnect(
  clientUrl,
  sendFn)
  {var _this = this;return _asyncToGenerator(function* () {var _querystring$parse =
      querystring.parse(
      /* $FlowFixMe: url might be null */
      url.parse(clientUrl).query);const bundleEntry = _querystring$parse.bundleEntry,platform = _querystring$parse.platform;


      // Create a new DeltaTransformer for each client. Once the clients are
      // modified to support Delta Bundles, they'll be able to pass the
      // DeltaBundleId param through the WS connection and we'll be able to share
      // the same DeltaTransformer between the WS connection and the HTTP one.
      const deltaBundler = _this._packagerServer.getDeltaBundler();var _ref =
      yield deltaBundler.getDeltaTransformer(
      getBundlingOptionsForHmr(bundleEntry, platform));const deltaTransformer = _ref.deltaTransformer;


      // Trigger an initial build to start up the DeltaTransformer.
      yield deltaTransformer.getDelta();

      // Listen to file changes.
      const client = { sendFn, deltaTransformer };
      deltaTransformer.on('change', _this._handleFileChange.bind(_this, client));

      return client;})();
  }

  onClientError(client, e) {
    this._reporter.update({
      type: 'hmr_client_error',
      error: e });

    this.onClientDisconnect(client);
  }

  onClientDisconnect(client) {
    // We can safely remove all listeners from the delta transformer since the
    // transformer is not shared between clients.
    client.deltaTransformer.removeAllListeners('change');
  }

  _handleFileChange(client) {var _this2 = this;return _asyncToGenerator(function* () {
      const processingHmrChange = log(
      createActionStartEntry({ action_name: 'Processing HMR change' }));


      client.sendFn(JSON.stringify({ type: 'update-start' }));
      const response = yield _this2._prepareResponse(client);

      client.sendFn(JSON.stringify(response));
      client.sendFn(JSON.stringify({ type: 'update-done' }));

      log(_extends({},
      createActionEndEntry(processingHmrChange), {
        outdated_modules: Array.isArray(response.body.modules) ?
        response.body.modules.length :
        null }));})();

  }

  _prepareResponse(client) {var _this3 = this;return _asyncToGenerator(function* () {
      let result;

      try {
        result = yield client.deltaTransformer.getDelta();
      } catch (error) {
        const formattedError = formatBundlingError(error);

        _this3._reporter.update({ type: 'bundling_error', error });

        return { type: 'error', body: formattedError };
      }
      const modules = [];

      for (const _ref2 of result.delta) {var _ref3 = _slicedToArray(_ref2, 2);const id = _ref3[0];const module = _ref3[1];
        // The Delta Bundle can have null objects: these correspond to deleted
        // modules, which we don't need to send to the client.
        if (module != null) {
          // When there are new modules added on the dependency tree, they are
          // appended on the Delta Bundle, but HMR needs to have them at the
          // beginning.
          modules.unshift({ id, code: module.code });
        }
      }

      return {
        type: 'update',
        body: {
          modules,
          inverseDependencies: result.inverseDependencies,
          sourceURLs: {},
          sourceMappingURLs: {} // TODO: handle Source Maps
        } };})();

  }}


module.exports = HmrServer;