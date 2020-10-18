---
order: 4
---

# Create posts with HTTP

Let's create a front-end for our blog application. In this guide we will be writing a client-side application in JavaScript that can create a wallet (public/private key pair), fetch a list of posts from our server, create posts and send to our server.

## `x/blog/client/rest/rest.go`

We'll be creating posts by sending `POST` requests to the same endpoint: `/blog/posts`. To add a handler add the following line to `func RegisterRoutes`:

```go
	r.HandleFunc("/blog/posts", createPostHandler(cliCtx)).Methods("POST")
```

Now let's create `createPostHandler` inside a new file you should create at `x/blog/client/rest/tx.go`.

## `x/blog/client/rest/tx.go`

We will need a `createPostReq` type that represents the request that we will be sending from the client:

<<< @/blog/blog/x/blog/client/rest/txPost.go{1-18}

`createPostHandler` first parses request parameters, performs basic validations, converts `Creator` field from a string into an SDK account address type then finally creates `MsgCreatePost` message.

<<< @/blog/blog/x/blog/client/rest/txPost.go{20-39}

## Setting up the client-side project

Inside a new empty directory create `package.json`:

```sh
npm init -y
```

Add dependencies:

```sh
npm add parcel-bundler local-cors-proxy axios @tendermint/sig
```

We'll be using [`parcel-bundler`](https://parceljs.org/) for bundling our dependencies and development server, [`local-cors-proxy`](https://github.com/garmeeh/local-cors-proxy) for providing a CORS proxy server as a development replacement for something like Nginx, [`axios`](https://github.com/axios/axios) for HTTP requests and [`@tendermint/sig`](https://github.com/tendermint/sig) for interacting with our application.

Replace `"scripts"` property in `package.json` with the following:

```json
"scripts": {
  "preserve": "lcp --proxyUrl http://localhost:1317 &",
  "serve": "parcel index.html",
  "postserve": "pkill -f lcp"
}
```

Responses from an HTTP server provided by `blogcli rest-server` do not provide `Access-Control-Allow-Origin` headers. This prevents client-side application running from different domains/ports from making successful requests due to most browsers' CORS policies. In production that can be managed by a web server like Nginx providing the right headers. In development we'll be using `local-cors-proxy` (`lcp`) to act as a proxy between our front-end and server.

The commands under `scripts` work as follows: `preserve` runs first and launches `lcp` in background. `serve` launches a development server and `postserve` kills `lcp` process after you shutdown the server.

## `index.html`

Our `index.html` file will contain only the minimal amount of markup:

<<< @/blog/blog/frontend/index.html

`<input/>` provides a text input for a title of a new post. `<button>` will submit a new post to our server. `<div id="post>` is an empty container which will be populated by a list of posts. We load `script.js` after all other markup has been loaded.

Now let's create `script.js`.

## `script.js`

Our app will depend on `axios` and `@tendermint/sig`:

<<< @/blog/blog/frontend/script.js{1,6}

<<< @/blog/blog/frontend/script.js{8}
<!-- 
```js
const API = "http://localhost:8010/proxy";
``` -->

HTTP-server runs on `localhost:1317`. Due to CORS we need to route our requests through a proxy: we're using `local-cors-proxy`, which uses the above URL by default.

First, we need to create an account: a public/private key pair, which we will use for signing our (create post) transactions before broadcasting them to the server. We'll be using a mnemonic, from which the keys will be generated.

<<< @/blog/blog/frontend/script.js{10}
<!-- 
```js
const mnemonic =
  "solid play vibrant paper clinic talent people employ august camp april reduce";
``` -->

Now let's generate a wallet:

<<< @/blog/blog/frontend/script.js{12}

<!-- ```js
const wallet = createWalletFromMnemonic(mnemonic);
``` -->

