import ReconnectingWebSocket from 'reconnecting-websocket'
import blockHelpers from '@/mixins/blocks/helpers'

export default {
  namespaced: true,
  state: {
    chainId: null,
    table: {
      highlightedBlock: {
        id: null,
        data: null
      }
    },
    entries: [],
    errorsQueue: []
  },
  getters: {
    highlightedBlock: state => state.table.highlightedBlock,
    blockEntries: state => state.entries,
    blockByHeight: state => height => state.entries.filter(block => block.height === height),
    chainId: state => state.chainId
  },
  mutations: {
    /**
     * 
     * 
     * @param {object|null} block
     * @param {string|null} block[].id
     * @param {object|null} block[].data
     * 
     * 
     */    
    setHighlightedBlock(state, block) {
      if (block == null || !block) {
        state.table.highlightedBlock = {
          id: null,
          data: null
        }
      } else {
        state.table.highlightedBlock = block
      }
    },
    /**
     * 
     * 
     * @param {boolean} tableState
     * 
     * 
     */    
    setTableSheetState(state, tableState) {
      state.table.isSheetActive = tableState
    },
    /**
     * 
     * 
     * @param {string} chainId
     * 
     * 
     */    
    setChainId(state, chainId) {
      if (!state.chainId || state.chainId.length<=0) state.chainId = chainId
    }, 
    /**
     * 
     * 
     * @param {object} block
     * TODO: define shape of block object
     * 
     * 
     */     
    addBlockEntry(state, block) {
      state.entries.unshift(block)
    },        
    addErrorBlock(state, {
      blockHeight,
      errLog
    }) {
      state.errorsQueue.push({blockHeight, errLog})
    },       
    addErrorTx(state, {
      blockHeight,
      txEncoded,
      errLog
    }) {
      let isBlockInQueue = false
      
      for (let errBlock of state.errorsQueue) {
        if (blockHeight === errBlock.blockHeight) {
          errBlock.txError = {
            txEncoded,
            errLog            
          }
          isBlockInQueue = true
          break          
        }      
      }

      if (!isBlockInQueue) {
        state.errorsQueue.push({blockHeight, txError: {
          txEncoded,
          errLog
        }})
      }
    }
  },
  actions: {
    addBlockEntry({ commit }, block) {
      commit('addBlockEntry', block)
    },
    initBlockConnection({ commit, getters }, { localEnv }) {
      const ws = new ReconnectingWebSocket(`ws://${localEnv.COSMOS_RPC}/websocket`)
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
  
          const blockErrCallback = (errLog) => commit('addErrorBlock', {
            blockHeight: header.height,
            errLog
          })
          const txErrCallback = (txEncoded, errLog) => commit('addErrorTx', {
            blockHeight: header.height,
            txEncoded,
            errLog
          })
  
          fetchBlockMeta(localEnv.COSMOS_RPC, header.height, blockErrCallback)
            .then(blockMeta => {
              blockHolder.setBlockMeta(blockMeta)
  
              if (txsData.txs && txsData.txs.length > 0) {
                const txsDecoded = txsData.txs
                  .map(txEncoded => fetchDecodedTx(localEnv.LCD, txEncoded, txErrCallback))
                
                txsDecoded.forEach(txRes => txRes.then(txResolved => {
                  blockHolder.setBlockTxsDecoded(txResolved)
                }))
              }    
              // this guards duplicated block pushed into blockEntries
              if (getters.blockByHeight(blockHolder.block.height).length<=0) {
                commit('addBlockEntry', blockHolder.block)
                commit('setChainId', blockHolder.block.header.chain_id)
              }
            })
        }         
      }      
    }
  }
}

