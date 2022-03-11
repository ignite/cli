// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

{{ range .Modules }}export { usePiniaStore as use{{ .Name }}PiniaStore, useModule as use{{ .Name }}Module } from './{{ .Path }}'
{{ end }}