// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

{{ range .Modules }}import {{ camelCaseUpperSta .Pkg.Name }} from './{{ .Pkg.Name }}'
{{ end }}

export default { 
  {{ range .Modules }}{{ camelCaseUpperSta .Pkg.Name }}: load({{ camelCaseUpperSta .Pkg.Name }}, '{{ .Pkg.Name }}'),
  {{ end }}
}


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