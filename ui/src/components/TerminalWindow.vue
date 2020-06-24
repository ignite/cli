<template>
  <div>
    <div class="container">
      <div class="row" v-if="messagesFiltered.length == 0">Waiting for blocks...</div>
      <div v-for="msg in messagesFiltered" :key="msg.timestamp">
        <div class="row" v-if="!empty(msg)">
          <div class="timestamp">Height: {{msg.msg.result.data.value.block.header.height}}</div>
          <vue-json-pretty :data="msg.msg.result" :deep="1" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  background-color: #111;
  color: rgba(255, 255, 255, 0.85);
  padding: 1rem;
  height: 400px;
  overflow-y: scroll;
  border-radius: 0.75rem;
}
.row {
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
  padding: 1rem 0;
  font-family: Monaco, Menlo, Consolas, Bitstream Vera Sans Mono, monospace;
  font-size: 14px;
  position: relative;
}
.timestamp {
  position: absolute;
  top: 1rem;
  right: 0;
  color: rgba(255, 255, 255, 0.25);
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
    VueJsonPretty
  },
  data: function() {
    return {
      messages: []
    };
  },
  methods: {
    empty(msg) {
      return (
        Object.keys(msg.msg.result).length === 0 &&
        msg.msg.result.constructor === Object
      );
    }
  },
  computed: {
    messagesFiltered() {
      return this.messages.filter(m => {
        return !this.empty(m);
      });
    }
  },
  created() {
    const ws = new ReconnectingWebSocket("ws://localhost:26657/websocket");
    ws.addEventListener("open", () => {
      ws.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "subscribe",
          id: "1",
          params: ["tm.event = 'NewBlock'"]
        })
      );
    });
    ws.onmessage = msg => {
      const message = {
        msg: JSON.parse(msg.data),
        timestamp: msg.timeStamp
      };
      this.messages = [message, ...this.messages];
    };
  }
};
</script>
