const blocks = {
  namespaced: true,
  state: {
    table: {
      highlightedBlock: {
        id: null,
        data: null
      }
    },
    entries: []
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
    }    
  },
  actions: {
    addBlockEntry({ commit }, block) {
      commit('addBlockEntry', block)
    }
  }
}

export default {
  namespaced: true,
  state: {
    tables: [
      // { id: null, isSheetActive: false }
    ]
  },
  getters: {
    targetTable: state => tableId => {
      const targetTable = state.tables.filter(table => table.id === tableId)[0]
      return targetTable ? targetTable : null
    },
    isTableSheetActive: (state, getters) => tableId => {
      const targetTable = getters.targetTable(tableId)
      return targetTable ? targetTable.isSheetActive : null
    }
  },
  mutations: {
    /**
     * @param {string} tableId
     */        
    createTable(state, tableId) {
      if (state.tables.filter(table => table.id === tableId).length>0) {
        console.warn(`TableId ${tableId} has been registered. Please register the table with another tableId.`)
        return 
      }
      
      state.tables.push({ id: tableId, isSheetActive: false })
    },
    /**
     * @param {object|null} payload
     * @param {string|null} payload[].tableId
     * @param {boolean|null} payload[].sheetState
     */        
    setTableSheetState(state, payload) {
      state.tables.filter(table => table.id === payload.tableId)[0]
        .isSheetActive = payload.sheetState
    }
  },
  actions: {},
  modules: {
    blocks
  }
}