export const ensureChainSetupMethod = "ensureChainSetup"
export const createPathMethod = "createPath"
export const getPathMethod = "getPath"
export const listPathsMethod = "listPaths"
export const getDefaultAccountMethod = "getDefaultAccount"
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance"

class EnsureChainSetupResponse {
  // id is the chain id of chain.
  id: string
}

// ensureChainSetup sets up a chain by its rpc address only if it is not set up already.
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

// createPath creates a path between the source chain and dest chain by their chain ids with given options.
// it returns a unique path id that represents the connection between these chains.
//
// createPath should only record the intention of connecting source and destion chains together
// and should not send any txs to these chains. this will later be done by linkChains().
export function createPath([srcID, dstID, options]: [string, string, ConnectOptions]): Path {
  return new Path;
}

// Path represents the connection between two chaons.
class Path {
  // id of the path.
  id: string

  // isLinked shows whether src and dst chains are connected on the chain with ibc txs.
  isLinked: boolean

  // src represents the source chain.
  src: PathEnd

  // dst represents the destionation chain.
  dst: PathEnd
}

// PathEnd represents a chain.
class PathEnd {
  channelID: string
  chainID: string
  portID: string
}

// getPath gets connection info between chains by path id.
export function getPath([id]: [string]): Path {
  return new Path;
}

// listPaths list all connections.
export function listPaths(): Path[] {
  return [new Path];
}

class Account {
  address: string
}

// getDefaultAccount gets the default account on chain by chain id.
export function getDefaultAccount([chainID]: [string]): Account {
  return new Account;
}

class Coin {
  denom: string
  amount: number
}

// getDefaultAccountBalance gets the balance of default account on chain by chain id.
export function getDefaultAccountBalance([chainID]: [string]): Coin[] {
  return [new Coin];
}
