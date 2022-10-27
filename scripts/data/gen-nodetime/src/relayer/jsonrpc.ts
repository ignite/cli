import { stdin } from "process";
import { JSONRPCServer, SimpleJSONRPCMethod} from "json-rpc-2.0";

type Handler = [string, SimpleJSONRPCMethod];

// run exposes the JSON-RPC server for given handlers through the standard streams.
export default async function run(handlers: Handler[]) {
  // init the rpc server.
  const server = new JSONRPCServer();

  // attach methods to the rpc server.
  for (const [name, func] of handlers) {
    server.addMethod(name, func);
  }

  // read the rpc call, invoke it and send a response.
  let jsonRequest: string = "";

  stdin.setEncoding("utf8");

  for await (const chunk of stdin) {
    jsonRequest += chunk;
  }

  const jsonResponse = await server.receiveJSON(jsonRequest);
  const response = JSON.stringify(jsonResponse);

  console.log(response);
}

