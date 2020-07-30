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
        <h2>Components</h2>
        <a
          :href="running.frontend && 'http://localhost:8080'"
          target="_blank"
          class="card"
          v-if="running.frontend || env.node_js"
          :style="{'background-color': running.frontend ? 'rgb(255, 234, 250)' : 'rgba(0,0,0,.05)', '--color-primary': running.frontend ? 'rgb(251, 80, 210)' : 'rgba(0,0,0,0.25)'}"
        >
          <logo-spaceship class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">User interface</div>
            <div class="card__desc__p">
              <span v-if="running.frontend">The front-end of your app. A Vue application is running on port <span class="card__desc__tag">8080</span></span>
              <span v-else>Loading...</code></span>
            </div>
          </div>
        </a>
        <a
          target="_blank"
          class="card"
          v-if="!running.frontend && !env.node_js"
          :style="{'background-color': running.frontend ? 'rgb(255, 234, 250)' : 'rgba(0,0,0,.05)', '--color-primary': running.frontend ? 'rgb(251, 80, 210)' : 'rgba(0,0,0,0.25)'}"
        >
          <logo-spaceship class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">User interface</div>
            <div class="card__desc__p">
              <span>Can't start UI, because Node.js is not found.</span>
            </div>
          </div>
        </a>
        <a
          href="http://localhost:1317"
          target="_blank"
          class="card"
          :style="{'background-color': running.api ? '#edefff' : 'rgba(0,0,0,.05)', '--color-primary': running.api ? 'rgb(80, 100, 251)' : 'rgba(0,0,0,0.25)'}"
        >
          <logo-cosmos-sdk class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">Cosmos</div>
            <div class="card__desc__p">
              <span v-if="running.api">The back-end of your app. Cosmos API is running locally on port <span class="card__desc__tag">1317</span></span>
              <span v-else>Loading...</span>
            </div>
          </div>
        </a>
        <a
          href="http://localhost:26657"
          target="_blank"
          class="card"
          :style="{'background-color': running.rpc ? '#edf8eb' : 'rgba(0,0,0,.05)', '--color-primary': running.rpc ? 'rgb(0, 187, 0)' : 'rgba(0,0,0,0.25)'}"
        >
          <logo-tendermint class="card__logo" />
          <div class="card__desc">
            <div class="card__desc__h1">Tendermint</div>
            <div class="card__desc__p">
              <span v-if="running.rpc">The consensus engine. Tendermint API is running on port <span class="card__desc__tag">26657</span></span>
              <span v-else>Loading...</span>
            </div>
          </div>
        </a>
        <h2>Next steps</h2>
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
        ~$: {{ env.chain_id }}cli tx {{ env.chain_id }} create-user Alice alice@example.org
        --from=user1
      </div>
      <div class="narrow">
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
  transition: all 0.2s;
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
import LogoSpaceship from "@/components/LogoSpaceship.vue";
import TerminalWindow from "@/components/TerminalWindow.vue";

export default {
  components: {
    LogoTendermint,
    LogoCosmosSdk,
    TerminalWindow,
    LogoSpaceship
  },
  data() {
    return {
      env: {
        chain_id: "{{chain_id}}",
        node_js: false
      },
      running: {
        rpc: true,
        api: true,
        frontend: false
      },
      timer: null
    };
  },
  methods: {
    async setStatusState() {
      try {
        const { data } = await axios.get("/status");
        const { status, env } = data;
        this.running = {
            rpc: status.is_consensus_engine_alive,
            api: status.is_my_app_backend_alive,
            frontend: status.is_my_app_frontend_alive,
        };
        this.env = env;
      } catch {
        this.running = {
            rpc: false,
            api: false,
            frontend: false,
        };
      }
    }
  },
  async created() {
    this.timer = setInterval(this.setStatusState.bind(this), 1000);
    try {
      await this.setStatusState();
    } catch {
      console.log("Can't fetch /env");
    }
  },
  beforeDestroy() {
    clearInterval(this.timer);
  }
};
</script>
