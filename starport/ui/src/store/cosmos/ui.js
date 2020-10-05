export default {
  namespaced: true,
  state: {
    blocksExplorerTableId: 'cosmosBlocksExplorer',
    tables: [
      // { id: null, isSheetActive: false }
    ]
  },
  getters: {
    blocksExplorerTableId: state => state.blocksExplorerTableId,
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
     * 
     * @param {object} state
     * @param {string} tableId
     * 
     */        
    createTable(state, tableId) {
      if (state.tables.filter(table => table.id === tableId).length>0) {
        console.warn(`TableId ${tableId} has been registered. Please register the table with another tableId.`)
        return 
      }
      
      state.tables.push({ id: tableId, isSheetActive: false })
    },
    /**
     * 
     * @param {object} state
     * @param {object|null} payload
     * @param {string|null} payload.tableId
     * @param {boolean|null} payload.sheetState
     * 
     */        
    setTableSheetState(state, payload) {
      state.tables.filter(table => table.id === payload.tableId)[0]
        .isSheetActive = payload.sheetState
    }
  },
  actions: {}
}