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

'use strict';

const meta = require('../../shared/output/meta');var _require =

require('./source-map');const createIndexMap = _require.createIndexMap;var _require2 =
require('./util');const getModuleCode = _require2.getModuleCode,concat = _require2.concat;



function asPlainBundle(_ref)





{let filename = _ref.filename,idsForPath = _ref.idsForPath,modules = _ref.modules,requireCalls = _ref.requireCalls,sourceMapPath = _ref.sourceMapPath;
  let code = '';
  let line = 0;
  const sections = [];
  const modIdForPath = x => idsForPath(x).moduleId;

  for (const module of concat(modules, requireCalls)) {const
    file = module.file;
    const moduleCode = getModuleCode(module, modIdForPath);

    code += moduleCode + '\n';
    if (file.map) {
      sections.push({
        map: file.map,
        offset: { column: 0, line } });

    }
    line += countLines(moduleCode);
  }

  if (sourceMapPath) {
    code += `//# sourceMappingURL=${sourceMapPath}`;
  }

  return {
    code,
    extraFiles: [[`${filename}.meta`, meta(code)]],
    map: createIndexMap({ file: filename, sections }) };

}

module.exports = asPlainBundle;

const reLine = /^/gm;
function countLines(string) {
  //$FlowFixMe This regular expression always matches
  return string.match(reLine).length;
}