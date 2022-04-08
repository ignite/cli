// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { OfflineDirectSigner } from '@cosmjs/proto-signing'
import { SigningStargateClient } from '@cosmjs/stargate'
import { Ignite, registry } from '../ts-client'
import { toRefs, ToRefs, reactive, UnwrapNestedRefs } from 'vue'

type State = UnwrapNestedRefs<Ignite>

type Response = {
    ignite: ToRefs<Ignite>
    signIn: (offlineSigner: OfflineDirectSigner) => void
    signOut: () => void
    inject: (instance: Ignite) => void
}

let _igniteGlobal: State

export default function (): Response {
    let signIn = async (offlineSigner: OfflineDirectSigner) => {
        let [acc] = await offlineSigner.getAccounts()

        let stargateClient = await SigningStargateClient.connectWithSigner(
            _igniteGlobal.env.rpcURL,
            offlineSigner,
            { registry }
        )
        let addr = acc.address

        _igniteGlobal.signer.client = stargateClient
        _igniteGlobal.signer.addr = addr

        {{ range .Modules }}
            _igniteGlobal.{{ camelCase .Name }}.withSigner(stargateClient, addr)
        {{ end }}
    }

    let signOut = () => {
        _igniteGlobal.signer.client = undefined
        _igniteGlobal.signer.addr = undefined


        {{ range .Modules }}
            _igniteGlobal.{{ camelCase .Name }}.noSigner()
        {{ end }}
    }

    let inject = (instance: Ignite) => {
        _igniteGlobal = reactive<Ignite>(instance)
    }

    return {
        inject,
        ignite: toRefs(_igniteGlobal as Ignite),
        signIn,
        signOut
    }
}
