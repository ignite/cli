import axios from 'axios'
import moment from 'moment'
import { capitalCase } from 'change-case'
import { sha256 } from 'js-sha256'

const getBlockTemplate = (header, txsData) => ({
  height: parseInt(header.height),
  header,
  txs: txsData.txs ? txsData.txs : 0,
  blockMeta: null,
  txsDecoded: []
})

export default {
  /**
   * 
   * 
   * @param {string} rpcUrl
   * @param {string|number} blockHeight
   * @param {function} errCallback
   *
   *  
   */      
  async fetchBlockMeta(rpcUrl, blockHeight, errCallback) {
    try {
      return await axios.get(`${rpcUrl}/block?height=${blockHeight}`)
    } catch (err) {
      console.error(err)
      if (errCallback) errCallback(err)
    }
  },
  /**
   * 
   * 
   * @param {string} rpcUrl
   * @param {function} errCallback
   *
   *  
   */      
  async fetchLatestBlock(apiUrl, errCallback) {
    try {
      return await axios.get(`${apiUrl}/blocks/latest`)
    } catch (err) {
      console.error(err)
      if (errCallback) errCallback(err)
    }
  },
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
  async fetchBlockchain({
    rpcUrl,
    minBlockHeight=undefined,
    maxBlockHeight=undefined,
    latestBlockHeight,
    maxStackCount=20,
    errCallback
  }) {
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
  },
  /**
   * 
   * 
   * @param {string} lcdUrl
   * @param {string} txEncoded
   * @param {function} errCallback
   *
   *  
   */      
  async fetchDecodedTx(lcdUrl, txEncoded, errCallback) {
    const hashedTx = sha256(Buffer.from(txEncoded, 'base64'))
    try {
      // return await axios.post(`${lcdUrl}/txs/decode`, { tx: txEncoded }) 
      return await axios.get(`${lcdUrl}/txs/${hashedTx}`) 
    } catch (err) {
      console.error(txEncoded, err)
      if (errCallback) errCallback(txEncoded, err)
    }        
  },   
  /**
   * 
   * 
   * @param {array} blocksStack 
   * 
   * 
   */
  getGapBlock(blocksStack) {
    for (let i=0; i<blocksStack.length; i++) {
      const currentBlock = blocksStack[i]
      const nextBlock = blocksStack[i+1]
      if (!nextBlock) continue
      
      if (parseInt(currentBlock.height) - parseInt(nextBlock.height) > 1) {
        return {
          block: currentBlock,
          index: i+1
        }
      }        
    }

    return null
  },
  /**
   * 
   * Container of methods for formatting block data
   * 
   */
  blockFormatter() {
    return {
      /**
       * 
       * 
       * @param {object} header
       * @param {object} txsData
       * TODO: define shape
       *
       *  
       */            
      setNewBlock(header, txsData) {
          const blockTemplate = getBlockTemplate(header, txsData)
          /**
           * 
           * 
           * @param {object} msg
           * TODO: define shape
           *
           *  
           */              
          const setBlockMeta = (msg) => {
            blockTemplate.blockMeta = msg.data?.result ? msg.data.result : msg
          }
          /**
           * 
           * 
           * @param {object} tx
           * TODO: define shape
           *
           *  
           */                        
          const setBlockTxsDecoded = (tx) => {
            blockTemplate.txsDecoded.push(tx.data)
          }
          /**
           * 
           * 
           * @param {function} fetchDecodedTx
           * TODO: define shape
           *
           *  
           */               
          const setBlockTxs = ({
            fetchDecodedTx,
            lcdUrl,
            txStackCallback,
            txErrCallback
          }) => {
            if (txsData.txs && txsData.txs.length > 0) {
              const txsDecoded = txsData.txs
                .map(txEncoded => fetchDecodedTx(lcdUrl, txEncoded, txErrCallback))

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
      },
      /**
       * 
       * 
       * @param {array} blocksStack
       * TODO: define shape of block object
       * 
       * 
       */    
      blockForTable(blocksStack) {
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
      },
      /**
       * 
       * 
       * @param {array} txs
       * TODO: define shape of block object
       * 
       * 
       */    
      txsForCard(txs, chainId) {
        return txs.map(tx => {
          const {
            code,
            gas_used,
            gas_wanted,
            logs,
            raw_log,
            tx: txObj,
            txhash
          } = tx

          const {
            fee,
            memo,
            msg
          } = txObj.value
  
          return {
            meta: this.txMeta({
              code,
              gas_used,
              gas_wanted,
              raw_log,
              fee,
              memo,  
              txhash            
            }),
            msgs: msg.map(msg => this.txMsg(msg, chainId)),
            tableData: {
              id: txhash,
              isActive: false
            },
          }
        })
      },      
      /**
       * 
       * 
       * @param {object} metaData
       * @param {string} metaData[].fee
       * @param {string} metaData[].memo
       * TODO: define shape of block object
       * 
       * 
       */                
      txMeta({
        code,
        gas_used,
        gas_wanted,
        raw_log,
        fee,
        memo,          
        txhash   
      }) {
        const fmtFee = fee.amount[0]

        return {
          'Status': !code ? 'Success' : 'Fail',
          'TxHash': txhash,
          'Gas (used / wanted)': `${gas_used} / ${gas_wanted}`,
          'Fee': fmtFee ? `${fmtFee.amount} ${fmtFee.denom}` : 'N/A',
          'Memo': memo && memo.length>0 ? memo : 'N/A',
        }
      },
      /**
       * @param {object} msg
       * @param {string} msg[].type
       * @param {object} msg[].value
       * @param {string} chainId
       * TODO: define shape of block object
       */                
      txMsg({ type, value }, chainId) {
        const amountDenomHolder = { amount: '', denom: '' }
        
        function setMsgHolder(msgs, holder) {
          for (const [key, msg] of Object.entries(msgs)) {          
            if (Array.isArray(msg)) {
              msg.forEach(subMsg => setMsgHolder(subMsg, holder))
              break
            }
            if (key !== 'amount' && key !== 'denom') {
              holder[capitalCase(key)] = msg
            } else {
              amountDenomHolder[key] = msg
            }
          }

          if (amountDenomHolder.amount) {
            holder['Amount'] = `${amountDenomHolder.amount} ${amountDenomHolder.denom}`
          }
        }        

        const msgHolder = {
          type: this.getMsgType(type, chainId)
        }

        setMsgHolder(value, msgHolder)

        return msgHolder
      },
      /**
       * 
       * 
       * @param {array} blocksStack
       * TODO: define shape of block object
       * 
       * 
       */          
      filterBlock(blocksStack) {
        const hideBlocksWithoutTxs = () => {
          return blocksStack.filter(block => block.txs && block.txs.length > 0)
        }

        return {
          hideBlocksWithoutTxs
        }
      },
      /**
       * 
       * 
       * @param {object} amountObj
       * TODO: define shape of amount object
       * 
       * 
       */             
      getAmount(amountObj) {
        return amountObj.amount 
          ? amountObj.amount+amountObj.denom
          : amountObj[0].amount+amountObj[0].denom
      },
      /**
       * 
       * 
       * @param {string} type
       * @param {string} prefix
       * 
       * 
       */             
      getMsgType(type, prefix) {
        return type.replace(prefix+'/', '')
      },
      txMsg({ type, value }, chainId) {
        const amountDenomHolder = { amount: '', denom: '' }
        
        function setMsgHolder(msgs, holder) {
          for (const [key, msg] of Object.entries(msgs)) {          
            if (Array.isArray(msg)) {
              msg.forEach(subMsg => setMsgHolder(subMsg, holder))
              break
            }
            if (key !== 'amount' && key !== 'denom') {
              holder[capitalCase(key)] = msg
            } else {
              amountDenomHolder[key] = msg
            }
          }

          if (amountDenomHolder.amount) {
            holder['Amount'] = `${amountDenomHolder.amount} ${amountDenomHolder.denom}`
          }
        }        

        const msgHolder = {
          type: this.getMsgType(type, chainId)
        }

        setMsgHolder(value, msgHolder)

        return msgHolder
      },           
    }
  }
}