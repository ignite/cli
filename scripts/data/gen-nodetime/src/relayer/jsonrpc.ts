import { stdin } from "process";
import { JSONRPCServer, SimpleJSONRPCMethod} from "json-rpc-2.0";

type Handler = [string, SimpleJSONRPCMethod];

// run exposes the JSON-RPC server for given handlers through the standard streams.
export default async function run(handlers: Handler[]) {
  // init the rpc server.
  const server = new JSONRPCServer();

  // attach methods to the rpc server.
  for (var [name, func] of handlers) {
    server.addMethod(name, func);
  }

  // read the rpc call, invoke it and send a response.
  let jsonreq: string = "";

  stdin.setEncoding("utf8");

  for await (const chunk of stdin) {
    jsonreq += chunk;
  }

  const jsonres = await server.receiveJSON(jsonreq);
  const res = JSON.stringify(jsonres);

  console.log(res);
}
