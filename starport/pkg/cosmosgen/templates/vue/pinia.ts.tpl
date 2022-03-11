// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { defineStore } from 'pinia'

import { 
{{ range .Module.Types }}{{ .Name }},
{{ end }}
} from "ts-client/{{ .Module.Pkg.Name }}/types";

type PiniaState = {
	{{ range .Module.Types }}{{ .Name }}All: {{ .Name }}[],
	{{ end }}
};

const piniaStore = {
  state: (): PiniaState => {
    return {
		{{ range .Module.Types }}{{ .Name }}All: [],
		{{ end }}
    }
  },
};

const usePiniaStore = defineStore('{{ .Module.Pkg.Name }}', piniaStore); 

export { usePiniaStore, PiniaState };