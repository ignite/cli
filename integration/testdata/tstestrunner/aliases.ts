import path from 'path'
import glob from 'glob'

// Alias type stores package name to path mappings
type Alias = {[key: string]: string}

// Collect aliases for all packages defined within a blockchain app
const collect = (chainPath: string): Alias => {
  const locations = [
    `${chainPath}/vue/node_modules`,
    `${chainPath}/vue/src/store/generated`,
  ]

  const options = {
    dot: false,
    nosort: true,
    nodir: true,
    ignore: [
      '**/test/**',
      '**/node_modules/**/node_modules/**',
    ],
  }

  const alias: Alias = {}

  // Search packages inside the locations and add them to the package aliases
  locations.forEach((location: string) => {
    glob.sync(`${location}/**/package.json`, options).forEach((match: string) => {
      const modulePath = path.dirname(match)

      if (location === modulePath) {
        return
      }

      let pkg: {name: string}

      try {
        pkg = require(match)
      } catch {
        return
      }

      if (pkg.name) {
        alias[pkg.name] = modulePath
      }
    })
  })

  return alias
}

export default { collect }
