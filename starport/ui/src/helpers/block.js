import axios from 'axios'
import moment from 'moment'


/**
 * 
 * 
 * Container of helper methods to work with blocks
 * 
 * 
 */
const blockHelpers = {}


/**
 * 
 * 
 * @param {string} rpcUrl
 * @param {string|number} blockHeight
 * @param {function} errCallback
 *
 *  
 */      
blockHelpers.getBlockByHeight = async (rpcUrl, blockHeight, errCallback) => {
  try {
    return await axios.get(`${rpcUrl}/block?height=${blockHeight}`)
  } catch (err) {
    console.error(err)
    if (errCallback) errCallback(err)
  }
}


/**
 * 
 * 
 * @param {string} rpcUrl
 * @param {function} errCallback
 *
 *  
 */      
blockHelpers.getLatestBlock = async (apiUrl, errCallback) => {
  try {
    return await axios.get(`${apiUrl}/blocks/latest`)
  } catch (err) {
    console.error(err)
    if (errCallback) errCallback(err)
  }
}


/**
 * 
 * 
 * @param {Object} payload
 * @param {string} payload.rpcUrl
 * @param {string|number} [payload.minBlockHeight=undefined]
 * @param {string|number} [payload.maxBlockHeight=undefined]
 * @param {string|number} payload.latestBlockHeight
 * @param {number} [payload.maxStackCount=20]
 * @param {function} payload.errCallback
 *
 *  
 */      
blockHelpers.getBlockchain = async ({
  rpcUrl,
  minBlockHeight=undefined,
  maxBlockHeight=undefined,
  latestBlockHeight,
  maxStackCount=20,
  errCallback
}) => {
  if (!minBlockHeight && !maxBlockHeight) {
    console.error('Please provide min or max block height value')
    return
  } 

  const fmtMinHeight = () => {
    if (maxBlockHeight) {
      return maxBlockHeight-1 - maxStackCount >= 0 
        ? maxBlockHeight-1 - maxStackCount
        : 0
    }
    return minBlockHeight
  }
  const fmtMaxHeight = () => {
    if (minBlockHeight) {
      return minBlockHeight + maxStackCount >= latestBlockHeight
        ? latestBlockHeight
        : minBlockHeight + maxStackCount
    }
    return maxBlockHeight-1
  }

  try {
    return await axios.get(`${rpcUrl}/blockchain?minHeight=${fmtMinHeight()}&maxHeight=${fmtMaxHeight()}`)
  } catch (err) {
    console.error(err)
    if (errCallback) errCallback(err)
  }
}


/**
 * 
 * 
 * Container of methods to format block data
 * 
 * 
 */
export const formatter = {}


/**
 * 
 * 
 * @param {object} header
 * @param {object} txsData
 *
 *  
 */            
formatter.setNewBlock = (header, txsData) => {
  const blockTemplate = {
    height: parseInt(header.height),
    header,
    txs: txsData.txs ? txsData.txs : 0,
    blockMeta: null,
    txsDecoded: []
  }

  const setBlockMeta = (msg) => {
    blockTemplate.blockMeta = msg.data?.result ? msg.data.result : msg
  }      
  const setBlockTxsDecoded = (tx) => {
    blockTemplate.txsDecoded.push(tx.data)
  }          
  const setBlockTxs = ({
    getDecodedTx,
    lcdUrl,
    txStackCallback,
    txErrCallback
  }) => {
    if (txsData.txs && txsData.txs.length > 0) {
      const txsDecoded = txsData.txs
        .map(txEncoded => getDecodedTx(lcdUrl, txEncoded, txErrCallback))

      txsDecoded.forEach(txRes => txRes.then(txResolved => {
        if (txResolved) setBlockTxsDecoded(txResolved)
        if (txStackCallback) txStackCallback(txResolved)
      }))
    }                
  }

  return {
    block: blockTemplate,
    setBlockMeta,
    setBlockTxsDecoded,
    setBlockTxs
  }
}


/**
 * 
 * 
 * @param {array} blocksStack
 * 
 * 
 */    
formatter.blockForTable = (blocksStack) => {
  if (blocksStack.length > 0) {
    return blocksStack.map((block) => {
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
          time,
          height: parseInt(height),
          proposer: `${proposer_address.slice(0,10)}...`,
          blockHash_sliced: `${hash.slice(0,15)}...`,
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
}


/**
 * 
 * 
 * @param {array} blocksStack
 * 
 * 
 */          
formatter.filterBlock = (blocksStack) => {
  const hideBlocksWithoutTxs = () => {
    return blocksStack.filter(block => block.txs && block.txs.length > 0)
  }

  return {
    hideBlocksWithoutTxs
  }
}


export default blockHelpers