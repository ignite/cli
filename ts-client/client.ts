import { GeneratedType } from "@cosmjs/proto-signing";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Env } from "./env";
import { UnionToIntersection, Return, Constructor } from "./helpers";
import { Module } from "./modules";

export class IgniteClient {
	static plugins: Module[] = [];
  env: Env;
  client: SigningStargateClient;
  registry: Array<[string, GeneratedType]>
  static plugin<T extends Module | Module[]>(plugin: T) {
    const currentPlugins = this.plugins;

    class AugmentedClient extends this {
      static plugins = currentPlugins.concat(plugin);
    }

    if (Array.isArray(plugin)) {
      type Extension = UnionToIntersection<Return<T>['module']>
      return AugmentedClient as typeof AugmentedClient & Constructor<Extension>;  
    }

    type Extension = Return<T>['module']
    return AugmentedClient as typeof AugmentedClient & Constructor<Extension>;
  }
  constructor(env: Env) {
    this.env = env;
    const classConstructor = this.constructor as typeof IgniteClient;
    classConstructor.plugins.forEach(plugin => {
      const pluginInstance = plugin(this);
      Object.assign(this, pluginInstance.module)
      if (this.registry) {
        this.registry = this.registry.concat(pluginInstance.registry)
      }
		});
		
  }
}