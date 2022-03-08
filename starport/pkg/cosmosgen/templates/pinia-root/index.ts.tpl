// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

{{ range .Modules }}export { default as use{{ .Name }}PiniaStore } from './{{ .Path }}'
{{ end }}