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
    maxEntriesCount: 20,
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
    addBlockEntry(state, { block, toInsert=true }) {
      if (toInsert) {
        state.entries.unshift(block)
      } else {
        state.entries.push(block)
      }
    },     
    /**
     * 
     * 
     * @param {object} block
     * TODO: define shape of block object
     * 
     * 
     */      
    popOverloadBlockEntry(state) {
      if (state.entries.length > state.maxEntriesCount) {
        state.entries.pop()
      }
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
    async setBlockMeta({ dispatch, commit, getters, rootGetters }, {
      header,
      blockMeta,
      txsData,
      toInsertBlockToFront=true
    }) {
      const localEnv = rootGetters['cosmos/localEnv']      
      const { fetchDecodedTx } = blockHelpers

      const blockFormatter = blockHelpers.blockFormatter()
      const blockHolder = blockFormatter.setNewBlock(header, txsData)

      const txErrCallback = (txEncoded, errLog) => commit('addErrorTx', {
        blockHeight: header.height,
        txEncoded,
        errLog
      })      
                      
      blockHolder.setBlockMeta(blockMeta)
      blockHolder.setBlockTxs(fetchDecodedTx, localEnv.LCD, txErrCallback)
      
      // this guards duplicated block pushed into blockEntries
      if (getters.blockByHeight(blockHolder.block.height).length<=0) {
        dispatch('addBlockEntry', {
          block: blockHolder.block,
          toInsert: toInsertBlockToFront
        })
        commit('setChainId', blockHolder.block.header.chain_id)
      }     
    },
    addBlockEntry({ commit }, { block, toInsert=true }) {
      commit('popOverloadBlockEntry') // 1. Pop entry with index > 20
      commit('addBlockEntry', { block, toInsert })  // 2. Push entry into stack
    },
    initBlockConnection({ commit, dispatch, getters, rootGetters }) {
      const localEnv = rootGetters['cosmos/localEnv']
      const ws = new ReconnectingWebSocket(`ws://${localEnv.COSMOS_RPC}/websocket`)
  
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
  
          const { fetchBlockMeta, fetchBlockchain } = blockHelpers
  
          const blockErrCallback = (errLog) => commit('addErrorBlock', {
            blockHeight: header.height,
            errLog
          })

          // 1. Fetch previous 20 blocks initially (if there's any)
          if (getters.blockEntries.length <= 0) {
            fetchBlockchain(localEnv.COSMOS_RPC, header.height-1)
              .then(blockchainRes => {
                const blockchain = blockchainRes.data.result.block_metas

                const promiseLoop = async _ => {
                  for (let i=0; i<blockchain.length; i++) {
                    const { header: prevHeader } = blockchain[i]

                    await fetchBlockMeta(localEnv.COSMOS_RPC, prevHeader.height, blockErrCallback)
                      .then(blockMeta => {
                        dispatch('setBlockMeta', {
                          header: prevHeader,
                          blockMeta,
                          txsData: blockMeta.data.result.block.data,
                          toInsertBlockToFront: false
                        })                      
                      })                         
                  }                  
                }
                promiseLoop()

              })
          }          
          
          // 2. Regular block fetching
          fetchBlockMeta(localEnv.COSMOS_RPC, header.height, blockErrCallback)
            .then(blockMeta => {
              dispatch('setBlockMeta', {
                header,
                blockMeta,
                txsData
              })                      
            })
        }         
      }      
    }
  }
}

