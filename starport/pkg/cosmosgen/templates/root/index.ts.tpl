// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { OfflineSigner, Registry } from "@cosmjs/proto-signing";
import { SigningStargateClient } from "@cosmjs/stargate";

{{ range .Modules }}import { Module as {{ .Name }}, msgTypes as {{ .Name }}MsgTypes } from './{{ .Path }}'
{{ end }}

const registry = new Registry(<any>[
  {{ range .Modules }}...{{ .Name }}MsgTypes,
  {{ end }}
])

class Signer {
  private _offlineSigner: OfflineSigner;
  private _client?: SigningStargateClient;
  private _addr: string;

  constructor(rpcAddr: string, offlineSigner: OfflineSigner) {
    this._offlineSigner = offlineSigner;
    this._addr = rpcAddr;
  }

  public async init() {
    this._client = await SigningStargateClient.connectWithSigner(
      this._addr,
      this._offlineSigner,
      { registry }
    );
  }

  get signer() {
    return this._client;
  }
}

interface Environment {
  chainID: string;
  chainName: string;
  apiURL: string;
  rpcURL: string;
  wsURL: string;
}

interface IgniteParams {
  env: Environment;
  signer: OfflineSigner;
  address: string;
}

class Ignite {
  private _env: Environment;
  private _signer: Signer;
  private _address: string;

{{ range .Modules }}public {{ .Name }}: {{ .Name }};
{{ end }}

  constructor({ env, signer, address }: IgniteParams) {
    this._env = env;
    this._address = address;
    this._signer = new Signer(env.rpcURL, signer);
  }

  public async init() {
   await this._signer.init();

     let client: SigningStargateClient = this._signer
       .signer as SigningStargateClient;

{{ range .Modules }}this.{{ .Name }} = new {{ .Name }}(
     client,
    this._address,
     this._env.apiURL
);
{{ end }}
   }

  get env(): Environment {
    return this._env;
  }
}

export {
    registry,
    Ignite,
}
