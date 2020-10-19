import blockHelpers from '@/mixins/blocks/helpers'

export default {
  namespaced: true,
  state: {
    stack: [],
  },
  getters: {
    txsStack: state => state.stack,
    txByHash: state => hash => state.stack.filter(tx => tx.txhash === hash),
    txByEncodedHash: state => hash => state.stack.filter(tx => tx.txEncoded && tx.txEncoded === hash),
  },
  mutations: {
    addTxEntry(state, { tx }) {
      /**
       * 
       // 1. No txs in stack yet
       * 
       */
      if (state.stack.length===0) {
        state.stack.push(tx)
        return
      }

      /**
       * 
       // 2. Txs already exist in the stack
       * 
       */      
      for (let txIndex=0; txIndex<state.stack.length; txIndex++) {
        const currentTxVal = state.stack[txIndex]
        const nextTxVal = state.stack[txIndex+1]
        
        // Push tx to the end of the stack
        if (!nextTxVal) {
          state.stack.push(tx)
          break
        }

        const txHeight = parseInt(tx.height)
        const currentTxValHeight = parseInt(currentTxVal.height)
        const nextTxValHeight = parseInt(nextTxVal.height)
        // Add tx to the start of the stack
        if (txHeight>currentTxValHeight && txIndex===0) {
          state.stack.unshift(tx)
          break
        }
        // Insert tx to the stack
        if (currentTxValHeight>txHeight && txHeight>nextTxValHeight) {
          state.stack.splice(txIndex+1, 0, tx)
          break
        }
      }
    }
  },
  actions: {
    /**
     * Add tx into txsStack
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {object} payload.txData
     * 
     * 
     */        
    addTxEntry: {
      root: true,
      handler({ commit, getters }, txData) {
        /*
         *
         // If tx is null, it's not decoded successfully,
         // and triggered from `addErrorTx` in `blocks` store mutations.
         *
         */
        const fmtTxData = txData.tx === null ? 
          {
            height: txData.height,
            txEncoded: txData.txEncoded
          } : txData.tx.data
        
        const isTxInStack = txData.tx === null
          ? getters.txByEncodedHash(txData.txEncoded).length>0
          : getters.txByHash(txData.tx.data.txhash).length>0

        if (!isTxInStack) commit('addTxEntry', { tx: fmtTxData })
      }
    },
    /**
     * Add tx into txsStack
     * 
     * @param {object} store 
     * @param {object} payload
     * @param {object} payload.txData
     * 
     * 
     */     
    initTxsStack({ dispatch, rootGetters }) {
      const { fetchBlockMeta, fetchBlockchain, fetchLatestBlock, fetchDecodedTx } = blockHelpers
      const appEnv = rootGetters['cosmos/appEnv']

      fetchLatestBlock(appEnv.API)
        .then(latestBlock => {
          const blockHeight = latestBlock.data.block.header.height

          /**
           * 
           * ⚠️ TODO: refactor code (messy)
           * 
           */
          for (let height=blockHeight; height>0; height-=20) {
            fetchBlockchain({
              rpcUrl: appEnv.RPC,
              minBlockHeight: undefined,
              maxBlockHeight: height,
              latestBlockHeight: height
            }).then(blockchainRes => {
                const blockchain = blockchainRes.data.result.block_metas
      
                const promiseLoop = async _ => {
                  for (let i=0; i<blockchain.length; i++) {
                    const { header: prevHeader } = blockchain[i]
      
                    await fetchBlockMeta(appEnv.RPC, prevHeader.height)
                      .then(blockMeta => {
                        const blockTxs = blockMeta.data.result.block.data.txs
  
                        if (blockTxs && blockTxs.length>0) {
                          const txsDecoded = blockTxs
                            .map(txEncoded => fetchDecodedTx(
                              appEnv.API,
                              txEncoded,
                              (txEncoded, errLog) => dispatch('txErrCallback', {
                                blockHeight: header.height,
                                txEncoded,
                                errLog
                              }, { root: true })
                            ))
                        
                          txsDecoded.forEach(txRes => txRes.then(txResolved => {
                            dispatch('addTxEntry', { tx: txResolved }, { root: true })
                          }))
                        }
                      })                         
                  }                  
                }
                promiseLoop()
              })  
          }
  
        })
    }
  }
}