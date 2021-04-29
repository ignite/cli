export const ensureChainSetupMethod = "ensureChainSetupMethod"
export const connectChainsMethod = "connectChains"
export const getPathMethod = "getPath"
export const listPathsMethod = "listPaths"
export const getDefaultAccountMethod = "getDefaultAccount"
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance"

class EnsureChainSetupResponse {
  id: string
}

export function ensureChainSetup([rpcAddr]: [string]): EnsureChainSetupResponse {
  return new EnsureChainSetupResponse;
}

interface ConnectOptions {
  sourcePort: string
  sourceVersion: string
  targetPort: string
  targetVersion: string
  ordering: string
}

class ConnectChainsResponse {
  id: string
}

export function connectChains([srcID, dstID, options]: [string, string, ConnectOptions]): ConnectChainsResponse {
  return new ConnectChainsResponse;
}

class Path {
  id: string
  isLinked: boolean
  src: PathEnd
  dst: PathEnd
}

class PathEnd {
  channelID: string
  chainID: string
  portID: string
}

export function getPath([id]: [string]): Path {
  return new Path;
}

export function listPaths(): Path[] {
  return [new Path];
}

class Account {
  address: string
}

export function getDefaultAccount([chainID]: [string]): Account {
  return new Account;
}

class Coin {
  denom: string
  amount: number
}

export function getDefaultAccountBalance([chainID]: [string]): Coin[] {
  return [new Coin];
}
