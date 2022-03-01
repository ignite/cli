// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient, DeliverTxResponse } from "@cosmjs/stargate";
import { Registry } from "@cosmjs/proto-signing";
import { Api } from "./rest";

{{ range .Module.Msgs }}import { {{ .Name }} } from "./types/{{ resolveFile .FilePath }}";
{{ end }}

/*********************
 *       MODULE      *
 *********************/

{{ range .Module.Msgs }}
type send{{ .Name }}Params = {
  value: {{ .Name }},
  fee?: StdFee,
  memo?: string
};
{{ end }}
{{ range .Module.Msgs }}
type {{ camelCase .Name }}Params = {
  value: {{ .Name }},
};
{{ end }}

class Module extends Api<any> {
	private _client: SigningStargateClient;
	private _address: string;

   types = [
    {{ range .Module.Msgs }}["/{{ .URI }}", {{ .Name }}],
    {{ end }}
  ];
  registry = new Registry(<any>this.types);

  	constructor(client: SigningStargateClient, address: string, baseUrl: string) {
		super({
			baseUrl
		})

		this._client = client;
		this._address = address;
	}

	{{ range .Module.Msgs }}
	async send{{ .Name }}({ value, fee, memo }: send{{ .Name }}Params): Promise<DeliverTxResponse> {
		try {
			let msg = { typeUrl: "/{{ .URI }}", value: {{ .Name }}.fromPartial( value ) }
			return await this._client.signAndBroadcast(this._address, [msg], fee ? fee : { amount: [], gas: '200000' }, memo)
		} catch (e: any) {
			throw new Error('TxClient:{{ .Name }}:Send Could not broadcast Tx: '+ e.message)
		}
	}
	{{ end }}
	{{ range .Module.Msgs }}
	{{ camelCase .Name }}({ value }: {{ camelCase .Name }}Params): { typeUrl: string; value: {{ .Name }} } {
		try {
			return { typeUrl: "/{{ .URI }}", value: {{ .Name }}.fromPartial( value ) }  
		} catch (e: any) {
			throw new Error('TxClient:{{ .Name }}:Create Could not create message: ' + e.message)
		}
	}
	{{ end }}
};

/*********************
 *        VUE        *
 *********************/

export {
  Module
};
