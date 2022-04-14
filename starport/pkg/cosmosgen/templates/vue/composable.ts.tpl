// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
import { unref } from 'vue'
import Module from "{{ .User }}-{{ .Repo }}-ts-client/{{ .Module.Pkg.Name }}/module";
import use{{ capitalCase .Repo }} from '../use'
		
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
  let { {{ .Repo }} } = use{{ capitalCase .Repo }}()

  let {
	{{ range .Module.Msgs }}
	send{{ .Name }},
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }},
  {{ end }}
  } = unref({{ .Repo }}.{{ camelCaseLowerSta .Module.Pkg.Name }})

  {{ $ModuleName := camelCaseLowerSta .Module.Pkg.Name }}
  {{ $Repo := .Repo }}
  {{ range .Module.Msgs }}
	send{{ .Name }} = send{{ .Name }}.bind({{ $Repo }}.{{ $ModuleName }})
  {{ end }}
  {{ range .Module.HTTPQueries }}
  {{ camelCase .FullName }} = {{ camelCase .FullName }}.bind({{ $Repo }}.{{ $ModuleName }})
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
