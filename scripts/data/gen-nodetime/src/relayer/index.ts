import run from "./jsonrpc";

import Relayer from "./lib/relayer";

const relayer = new Relayer(parseInt(process.argv[2]));

run([
	["link", relayer.link.bind(relayer)],
	["start", relayer.start.bind(relayer)],
]);
