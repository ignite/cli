import { SigningStargateClient } from '@cosmjs/stargate'
import { OfflineDirectSigner } from '@cosmjs/proto-signing'
import {
    AppCurrency,
    Bech32Config,
    BIP44,
    Currency,
    Key
} from '@keplr-wallet/types'
import { V1Beta1QueryParamsResponse } from './cosmos.staking.v1beta1/rest'
import { V1Beta1QueryTotalSupplyResponse } from './cosmos.bank.v1beta1/rest'
import axios from 'axios'
import ReconnectingWebSocket from 'reconnecting-websocket'
import EventEmitter from 'eventemitter3'

interface Env {
    chainID?: string
    chainName?: string
    apiURL: string
    rpcURL: string
    wsURL: string
    prefix?: string
    status?: {
        apiConnected?: boolean
        rpcConnected?: boolean
        wsConnected?: boolean
    }
}
function plugEnv(initial: Env): {
    env: Env
} {
    return {
        env: {
            ...initial
        }
    }
}

function plugSigner(): {
    signer: {
        client: SigningStargateClient
        addr: string
    }
} {
    return {
        signer: {
            client: undefined,
            addr: undefined
        }
    }
}

function plugKeplr(): {
    keplr: {
        connect: (
            onSuccessCb: () => void,
            onErrorCb: () => void,
            params: {
                stakinParams: V1Beta1QueryParamsResponse
                tokens: V1Beta1QueryTotalSupplyResponse
                env: Env
            }
        ) => void
        isAvailable: () => boolean
        getOfflineSigner: (chainId: string) => OfflineDirectSigner
        getAccParams: (chainId: string) => Promise<Key>
        listenToAccChange: (cb: EventListener) => void
    }
} {
    return {
        keplr: {
            connect: async (onSuccessCb: () => void, onErrorCb: () => void, p) => {
                try {
                    let staking = p.stakinParams
                    let tokens = p.tokens

                    let chainId: string = p.env.chainID || 'ee'
                    let chainName: string = p.env.chainName || 'Ee'
                    let rpc: string = p.env.rpcURL || ''
                    let rest: string = p.env.apiURL || ''
                    let addrPrefix: string = p.env.prefix || ''

                    let stakeCurrency: Currency = {
                        coinDenom: staking.params?.bond_denom?.toUpperCase() || '',
                        coinMinimalDenom: staking.params?.bond_denom || '',
                        coinDecimals: 0
                    }

                    let bip44: BIP44 = {
                        coinType: 118
                    }

                    let bech32Config: Bech32Config = {
                        bech32PrefixAccAddr: addrPrefix,
                        bech32PrefixAccPub: addrPrefix + 'pub',
                        bech32PrefixValAddr: addrPrefix + 'valoper',
                        bech32PrefixValPub: addrPrefix + 'valoperpub',
                        bech32PrefixConsAddr: addrPrefix + 'valcons',
                        bech32PrefixConsPub: addrPrefix + 'valconspub'
                    }

                    let currencies: AppCurrency[] = tokens.supply?.map((c) => {
                        const y: any = {
                            amount: '0',
                            denom: '',
                            coinDenom: '',
                            coinMinimalDenom: '',
                            coinDecimals: 0
                        }
                        y.coinDenom = c.denom?.toUpperCase() || ''
                        y.coinMinimalDenom = c.denom || ''
                        y.coinDecimals = 0

                        return y
                    }) as AppCurrency[]

                    let feeCurrencies: AppCurrency[] = tokens.supply?.map((c) => {
                        const y: any = {
                            amount: '0',
                            denom: '',
                            coinDenom: '',
                            coinMinimalDenom: '',
                            coinDecimals: 0
                        }
                        y.coinDenom = c.denom?.toUpperCase() || ''
                        y.coinMinimalDenom = c.denom || ''
                        y.coinDecimals = 0

                        return y
                    }) as AppCurrency[]

                    let coinType = 118

                    let gasPriceStep = {
                        low: 0.01,
                        average: 0.025,
                        high: 0.04
                    }

                    if (chainId) {
                        await window.keplr.experimentalSuggestChain({
                            chainId,
                            chainName,
                            rpc,
                            rest,
                            stakeCurrency,
                            bip44,
                            bech32Config,
                            currencies,
                            feeCurrencies,
                            coinType,
                            gasPriceStep
                        })

                        window.keplr.defaultOptions = {
                            sign: {
                                preferNoSetFee: true,
                                preferNoSetMemo: true
                            }
                        }

                        await window.keplr.enable(chainId)
                        onSuccessCb()
                    } else {
                        console.error('Cannot access chain data')
                        onErrorCb()
                    }
                } catch (e) {
                    console.error(e)
                    onErrorCb()
                }
            },

            isAvailable: () => {
                // @ts-ignore
                return !!window.keplr
            },

            getOfflineSigner: (chainId: string) =>
                // @ts-ignore
                window.keplr.getOfflineSigner(chainId),

            getAccParams: async (chainId: string) =>
                // @ts-ignore
                await window.keplr.getKey(chainId),

            listenToAccChange: (cb: EventListener) => {
                window.addEventListener('keplr_keystorechange', cb)
            }
        }
    }
}

