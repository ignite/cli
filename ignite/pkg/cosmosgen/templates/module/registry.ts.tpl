import { GeneratedType } from "@cosmjs/proto-signing";
{{ range .Module.Msgs }}import { {{ .Name }} } from "{{ resolveFile .FilePath }}";
{{ end }}
const msgTypes: Array<[string, GeneratedType]>  = [
    {{ range .Module.Msgs }}["/{{ .URI }}", {{ .Name }}],
    {{ end }}
];

export { msgTypes }