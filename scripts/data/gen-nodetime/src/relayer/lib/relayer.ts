import os from "os";
import fs from "fs";
import { join } from "path";

import yaml from "js-yaml";

import { Bip39, Random } from "@cosmjs/crypto";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { GasPrice } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import { stringToPath } from "@cosmjs/crypto";
import { Coin } from "@cosmjs/stargate";
import { Link, IbcClient } from "@confio/relayer/build";
import { orderFromJSON } from "@confio/relayer/build/codec/ibc/core/channel/v1/channel";

import Errors from "./errors";
import ConsoleLogger from "./logger";

export const ensureChainSetupMethod = "ensureChainSetup";
export const createPathMethod = "createPath";
export const getPathMethod = "getPath";
export const listPathsMethod = "listPaths";
export const getDefaultAccountMethod = "getDefaultAccount";
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance";
export const linkMethod = "link";
export const startMethod = "start";
export const infoMethod = "info";

const IBCSetupGas = 2256000;
const defaultMaxAge = 86400;
const configDir = ".starport/relayer";
const configName = "config.yaml";

type Account = {
	address: string;
};

type ChainConfig = {
	chainId: string;
	rpcAddr: string;
	addressPrefix: string;
	gasPrice: string;
};

type ConnectOptions = {
	sourcePort: string;
	sourceVersion: string;
	targetPort: string;
	targetVersion: string;
	ordering: "ORDER_UNORDERED" | "ORDER_ORDERED";
};

type Connections = {
	srcConnection: string;
	destConnection: string;
};

type Endpoint = {
	channelID?: string;
	chainID: string;
	portID: string;
};

type Path = {
	id: string;
	isLinked: boolean;
	src: Endpoint;
	dst: Endpoint;
};

type PacketHeights = {
	packetHeightA: number;
	packetHeightB: number;
	ackHeightA: number;
	ackHeightB: number;
};

type PathConfig = {
	path: Path;
	options?: ConnectOptions;
	connections?: Connections;
	relayerData?: PacketHeights;
};

type RelayerConfig = {
	mnemonic: string;
	chains?: Array<ChainConfig>;
	paths?: Array<PathConfig>;
};

interface ChainSetupOptions {
	gasPrice: string;
	addressPrefix: string;
}

type EnsureChainSetupResponse = {
	// id is chain id.
	id: string;
};
type LinkError = {
	pathName: string;
	error: string;
};
type LinkResponse = {
	linkedPaths: string[];
	alreadyLinkedPaths: string[];
	failedToLinkPaths: LinkError[];
};
type LinkStatus = {
	status: boolean;
	pathName: string;
	error?: string;
};
type StartResponse = {};

type InfoResponse = {
	configPath: string;
};

export default class Relayer {
	public config: RelayerConfig;
	public homedir: string;
	public configFile: string;
	private pollTime: number;
	private relayers: Map<string, ReturnType<typeof setInterval>>;

	constructor(configFile: string = configName, pollTime = 5000) {
		this.homedir = os.homedir();
		this.pollTime = pollTime;
		this.configFile = configFile;
		this.relayers = new Map();
		this.config = this.readOrCreateConfig();

		const nestedProxy = {
			set: (target, prop, value) => {
				target[prop] = value;
				this.writeConfig(this.config);
				return true;
			},

			get: (target, prop) => {
				if (typeof target[prop] === "object" && target[prop] !== null) {
					return new Proxy(target[prop], nestedProxy);
				} else {
					return target[prop];
				}
			},
		};

		this.config = new Proxy(this.config, nestedProxy);
	}

	ensureConfigDirCreated() {
		try {
			if (!fs.existsSync(this.getConfigFolder())) {
				fs.mkdirSync(this.getConfigFolder(), { recursive: true });
			}
		} catch (e) {
			throw Errors.ConfigFolderFailed(e);
		}
	}

	getConfigFolder() {
		return join(this.homedir, configDir);
	}

	getConfigPath() {
		return join(this.getConfigFolder(), this.configFile);
	}

	readOrCreateConfig() {
		this.ensureConfigDirCreated();

		try {
			if (fs.existsSync(this.getConfigPath())) {
				let configFile = fs.readFileSync(this.getConfigPath(), "utf8");
				return yaml.load(configFile);
			}

			let config = {
				mnemonic: Bip39.encode(Random.getBytes(32)).toString(),
			};

			let configFile = yaml.dump(config);
			fs.writeFileSync(this.getConfigPath(), configFile, "utf8");

			return config;
		} catch (e) {
			throw Errors.ConfigReadFailed(e);
		}
	}

