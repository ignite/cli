import { IgniteClient } from "./client";

export type ModuleInterface = { [key: string]: any }
export type Module = (instance: IgniteClient) => { module: ModuleInterface, registry }
