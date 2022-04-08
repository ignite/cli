// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import Module from "../../client/{{ .Module.Pkg.Name }}/module";
import useIgnite from '../useIgnite'
import { unref } from 'vue'
		
{{ range .Module.Msgs }}type Send{{ .Name }}Type = typeof Module.prototype.send{{ .Name }}
{{ end }}
{{ range .Module.HTTPQueries }}type {{ .FullName }}Type = typeof Module.prototype.{{ camelCase .FullName }}
{{ end }}

type Response = {
  {{ range .Module.Msgs }}send{{ .Name }}: Send{{ .Name }}Type,
  {{ end }}
  {{ range .Module.HTTPQueries }}{{ camelCase .FullName }}: {{ .FullName }}Type
  {{ end }}
}

function useModule(): Response {
  // ignite
  let { ignite } = useIgnite()

  let {
	{{ range .Module.Msgs }}
	send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  } = unref(ignite.{{ camelCase .Module.Name }})

  {{ $ModuleName := camelCase .Module.Name }}
  {{ range .Module.Msgs }}
	send{{ .Name }} = send{{ .Name }}.bind(ignite.{{ $ModuleName }})
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }} = {{ camelCase .FullName }}.bind(ignite.{{ $ModuleName }})
  {{ end }}

  return {
  {{ range .Module.Msgs }}
  send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  }
}

export { useModule }
