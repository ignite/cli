import { OfflineSigner } from "@cosmjs/proto-signing";

export interface Env {
  chainID?: string
  signer?: OfflineSigner
  chainName?: string
  apiURL: string
  rpcURL: string
  wsURL: string
  prefix?: string
  status?: {
    apiConnected: boolean
    rpcConnected: boolean
    wsConnected: boolean
  }
}