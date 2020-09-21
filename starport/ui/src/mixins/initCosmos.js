import blockConnection from './blocks/initConnection'
import backendConnection from './backend/initConnection'

export default {
  mixins: [blockConnection, backendConnection]
}