<template>
  <div>
    <div class="container">
      <div class="narrow">
        <h1>Start<br />developing</h1>
        <p>
          Well done! You have successfully created and started a Cosmos
          application. Your app is built on top of Tendermint Core engine and
          uses Cosmos SDK.
        </p>
        <h2>Next steps</h2>
        <p>It's time to add features to it!</p>
        <p>
          Keep the terminal window open, open a new one and run the following
          command to create a new type:
        </p>
      </div>
      <div class="window">
        ~$: starport type user name email
      </div>
      <div class="narrow">
        <p>
          This creates a <code>User</code> type. Now you can create and view
          users in your application.
        </p>
      </div>
      <div class="window">
        ~$: {{ chain_id }}cli tx {{ chain_id }} create-user Alice alice@example.org
        --from=me
      </div>
      <div class="narrow">
        <h2>Block viewer</h2>
        <terminal-window />
        <h2>Components</h2>
        <a
          href="http://localhost:26657"
          target="_blank"
          class="card"
          style="background-color: #edf8eb; --color-primary: rgb(0, 187, 0)"
        >
          <logo-tendermint class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">Tendermint</div>
            <div class="card__desc__p">
              Tendermint RPC is running locally on port
              <span class="card__desc__tag">26657</span>
            </div>
          </div>
        </a>
        <a
          href="http://localhost:1317"
          target="_blank"
          class="card"
          style="background-color: #edefff; --color-primary: rgb(80, 100, 251)"
        >
          <logo-cosmos-sdk class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">Cosmos</div>
            <div class="card__desc__p">
              Cosmos REST API is running locally on port
              <span class="card__desc__tag">1317</span>
            </div>
          </div>
        </a>
        <h2 id="front-end">Front-end</h2>
        <p>Starport generated a simple web application inside <code>./ui</code> directory that can interact with your blockchain. In a new terminal window inside your project's directory run the following command:</p>
      </div>
      <div class="window">
        ~$: cd ui && npm i && npm run serve
      </div>
      <div class="narrow">
        <p>When the build is finished, open <a href="localhost:1234" target="_blank">a new browser tab</a> to see the application.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
a {
  color: inherit;
}
code {
  font-family: monospace;
}
.container {
  padding: 4rem 1rem 12rem;
  max-width: 900px;
  width: 100%;
  margin-left: auto;
  margin-right: auto;
  box-sizing: border-box;
}
.narrow {
  padding-left: 40px;
  padding-right: 40px;
  box-sizing: border-box;
}
h1 {
  font-weight: 800;
  font-size: 5rem;
  line-height: 0.9;
}
h2 {
  margin-top: 4rem;
  margin-bottom: 1rem;
  font-size: 1.25rem;
  font-weight: 900;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
p {
  font-size: 1.15rem;
  font-weight: 400;
  color: rgba(0, 0, 0, 0.65);
  line-height: 1.4;
  letter-spacing: 0.01em;
}
code {
  font-family: "Monaco", monospace;
  font-size: 1rem;
  font-weight: bold;
}
.window {
  font-family: "Monaco";
  box-shadow: 0 20px 50px 10px rgba(0, 0, 0, 0.05);
  padding: 3rem 2rem;
  margin-top: 3rem;
  margin-bottom: 3rem;
  border-radius: 0.75rem;
}
.card {
  padding: 1.5rem 2.5rem;
  border-radius: 0.75rem;
  display: flex;
  margin: 2rem 0;
  text-decoration: none;
  transition: transform 0.2s;
  color: inherit;
  align-items: center;
  color: var(--color-primary);
  letter-spacing: 0.01em;
}
.card:hover {
  transform: scale(1.01);
}
.card__logo {
  flex-shrink: 0;
  height: 4.5rem;
  width: 4.5rem;
  fill: var(--color-primary);
}
.card__desc {
  padding: 0 1rem;
}
.card__desc__h1 {
  font-size: 1.5rem;
  font-weight: 500;
  margin-bottom: 0.15rem;
}
.card__desc__p {
  opacity: 0.5;
  filter: brightness(85%);
}
.card__desc__tag {
  font-weight: 500;
  font-size: 0.9em;
  letter-spacing: 0.03em;
}
@media screen and (max-width: 980px) {
  h1 {
    font-size: 3rem;
  }
  .narrow {
    padding: 0;
  }
  .card {
    padding: 1.5rem 1.5rem;
  }
}
</style>

<script>
import axios from "axios";
import LogoTendermint from "@/components/LogoTendermint.vue";
import LogoCosmosSdk from "@/components/LogoCosmosSdk.vue";
import TerminalWindow from "@/components/TerminalWindow.vue";

export default {
  components: {
    LogoTendermint,
    LogoCosmosSdk,
    TerminalWindow
  },
  data() {
    return {
      chain_id: "{{chain_id}}"
    };
  },
  async created() {
    try {
      this.chain_id = (await axios.get("/chain_id")).data;
    } catch {
      console.log("Can't fetch /chain_id");
    }
  }
};
</script>
