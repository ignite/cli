import { txClient, queryClient, MissingWalletError } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

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

function mergeResults(value, next_values) {
	for (let prop of Object.keys(next_values)) {
		if (Array.isArray(next_values[prop])) {
			value[prop]=[...value[prop], ...next_values[prop]]
		}else{
			value[prop]=next_values[prop]
		}
	}
	return value
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
        {{ range .Module.HTTPQueries }}{{ .Name }}: {},
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
        {{ range .Module.HTTPQueries }}get{{ .Name }}: (state) => (params = {}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
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
		{{ range .Module.HTTPQueries }}async {{ .FullName }}({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value
				{{ if (index .Rules 0).HasBody}}
					{{ if (index .Rules 0).HasQuery}}				
				value = (await queryClient.{{ camelCase .FullName }}({{ range $i,$a :=(index .Rules 0).Params}} key.{{$a}}, {{end}} query, {...key})).data
					{{ else }}
				value = (await queryClient.{{ camelCase .FullName }}({{ range $i,$a :=(index .Rules 0).Params}} key.{{$a}}, {{end}} {...key})).data
					{{ end }}
				{{ else }}
				if (query) {
					value = (await queryClient.{{ camelCase .FullName }}({{ range $i,$a :=(index .Rules 0).Params}} key.{{$a}}, {{end}} query)).data
				}else{
					value = (await queryClient.{{ camelCase .FullName }}({{ range $i,$a :=(index .Rules 0).Params}}{{ if (gt $i 0)}}, {{ end}} key.{{$a}} {{end}})).data
				}
				{{ end }}
				{{ if (index .Rules 0).HasQuery}}
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await queryClient.{{ camelCase .FullName }}({{ range $i,$a :=(index .Rules 0).Params}} key.{{$a}}, {{end}}{...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					value = mergeResults(value,next_values);
				}
				{{ end }}
				commit('QUERY', { query: '{{ .Name }}', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: '{{ .FullName }}', payload: { options: { all }, params: {...key},query }})
				return getters['get{{.Name }}']( { params: {...key}, query}) ?? {}
			} catch (e) {
				let err = new SpVuexError('QueryClient:{{ .FullName }}', 'API Node Unavailable. Could not perform query.')
				err.original = e
				console.error(err)
				return {}
			}
		},
		{{ end }}
		{{ range .Module.Msgs }}async send{{ .Name }}({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.{{ camelCase .Name }}(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
  gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError ) {
					throw new SpVuexError('TxClient:{{ .Name }}:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					let err = new SpVuexError('TxClient:{{ .Name }}:Send', 'Could not broadcast Tx.')
					err.original = e
					throw err
				}
			}
		},
		{{ end }}
		{{ range .Module.Msgs }}async {{ .Name }}({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.{{ camelCase .Name }}(value)
				return msg
			} catch (e) {
				if (e.toString()=='wallet is required') {
					throw new SpVuexError('TxClient:{{ .Name }}:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					let err = new SpVuexError('TxClient:{{ .Name }}:Create', 'Could not create message.')
					err.original = e
					throw err
				}
			}
		},
		{{ end }}
	}
}
