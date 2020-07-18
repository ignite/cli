const path = require('path');
const home = require('home-path');

exports.starportDir = path.join(home(), ".starport");
exports.starportPath = path.join(exports.starportDir, "starport");
