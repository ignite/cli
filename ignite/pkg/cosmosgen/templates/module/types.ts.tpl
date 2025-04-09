{{ range .Module.Types }}import { {{ .Name }} } from "{{ resolveFile .FilePath }}"
{{ end }}

export {     
    {{ range .Module.Types }}{{ .Name }},
    {{ end }}
 }