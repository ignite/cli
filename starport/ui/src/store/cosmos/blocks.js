import axios from 'axios'
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
    maxBlockchainCount: 20,
    maxStackCount: 100,
    stack: [],
    stackChainRange: {
      highestHeight: null,
      lowestHeight: null
    },
    latestBlock: null,
    errorsQueue: []
  },
  getters: {
    highlightedBlock: state => state.table.highlightedBlock,
    blocksStack: state => state.stack,
    blockByHeight: state => height => state.stack.filter(block => parseInt(block.height) === parseInt(height)),
    blockByHash: state => hash => state.stack.filter(block => block.blockMeta.block_id.hash === hash),
    latestBlock: state => state.latestBlock,
    stackChainRange: state => state.stackChainRange,
    lastBlock: state => state.stack[state.stack.length-1],
    gapBlock: state => blockHelpers.getGapBlock(state.stack),
    chainId: state => state.chainId,
    errorsQueue: state => state.errorsQueue
  },
  mutations: {
    /**
     * 
     * 
     * @param {object} payload
     * @param {number|string} payload.highest
     * @param {number|string} payload.lowest
     * 
     * 
     */
    setStackChainRange(state, { highest, lowest }) {
      if (highest) {
        state.stackChainRange.highestHeight = parseInt(highest)
      }
      if (lowest) {
        state.stackChainRange.lowestHeight = parseInt(lowest)
      }
    },
    /**
     * 
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
     * 
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
     * 
     * 
     * 
     */    
    sortBlocksStack(state) {
      state.stack.sort((a,b) => b.height - a.height)
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
      if (!toInsert) {
        state.stack.push(block)
      } else {
        state.stack.unshift(block)
      }
    },     
    /**
     * Pop overloaded blocks in stack (if more than 500)
     * 
     * @param {object} state
     * @param {boolean} toPop - default is True
     * 
     * 
     */      
    popOverloadBlocks(state, {
      toPop,
      toPopOverBlockchainCount
    }={
      toPop: true,
      toPopOverBlockchainCount: false
    }) {
      if (toPopOverBlockchainCount) {
        state.stack.splice(state.maxBlockchainCount)        
        return 
      }
      if (state.stack.length > state.maxStackCount) {
        if (toPop) {
          state.stack.splice(state.maxStackCount)
        } else {
          state.stack.splice(0, state.stack.length-state.maxStackCount)
        }
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
      errLog,
      txStackCallback
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

        txStackCallback()
      }
    }
  },
  actions: {
    /**
     * 
     * 
     */
    async setStackChainRange({ commit, getters }, {
      highest,
      lowest
    }={
      highest: null,
      lowest: null,
    }) {
      const stack = getters.blocksStack
      const highestHeight = highest ? highest : stack[0]?.height
      const lowestHeight = lowest ? lowest : stack[stack.length-1]?.height

      commit('setStackChainRange', {
        highest: highestHeight ? parseInt(highestHeight) : null, 
        lowest: lowestHeight ? parseInt(lowestHeight) : null
      })        
    },
    /**
     * 
     * 
     * 
     */
    txErrCallback: {
      root: true,
      handler({ commit, dispatch }, {
        blockHeight,
        txEncoded,
        errLog,
      }) {
        commit('addErrorTx', {
          blockHeight,
          txEncoded,
          errLog,
          txStackCallback: () => dispatch('addTxEntry', {
            tx: null,
            height,
            txEncoded
          }, { root: true })
        })          
      }    
    },
    /**
     * 
     * Fetch blocks (20 of which) from RPC endpoint
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {number} payload.blockHeight
     * @param {boolean} payload.toGetLowerBlocks - to get older or newer blocks
     * 
     * 
     */
    async getBlockchain({ commit, dispatch, getters, rootGetters }, {
      blockHeight,
      toGetLowerBlocks=true,
      toReset=false
    }) {
      const appEnv = rootGetters['cosmos/appEnv']      
      const { fetchBlockMeta, fetchBlockchain } = blockHelpers
      const latestBlock = getters.latestBlock    

      const blockErrCallback = (errLog) => commit('addErrorBlock', {
        blockHeight: header.height,
        errLog
      })

      /**
       * 
       // minHeight to maxHeight is the range of blockchain to fetch.
       // To get older blocks (blocks with lower heights),
       // set minHeight to `undefined`, because it's dependent on maxHeight (maxHeight-20).
       // And vise versa.
       * 
       */
      const minBlockHeight = toGetLowerBlocks ? undefined : parseInt(blockHeight)
      const maxBlockHeight = toGetLowerBlocks ? parseInt(blockHeight) : undefined
 
      const toFetchBlockchain = async (min, max, toInsert=false, toReset=false) => fetchBlockchain({
        rpcUrl: appEnv.RPC,
        minBlockHeight: min,
        maxBlockHeight: max,
        latestBlockHeight: latestBlock ? latestBlock.height : null
      }).then(blockchainRes => {
          const blockchain = blockchainRes.data.result.block_metas
          const toReverse = toReset ? true : toInsert          
          const fmtBlockchainRes = toReverse ? blockchain.reverse() : blockchain          

          return async _ => {
            for (let i=0; i<fmtBlockchainRes.length; i++) {
              const { header: prevHeader } = fmtBlockchainRes[i]
              await fetchBlockMeta(appEnv.RPC, prevHeader.height, blockErrCallback)
                .then(blockMeta => {
                  dispatch('setBlockMeta', {
                    header: prevHeader,
                    blockMeta,
                    txsData: blockMeta.data.result.block.data,
                    toInsertBlockToFront: toInsert,
                    toReset
                  })   
                })     
            }                  
          }
        })     
      
  
      await toFetchBlockchain(minBlockHeight, maxBlockHeight, !toGetLowerBlocks, toReset)
        .then(promiseLoop => promiseLoop()
          .then(() => {
            const isToPop = toReset ? true : !toGetLowerBlocks
            commit('popOverloadBlocks', {
              toPop: isToPop, 
              toPopOverBlockchainCount: toReset
            })
            dispatch('setStackChainRange')
          })
        )      
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
    async setBlockMeta({ dispatch, commit, getters, rootGetters, rootCommit }, {
      header,
      blockMeta,
      txsData,
      toInsertBlockToFront=false,
      isValidLatestBlock=false,
      toReset=false
    }) {
      const appEnv = rootGetters['cosmos/appEnv']      
      const { fetchDecodedTx } = blockHelpers

      const blockFormatter = blockHelpers.blockFormatter()
      const blockHolder = blockFormatter.setNewBlock(header, txsData)
                      
      blockHolder.setBlockMeta(blockMeta)
      blockHolder.setBlockTxs({
        fetchDecodedTx,
        lcdUrl: appEnv.API,
        txStackCallback: (tx) => dispatch('addTxEntry', { tx }, { root: true }),
        txErrCallback: (txEncoded, errLog) => dispatch('txErrCallback', {
          blockHeight: header.height,
          txEncoded,
          errLog
        }, { root: true })
      })
      
      // this guards duplicated block pushed into blocksStack
      if (getters.blockByHeight(blockHolder.block.height).length<=0) {
        /*
         *
         // 2. Check block position
         *
         */    
        const newBlockPosition = (() => {
          const { highestHeight, lowestHeight } = getters.stackChainRange
          const newBlockHeight = parseInt(blockHolder.block.height)

          let isHigher=false,
              isLower=false,
              isAdjacent=false

          if (!highestHeight && !lowestHeight) {
            isHigher = true
            isAdjacent = true
          } else if (newBlockHeight>highestHeight) {
            isHigher = true
            isAdjacent = !(newBlockHeight-highestHeight>1)
          } else if (newBlockHeight<lowestHeight) {
            isLower = true
            isAdjacent = !(lowestHeight-newBlockHeight>1)
          }

          return { isHigher, isLower, isAdjacent }
        })()

        /*
         *
         // 1. Save the latest block (if the block is coming from WS connection)
         *
         */  
        if (isValidLatestBlock) {
          commit('setLatestBlock', blockHolder.block)
        }                
        
        /*
         *
         // 3. Add block to stack (toReset is to travel to top of the chain)
         *
         */     
        if (newBlockPosition.isAdjacent || toReset) {
          commit('addBlockEntry', {
            block: blockHolder.block,
            toInsert: toReset ? true : newBlockPosition.isHigher
          })          
          dispatch('setStackChainRange', {
            highest: newBlockPosition.isHigher ? blockHolder.block.height : null,
            lowest: newBlockPosition.isLower ? blockHolder.block.height : null,
          })
        }    
        /*
         *
         // 4. Set application's chainId
         *
         */        
        commit('setChainId', blockHolder.block.header.chain_id)
      }     
    },
    /**
     * Initiate WS connection subscribes to LCD endpoint
     * 
     * @param {object} store 
     * 
     * 
     */        
    async initBlockConnection({ commit, dispatch, getters, rootGetters }) {
      const appEnv = rootGetters['cosmos/appEnv']     
      console.log(appEnv.STARPORT_APP)
      const { data } = await axios.get(`${appEnv.STARPORT_APP}/status`)      
      const GITPOD = data.env.vue_app_custom_url && new URL(data.env.vue_app_custom_url)
      const wsUrl = GITPOD
        ? process.env.VUE_APP_WS_TENDERMINT || (GITPOD && `wss://26657-${GITPOD.hostname}/websocket`)
        : 'ws://localhost:26657/websocket'

      const ws = new ReconnectingWebSocket(wsUrl) 
  
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
      
      ws.onmessage = async (msg) => {
        const { result } = JSON.parse(msg.data)
  
        if (result.data && result.events) {
          const { data } = result        
          const { data: txsData, header } = data.value.block
          const { fetchBlockMeta, fetchDecodedTx } = blockHelpers

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
            await dispatch('getBlockchain', { blockHeight: header.height })
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

            /**
             * 
             // 3. Check and resolve errors queue
             * 
             */
            if (getters.errorsQueue.length > 0) {
              getters.errorsQueue.forEach((errObj, index) => {
                const errBlockInStack = getters.blockByHeight(errObj.blockHeight)[0]

                if (errBlockInStack) {
                  if (errObj.txError && errObj.txError.txEncoded) {
                    fetchDecodedTx(appEnv.API, errObj.txError.txEncoded)
                      .then(txRes => {
                        const isTxAlreadyDecoded = errBlockInStack.txsDecoded
                          .filter(tx => tx.txhash === txRes.data.txhash).length>0
                        
                        if (!isTxAlreadyDecoded) {
                          errBlockInStack.txsDecoded.push(txRes.data)
                        }
                        getters.errorsQueue.splice(index,1)
                        console.info(`✨TX fetching error ${txRes.data.txhash} was resolved.`)
                      })
                  }
                }

              })
            }
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
        commit('setHighlightedBlock', { block })
        // await dispatch('setHighlightedBlockMeta', { block })
        //   .then(() => commit('setHighlightedBlock', { block }))
      }
    },    
  }
}

