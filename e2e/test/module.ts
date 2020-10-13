import test from "ava"
import * as tmp from "tmp"
import * as path from "path"
import * as util from "util"

const exec = util.promisify(require('child_process').exec)

test('generate an app with a empty module and verify', async t => {
  const workdir = tmp.dirSync()
  const appPath = path.join(workdir.name, "e2e")

  // generate the app. 
  await exec('starport app github.com/e2e/e2e', { cwd: workdir.name })

  // create a module
  await exec('starport module create example', { cwd: appPath })

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

  // cannot create a module that already exists 
  try {
    await exec('starport module create example', { cwd: appPath })
    t.fail("cannot cannot create an already existing module")
  } catch(e) {
  }

})

test('generate an app with stargate version with a empty module and verify', async t => {
    const workdir = tmp.dirSync()
    const appPath = path.join(workdir.name, "e2e")
  
    // generate the app. 
    await exec('starport app github.com/e2e/e2e --sdk-version stargate', { cwd: workdir.name })
  
    // create a module
    await exec('starport module create example', { cwd: appPath })
  
    // check if app builds and tests passes.
    const gotest = await exec('go test ./...', { cwd: appPath })
    console.log(gotest.stdout.toString())
    t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

    // cannot create a module that already exists 
    try {
      await exec('starport module create example', { cwd: appPath })
      t.fail("cannot cannot create an already existing module")
    } catch(e) {
    }
  })
  