export const ensureChainSetupMethod = "ensureChainSetup";
export const createPathMethod = "createPath";
export const getPathMethod = "getPath";
export const listPathsMethod = "listPaths";
export const getDefaultAccountMethod = "getDefaultAccount";
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance";
import { readOrCreateConfig, writeConfig } from "./persistence";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { GasPrice } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import { stringToPath } from "@cosmjs/crypto";
import { IbcClient } from "@confio/relayer/build";
import { Coin } from "@cosmjs/stargate";
type EnsureChainSetupResponse = {
	// id is the chain id of chain.
	id: string;
};
interface Chain {
	chainId: string;
	rpcAddr: string;
	addrPrefix: string;
	gasPrice: string;
}
// ensureChainSetup sets up a chain by its rpc address only if it is not set up already.
export async function ensureChainSetup([rpcAddr, gasPrice, addrPrefix]: [
	string,
	string,
	string
]): Promise<EnsureChainSetupResponse> {
	try {
		const config = readOrCreateConfig();
		const tmClient = await Tendermint34Client.connect(rpcAddr);
		let status = await tmClient.status();
		const chain = {
			chainId: status.nodeInfo.network,
			rpcAddr,
			addrPrefix,
			gasPrice,
		};
		if (!config.chains) {
			config.chains = [];
			config.chains.push(chain);
			writeConfig(config);
			return { id: chain.chainId };
		}
		if (
			config.chains &&
			config.chains.find(
				(x) => x.chainId == chain.chainId && x.rpcAddr == chain.rpcAddr
			)
		) {
			throw new Error("chain already exists");
		} else {
			config.chains.push(chain);
			writeConfig(config);
			return { id: chain.chainId };
		}
	} catch (e) {
		throw new Error("Could not setup chain: " + e);
	}
}

interface ConnectOptions {
	sourcePort: string;
	sourceVersion: string;
	targetPort: string;
	targetVersion: string;
	ordering: string;
}

// createPath creates a path between the source chain and dest chain by their chain ids with given options.
// it returns a unique path id that represents the connection between these chains.
//
// createPath should only record the intention of connecting source and destion chains together
// and should not send any txs to these chains. this will later be done by link().
export function createPath([srcID, dstID, options]: [
	string,
	string,
	ConnectOptions
]): Path {
	try {
		const config = readOrCreateConfig();
		let path = {
			id: srcID + "-" + dstID,
			isLinked: false,
			src: {
				chainID: srcID,
				portID: options.sourcePort,
			},
			dst: {
				chainID: dstID,
				portID: options.targetPort,
			},
		};
		if (!config.paths) {
			config.paths = [];
		}
		config.paths.push({ path, options });
		writeConfig(config);
		return path;
	} catch (e) {
		throw new Error("Could not create path:" + e);
	}
}

// Path represents the connection between two chaons.
interface Path {
	// id of the path.
	id: string;

	// isLinked shows whether src and dst chains are connected on the chain with ibc txs.
	isLinked: boolean;

	// src represents the source chain.
	src: PathEnd;

	// dst represents the destionation chain.
	dst: PathEnd;
}

// PathEnd represents a chain.
interface PathEnd {
	channelID?: string;
	chainID: string;
	portID: string;
}
export interface FullPath {
	path: Path;
	options: ConnectOptions;
	connections?: {
		srcConnection: string;
		destConnection: string;
	};
	relayerData?: {
		packetHeightA?: number;
		packetHeightB?: number;
		ackHeightA?: number;
		ackHeightB?: number;
	} | null;
}
// getPath gets connection info between chains by path id.
export function getPath(id: [string]): Path {
	const config = readOrCreateConfig();
	if (config.paths) {
		let path = config.paths.find((x) => x.path.id == id);
		if (path) {
			return path.path;
		} else {
			throw new Error("Path does not exist");
		}
	} else {
		throw new Error("Path does not exist");
	}
}
export function getFullPath(id: string): FullPath {
	const config = readOrCreateConfig();
	if (config.paths) {
		let path = config.paths.find((x) => x.path.id == id);
		if (path) {
			return path;
		} else {
			throw new Error("Path does not exist");
		}
	} else {
		throw new Error("Path does not exist");
	}
}
// listPaths list all connections.
export function listPaths(): Path[] {
	const config = readOrCreateConfig();
	if (config.paths) {
		let paths = config.paths.map((x) => x.path);
		return paths;
	} else {
		throw new Error("No paths defined");
	}
}

interface Account {
	address: string;
}

// getDefaultAccount gets the default account on chain by chain id.
export async function getDefaultAccount([chainID]: string[]): Promise<Account> {
	const config = readOrCreateConfig();
	const chain = config.chains.find((x) => x.chainId == chainID);
	if (chain) {
		let signer = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
			hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
			prefix: chain.addrPrefix,
		});
		const [account] = await signer.getAccounts();
		return {
			address: account.address,
		};
	} else {
		throw new Error("Chain not found: " + chainID);
	}
}

// getDefaultAccountBalance gets the balance of default account on chain by chain id.
export async function getDefaultAccountBalance([chainID]: string[]): Promise<
	Coin[]
> {
	const config = readOrCreateConfig();
	const chain = config.chains.find((x) => x.chainId == chainID);
	if (chain) {
		let signer = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
			hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
			prefix: chain.addrPrefix,
		});
		const [account] = await signer.getAccounts();

		const client = await IbcClient.connectWithSigner(
			chain.rpcAddr,
			signer,
			account.address,
			{
				prefix: chain.addrPrefix,
				gasPrice: GasPrice.fromString(chain.gasPrice),
			}
		);
		return await client.query.bank.allBalances(account.address);
	} else {
		throw new Error("Chain not found: " + chainID);
	}
}
