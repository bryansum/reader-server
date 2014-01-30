var ejs = require('ejs'),
    path = require('path'),
    fs = require('fs');

var tmpl = {};

tmpl.view = function(path, opts) {
  return ejs.render(fs.readFileSync(path, 'utf8'), opts);
}

tmpl.layoutPath = function(name, ctx) {
  return path.join(ctx.appDir, name + '.ejs');
}

tmpl.viewPath = function(name, ctx) {
  return path.join(ctx.viewDir, name + '.ejs');
}

module.exports = tmpl;
