const tmp = require('tmp');
const { spawn } = require('child_process');
const { starportPath } = require('./conf');

async function main() {
  const args = process.argv;
  args.shift(); // node
  args.shift(); // starport
  spawn(starportPath, args, {
    stdio: [process.stdin, process.stdout, process.stderr],
  });
}

(async function () {
  try {
    await main();
  } catch (e) {
    console.log(e.message);
    process.exit(1);
  }
})();

