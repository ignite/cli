// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { Ignite } from "@ignt/client";
import Module from "@ignt/client/{{ .Module.Pkg.Name }}/module";
		
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

type Params = {
  ignite: Ignite;
}

function useModule({ ignite }: Params): Response {
  let {
	{{ range .Module.Msgs }}
	send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  } = ignite.{{ .Module.Name }}

  {{ $ModuleName := .Module.Name }}
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
