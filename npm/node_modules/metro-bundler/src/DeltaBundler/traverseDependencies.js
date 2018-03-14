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














/**
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * Dependency Traversal logic for the Delta Bundler. This method calculates
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * the modules that should be included in the bundle by traversing the
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * dependency graph.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * Instead of traversing the whole graph each time, it just calculates the
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * difference between runs by only traversing the added/removed dependencies.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * To do so, it uses the passed `edges` paramater, which is a data structure
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * that contains the whole status of the dependency graph. During the
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * recalculation of the dependencies, it mutates the edges graph.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       *
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * The paths parameter contains the absolute paths of the root files that the
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * method should traverse. Normally, these paths should be the modified files
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       * since the last traversal.
                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       */let traverseDependencies = (() => {var _ref = _asyncToGenerator(
  function* (
  paths,
  dependencyGraph,
  transformOptions,
  edges)

  {let onProgress = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : function () {};
    const changes = yield Promise.all(
    paths.map(function (path) {return (
        traverseDependenciesForSingleFile(
        path,
        dependencyGraph,
        transformOptions,
        edges,
        onProgress));}));




    const added = new Set();
    const deleted = new Set();

    for (const change of changes) {
      for (const path of change.added) {
        added.add(path);
      }
      for (const path of change.deleted) {
        deleted.add(path);
      }
    }

    return {
      added,
      deleted };

  });return function traverseDependencies(_x, _x2, _x3, _x4) {return _ref.apply(this, arguments);};})();let initialTraverseDependencies = (() => {var _ref2 = _asyncToGenerator(

  function* (
  path,
  dependencyGraph,
  transformOptions,
  edges)

  {let onProgress = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : function () {};
    createEdge(path, edges);

    return yield traverseDependenciesForSingleFile(
    path,
    dependencyGraph,
    transformOptions,
    edges,
    onProgress);

  });return function initialTraverseDependencies(_x6, _x7, _x8, _x9) {return _ref2.apply(this, arguments);};})();let traverseDependenciesForSingleFile = (() => {var _ref3 = _asyncToGenerator(

  function* (
  path,
  dependencyGraph,
  transformOptions,
  edges)

  {let onProgress = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : function () {};
    const edge = edges.get(path);

    // If the passed edge does not exist does not exist in the graph, ignore it.
    if (!edge) {
      return { added: new Set(), deleted: new Set() };
    }

    // Get the absolute path of all sub-dependencies (some of them could have been
    // moved but maintain the same relative path).
    const currentDependencies = resolveDependencies(
    path, (
    yield dependencyGraph.getShallowDependencies(path, transformOptions)),
    dependencyGraph,
    transformOptions);


    const previousDependencies = new Set(edge.dependencies.values());

    const nonNullEdge = edge;

    let numProcessed = 0;
    let total = currentDependencies.size;

    const deleted = Array.from(edge.dependencies.entries()).
    map(function (_ref4) {var _ref5 = _slicedToArray(_ref4, 2);let relativePath = _ref5[0],absolutePath = _ref5[1];
      if (!currentDependencies.has(absolutePath)) {
        return removeDependency(nonNullEdge, relativePath, edges);
      } else {
        return undefined;
      }
    }).
    filter(Boolean);

    // Check all the module dependencies and start traversing the tree from each
    // added and removed dependency, to get all the modules that have to be added
    // and removed from the dependency graph.
    const added = yield Promise.all(
    Array.from(
    currentDependencies).
    map((() => {var _ref6 = _asyncToGenerator(function* (_ref7) {var _ref8 = _slicedToArray(_ref7, 2);let absolutePath = _ref8[0],relativePath = _ref8[1];
        let newDependencies;

        if (!previousDependencies.has(absolutePath)) {
          newDependencies = yield addDependency(
          nonNullEdge,
          relativePath,
          dependencyGraph,
          transformOptions,
          edges,
          function () {
            total++;
            onProgress(numProcessed, total);
          },
          function () {
            numProcessed++;
            onProgress(numProcessed, total);
          });

        } else {
          newDependencies = new Set();
        }

        return newDependencies;
      });return function (_x16) {return _ref6.apply(this, arguments);};})()));


    return {
      added: flatten(added),
      deleted: flatten(deleted) };

  });return function traverseDependenciesForSingleFile(_x11, _x12, _x13, _x14) {return _ref3.apply(this, arguments);};})();let addDependency = (() => {var _ref9 = _asyncToGenerator(

  function* (
  parentEdge,
  relativePath,
  dependencyGraph,
  transformOptions,
  edges,
  onDependencyAdd,
  onDependencyAdded)
  {
    const parentModule = dependencyGraph.getModuleForPath(parentEdge.path);
    const module = dependencyGraph.resolveDependency(
    parentModule,
    relativePath,
    transformOptions.platform);


    // Update the parent edge to keep track of the new dependency.
    parentEdge.dependencies.set(relativePath, module.path);

    let dependencyEdge = edges.get(module.path);

    // The new dependency was already in the graph, we don't need to do anything.
    if (dependencyEdge) {
      dependencyEdge.inverseDependencies.add(parentEdge.path);

      return new Set();
    }

    onDependencyAdd();

    // Create the new edge and traverse all its subdependencies, looking for new
    // subdependencies recursively.
    dependencyEdge = createEdge(module.path, edges);
    dependencyEdge.inverseDependencies.add(parentEdge.path);

    const addedDependencies = new Set([dependencyEdge.path]);

    const shallowDeps = yield dependencyGraph.getShallowDependencies(
    dependencyEdge.path,
    transformOptions);


    const nonNullDependencyEdge = dependencyEdge;

    const added = yield Promise.all(
    shallowDeps.map(function (dep) {return (
        addDependency(
        nonNullDependencyEdge,
        dep,
        dependencyGraph,
        transformOptions,
        edges,
        onDependencyAdd,
        onDependencyAdded));}));




    for (const newDependency of flatten(added)) {
      addedDependencies.add(newDependency);
    }

    onDependencyAdded();

    return addedDependencies;
  });return function addDependency(_x17, _x18, _x19, _x20, _x21, _x22, _x23) {return _ref9.apply(this, arguments);};})();function _asyncToGenerator(fn) {return function () {var gen = fn.apply(this, arguments);return new Promise(function (resolve, reject) {function step(key, arg) {try {var info = gen[key](arg);var value = info.value;} catch (error) {reject(error);return;}if (info.done) {resolve(value);} else {return Promise.resolve(value).then(function (value) {step("next", value);}, function (err) {step("throw", err);});}}return step("next");});};}

