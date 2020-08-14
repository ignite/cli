import test from "ava";
import * as tmp from "tmp";
import * as path from "path";
import * as fs from "fs";
import * as util from "util";
import { ChildProcess, spawn } from 'child_process';
import * as retry from'async-retry';
import * as fetch from 'node-fetch';
const exec = util.promisify(require('child_process').exec);

const cosmosAddr = "http://localhost:1317"

test('generate an app and verify', async t => {
  const workdir = tmp.dirSync();

  // check if app can be generated.
  await exec('starport app github.com/e2e/e2e', {
    cwd: workdir.name,
  });
  t.true(fs.existsSync(path.join(workdir.name, "e2e/config.yml")), "cannot locate a config.yml")

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', {
    cwd: path.join(workdir.name, "e2e"),
  })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")
});

test('generate an app with CosmWasm and verify', async t => {
  const workdir = tmp.dirSync();
  const appPath = path.join(workdir.name, "e2e");

  // generate the app. 
  await exec('starport app github.com/e2e/e2e', { cwd: workdir.name });

  // add CosmWasm. 
  await exec('starport add wasm', { cwd: appPath });

  // cannot add CosmWasm again. 
  try {
    await exec('starport add wasm', { cwd: appPath });
    t.fail("cannot add wasm twice")
  } catch(e) {
  }

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")
});

test('serve app with CosmWasm', async t => {
  const workdir = tmp.dirSync();
  const appPath = path.join(workdir.name, "e2e");

  // generate the app. 
  await exec('starport app github.com/e2e/e2e', { cwd: workdir.name });

  // add CosmWasm. 
  await exec('starport add wasm', { cwd: appPath });

  // serve should be OK and there should not be any runtime errors during serve.
  await new Promise(async (resolve, reject) => {
    let stderr = "";
    const serveproc = spawn("starport", ["serve", "--verbose"], { cwd: appPath });

    // collect logs from the error stream.
    serveproc.stderr.on('data', (data) => {
      stderr += data.toString();
    });

    // watch serve exit.
    serveproc.on("exit", (code) => {
      if (code > 0) {
        reject(`cannot serve:\n${stderr}`);
      }
    });

    // try to reach to App through API to see if it works properly. 
    await retry(async (bail, no) => {
      const res = await fetch(`${cosmosAddr}`)
      if (res.status != 404) {
        bail(new Error('app is not reachable'))
      }
    }, {
      retries: 15,
      minTimeout: 1000 * 5,
    });
    resolve();
  })

  t.pass("app is reachable");
});
