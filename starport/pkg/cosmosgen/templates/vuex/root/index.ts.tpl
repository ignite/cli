// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { Registry, GeneratedType } from "@cosmjs/proto-signing";
{{ range .Modules }}import {{ .FullName }} from './{{ .FullPath }}'
import {MsgTypes as {{ .FullName}}MsgTypes} from './{{ .FullPath }}/module'
{{ end }}

export default { 
  {{ range .Modules }}{{ .FullName }}: load({{ .FullName }}, '{{ .Path }}'),
  {{ end }}
}

export const registry = new Registry(<any>[{{ range $j,$mod :=.Modules }}{{- if (gt $j 0) -}}, {{ end }}...{{$mod.FullName}}MsgTypes{{end}}]);

function load(mod, fullns) {
    return function init(store) {        
        if (store.hasModule([fullns])) {
            throw new Error('Duplicate module name detected: '+ fullns)
        }else{
            store.registerModule([fullns], mod)
            store.subscribe((mutation) => {
                if (mutation.type == 'common/env/INITIALIZE_WS_COMPLETE') {
                    store.dispatch(fullns+ '/init', null, {
                        root: true
                    })
                }
            })
        }
    }
}
