/* eslint-disable @typescript-eslint/no-unused-vars */
import { useQuery, type UseQueryOptions, useInfiniteQuery, type UseInfiniteQueryOptions } from "@tanstack/{{- .FrontendType -}}-query";
import { useClient } from '../useClient';

export default function use{{ camelCaseUpperSta $.Module.Pkg.Name }}() {
  const client = useClient();

  {{- range .Module.HTTPQueries -}}
    {{- $FullName := .FullName -}}
    {{- $Name := .Name -}}
    {{- if .Paginated -}}
    {{- range $i,$rule := .Rules -}}
      {{- $n := "" -}}
      {{- if (gt $i 0) -}}
        {{- $n = inc $i -}}
      {{- end }}
  const {{ $FullName }}{{ $n }} = (
      {{- if $rule.Params -}}
      {{- range $j,$a :=$rule.Params -}}
        {{- if (gt $j 0) -}}, {{ end -}}
        {{- $a -}}: string
      {{- end -}}
      , {{ end -}}
      {{- if $rule.HasQuery -}}        
        query: any, 
      {{- end }} options: any, perPage: number) => {
    const key = { type: '{{ $FullName }}{{ $n }}', {{ range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{ $a -}}
          {{- end -}} {{- if $rule.HasQuery -}}
            {{- if $rule.Params -}}, {{ end -}}
            query
          {{- end }} };    
    return useInfiniteQuery([key], ({pageParam = 1}: { pageParam?: number}) => {
       {{- if or $rule.HasQuery $rule.Params}}
      const { {{- if $rule.Params -}}{{- range $j,$a :=$rule.Params -}}
            {{- if (gt $j 0) -}}, {{ end }} {{ $a -}}
          {{- end -}}{{- if $rule.HasQuery -}}, {{- end -}}{{- end -}}{{- if $rule.HasQuery -}}query{{- end }} } = key{{ end }}

      query['pagination.limit']=perPage;
      query['pagination.offset']= (pageParam-1)*perPage;
      query['pagination.count_total']= true;
      return  client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCase $FullName -}}
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
            ).then( res => ({...res.data,pageParam}) );
    }, {...options,
      getNextPageParam: (lastPage, allPages) => { if ((lastPage.pagination?.total ?? 0) >((lastPage.pageParam ?? 0) * perPage)) {return lastPage.pageParam+1 } else {return undefined}},
      getPreviousPageParam: (firstPage, allPages) => { if (firstPage.pageParam==1) { return undefined } else { return firstPage.pageParam-1}}
    }
    );
  }
  {{ end -}}
  {{- else -}}
    {{- range $i,$rule := .Rules -}}
      {{- $n := "" -}}
      {{- if (gt $i 0) -}}
        {{- $n = inc $i -}}
      {{- end }}
  const {{ $FullName }}{{ $n }} = (
      {{- if $rule.Params -}}
      {{- range $j,$a :=$rule.Params -}}
        {{- if (gt $j 0) -}}, {{ end -}}
        {{- $a -}}: string
      {{- end -}}
      , {{ end -}}
      {{- if $rule.HasQuery -}}        
        query: any, 
      {{- end }} options: any) => {
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
      return  client.{{ camelCaseUpperSta $.Module.Pkg.Name }}.query.{{ camelCase $FullName -}}
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
    }, options);
  }
  {{ end -}}
  {{- end -}}
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
