import axios from 'axios'
import moment from 'moment'

const getBlockTemplate = (header, txsData) => ({
  height: header.height,
  header,
  txs: txsData.txs,
  blockMeta: null,
  txsDecoded: []
})

export default {
  async fetchBlockMeta(cosmosUrl, blockHeight, errCallback) {
    try {
      return await axios.get(`http://${cosmosUrl}/block?${blockHeight}`)
    } catch (err) {
      console.error(err)
      errCallback(err)
    }
  },
  async fetchDecodedTx(lcdUrl, txEncoded, errCallback) {
    try {
      return await axios.post(`http://${lcdUrl}/txs/decode`, { tx: txEncoded }) 
    } catch (err) {
      console.error(txEncoded, err)
      errCallback(txEncoded, err)
    }        
  },   
  blockFormatter() {
    return {
      /**
       * @param {object} header
       * @param {object} txsData
       * TODO: define shape
       */            
      setNewBlock(header, txsData) {
          const blockTemplate = getBlockTemplate(header, txsData)
          const setBlockMeta = (msg) => {
            blockTemplate.blockMeta = msg.data.result.block_meta ? msg.data.result.block_meta : msg.data.result 
          }
          const setBlockTxsDecoded = (tx) => {
            blockTemplate.txsDecoded.push(tx.data.result)
          }

          return {
            block: blockTemplate,
            setBlockMeta,
            setBlockTxsDecoded
          }
      },
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
                txs: block.txs ? block.txs.length : 0          
              },
              tableData: {
                id: height,
                isActive: false
              },
              txs: block.txsDecoded
            }          
          })        
        }
      },
      /**
       * @param {array} blockEntries
       * TODO: define shape of block object
       */          
      filterBlock(blockEntries) {
        const hideBlocksWithoutTxs = () => {
          return blockEntries.filter(block => block.txs && block.txs.length > 0)
        }

        return {
          hideBlocksWithoutTxs
        }
      }
    }
  }
}