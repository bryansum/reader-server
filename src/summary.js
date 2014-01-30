// http://c2.com/doc/SignatureSurvey/

module.exports = function summary(code) {
  return code.replace(/[^{};"]/g, '');
};
