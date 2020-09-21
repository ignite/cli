export default {
  namespaced: true,
  state: {
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
    blockByHeight: state => height => state.entries.filter(block => block.height === height)
  },
  mutations: {
    /**
     * @param {object|null} block
     * @param {string|null} block[].id
     * @param {object|null} block[].data
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
     * @param {boolean} tableState
     */    
    setTableSheetState(state, tableState) {
      state.table.isSheetActive = tableState
    },
    /**
     * @param {object} block
     * TODO: define shape of block object
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
    }
  }
}

