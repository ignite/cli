// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
{{ range .Modules }}export { useModule as use{{ .Name }} } from './{{ .Path }}'
{{ end }}
export { default as useIgnite } from './useIgnite'