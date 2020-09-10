<template>

  <div class="main">
    <div class="hero">
      <h4>Congratulations!</h4>
      <h2>Your Blockchain is live</h2>
    </div>
    <div class="texts-field">
      <p>You've just scaffolded and launched a Cosmos application.</p>
      <p>This application is a blockchain with its own tokens, accounts, governance and more. You can create your very own token and distribute to anyone you want. You can create governance proposals and vote on them. You can create custom data types ($EXAMPLES) and define operations on these types, as well as display this data on a web application.</p>

      <h3>01. Brief intro to architecture</h3>
      <p>This blockchain application is built with the Cosmos SDK. The SDK is a modular framework for building blockchains. It includes modules for auth, bank, staking, governance, and more. Every feature is packaged as a separate module that can interact with other modules. Starport has actually generated a module for you, which you can use to start developing your own features.</p>
      <a href="#">Learn more about Cosmos SDK</a>

      <h3>02. Configure your genesis</h3>
      <p>Cosmos blockchain applications require a config file (genesis) that contains the initial state of your blockchain. Genesis can be configured with a config.yml file, in which you can set accounts, token balances, validators, etc.</p>
      <div class="flex -align-center">
        <code>Go to `my-blog/config.yml` to configure your config file.</code>
        <a href="#">Config guide</a>
      </div>
      <p>Right now your blockchain runs on a single server node on your machine, but in production it'll be deployed on a network of machines around the world.</p>

      <h3>03. Next Steps</h3>
      <div class="flex -justify-between">
        <div class="card">
          <div class="card__main">
            <h5>Tokens</h5>
            <p>Create your very own token and distribute to anyone.</p>
          </div>
          <div class="card__sub">
            <a href="#">Start tutorial</a>
          </div>
        </div>
        <div class="card">
          <div class="card__main">
            <h5>Governance</h5>
            <p>Create governance proposals and vote on them.</p>
          </div>
          <div class="card__sub">
            <a href="#">Start tutorial</a>
          </div>
        </div>
        <div class="card">
          <div class="card__main">
            <h5>Custom data types</h5>
            <p>Create custom data types ($EXAMPLES), define operations, and display data on a web app.</p>
          </div>
          <div class="card__sub">
            <a href="#">Start tutorial</a>
          </div>
        </div>
      </div>
    </div>
  </div>
  
</template>

<script>
export default {

}
</script>

<style scoped>

.main {
  margin-bottom: 6rem;
}

.hero {
  margin-bottom: 4rem;
}
.hero h2 {
  font-size: 3rem;
  font-weight: 900;
}
.hero h4 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.texts-field {
  max-width: 796px;
}
.texts-field p {
  font-size: 1.125rem;
  line-height: 1.5rem;
  margin-bottom: 1rem;
  color: var(--c-txt-secondary);
}
.texts-field p ~ h3 {
  margin-top: 4rem;
}
.texts-field h3 {
  font-size: 1.8rem;
  font-weight: 800;
  margin-bottom: 1.5rem;
}
.texts-field .flex {
  display: flex;
  margin: 1.85rem 0 1.85rem 0;  
}
.texts-field .flex.-align-center {
  align-items: center;
}
.texts-field .flex.-justify-between {
  justify-content: space-between;
}
.texts-field p ~ a {
  margin-top: 1.5rem;
  display: inline-block;
}
.texts-field a {
  color: #66a1ff;
}
.texts-field code {
  display: inline-block;
  background: var(--c-bg-secondary);
  padding: 1.5rem 1.85rem;
}
.texts-field code + a {
  margin-left: 2rem;
}
.texts-field .card {
  background: var(--c-bg-secondary);
  padding: 1.125rem 1.25rem;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.texts-field .flex .card {
  width: calc(100% / 3 - 1rem);
}
.texts-field .card h5 {
  font-size: 1.125rem;
  margin-bottom: 0.85rem;
}
.texts-field .card p {
  font-size: 1rem;
  line-height: 1.25rem;
}
.texts-field .card .card__main {
  margin-bottom: 3rem;
}

</style>