// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { defineStore } from 'pinia'

{{ range .Module.Types }}import { {{ .Name }} } from "./types/{{ resolveFile .FilePath }}"
{{ end }}

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
  }
};

const usePiniaStore = defineStore('module', piniaStore); 

export default usePiniaStore;
