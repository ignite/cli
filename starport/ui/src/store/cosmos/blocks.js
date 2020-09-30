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
     * Highlight the block selected in BlockTable
     * and keep the block in the store.
     * 
     * @param {object} state
     * @param {Object} payload
     * @param {object|null} payload.block
     * @param {string|null} payload.block.id
     * @param {object|null} payload.block.data
     * 
     * 
     */    
    setHighlightedBlock(state, { block }) {
      if (block == null || !block) {
        state.table.highlightedBlock = {
          id: null,
          data: null
        }
      } else {
        state.table.highlightedBlock = {
          ...state.table.highlightedBlock,
          ...block
        }
      }
    },
    /**
     * Set the state of table's side sheet to true/false
     * 
     * @param {object} state
     * @param {boolean} [tableState=false]
     * 
     * 
     */    
    setTableSheetState(state, tableState=false) {
      state.table.isSheetActive = tableState
    },
    /**
     * Set chainId of the app (if there's no existing one yet)
     * 
     * @param {object} state
     * @param {string} chainId
     * 
     * 
     */    
    setChainId(state, chainId) {
      if (!state.chainId || state.chainId.length<=0) state.chainId = chainId
    }, 
    /**
     * 
     * @param {object} state
     * @param {object} payload
     * @param {object} payload.block - The block to add into stack
     * @param {boolean} [payload.toInsert=true] - To push or unshift block into stack
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
     * Pop overloaded blocks in stack (if more than 500)
     * 
     * @param {object} state
     * 
     * 
     */      
    popOverloadBlocks(state) {
      if (state.stack.length > state.maxStackCount) {
        state.stack.splice(state.maxStackCount)
      }
    },       
    /**
     * Store the latest block fetched from WS connection
     * 
     * @param {object} state
     * @param {object} block
     * 
     * 
     */   
    setLatestBlock(state, block) {
      state.latestBlock = block
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
    /**
     * Fetch blocks (20 of which) from RPC endpoint
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {number} payload.blockHeight
     * @param {boolean} payload.toGetOlderBlocks - to get older or newer blocks
     * 
     * 
     */
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
    /**
     * Format the fetched block and add it into store's stack
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {object} payload.header
     * @param {object} payload.blockMeta
     * @param {object} payload.txsData
     * @param {boolean} payload.toInsertBlockToFront
     * @param {boolean} payload.isValidLatestBlock
     * 
     * 
     */    
    async setBlockMeta({ commit, getters, rootGetters }, {
      header,
      blockMeta,
      txsData,
      toInsertBlockToFront=true,
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
      blockHolder.setBlockTxs(fetchDecodedTx, appEnv.API, txErrCallback)
      
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
         // 2. Set application's chainId
         *
         */        
        commit('setChainId', blockHolder.block.header.chain_id)
        /*
         *
         // 3. Save the latest block (if the block is coming from WS connection)
         *
         */  
        if (isValidLatestBlock) {
          commit('setLatestBlock', blockHolder.block)
        }
      }     
    },
    /**
     * Initiate WS connection subscribes to LCD endpoint
     * 
     * @param {object} store 
     * 
     * 
     */        
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

          /**
           * 
           // 1. Fetch previous 20 blocks initially (if there's any) 
           * 
           */
          if (getters.blocksStack.length <= 0) {
            dispatch('getBlockchain', { blockHeight: header.height })
          }          
          
          /** 
           * 
           // 2. Regular block fetching
           * 
           */
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
    },
    /**
     * Fetch raw block's meta for highlighted block
     * and add rawJson data into highlightedBlock
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {object} payload.block
     * 
     * 
     */       
    async setHighlightedBlockMeta({ state, rootGetters }, { block }) {
      blockHelpers
        .fetchBlockMeta(rootGetters['cosmos/appEnv'].RPC, block.data.blockMsg.height)
        .then(blockMeta => state.table.highlightedBlock.rawJson = blockMeta)
    },
    /**
     * Set highlightedBlock to be null or active with block's info
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {object} payload.block
     * 
     * 
     */          
    async setHighlightedBlock({ dispatch, commit }, { block }) {
      if ( block == null || !block ) {
        commit('setHighlightedBlock', { block: null })
      } else {
        await dispatch('setHighlightedBlockMeta', { block })
          .then(() => commit('setHighlightedBlock', { block }))
      }
    },    
  }
}

