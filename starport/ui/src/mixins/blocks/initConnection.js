import { mapGetters, mapActions, mapMutations } from 'vuex'
import ReconnectingWebSocket from "reconnecting-websocket"

import blockHelpers from './helpers'

export default {
  computed: {
    ...mapGetters('cosmos', [ 'localEnv' ]),
    ...mapGetters('cosmos/blocks', [ 'blockByHeight' ]),
  },
  methods: {
    ...mapMutations('cosmos/blocks', [ 'addErrorBlock', 'addErrorTx' ]),
    ...mapActions('cosmos/blocks', [ 'addBlockEntry' ]),
  },
  created() {
    const ws = new ReconnectingWebSocket(`ws://${this.localEnv.COSMOS_RPC}/websocket`)
    // const ws = new ReconnectingWebSocket(`wss://${this.localEnv.COSMOS_RPC}:443/websocket`, [], { WebSocket: WebSocket })    

    ws.onopen = function() {
      ws.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "subscribe",
          id: "1",
          params: ["tm.event = 'NewBlock'"]
        })
      )
    }
    
    ws.onmessage = (msg) => {
      const { result } = JSON.parse(msg.data)

      if (result.data && result.events) {
        const { data } = result        
        const { data: txsData, header } = data.value.block

        const { fetchBlockMeta, fetchDecodedTx } = blockHelpers
        const blockFormatter = blockHelpers.blockFormatter()
        const blockHolder = blockFormatter.setNewBlock(header, txsData)

        const blockErrCallback = (errLog) => this.addErrorBlock({
          blockHeight: header.height,
          errLog
        })
        const txErrCallback = (txEncoded, errLog) => this.addErrorTx({
          blockHeight: header.height,
          txEncoded,
          errLog
        })

        fetchBlockMeta(this.localEnv.COSMOS_RPC, header.height, blockErrCallback)
          .then(blockMeta => {
            blockHolder.setBlockMeta(blockMeta)

            if (txsData.txs && txsData.txs.length > 0) {
              const txsDecoded = txsData.txs
                .map(txEncoded => fetchDecodedTx(this.localEnv.LCD, txEncoded, txErrCallback))
              
              txsDecoded.forEach(txRes => txRes.then(txResolved => {
                blockHolder.setBlockTxsDecoded(txResolved)
              }))
            }    
            // this guards duplicated block pushed into blockEntries
            if (this.blockByHeight(blockHolder.block.height).length<=0) {
              this.addBlockEntry(blockHolder.block)
            }
          })
      }         
    }
  }
}