If you use the same mnemonic (think of it as a password), you will get the same wallet address: in out case `wallet.address` is `cosmos152gzu3vzf7g9tu46vszgpac24lwr48vc8k8kkh`. If you want to generate unique mnemonics for your users, you can use `bip39` [package](https://www.npmjs.com/package/bip39).

Our app will do two things: it will fetch a list of posts and create posts when "Create post" button is clicked.

<<< @/blog/blog/frontend/script.js{14-17}

<!-- ```js
const init = () => {
  fetchPosts();
  document.getElementById("button").addEventListener("click", createPost);
};
``` -->

For `init()` to work we need to define two functions: `fetchPosts` and `createPost`.

<<< @/blog/blog/frontend/script.js{19-24}

<!-- ```js
const fetchPosts = () => {
  axios.get(`${API}/blog/posts`).then(({ data }) => {
    const posts = JSON.stringify(data.result);
    document.getElementById("posts").innerText = posts;
  });
};
``` -->

`fetchPosts` makes an HTTP GET request to `/blog/posts` and inserts the posts into the `<div id="posts">` container.

Next, we need to define the `createPost` function. It will take the value from the text input, fetch the required account parameters from the server, fetch an unsigned transaction from the server, sign it using our private key from the wallet, broadcast it back to the server and fetch the list of posts again.

<<< @/blog/blog/frontend/script.js{26-71}

<!-- 
```js
const createPost = () => {
  // Getting a post title value from text input
  const title = document.getElementById("input").value;
  // Fetching account parameters: sequence and account_number
  axios.get(`${API}/auth/accounts/${wallet.address}`).then(({ data }) => {
    const account = data.result.value;
    const chain_id = "blog";
    const meta = {
      // Making sure both sequence and account_number are strings
      // Sequence number changes every time we submit a new transaction
      sequence: `${account.sequence}`,
      // Account number stays the same
      account_number: `${account.account_number}`,
      chain_id,
    };
    const req = {
      base_req: {
        chain_id,
        from: wallet.address,
      },
      creator: wallet.address,
      title,
    };
    // Fetching am unsigned transaction
    axios.post(`${API}/blog/posts`, req).then(({ data }) => {
      const tx = data.value;
      // Signing the transaction with the private key and meta info
      const stdTx = signTx(tx, meta, wallet);
      // Preparing transaction for broadcasting
      // "block" will make sure we get the response after
      // the transaction is included in the block
      const txBroadcast = createBroadcastTx(stdTx, "block");
      const params = {
        headers: {
          "Content-Type": "application/json",
        },
      };
      // Sending our post to be processed by the server
      axios.post(`${API}/txs`, txBroadcast, params).then(() => {
        // Fetch a new list of posts after we've successfully
        // created a new post
        fetchPosts();
      });
    });
  });
};
``` -->

Run init to initialize our app:

<<< @/blog/blog/frontend/script.js{73}

<!-- ```js
init();
``` -->

## Creating an account

Before we start using our application we need to make sure that the account we have generated exists on our chain. To do so we will send a nominal amount of tokens from an existing account to a new one.

```sh
blogcli keys show user1
```

You will get information about one of the existing accounts (your values will be different):

```json
{
  "name": "user1",
  "type": "local",
  "address": "cosmos1wt47yve6l29yjtxtsajhltr2vqhf7mpw5n6fx6",
  "pubkey": "cosmospub1addwnpepq03v7d6q4yt4nalj74elq8l5498wd9krcx92mxudkarj8aapy0qjvfaga8z"
}
```

Transfer some tokens from this account to the new one:

```sh
blogcli tx send $(blogcli keys show user1 -a) cosmos152gzu3vzf7g9tu46vszgpac24lwr48vc8k8kkh 10token --from=user1
```

Notice that the sender address can be queried automatically using the sub-command `$(blogcli keys show user1 -a)` with the flag `-a` to show just the address and the receiver account address `cosmos152gzu3vzf7g9tu46vszgpac24lwr48vc8k8kkh` is the one we have generated from the mnemonic in the browser.

In this guide we're activating accounts manually, but in production apps you might want to do it as part of a signing up process.

Now that we've created the app and set up our account, let's run it!

```
npm run serve
```

Open `http://localhost:1234/` and try creating new posts. New posts should show up in the list after 2-4 seconds after submission. Notice that posts persist after you refresh the page.

Congratulations! You have successfully created your first client-side app that interacts with our custom blockchain.
