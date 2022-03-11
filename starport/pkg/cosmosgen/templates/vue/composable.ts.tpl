// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { Store } from 'pinia'
import { usePiniaStore, PiniaState } from './'

import { Ignite } from "ts-client";
import Module from "ts-client/{{ .Module.Pkg.Name }}/module";
		
{{ range .Module.Msgs }}type Send{{ .Name }}Type = typeof Module.prototype.send{{ .Name }}
{{ end }}
{{ range .Module.HTTPQueries }}type {{ .FullName }}Type = typeof Module.prototype.{{ camelCase .FullName }}
{{ end }}

type Response = {
  $s: Store<'{{ .Module.Pkg.Name }}', PiniaState, {}, {}>
  {{ range .Module.Msgs }}send{{ .Name }}: Send{{ .Name }}Type,
  {{ end }}
  {{ range .Module.HTTPQueries }}{{ camelCase .FullName }}: {{ .FullName }}Type
  {{ end }}
}

type Params = {
  $ignt: Ignite;
}

function useModule({ $ignt }: Params): Response {
  let $s = usePiniaStore()

  let {
	{{ range .Module.Msgs }}
	send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  } = $ignt.{{ .Module.Name }} 

  return {
	$s,
  {{ range .Module.Msgs }}
  send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  }
}

export { useModule }
