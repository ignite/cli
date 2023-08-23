import run from "./jsonrpc";
import { LogLevels } from "./lib/logger";

import Relayer from "./lib/relayer";

const logLevel = parseInt(process.argv[2]);
const relayer = new Relayer(isNaN(logLevel) ? LogLevels.INFO: logLevel);

run([
	["link", relayer.link.bind(relayer)],
	["start", relayer.start.bind(relayer)],
]);
