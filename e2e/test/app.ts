import test from "ava";
import * as tmp from "tmp";
import * as path from "path";
import * as fs from "fs";
import * as util from "util";
const exec = util.promisify(require('child_process').exec);

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
