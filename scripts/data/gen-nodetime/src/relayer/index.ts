import run from "./jsonrpc";

import {
	linkMethod,
	startMethod,
	ensureChainSetupMethod,
	createPathMethod,
	getPathMethod,
	listPathsMethod,
	getDefaultAccountMethod,
	getDefaultAccountBalanceMethod,
} from "./lib/relayer";
import Relayer from "./lib/relayer";

const relayer = new Relayer();

run([
	[linkMethod, relayer.link.bind(relayer)],
	[startMethod, relayer.start.bind(relayer)],
	[ensureChainSetupMethod, relayer.ensureChainSetup.bind(relayer)],
	[createPathMethod, relayer.createPath.bind(relayer)],
	[getPathMethod, relayer.getPath.bind(relayer)],
	[listPathsMethod, relayer.listPaths.bind(relayer)],
	[getDefaultAccountMethod, relayer.getDefaultAccount.bind(relayer)],
	[
		getDefaultAccountBalanceMethod,
		relayer.getDefaultAccountBalance.bind(relayer),
	],
]);
