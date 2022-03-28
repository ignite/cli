// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { Registry } from '@cosmjs/proto-signing'
import { SigningStargateClient } from '@cosmjs/stargate'
import axios from 'axios'
import { EventEmitter } from 'events'
import ReconnectingWebSocket from 'reconnecting-websocket'
{{ range .Modules }}import { Module as {{ .Name }}, msgTypes as {{ .Name }}MsgTypes } from './{{ .Path }}'
{{ end }}

const registry = new Registry([
  {{ range .Modules }}...{{ .Name }}MsgTypes,
  {{ end }}
])

interface Environment {
  chainID?: string;
  chainName?: string;
  apiURL: string;
  rpcURL: string;
  wsURL: string;
  prefix?: string
}

interface EnvironmentStatus {
  apiConnected: boolean
  rpcConnected: boolean
  wsConnected: boolean
}

interface IgniteParams {
  env: Environment;
}

class Ignite {
    private _env: Environment
    private _envStatus: EnvironmentStatus = {
        apiConnected: false,
        rpcConnected: false,
        wsConnected: false
    }
    private _ws?: WebSocketClient
    private _signer?: SigningStargateClient
    private _addr?: string

    {{ range .Modules }}public {{ .Name }}: {{ .Name }};
    {{ end }}

    constructor({ env }: IgniteParams) {
        this._env = env;
        {{ range .Modules }}this.{{ .Name }} = new {{ .Name }}(
         this._env.apiURL
        );
        {{ end }}
    }

    public signIn(client: SigningStargateClient, addr: string): Ignite {
        this._signer = client
        this._addr = addr

        {{ range .Modules }}this.{{ .Name }}.withSigner(client, addr)
        {{ end }}

        return this
    }

    public signOut() {
        this._signer = undefined
        this._addr = undefined
    }

    public connectWS(ws?: WebSocketClient): Ignite {
        if (this._ws) {
            this._ws.close()
        }

        if (ws) {
            this._ws = ws
        } else {
            this._ws = new WebSocketClient({ env: this._env })
        }

        // @ts-ignore
        this._ws.on('chain-id', (id) => {
            this._env.chainID = id
        })
        // @ts-ignore
        this._ws.on('chain-name', (name) => {
            this._env.chainName = name
        })
        // @ts-ignore
        this._ws.on('api-status', (connected) => {
            this._envStatus.apiConnected = connected
        })
        // @ts-ignore
        this._ws.on('rpc-status', (connected) => {
            this._envStatus.rpcConnected = connected
        })
        // @ts-ignore
        this._ws.on('ws-status', (connected) => {
            this._envStatus.wsConnected = connected
        })

        this._ws.connect()

        return this
    }

    get ws(): WebSocketClient | undefined {
        return this._ws
    }

    get env(): Environment {
        return this._env
    }

    get envStatus(): EnvironmentStatus {
        return this._envStatus
    }

    get signer(): SigningStargateClient | undefined {
        return this._signer
    }

    get addr(): string | undefined {
        return this._addr
    }
}

interface WebSocketParams {
  refresh?: number
  env: Environment
}

class WebSocketClient extends EventEmitter {
    private _refresh: number
    private _env: Environment
    private _socket: ReconnectingWebSocket
    private _timer: number

    constructor({ env, refresh = 5000 }: WebSocketParams) {
        super()
        this._env = env
        this._refresh = refresh
    }

    public close() {
        clearInterval(this._timer)
        this._timer = undefined
        this._socket.close()
    }

    public connect() {
        this._socket = new ReconnectingWebSocket(this._env.wsURL)
        this.sendBeacon()
        this._timer = setInterval(() => this.sendBeacon(), this._refresh)

        this._socket.onopen = this.onOpen.bind(this)
        this._socket.onmessage = this.onMessage.bind(this)
        this._socket.onerror = this.onError.bind(this)
        this._socket.onclose = this.onClose.bind(this)
    }

    private async sendBeacon(): Promise<void> {
        if (this._env.apiURL) {
            try {
                const status: any = await axios.get(this._env.apiURL + '/node_info')
                this.emit('chain-id', status.data.node_info.network)
                status.data.application_version.name
                    ? this.emit('chain-name', status.data.application_version.name)
                    : this.emit('chain-name', status.data.node_info.network)
                this.emit('api-status', true)
            } catch (error) {
                if (!error.response) {
                    this.emit('api-status', false)
                    console.error(new Error('WebSocketClient: API Node unavailable'))
                } else {
                    this.emit('api-status', true)
                }
            }
        }
        if (this._env.rpcURL) {
            try {
                await axios.get(this._env.rpcURL)
                this.emit('rpc-status', true)
            } catch (error) {
                if (!error.response) {
                    console.error(new Error('WebSocketClient: RPC Node unavailable'))
                    this.emit('rpc-status', false)
                } else {
                    this.emit('rpc-status', true)
                }
            }
        }
    }

    private onError(): void {
        console.error(new Error('WebSocketClient: Could not connect to websocket.'))
    }

    private onClose(): void {
        this.emit('ws-status', false)
    }

    private onOpen(): void {
        this.emit('ws-status', true)
        this._socket.send(
            JSON.stringify({
                jsonrpc: '2.0',
                method: 'subscribe',
                id: '1',
                params: ["tm.event = 'NewBlock'"]
            })
        )
    }

    private onMessage(msg): void {
        const result = JSON.parse(msg.data).result
        if (result.data && result.data.type === 'tendermint/event/NewBlock') {
            this.emit('newblock', JSON.parse(msg.data).result)
        }
    }
}

function createIgnite({ env }: IgniteParams): Ignite {
    return new Ignite({
        env
    })
}

export { createIgnite, Environment, Ignite, registry, WebSocketClient }