	writeConfig(config) {
		try {
			let configFile = yaml.dump(config);
			fs.writeFileSync(this.getConfigPath(), configFile, "utf8");
		} catch (e) {
			throw Errors.ConfigWriteFailed(e);
		}
	}

	// ensureChainSetup is an RPC shim func to be executed by the Go API.
	public async ensureChainSetup([rpcAddr, { addressPrefix, gasPrice }]: [
		string,
		ChainSetupOptions
	]): Promise<EnsureChainSetupResponse> {
		try {
			const tmClient = await Tendermint34Client.connect(rpcAddr);
			const status = await tmClient.status();

			const chain = {
				chainId: status.nodeInfo.network,
				rpcAddr,
				addressPrefix,
				gasPrice,
			};

			if (!this.config.chains) {
				this.config.chains = [chain];
				return { id: chain.chainId };
			}

			const endpointExistsWithDifferentChainID =
				this.config.chains &&
				this.config.chains.find(
					(x) => x.chainId != chain.chainId && x.rpcAddr == chain.rpcAddr
				);

			const chainExistsWithSameEndpoint =
				this.config.chains &&
				this.config.chains.find(
					(x) => x.chainId == chain.chainId && x.rpcAddr == chain.rpcAddr
				);

			if (endpointExistsWithDifferentChainID)
				throw Errors.EndpointExistsWithDifferentChainID;

			if (chainExistsWithSameEndpoint) return { id: chain.chainId };

			this.config.chains.push(chain);

			return { id: chain.chainId };
		} catch (e) {
			throw Errors.ChainSetupFailed(e);
		}
	}

