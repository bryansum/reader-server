var walk = require('./walk'),
    fileStats = require('./file-stats');

module.exports = function(dir) {
  var stats = {},
      tree = [];

  function filter(name, stat) {
    return stat.isDirectory() && name === '.git';
  }

  stats.tree = function() {
    tree = walk(dir, function(name, stat) {
      if (filter(name, stat)) return false;
      else return {name: name, size: stat.size};
    });
    return tree;
  };

  stats.file = function(name) {
    return fileStats(dir, name);
  }

  return stats;
};