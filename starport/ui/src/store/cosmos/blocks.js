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
    maxStackCount: 500,
    stack: [],
    latestBlock: null,
    errorsQueue: []
  },
  getters: {
    highlightedBlock: state => state.table.highlightedBlock,
    blocksStack: state => state.stack,
    blockByHeight: state => height => state.stack.filter(block => block.height === height),
    latestBlock: state => state.latestBlock,
    lastBlock: state => state.stack[state.stack.length-1],
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
     * @param {boolean} toInsert
     * TODO: define shape of block object
     * 
     * 
     */     
    addBlockEntry(state, { block, toInsert=true }) {
      if (toInsert) {
        state.stack.unshift(block)
      } else {
        state.stack.push(block)
      }
    },     
    /**
     * 
     * 
     * 
     * 
     */      
    popOverloadBlocks(state) {
      if (state.stack.length > state.maxStackCount) {
        state.stack.splice(state.maxStackCount)
      }
    },       
    /**
     * 
     * 
     * @param {object} state
     * @param {object} block
     * TODO: define shape of block object
     * 
     * 
     */   
    setLatestBlock(state, block) {
      state.latestBlock = block
    },
    // addErrorBlock(state, {
    //   blockHeight,
    //   errLog
    // }) {
    //   state.errorsQueue.push({blockHeight, errLog})
    // },       
    // addErrorTx(state, {
    //   blockHeight,
    //   txEncoded,
    //   errLog
    // }) {
    //   let isBlockInQueue = false
      
    //   for (let errBlock of state.errorsQueue) {
    //     if (blockHeight === errBlock.blockHeight) {
    //       errBlock.txError = {
    //         txEncoded,
    //         errLog            
    //       }
    //       isBlockInQueue = true
    //       break          
    //     }      
    //   }

    //   if (!isBlockInQueue) {
    //     state.errorsQueue.push({blockHeight, txError: {
    //       txEncoded,
    //       errLog
    //     }})
    //   }
    // }
  },
  actions: {
    async getBlockchain({ dispatch, getters, rootGetters }, {
      blockHeight,
      toGetOlderBlocks=true
    }) {
      const appEnv = rootGetters['cosmos/appEnv']      
      const { fetchBlockMeta, fetchBlockchain } = blockHelpers
      const latestBlock = getters.latestBlock    

      const blockErrCallback = (errLog) => commit('addErrorBlock', {
        blockHeight: header.height,
        errLog
      })

      fetchBlockchain({
        rpcUrl: appEnv.RPC,
        minBlockHeight: undefined,
        maxBlockHeight: blockHeight,
        latestBlockHeight: latestBlock ? latestBlock.height : null
      })
        .then(blockchainRes => {
          const blockchain = blockchainRes.data.result.block_metas

          const promiseLoop = async _ => {
            for (let i=0; i<blockchain.length; i++) {
              const { header: prevHeader } = blockchain[i]

              await fetchBlockMeta(appEnv.RPC, prevHeader.height, blockErrCallback)
                .then(blockMeta => {
                  dispatch('setBlockMeta', {
                    header: prevHeader,
                    blockMeta,
                    txsData: blockMeta.data.result.block.data,
                    toInsertBlockToFront: false,
                    toPopOverloadBlocks: !toGetOlderBlocks
                  })                      
                })                         
            }                  
          }
          promiseLoop()

        })
    },
    async setBlockMeta({ commit, getters, rootGetters }, {
      header,
      blockMeta,
      txsData,
      toInsertBlockToFront=true,
      toPopOverloadBlocks=true,
      isValidLatestBlock=false
    }) {
      const appEnv = rootGetters['cosmos/appEnv']      
      const { fetchDecodedTx } = blockHelpers

      const blockFormatter = blockHelpers.blockFormatter()
      const blockHolder = blockFormatter.setNewBlock(header, txsData)

      const txErrCallback = (txEncoded, errLog) => commit('addErrorTx', {
        blockHeight: header.height,
        txEncoded,
        errLog
      })      
                      
      blockHolder.setBlockMeta(blockMeta)
      blockHolder.setBlockTxs(fetchDecodedTx, appEnv.LCD, txErrCallback)
      
      // this guards duplicated block pushed into blocksStack
      if (getters.blockByHeight(blockHolder.block.height).length<=0) {
        /*
         *
         // 1. Add block to stack
         *
         */        
        commit('addBlockEntry', {
          block: blockHolder.block,
          toInsert: toInsertBlockToFront
        })
        /*
         *
         // 3. Set application's chainId
         *
         */        
        commit('setChainId', blockHolder.block.header.chain_id)
        /*
         *
         // 4. Save the latest block (if the block is coming from WS connection)
         *
         */  
        if (isValidLatestBlock) {
          commit('setLatestBlock', blockHolder.block)
        }
      }     
    },
    initBlockConnection({ commit, dispatch, getters, rootGetters }) {
      const appEnv = rootGetters['cosmos/appEnv']
      const ws = new ReconnectingWebSocket(appEnv.WS) 
  
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
          const { fetchBlockMeta } = blockHelpers

          const blockErrCallback = (errLog) => commit('addErrorBlock', {
            blockHeight: header.height,
            errLog
          })          

          // 1. Fetch previous 20 blocks initially (if there's any)
          if (getters.blocksStack.length <= 0) {
            dispatch('getBlockchain', { blockHeight: header.height })
          }          
          
          // 2. Regular block fetching
          fetchBlockMeta(appEnv.RPC, header.height, blockErrCallback)
            .then(blockMeta => {
              dispatch('setBlockMeta', {
                header,
                blockMeta,
                txsData,
                isValidLatestBlock: true
              })                      
            })
        }         
      }      
    }
  }
}

