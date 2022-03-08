// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

{{ range .Modules }}export { default as use{{ .FullName }}PiniaStore } from './{{ .Path }}'
{{ end }}