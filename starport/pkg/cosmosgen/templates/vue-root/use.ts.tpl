// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { OfflineDirectSigner } from '@cosmjs/proto-signing'
import { SigningStargateClient } from '@cosmjs/stargate'
import { {{ capitalCase .Repo }}, registry } from '{{ .User }}-{{ .Repo }}-ts-client'
import { toRefs, ToRefs, reactive, UnwrapNestedRefs } from 'vue'

type State = UnwrapNestedRefs<{{ capitalCase .Repo }}>

type Response = {
    {{ .Repo }}: ToRefs<{{ capitalCase .Repo }}>
    signIn: (offlineSigner: OfflineDirectSigner) => void
    signOut: () => void
    inject: (instance: {{ capitalCase .Repo }}) => void
}

let _{{ .Repo }}Global: State

export default function (): Response {
    let signIn = async (offlineSigner: OfflineDirectSigner) => {
        let [acc] = await offlineSigner.getAccounts()

        let stargateClient = await SigningStargateClient.connectWithSigner(
            _{{ .Repo }}Global.env.rpcURL,
            offlineSigner,
            { registry }
        )
        let addr = acc.address

        _{{ .Repo }}Global.signer.client = stargateClient
        _{{ .Repo }}Global.signer.addr = addr

        {{ $Repo := .Repo }}
        {{ range .Modules }}
        _{{ $Repo }}Global.{{ camelCaseLowerSta .Pkg.Name }}.withSigner(stargateClient, addr)
        {{ end }}
    }

    let signOut = () => {
        _{{ .Repo }}Global.signer.client = undefined
        _{{ .Repo }}Global.signer.addr = undefined

        {{ $Repo := .Repo }}
        {{ range .Modules }}
        _{{ $Repo }}Global.{{ camelCaseLowerSta .Pkg.Name }}.noSigner()
        {{ end }}
    }

    let inject = (instance: {{ capitalCase .Repo }}) => {
        _{{ .Repo }}Global = reactive<{{ capitalCase .Repo }}>(instance)
    }

    return {
        inject,
        {{ .Repo }}: toRefs(_{{ .Repo }}Global as {{ capitalCase .Repo }}),
        signIn,
        signOut
    }
}
