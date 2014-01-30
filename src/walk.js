var fs = require('fs'),
    path = require('path');

function walk(root, dir, cb) {
  var results = [];
  var list = fs.readdirSync(path.join(root, dir));
  list.forEach(function(file) {
      file = path.join(dir, file);
      var stat = fs.statSync(path.join(root, file)),
          keep = cb(file, stat);
      if (keep !== false) {
        if (stat && stat.isDirectory()) results = results.concat(walk(root, file, cb));
        else results.push(keep);
      }
  });
  return results;
}

/**
  * Walks a directory.
  * @param {string} dir - Root directory to walk. All file results will be relative to this.
  * @param {function(fname, fs.Stats)} cb - Callback with file name and stats.
  */
module.exports = function(dir, cb) {
  return walk(dir, '.', cb)
}
