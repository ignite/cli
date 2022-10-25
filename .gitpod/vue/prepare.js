const pkgjson=require('./package.json');

for (let pkg in pkgjson.dependencies) {
  if (pkgjson.dependencies[pkg].startsWith('file:')) {
    console.error('\x1b[31m%s\x1b[0m', `Package '${pkg}' located at '${pkgjson.dependencies[pkg].replace('file:', '')}' needs to be published and your package.json file updated.`);
	}
}