function plugWebsocket(env: Env): {
    ws: {
        ee: () => EventEmitter
        close: () => void
        connect: () => void
    }
} {
    let _refresh: number = 5000
    let _socket: ReconnectingWebSocket
    let _timer: number
    let _ee = new EventEmitter()

    let ping = async () => {
        if (env.apiURL) {
            try {
                const status: any = await axios.get(env.apiURL + '/node_info')
                _ee.emit('ws-newblock', 'ws-chain-id', status.data.node_info.network)

                status.data.application_version.name
                    ? _ee.emit('ws-chain-name', status.data.application_version.name)
                    : _ee.emit('ws-chain-name', status.data.node_info.network)

                _ee.emit('ws-api-status')
            } catch (error) {
                if (!error.response) {
                    _ee.emit('ws-api-status')
                    console.error(new Error('WebSocketClient: API Node unavailable'))
                } else {
                    _ee.emit('ws-api-status')
                }
            }
        }
        if (env.rpcURL) {
            try {
                await axios.get(env.rpcURL)
                _ee.emit('ws-rpc-status')
            } catch (error) {
                if (!error.response) {
                    console.error(new Error('WebSocketClient: RPC Node unavailable'))
                    _ee.emit('ws-rpc-status')
                } else {
                    _ee.emit('ws-rpc-status')
                }
            }
        }
    }

    let onError = () => {
        console.error(new Error('WebSocketClient: Could not connect to websocket.'))
    }

    let onClose = () => {
        _ee.emit('ws-close')
    }

    let onOpen = () => {
        _ee.emit('ws-open')
        _socket.send(
            JSON.stringify({
                jsonrpc: '2.0',
                method: 'subscribe',
                id: '1',
                params: ["tm.event = 'NewBlock'"]
            })
        )
    }

    let onMessage = (msg) => {
        const result = JSON.parse(msg.data).result
        if (result.data && result.data.type === 'tendermint/event/NewBlock') {
            _ee.emit('ws-newblock', JSON.parse(msg.data).result)
        }
    }

    return {
        ws: {
            ee: () => _ee,

            close: () => {
                clearInterval(_timer)
                _timer = undefined
                _socket.close()
            },

            connect: () => {
                _socket = new ReconnectingWebSocket(env.wsURL)
                ping()
                _timer = setInterval(() => ping(), _refresh)
                _socket.onopen = onOpen.bind(this)
                _socket.onmessage = onMessage.bind(this)
                _socket.onerror = onError.bind(this)
                _socket.onclose = onClose.bind(this)
            }
        }
    }
}

export { plugEnv, plugWebsocket, plugSigner, plugKeplr }
