export const linkMethod = "link"

interface Response {
  // linkedPaths is a list of paths that are linked when link() called.
  linkedPaths: string[]

  // alreadyLinkedPaths is a list of paths that already linked on chain.
  alreadyLinkedPaths: string[]
}

// link connects src and dst chains by their paths on chain with ibc txs.
export function link(paths: string[]): Response {
  throw new Error("link() not implemented");
}
