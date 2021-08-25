import run from "./jsonrpc";

import Relayer from "./lib/relayer";

const relayer = new Relayer();

run([
	["link", relayer.link.bind(relayer)],
	["start", relayer.start.bind(relayer)],
]);
