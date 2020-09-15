const blocks = {
  namespaced: true,
  state: {
    table: {
      isSheetActive: false,
      highlightedBlock: {
        id: null,
        data: null
      }
    }
  },
  getters: {
    highlightedBlock: state => state.table.highlightedBlock,
    isTableSheetActive: state => state.table.isSheetActive
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
    }    
  },
  actions: {}
}

export default {
  namespaced: true,
  state: {},
  getters: {},
  mutations: {},
  actions: {},
  modules: {
    blocks
  }
}