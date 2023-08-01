{{ range .Module.Types }}import { {{ .Name }} } from "./types/{{ resolveFile .FilePath }}"
{{ end }}

export {     
    {{ range .Module.Types }}{{ .Name }},
    {{ end }}
 }