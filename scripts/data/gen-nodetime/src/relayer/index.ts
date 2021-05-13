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
} from "./lib/xrelayer";
import XRelayer from "./lib/xrelayer";

const relayer = new XRelayer();
const link = relayer.link.bind(relayer);
const start = relayer.start.bind(relayer);
const ensureChainSetup = relayer.ensureChainSetup.bind(relayer);
const createPath = relayer.createPath.bind(relayer);
const getPath = relayer.getPath.bind(relayer);
const listPaths = relayer.listPaths.bind(relayer);
const getDefaultAccount = relayer.getDefaultAccount.bind(relayer);
const getDefaultAccountBalance = relayer.getDefaultAccountBalance.bind(relayer);

run([
	[linkMethod, link],
	[startMethod, start],
	[ensureChainSetupMethod, ensureChainSetup],
	[createPathMethod, createPath],
	[getPathMethod, getPath],
	[listPathsMethod, listPaths],
	[getDefaultAccountMethod, getDefaultAccount],
	[getDefaultAccountBalanceMethod, getDefaultAccountBalance],
]);
