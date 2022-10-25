const pkgjson=require('./package.json');
var exec = require('child_process').exec;

for (let pkg in pkgjson.dependencies) {
  if (pkgjson.dependencies[pkg].startsWith('file:')) {
    exec(`cd ./node_modules/${pkg} && npm install`);
	}
}