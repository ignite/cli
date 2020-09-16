import { mapGetters, mapActions } from 'vuex'
import axios from "axios"
import ReconnectingWebSocket from "reconnecting-websocket"

export default {
  data() {
    return {
      TEMP_ENV: {
        // COSMOS_RPC: 'rpc.nylira.net',
        COSMOS_RPC: 'localhost:26657',
        LCD: 'localhost:1317'
      }
    }
  },
  computed: {
    ...mapGetters('cosmos/blocks', [ 'blockByHeight' ]),
  },
  methods: {
    ...mapActions('cosmos/blocks', [ 'addBlockEntry' ]),
  },
  mounted() {
    // let ws = new ReconnectingWebSocket(`wss://${this.TEMP_ENV.COSMOS_RPC}:443/websocket`, [], { WebSocket: WebSocket })
    const ws = new ReconnectingWebSocket(`ws://${this.TEMP_ENV.COSMOS_RPC}/websocket`)

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

      /* TODO: move block processors into ./helpers mixins */
      if (result.data && result.events) {
        const { data, events } = result        
        const { data: txsData, header } = data.value.block

        console.log(result)

        async function fetchBlockMeta(cosmosUrl) {
          try {
            return await axios.get(`http://${cosmosUrl}/block?${header.height}`)
          } catch (err) {
            console.error(err)
          }
        }
        async function fetchDecodedTx(txEncoded, lcdUrl) {
          try {
            return await axios.post(`http://${lcdUrl}/txs/decode`, { tx: txEncoded }) 
          } catch (err) {
            console.error(txEncoded, err)
          }        
        }   
        
        /* TODO: Proposer address is in HEX format? Decoding API is required? */
        // async function fetchValidator() {
        //   try {
        //     console.log(header)
        //     return await axios.get(`https://lcd.nylira.net/staking/validators/${header.proposer_address}`)
        //   } catch (err) {
        //     console.error(err)
        //   }   
        // }
        // fetchValidator().then(validator => console.log(validator))

        const blockHolder = {
          height: '',
          header,
          txs: txsData.txs,
          blockMeta: null,
          txsDecoded: []
        }

        fetchBlockMeta(this.TEMP_ENV.COSMOS_RPC)
          .then(blockMeta => {
            blockHolder.height = blockMeta.data.result.block_meta.header.height
            blockHolder.blockMeta = blockMeta.data.result.block_meta

            if (txsData.txs && txsData.txs.length > 0) {
              const txsDecoded = txsData.txs.map(txEncoded => fetchDecodedTx(txEncoded, this.TEMP_ENV.LCD))
              
              txsDecoded.forEach(txRes => txRes.then(txResolved => {
                blockHolder.txsDecoded.push(txResolved.data.result)
              }))
            }    

            // this guards duplicated block pushed into blockEntries
            if (this.blockByHeight(blockHolder.height).length<=0) {
              this.addBlockEntry(blockHolder)
            }
          })
   
      }         
    }
  }
}