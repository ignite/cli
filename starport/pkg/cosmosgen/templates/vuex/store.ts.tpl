import { txClient, queryClient } from './module'

{{ range .Module.Types }}import { {{ .Name }} } from "./module/types/{{ resolveFile .FilePath }}"
{{ end }}

async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['common/wallet/signer'], {
		addr: vuexGetters['common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['common/env/apiCosmos']
	})
}

function getStructure(template) {
	let structure = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field: any = {}
		field.name = key
		field.type = typeof value
		structure.fields.push(field)
	}
	return structure
}

const getDefaultState = () => {
	return {
        {{ range .Module.Queries }}{{ .Name }}: {},
        {{ end }}
        _Structure: {
            {{ range .Module.Types }}{{ .Name }}: getStructure({{ .Name }}.fromPartial({})),
            {{ end }}
		},
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(subscription)
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(subscription)
		}
	},
	getters: {
        {{ range .Module.Queries }}get{{ .Name }}: (state) => (params = {}) => {
			return state.{{ .Name }}[JSON.stringify(params)] ?? {}
		},
        {{ end }}
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('init')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach((subscription) => {
				dispatch(subscription.action, subscription.payload)
			})
		},
        {{ range .Module.Queries }}async {{ .FullName }}({ commit, rootGetters }, { subscribe = false, ...key }) {
			try {
				const value = (await (await initQueryClient(rootGetters)).{{ camelCase .FullName }}.apply(null, Object.values(key))).data
				commit('QUERY', { query: '{{ .Name }}', key, value })
				if (subscribe) commit('SUBSCRIBE', { action: '{{ .FullName }}', payload: key })
			} catch (e) {
				console.log('Query Failed: API node unavailable')
			}
		},
        {{ end }}
        {{ range .Module.Msgs }}async {{ .Name }}({ rootGetters }, { value }) {
			try {
				const msg = await (await initTxClient(rootGetters)).{{ camelCase .Name }}(value)
				await (await initTxClient(rootGetters)).signAndBroadcast([msg])
			} catch (e) {
				throw 'Failed to broadcast transaction: ' + e
			}
		},
        {{ end }}
	}
}
