import useSWRV from 'swrv'
import { useClient } from '../useClient';
import type { Ref } from 'vue'

type SwrvReturn<T> = {
  data: Ref<T>;
  error: Ref<unknown>;
};
type SwrvHelper<T extends (...args: any) => any> = SwrvReturn<
  Awaited<ReturnType<T>>["data"]
>;

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
    const key = JSON.stringify({ type: '{{ $FullName }}{{ $n }}', {{ range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{ $a -}}
          {{- end -}} {{- if $rule.HasQuery -}}
            {{- if $rule.Params -}}, {{ end -}}
            query
          {{- end }} });
    return useSWRV(key, ({{- if or $rule.HasQuery $rule.Params -}}key{{- else -}}_{{- end -}}) => {
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
    }) as SwrvHelper<typeof client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCaseSta $FullName -}}
        {{- $n -}}>
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