	// createPath is an RPC shim func to be executed by the Go API.
	public createPath([srcID, dstID, options]: [
		string,
		string,
		ConnectOptions
	]): Path {
		try {
			let path = {
				id: `${srcID}-${dstID}`,
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

			if (!this.config.paths) this.config.paths = [];

			this.config.paths.push({ path, options });

			return path;
		} catch (e) {
			throw Errors.PathSetupFailed(e);
		}
	}

	// getPath is an RPC shim func to be executed by the Go API.
	public getPath([id]: [string]): Path {
		if (this.config.paths) {
			let path = this.config.paths.find((x) => x.path.id == id);
			if (path) return path.path;
		}
		throw Errors.PathNotExists;
	}

	// listPaths is an RPC shim func to be executed by the Go API.
	public listPaths(): Path[] {
		if (this.config.paths) {
			let paths = this.config.paths.map((x) => x.path);
			return paths;
		}
		throw Errors.PathsNotDefined;
	}

	// getDefaultAccount is an RPC shim func to be executed by the Go API.
	public async getDefaultAccount([chainID]: [string]): Promise<Account> {
		const chain = this.chainById(chainID);
		if (chain) {
			let client = await this.getIBCClient(chain);
			return {
				address: client.senderAddress,
			};
		}
		throw Errors.ChainNotFound(chainID);
	}

	// getDefaultAccountBalance is an RPC shim func to be executed by the Go API.
	public async getDefaultAccountBalance([chainID]: [string]): Promise<Coin[]> {
		const chain = this.chainById(chainID);
		if (chain) {
			let client = await this.getIBCClient(chain);
			return await client.query.bank.allBalances(client.senderAddress);
		}
		throw Errors.ChainNotFound(chainID);
	}

	// link is an RPC shim func to be executed by the Go API.
	public async link([paths]: [string[]]): Promise<LinkResponse> {
		if (this.config.paths) {
			let response: LinkResponse = {
				linkedPaths: [],
				alreadyLinkedPaths: [],
				failedToLinkPaths: [],
			};
			const links = [];
			for (let pathName of paths) {
				const path = this.pathById(pathName);

				if (path?.path.isLinked) {
					response.alreadyLinkedPaths.push(pathName);
					continue;
				}
				links.push(this.createLink(path));
			}
			const results = await Promise.all(links);
			for (let result of results) {
				if (result.status) {
					response.linkedPaths.push(result.pathName);
				} else {
					response.failedToLinkPaths.push({
						pathName: result.pathName,
						error: result.error,
					});
				}
			}
			return response;
		}

		throw Errors.PathsNotDefined;
	}

	// start is an RPC shim func to be executed by the Go API.
	public async start([paths]: [string[]]): Promise<StartResponse> {
		if (this.config.paths) {
			for (let pathName of paths) {
				const path = this.pathById(pathName);

				if (path?.path.isLinked) {
					const link = await this.getLink(path);
					this.relayers.set(
						pathName,
						setInterval(async () => {
							let heights = this.pathById(pathName).relayerData;
							let newHeights = await this.relayPackets(link, heights);
							this.pathById(pathName).relayerData = newHeights;
						}, this.pollTime)
					);

					continue;
				}

				throw Errors.PathNotLinked;
			}

			return {};
		}

		throw Errors.PathsNotDefined;
	}

	// info is an RPC shim func to be executed by the Go API.
	public async info(): Promise<InfoResponse> {
		return { configPath: this.getConfigPath() };
	}

	protected chainById(chainID: string): ChainConfig {
		return this.config.chains
			? this.config.chains.find((x) => x.chainId == chainID)
			: null;
	}

	protected pathById(pathID: string): PathConfig {
		return this.config.paths
			? this.config.paths.find((x) => x.path.id == pathID)
			: null;
	}
	protected async balanceCheck(chain: ChainConfig): Promise<boolean> {
		let chainBalances = await this.getDefaultAccountBalance([chain.chainId]);
		let chainGP = GasPrice.fromString(chain.gasPrice);
		if (!chainBalances.find((x) => x.denom == chainGP.denom)) return false;

		return !chainBalances.find(
			(x) =>
				x.denom == chainGP.denom &&
				parseInt(x.amount) < chainGP.amount.toFloatApproximation() * IBCSetupGas
		);
	}

	protected notEnoughBalanceError(chain, gasPrice) {
		const { chainId } = chain;
		const { amount, denom } = gasPrice;
		const calcAmount = amount.toFloatApproximation() * IBCSetupGas;

		return Errors.NotEnoughBalance(`${calcAmount} ${denom} (${chainId})`);
	}

	protected async createLink({
		path,
		options,
	}: PathConfig): Promise<LinkStatus> {
		let chainA = this.chainById(path.src.chainID);
		let chainB = this.chainById(path.dst.chainID);
		let chainAGP = GasPrice.fromString(chainA.gasPrice);
		let chainBGP = GasPrice.fromString(chainA.gasPrice);

		if (!(await this.balanceCheck(chainA)))
			return {
				status: false,
				pathName: path.id,
				error: this.notEnoughBalanceError(chainA, chainAGP).message,
			};

		if (!(await this.balanceCheck(chainB)))
			return {
				status: false,
				pathName: path.id,
				error: this.notEnoughBalanceError(chainB, chainBGP).message,
			};

		// create IBC clients.
		const clientA = await this.getIBCClient(chainA);
		const clientB = await this.getIBCClient(chainB);
		try {
			const link = await Link.createWithNewConnections(clientA, clientB);

			const channels = await link.createChannel(
				"A",
				options.sourcePort,
				options.targetPort,
				orderFromJSON(options.ordering),
				options.targetVersion
			);

			let configPath = this.pathById(path.id);
			configPath.path.src.channelID = channels.src.channelId;
			configPath.path.dst.channelID = channels.dest.channelId;
			configPath.path.isLinked = true;
			configPath.connections = {
				srcConnection: link.endA.connectionID,
				destConnection: link.endB.connectionID,
			};
			configPath.relayerData = null;
			return {
				status: true,
				pathName: path.id,
			};
		} catch (e) {
			return {
				status: false,
				pathName: path.id,
				error: e,
			};
		}
	}

	protected async getLink({ path, connections }: PathConfig): Promise<Link> {
		let chainA = this.chainById(path.src.chainID);
		let chainB = this.chainById(path.dst.chainID);

		// create IBC clients.
		const clientA = await this.getIBCClient(chainA);
		const clientB = await this.getIBCClient(chainB);

		const link = Link.createWithExistingConnections(
			clientA,
			clientB,
			connections.srcConnection,
			connections.destConnection,
			new ConsoleLogger()
		);

		return link;
	}

	protected async getIBCClient(chain: ChainConfig): Promise<IbcClient> {
		let chainGP = GasPrice.fromString(chain.gasPrice);
		let signer = await DirectSecp256k1HdWallet.fromMnemonic(
			this.config.mnemonic,
			{
				hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
				prefix: chain.addressPrefix,
			}
		);

		const [account] = await signer.getAccounts();

		const client = await IbcClient.connectWithSigner(
			chain.rpcAddr,
			signer,
			account.address,
			{
				prefix: chain.addressPrefix,
				gasPrice: chainGP,
			}
		);

		return client;
	}

	protected async relayPackets(
		link,
		relayHeights,
		options = { maxAgeDest: defaultMaxAge, maxAgeSrc: defaultMaxAge }
	) {
		try {
			const heights = await link.checkAndRelayPacketsAndAcks(
				relayHeights ?? {},
				2,
				6
			);

			await link.updateClientIfStale("A", options.maxAgeDest);
			await link.updateClientIfStale("B", options.maxAgeSrc);

			return heights;
		} catch (e) {
			throw Errors.RelayPacketError;
		}
	}
}
