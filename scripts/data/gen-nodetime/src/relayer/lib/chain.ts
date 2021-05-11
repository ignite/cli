export const ensureChainSetupMethod = "ensureChainSetup"
export const createPathMethod = "createPath"
export const getPathMethod = "getPath"
export const listPathsMethod = "listPaths"
export const getDefaultAccountMethod = "getDefaultAccount"
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance"

interface ChainSetupOptions {
  gasPrice: string
}

class EnsureChainSetupResponse {
  // id is the chain id of chain.
  id: string
}

// ensureChainSetup sets up a chain by its rpc address only if it is not set up already.
export function ensureChainSetup([rpcAddr, options]: [string, ChainSetupOptions]): EnsureChainSetupResponse {
  throw new Error("ensureChainSetup() not implemented");
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
// and should not send any txs to these chains. this will later be done by link().
export function createPath([srcID, dstID, options]: [string, string, ConnectOptions]): Path {
  throw new Error("createPath() not implemented");
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
  throw new Error("getPath() not implemented");
}

// listPaths list all connections.
export async function listPaths(): Promise<Path[]> {
  throw new Error("listPaths() not implemented");
}

class Account {
  address: string
}

// getDefaultAccount gets the default account on chain by chain id.
export function getDefaultAccount([chainID]: [string]): Account {
  throw new Error("getDefaultAccount() not implemented");
}

class Coin {
  denom: string
  amount: number
}

// getDefaultAccountBalance gets the balance of default account on chain by chain id.
export function getDefaultAccountBalance([chainID]: [string]): Coin[] {
  throw new Error("getDefaultAccountBalance() not implemented");
}
