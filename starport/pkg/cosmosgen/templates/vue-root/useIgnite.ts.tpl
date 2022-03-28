// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { OfflineDirectSigner } from '@cosmjs/proto-signing'
import { SigningStargateClient } from '@cosmjs/stargate'
import { createIgnite, Ignite, registry, Environment } from '@ignt/client'
import { reactive, Ref, toRefs, ToRefs } from 'vue'

type State = {
    ignite: Ref<Ignite>
}

type Response = {
    state: ToRefs<State>
    signIn: (offlineSigner: OfflineDirectSigner) => void
    signOut: () => void
}

type Params = {
    env: Environment
    autoConnectWS: boolean
}

// singleton state
const state = reactive({} as State)

export default function (
    p: Params = {
        env: {
            apiURL: 'http://localhost:1317',
            rpcURL: 'http://localhost:26657',
            wsURL: 'ws://localhost:26657/websocket',
            prefix: 'cosmos'
        },
        autoConnectWS: true
    }
): Response {
    if (!state.ignite) {
        if (!p.env) {
            throw new Error('Ignite: Unable to create instance without env')
        }

        state.ignite = createIgnite({
            env: p.env
        }) as Ignite

        if (p.autoConnectWS) {
            state.ignite.connectWS()
        }
    }

    let signIn = async (offlineSigner: OfflineDirectSigner) => {
        let [acc] = await offlineSigner.getAccounts()

        let stargateClient = await SigningStargateClient.connectWithSigner(
            state.ignite.env.rpcURL,
            offlineSigner,
            { registry }
        )

        state.ignite.signIn(stargateClient, acc.address)
    }

    let signOut = () => {
        state.ignite.signOut()
    }

    return {
        state: toRefs(state),
        signIn,
        signOut
    }
}
