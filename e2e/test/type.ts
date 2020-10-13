import test from "ava"
import * as tmp from "tmp"
import * as path from "path"
import * as util from "util"

const exec = util.promisify(require('child_process').exec)

test('generate an app with a type and verify', async t => {
  const workdir = tmp.dirSync()
  const appPath = path.join(workdir.name, "e2e")

  // generate the app. 
  await exec('starport app github.com/e2e/e2e', { cwd: workdir.name })

  // add a type. 
  await exec('starport type user name email', { cwd: appPath })

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

  // cannot create an existing type
  try {
    await exec('starport type user name email', { cwd: appPath })
    t.fail("should prevent creating existing types")
  } catch(e) {
  }
})

test('generate an app using stargate version with a type and verify', async t => {
  const workdir = tmp.dirSync()
  const appPath = path.join(workdir.name, "e2e")

  // generate the app. 
  await exec('starport app --sdk-version stargate github.com/e2e/e2e', { cwd: workdir.name })

  // add a type. 
  await exec('starport type user name email', { cwd: appPath })

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

  // cannot create an existing type
  try {
    await exec('starport type user name email', { cwd: appPath })
    t.fail("should prevent creating existing types")
  } catch(e) {
  }
})

test('creating types in a custom module', async t => {
  const workdir = tmp.dirSync();
  const appPath = path.join(workdir.name, "e2e")

  // generate the app. 
  await exec('starport app github.com/e2e/e2e', { cwd: workdir.name })

  // create a module
  await exec('starport module create example', { cwd: appPath })

  // add a type. 
  await exec('starport type user name email --module example', { cwd: appPath })

  // can create the same type in the app's module
  await exec('starport type user name email', { cwd: appPath })

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

  // cannot create a type in a non existing module 
  try {
    await exec('starport type user name email --module idontexist', { cwd: appPath })
    t.fail("cannot cannot create type in a non existant module")
  } catch(e) {
  }

  // cannot create an existing type
  try {
    await exec('starport type user name email --module example', { cwd: appPath })
    t.fail("should prevent creating existing types")
  } catch(e) {
  }
})

test('creating types in a custom module using stargate version ', async t => {
  const workdir = tmp.dirSync();
  const appPath = path.join(workdir.name, "e2e")

  // generate the app. 
  await exec('starport app --sdk-version stargate github.com/e2e/e2e', { cwd: workdir.name })

  // create a module
  await exec('starport module create example', { cwd: appPath })

  // add a type. 
  await exec('starport type user name email --module example', { cwd: appPath })

  // can create the same type in the app's module
  await exec('starport type user name email', { cwd: appPath })

  // check if app builds and tests passes.
  const gotest = await exec('go test ./...', { cwd: appPath })
  console.log(gotest.stdout.toString())
  t.falsy(gotest.error,  "app cannot be build or tests aren't passing")

  // cannot create a type in a non existing module 
  try {
    await exec('starport type user name email --module idontexist', { cwd: appPath })
    t.fail("cannot cannot create type in a non existant module")
  } catch(e) {
  }

  // cannot create an existing type
  try {
    await exec('starport type user name email --module example', { cwd: appPath })
    t.fail("should prevent creating existing types")
  } catch(e) {
  }
})