/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
'use strict';var _require =

require('console');const Console = _require.Console;var _require2 =
require('stream');const Writable = _require2.Writable;

const write = (_, __, callback) => callback();
module.exports = new Console(new Writable({ write, writev: write }));