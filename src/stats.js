var walk = require('./walk'),
    fileStats = require('./file-stats');

module.exports = function(dir, config) {
  var stats = {};

  function filter(name, stat) {
    return stat.isDirectory() && name === '.git';
  }

  stats.main = function() {
    return fileStats(dir, config.main);
  }

  stats.tree = function() {
    return walk(dir, function(name, stat) {
      if (filter(name, stat)) return false;
      else return {name: name, size: stat.size};
    });
  };

  stats.file = function(name) {
    return fileStats(dir, name);
  }

  return stats;
};
