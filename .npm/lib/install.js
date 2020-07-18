const fs = require("fs");
const tmp = require("tmp");
const decompress = require("decompress");
const download = require("download");
const ora = require("ora");
const pkg = require("../package.json");
const { starportDir, starportPath } = require("./conf");

const baseUrl = "http://github.com/tendermint/starport/releases/download";
const urls = {
  linux(version) {
    return `${baseUrl}/v${version}/starport_${version}_linux_amd64.tar.gz`;
  },
  darwin(version) {
    return `${baseUrl}/v${version}/starport_${version}_darwin_amd64.tar.gz`;
  },
};

async function main() {
  const urlf = urls[process.platform];
  if (!urlf) {
    throw new Error(`unsupported platform '${process.platform}'`);
  }
  const url = urlf(pkg.version);
  const tmptar = tmp.fileSync();
  if (!fs.existsSync(starportDir)) {
    fs.mkdirSync(starportDir);
  }
  if (fs.existsSync(starportPath)) {
    fs.unlinkSync(starportPath);
  }
  const sp = ora("💫 Installing Starport...").start();
  sp.color = "yellow";
  try {
    fs.writeFileSync(tmptar.name, await download(url));
    await decompress(tmptar.name, starportDir);
  } catch ({ message }) {
    throw new Error(`cannot install starport err:\n\t${message}`);
  } finally {
    sp.stop();
    tmptar.removeCallback();
  }
}

(async function () {
  try {
    await main();
  } catch (e) {
    console.log(e.message);
    process.exit(1);
  }
})();
