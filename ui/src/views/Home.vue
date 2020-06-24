<template>
  <div>
    <div class="container">
      <div class="row" v-for="msg in messages" :key="msg.timestamp">
        <vue-json-pretty :data="msg" :deep="1" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  background-color: #111;
  color: rgba(255, 255, 255, 0.85);
  padding: 1rem;
}
.row {
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
  padding: 1rem 0;
}
.vjs-tree,
.vjs-key {
  font-size: 1rem;
}
</style>

<script>
import VueJsonPretty from "vue-json-pretty";
import ReconnectingWebSocket from "reconnecting-websocket";

export default {
  components: {
    VueJsonPretty,
  },
  data: function() {
    return {
      messages: [],
    };
  },
  created() {
    const ws = new ReconnectingWebSocket("ws://localhost:26657/websocket");
    ws.addEventListener("open", () => {
      ws.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "subscribe",
          id: "1",
          params: ["tm.event = 'NewBlock'"],
        })
      );
    });
    ws.onmessage = (msg) => {
      const message = {
        msg: JSON.parse(msg.data),
        timestamp: msg.timeStamp,
      };
      this.messages = [message, ...this.messages];
    };
  },
};
</script>
