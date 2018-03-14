/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 *
 * 
 */
'use strict';

const optimizeModule = require('./worker/optimize-module');
const transformModule = require('./worker/transform-module');
const wrapWorkerFn = require('./worker/wrap-worker-fn');





exports.optimizeModule =
wrapWorkerFn(optimizeModule);
exports.transformModule =
wrapWorkerFn(transformModule);