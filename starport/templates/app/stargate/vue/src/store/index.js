import { createStore, createLogger } from 'vuex'
//import { createStore } from 'vuex'
import init from './config'

const store = createStore({
	state() {
		return {}
	},
	mutations: {},
	actions: {},
	plugins: [createLogger()]
})
init(store)
export default store
