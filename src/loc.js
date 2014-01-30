module.exports = function loc(code) {
  var count = 0
  for (var i = 0; i < code.length; i++) {
    if (code[i] === '\n') count++
  }
  return count
};
