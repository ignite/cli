// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
{{ range .Modules }}export { useModule as use{{ .Name }}Module } from './{{ .Path }}'
{{ end }}
export { default as useIgnite } from './useIgnite'