function removeDependency(
parentEdge,
relativePath,
edges)
{
  // Find the actual edge that represents the removed dependency. We do this
  // from the egdes data structure, since the file may have been deleted
  // already.
  const dependencyEdge = resolveEdge(parentEdge, relativePath, edges);
  if (!dependencyEdge) {
    return new Set();
  }

  parentEdge.dependencies.delete(relativePath);
  dependencyEdge.inverseDependencies.delete(parentEdge.path);

  // This module is still used by another modules, so we cannot remove it from
  // the bundle.
  if (dependencyEdge.inverseDependencies.size) {
    return new Set();
  }

  const removedDependencies = new Set([dependencyEdge.path]);

  // Now we need to iterate through the module dependencies in order to
  // clean up everything (we cannot read the module because it may have
  // been deleted).
  for (const subDependency of dependencyEdge.dependencies.keys()) {
    const removed = removeDependency(dependencyEdge, subDependency, edges);

    for (const removedDependency of removed.values()) {
      removedDependencies.add(removedDependency);
    }
  }

  // This module is not used anywhere else!! we can clear it from the bundle
  destroyEdge(dependencyEdge, edges);

  return removedDependencies;
}

function createEdge(path, edges) {
  const edge = {
    dependencies: new Map(),
    inverseDependencies: new Set(),
    path };

  edges.set(path, edge);

  return edge;
}

function destroyEdge(edge, edges) {
  edges.delete(edge.path);
}

function resolveEdge(
parentEdge,
relativePath,
edges)
{
  const absolutePath = parentEdge.dependencies.get(relativePath);
  if (!absolutePath) {
    return null;
  }

  return edges.get(absolutePath);
}

function resolveDependencies(
parentPath,
dependencies,
dependencyGraph,
transformOptions)
{
  const parentModule = dependencyGraph.getModuleForPath(parentPath);

  return new Map(
  dependencies.map(relativePath => [
  dependencyGraph.resolveDependency(
  parentModule,
  relativePath,
  transformOptions.platform).
  path,
  relativePath]));


}

function flatten(input) {
  const output = new Set();

  for (const items of input) {
    for (const item of items) {
      output.add(item);
    }
  }

  return output;
}

module.exports = {
  initialTraverseDependencies,
  traverseDependencies };