// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.
{{ range .Modules }}export { useModule as use{{ camelCaseUpperSta .Pkg.Name }} } from './{{ .Pkg.Name }}'
{{ end }}
export { default as use{{ capitalCase .Repo }} } from './use'