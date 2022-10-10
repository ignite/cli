/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery } from "@tanstack/vue-query";
import { useClient } from '../useClient';
import type { Ref } from 'vue'

export default function use{{ camelCaseUpperSta $.Module.Pkg.Name }}() {
  const client = useClient();

  {{- range .Module.HTTPQueries -}}
    {{- $FullName := .FullName -}}
    {{- $Name := .Name -}}
    {{- range $i,$rule := .Rules -}}
      {{- $n := "" -}}
      {{- if (gt $i 0) -}}
        {{- $n = inc $i -}}
      {{- end }}
  const {{ $FullName }}{{ $n }} = (
      {{- range $j,$a :=$rule.Params -}}
        {{- if (gt $j 0) -}}, {{ end -}}
        {{- $a -}}: string
      {{- end -}}
      {{- if $rule.HasQuery -}}
        {{- if $rule.Params -}}, {{ end -}}
        query?: unknown 
      {{- end }}) => {
    const key = { type: '{{ $FullName }}{{ $n }}', {{ range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{ $a -}}
          {{- end -}} {{- if $rule.HasQuery -}}
            {{- if $rule.Params -}}, {{ end -}}
            query
          {{- end }} };    
    return useQuery([key], () => {
       {{- if or $rule.HasQuery $rule.Params}}
      const { {{- if $rule.Params -}}{{- range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{ $a -}}
          {{- end -}}{{- if $rule.HasQuery -}}, {{- end -}}{{- end -}}{{- if $rule.HasQuery -}}query{{- end }} } = key{{ end }}
      return  client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCaseSta $FullName -}}
        {{- $n -}}({{- range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{- $a -}}
          {{- end -}}
          {{- if $rule.HasQuery -}}
            {{- if $rule.Params -}}, {{ end -}}
            query ?? undefined
          {{- end -}}
          {{- if $rule.HasBody -}}
            {{- if or $rule.HasQuery $rule.Params}},{{ end -}}
              {...key}
            {{- end -}}
            ).then( res => res.data );
    })
  }
  {{ end -}}
  {{- end }}
  return {
  {{- range .Module.HTTPQueries -}}
  {{- $FullName := .FullName -}}
  {{- $Name := .Name -}}
  {{- range $i,$rule := .Rules -}}
  {{- $n := "" -}}
  {{- if (gt $i 0) -}}
  {{- $n = inc $i -}}
  {{- end -}}
    {{ $FullName }}{{ $n }},
  {{- end -}}
  {{- end }}
  }
}