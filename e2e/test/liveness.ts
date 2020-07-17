// liveness test checks if the developer and app pages are reachable and
// can unlock wallet on the app and can see the available tokens from the browser.

import test from 'ava';
import * as puppeteer from 'puppeteer';
import * as tmp from 'tmp';
import * as util from 'util';
import { ChildProcess, spawn } from 'child_process';
const exec = util.promisify(require('child_process').exec);

const devAddr = 'http://localhost:12345';
const appAddr = 'http://localhost:8080';

const timeout = 1000 * 15;

test('dev page and wallet unlock', async t => {
  const served = await serve(t);
  const browser = await puppeteer.launch({dumpio: true, args: ['--no-sandbox']});

  await sleep(15000)

  try {
    // check appearance of activated go to app card on developer page.
    let page = await browser.newPage();
    await page.goto(devAddr);
    await page.waitForSelector(`a.card[href="${appAddr}"]`, { timeout })
    await page.close();

    // check tokens in the wallet on app page.
    page = await browser.newPage();
page.on('console', consoleObj => console.log(consoleObj.text()));
    await page.goto(appAddr);
    console.log(served.mnemonic)
    await page.type('.password__input', served.mnemonic);
    await page.click('.button');
    setInterval(() => {
page.screenshot({path: 'example.png'});
    }, 3333)
    await page.waitForSelector('.coin__amount', { timeout })
    const tokenText = await page.evaluate(() => (<HTMLElement>document.querySelector('.coin__amount')).innerText);
    t.is(tokenText, "500 TOKEN");
  } catch ({ message }) {
    t.fail(message);
  } finally {
    browser.close();
    served.proc.kill();
  }
});

interface Serve {
  proc: ChildProcess
  mnemonic: string
}

async function serve(t): Promise<Serve> {
  const workdir = tmp.dirSync();
  const o = { cwd: workdir.name };
  // create the app and serve.
  await exec('starport app github.com/a/b', o);
  const serveproc = spawn('starport', ['serve', '-p', 'b/'], o);
  serveproc.on("exit", (code) => {
    if (code !== null && code !== 0) {
      t.fail("serve command failed");
    }
  });
  // get mnemonic from the serve logs and wait developer
  // server to start.
  let mnemonic;
  for await (const log of serveproc.stdout) {
    const line = log.toString();
    if (line.includes("mnemonic")) {
      mnemonic = line.split(":")[1].trim();
    }
    if (line.includes(devAddr)) {
      break;
    }
  };
  if (!mnemonic) {
    throw new Error("cannot get mnemonic from serve logs");
  }
  return {
    proc: serveproc,
    mnemonic: mnemonic,
  };
}
function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
