import moment from 'moment'

export default {
  methods: {
    $_blockFormatter() {
      return {
        /**
         * @param {array} blockEntries
         * TODO: define shape of block object
         */    
        blockForTable(blockEntries) {
          if (blockEntries.length > 0) {
            return blockEntries.map((block) => {
              const {
                time,
                height,
                proposer_address,
                num_txs
              } = block.header

              const {
                hash
              } = block.blockMeta.block_id

              return {
                blockMsg: {
                  time_formatted: moment(time).fromNow(true),
                  time: time,
                  height,
                  proposer: `${proposer_address.slice(0,10)}...`,
                  blockHash_sliced: `${hash.slice(0,30)}...`,
                  blockHash: hash,
                  txs: num_txs          
                },
                tableData: {
                  id: height,
                  isActive: false
                },
                txs: block.txsDecoded
              }          
            })        
          }
        }        
      }
    }
  }
}