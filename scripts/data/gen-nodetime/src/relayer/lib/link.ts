export const linkMethod = "link"

interface Response {
  linkedPaths: string[]
  alreadyLinkedPaths: string[]
}

export default function(paths: string[]): Response {
  return {
    linkedPaths: [],
    alreadyLinkedPaths: [],
  }
}
