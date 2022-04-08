// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { Registry } from '@cosmjs/proto-signing'

import { plugEnv, plugKeplr, plugSigner, plugWebsocket } from './plugins'

{{ range .Modules }}import { Module as {{ .Name }}, msgTypes as {{ .Name }}MsgTypes } from './{{ .Path }}'
{{ end }}

const registry = new Registry([
  {{ range .Modules }}...{{ .Name }}MsgTypes,
  {{ end }}
])

type Environment = ReturnType<typeof plugEnv>['env']

{{ range .Modules }}
function plug{{ .Name }}(env: Environment): {
    {{ camelCase .Name }}: {{ .Name }}
} {
    return {
        {{ camelCase .Name }}: new {{ .Name }}(env.apiURL)
    }
}
{{ end }}

type Ignite = {{ range .Modules }}ReturnType<typeof plug{{ .Name }}> & {{ end }}

ReturnType<typeof plugSigner> &
ReturnType<typeof plugKeplr> &
ReturnType<typeof plugWebsocket> &
ReturnType<typeof plugEnv>

let createIgnite = (p: { env: Environment }) => _use(
{
    {{ range .Modules }}
        ...plug{{ .Name }}(p.env),
    {{ end }}
        ...plugSigner(),

        ...plugKeplr(),

        ...plugWebsocket(p.env),

        ...plugEnv(p.env)

    }
)

function _use<T>(elements: T): { [K in keyof T]: T[K] } {
    return Object.assign({}, elements)
}

export * from './plugins'
export {
    {{ range .Modules }}
        plug{{ .Name }},
    {{ end }}
    createIgnite,
    Environment,
    registry,
    Ignite,
    _use
}
