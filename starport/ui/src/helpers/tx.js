import axios from 'axios'
import { sha256 } from 'js-sha256'

const txHelpers = {}

/**
 * 
 * 
 * @param {string} lcdUrl
 * @param {string} txEncoded
 * @param {function} errCallback
 *
 *  
 */      
txHelpers.getDecodedTx = async (lcdUrl, txEncoded, errCallback) => {
  const hashedTx = sha256(Buffer.from(txEncoded, 'base64'))
  try {
    // return await axios.post(`${lcdUrl}/txs/decode`, { tx: txEncoded }) 
    return await axios.get(`${lcdUrl}/txs/${hashedTx}`) 
  } catch (err) {
    console.error(txEncoded, err)
    if (errCallback) errCallback(txEncoded, err)
  }        
}

export default txHelpers