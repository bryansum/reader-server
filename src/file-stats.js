var fs = require('fs'),
    path = require('path');

module.exports = function(root, file, stat) {
  var fpath = path.join(root, file),
      code = fs.readFileSync(fpath, 'utf8'),
      stats = { code: code };

  if (!stat) stat = fs.statSync(fpath);

  stats.dir = path.dirname(file);
  stats.name = path.basename(file);
  stats.bytes = stat.size;

  stats.loc = (function() {
    var count = 0;
    for (var i = 0; i < code.length; i++) {
      if (code[i] === '\n') count++;
    }
    return count;
  })();

  // http://c2.com/doc/SignatureSurvey/
  stats.signature = code.replace(/[^{};"]/g, '');

  return stats;
}