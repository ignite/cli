import { Client, registry, MissingWalletError } from '{{ .PackageNS }}-client-ts'

{{ range .Module.Types }}import { {{ .Name }} } from "{{ $.PackageNS }}-client-ts/{{ $.Module.Pkg.Name }}/types"
{{ end }}

export { {{ range $i,$type:=.Module.Types }}{{ if (gt $i 0) }}, {{ end }}{{ $type.Name }}{{ end }} };

function initClient(vuexGetters) {
	return new Client(vuexGetters['common/env/getEnv'], vuexGetters['common/wallet/signer'])
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

type Field = {
	name: string;
	type: unknown;
}
function getStructure(template) {
	let structure: {fields: Field[]} = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field = { name: key, type: typeof value }
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
		_Registry: registry,
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
			state._Subscriptions.add(JSON.stringify(subscription))
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(JSON.stringify(subscription))
		}
	},
	getters: {
				{{ range .Module.HTTPQueries }}get{{ .Name }}: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.{{ .Name }}[JSON.stringify(params)] ?? {}
		},
				{{ end }}
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		},
		getRegistry: (state) => {
			return state._Registry
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('Vuex module: {{ .Module.Pkg.Name }} initialized!')
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
			state._Subscriptions.forEach(async (subscription) => {
				try {
					const sub=JSON.parse(subscription)
					await dispatch(sub.action, sub.payload)
				}catch(e) {
					throw new Error('Subscriptions: ' + e.message)
				}
			})
		},
		{{ range .Module.HTTPQueries }}
		{{ $FullName := .FullName }}
		{{ $Name := .Name }}
		{{ range $i,$rule := .Rules}} 		
		{{ $n := "" }}
		{{ if (gt $i 0) }}
		{{ $n = inc $i }}
		{{ end}}
		async {{ $FullName }}{{ $n }}({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCase $FullName -}}
				{{- $n -}}(
					{{- range $j,$a :=$rule.Params -}}
						{{- if (gt $j 0) -}}, {{ end }} key.{{ $a -}}
					{{- end -}}
					{{- if $rule.HasQuery -}}
						{{- if $rule.Params -}}, {{ end -}}
						query ?? undefined
					{{- end -}}
					{{- if $rule.HasBody -}}
						{{- if or $rule.HasQuery $rule.Params}},{{ end -}}
							{...key}
						{{- end -}}
					 )).data
				
					{{ if $rule.HasQuery }}
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCase $FullName -}}
					{{- $n -}}(
						{{- range $j,$a :=$rule.Params }} key.{{$a}}, {{ end -}}{...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any
						{{- if $rule.HasBody -}}, {...key}
						{{- end -}}
						)).data
					value = mergeResults(value, next_values);
				}
					{{- end }}
				commit('QUERY', { query: '{{ $Name }}', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: '{{ $FullName }}{{ $n }}', payload: { options: { all }, params: {...key},query }})
				return getters['get{{ $Name }}']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:{{ $FullName }}{{ $n }} API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		{{ end }}
		{{ end }}
		{{ range .Module.Msgs }}async send{{ .Name }}({ rootGetters }, { value, fee = {amount: [], gas: "200000"}, memo = '' }) {
			try {
				const client=await initClient(rootGetters)
				const fullFee = Array.isArray(fee)  ? {amount: fee, gas: "200000"} :fee;
				const result = await client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.tx.send{{ .Name }}({ value, fee: fullFee, memo })
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new Error('TxClient:{{ .Name }}:Init Could not initialize signing client. Wallet is required.')
				}else{
					throw new Error('TxClient:{{ .Name }}:Send Could not broadcast Tx: '+ e.message)
				}
			}
		},
		{{ end }}
		{{ range .Module.Msgs }}async {{ .Name }}({ rootGetters }, { value }) {
			try {
				const client=initClient(rootGetters)
				const msg = await client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.tx.{{ camelCase .Name }}({value})
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new Error('TxClient:{{ .Name }}:Init Could not initialize signing client. Wallet is required.')
				} else{
					throw new Error('TxClient:{{ .Name }}:Create Could not create message: ' + e.message)
				}
			}
		},
		{{ end }}
	}
}