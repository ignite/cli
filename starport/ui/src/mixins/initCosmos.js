import blockConnection from './blocks/initConnection'
import backendConnection from './backend/initConnection'

export default {
  mixins: [backendConnection, blockConnection]
}