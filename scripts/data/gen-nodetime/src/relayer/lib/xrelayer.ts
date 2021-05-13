import os from "os";
import yaml from "js-yaml";
import fs from "fs";
import { Bip39, Random } from "@cosmjs/crypto";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { GasPrice } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import { stringToPath } from "@cosmjs/crypto";
import { Coin } from "@cosmjs/stargate";
import { Link, IbcClient } from "@confio/relayer/build";
import { orderFromJSON } from "@confio/relayer/build/codec/ibc/core/channel/v1/channel";
export const ensureChainSetupMethod = "ensureChainSetup";
export const createPathMethod = "createPath";
export const getPathMethod = "getPath";
export const listPathsMethod = "listPaths";
export const getDefaultAccountMethod = "getDefaultAccount";
export const getDefaultAccountBalanceMethod = "getDefaultAccountBalance";
export const linkMethod = "link";
export const startMethod = "start";

type EnsureChainSetupResponse = {
	// id is the chain id of chain.
	id: string;
};
type Account = {
	address: string;
};
type ChainConfig = {
	chainId: string;
	rpcAddr: string;
	addrPrefix: string;
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
type LinkResponse = {
	linkedPaths: string[];
	alreadyLinkedPaths: string[];
};
type StartResponse = {};
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

type LogMethod = (message: string) => Logger;

interface Logger {
	error: LogMethod;
	warn: LogMethod;
	info: LogMethod;
	verbose: LogMethod;
	debug: LogMethod;
}

class ConsoleLogger {
	public readonly error: LogMethod;
	public readonly warn: LogMethod;
	public readonly info: LogMethod;
	public readonly verbose: LogMethod;
	public readonly debug: LogMethod;

	constructor() {
		this.error = (msg) => {
			console.log(msg);
			return this;
		};
		this.warn = (msg) => {
			console.log(msg);
			return this;
		};
		this.info = (msg) => {
			console.log(msg);
			return this;
		};
		this.verbose = (msg) => {
			console.log(msg);
			return this;
		};
		this.debug = (msg) => {
			console.log(msg);
			return this;
		};
	}
}

export default class XRelayer {
	private _config: RelayerConfig;
	public config: RelayerConfig;
	public homedir: string;
	private relayers: Map<string, ReturnType<typeof setInterval>>;

	constructor() {
		this.homedir = os.homedir();
		this.relayers = new Map();
		this._config = this.readOrCreateConfig();
		const nestedProxy = {
			set: (target, prop, value) => {
				target[prop] = value;
				this.writeConfig(this._config);
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
		this.config = new Proxy(this._config, nestedProxy);
	}
	createConfigFolder() {
		const configPath = this.homedir + "/.ts-relayer";
		try {
			if (!fs.existsSync(configPath)) {
				fs.mkdirSync(configPath);
			}
		} catch (e) {
			throw new Error("Could not create config folder: " + e);
		}
	}
	readOrCreateConfig() {
		this.createConfigFolder();
		try {
			if (fs.existsSync(this.homedir + "/.ts-relayer/config.yaml")) {
				let configFile = fs.readFileSync(
					this.homedir + "/.ts-relayer/config.yaml",
					"utf8"
				);
				return yaml.load(configFile);
			} else {
				let config = {
					mnemonic: Bip39.encode(Random.getBytes(32)).toString(),
				};
				let configFile = yaml.dump(config);
				fs.writeFileSync(
					this.homedir + "/.ts-relayer/config.yaml",
					configFile,
					"utf8"
				);
				return config;
			}
		} catch (e) {
			throw new Error("Failed reading config: " + e);
		}
	}
	writeConfig(config) {
		try {
			let configFile = yaml.dump(config);
			fs.writeFileSync(
				this.homedir + "/.ts-relayer/config.yaml",
				configFile,
				"utf8"
			);
		} catch (e) {
			throw new Error("Failed writing config: " + e);
		}
	}

	// RPC Handlers
	public async ensureChainSetup([rpcAddr, gasPrice, addrPrefix]: [
		string,
		string,
		string
	]): Promise<EnsureChainSetupResponse> {
		try {
			const tmClient = await Tendermint34Client.connect(rpcAddr);
			let status = await tmClient.status();
			const chain = {
				chainId: status.nodeInfo.network,
				rpcAddr,
				addrPrefix,
				gasPrice,
			};
			if (!this.config.chains) {
				this.config.chains = [chain];
				return { id: chain.chainId };
			}
			if (
				this.config.chains &&
				this.config.chains.find(
					(x) => x.chainId != chain.chainId && x.rpcAddr == chain.rpcAddr
				)
			) {
				throw new Error(
					"RPC endpoint already exists with a different chain id"
				);
			}
			if (
				this.config.chains &&
				this.config.chains.find(
					(x) => x.chainId == chain.chainId && x.rpcAddr == chain.rpcAddr
				)
			) {
				return { id: chain.chainId };
			} else {
				this.config.chains.push(chain);
				return { id: chain.chainId };
			}
		} catch (e) {
			throw new Error("Could not setup chain: " + e);
		}
	}
	public createPath([srcID, dstID, options]: [
		string,
		string,
		ConnectOptions
	]): Path {
		try {
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
			if (!this.config.paths) {
				this.config.paths = [];
			}
			this.config.paths.push({ path, options });
			return path;
		} catch (e) {
			throw new Error("Could not create path:" + e);
		}
	}
	public getPath([id]: [string]): Path {
		if (this.config.paths) {
			let path = this.config.paths.find((x) => x.path.id == id);
			if (path) {
				return path.path;
			} else {
				throw new Error("Path does not exist");
			}
		} else {
			throw new Error("Path does not exist");
		}
	}
	public listPaths(): Path[] {
		if (this.config.paths) {
			let paths = this.config.paths.map((x) => x.path);
			return paths;
		} else {
			throw new Error("No paths defined");
		}
	}
	public async getDefaultAccount([chainID]: string[]): Promise<Account> {
		const chain = this.chainById(chainID);
		if (chain) {
			let client = await this.getIBCClient(chain);
			return {
				address: client.senderAddress,
			};
		} else {
			throw new Error("Chain not found: " + chainID);
		}
	}
	public async getDefaultAccountBalance([chainID]: string[]): Promise<Coin[]> {
		const chain = this.chainById(chainID);
		if (chain) {
			let client = await this.getIBCClient(chain);
			return await client.query.bank.allBalances(client.senderAddress);
		} else {
			throw new Error("Chain not found: " + chainID);
		}
	}
	public async link(paths: string[]): Promise<LinkResponse> {
		if (this.config.paths) {
			let response: LinkResponse = {
				linkedPaths: [],
				alreadyLinkedPaths: [],
			};
			for (let pathName of paths) {
				const path = this.pathById(pathName);
				if (path?.path.isLinked) {
					response.alreadyLinkedPaths.push(pathName);
				} else {
					try {
						await this.createLink(path);
						response.linkedPaths.push(pathName);
					} catch (e) {
						throw new Error("Could not link path: " + pathName + ": " + e);
					}
				}
			}
			return response;
		} else {
			throw new Error("No paths defined");
		}
	}
	public async start(paths: string[]): Promise<StartResponse> {
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
						}, 5000)
					);
				} else {
					throw new Error("Path: " + pathName + " is not linked.");
				}
			}
			return {};
		} else {
			throw new Error("No paths defined");
		}
	}
	// Helper functions
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
		if (!chainBalances.find((x) => x.denom == chainGP.denom)) {
			return false;
		}
		if (
			chainBalances.find(
				(x) =>
					x.denom == chainGP.denom &&
					parseInt(x.amount) < chainGP.amount.toFloatApproximation() * 2256000
			)
		) {
			return false;
		}
		return true;
	}
	protected async createLink({ path, options }: PathConfig): Promise<void> {
		let chainA = this.chainById(path.src.chainID);
		let chainB = this.chainById(path.dst.chainID);
		let chainAGP = GasPrice.fromString(chainA.gasPrice);
		let chainBGP = GasPrice.fromString(chainA.gasPrice);
		if (!(await this.balanceCheck(chainA))) {
			throw new Error(
				"Not enough balance available on '" +
					chainA.chainId +
					"'. You need at least " +
					chainAGP.amount.toFloatApproximation() * 2256000 +
					chainAGP.denom
			);
		}
		if (!(await this.balanceCheck(chainB))) {
			throw new Error(
				"Not enough balance available on '" +
					chainB.chainId +
					"'. You need at least " +
					chainBGP.amount.toFloatApproximation() * 2256000 +
					chainBGP.denom
			);
		}
		// Create IBC Client for chain A

		const clientA = await this.getIBCClient(chainA);
		// Create IBC Client for chain B

		const clientB = await this.getIBCClient(chainB);

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
	}
	protected async getLink({ path, connections }: PathConfig): Promise<Link> {
		let chainA = this.chainById(path.src.chainID);
		let chainB = this.chainById(path.dst.chainID);
		// Create IBC Client for chain A

		const clientA = await this.getIBCClient(chainA);
		// Create IBC Client for chain B

		const clientB = await this.getIBCClient(chainB);

		const logger = new ConsoleLogger();
		const link = Link.createWithExistingConnections(
			clientA,
			clientB,
			connections.srcConnection,
			connections.destConnection,
			logger
		);
		return link;
	}
	protected async getIBCClient(chain: ChainConfig): Promise<IbcClient> {
		let chainGP = GasPrice.fromString(chain.gasPrice);
		let signer = await DirectSecp256k1HdWallet.fromMnemonic(
			this.config.mnemonic,
			{
				hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
				prefix: chain.addrPrefix,
			}
		);
		const [account] = await signer.getAccounts();

		const client = await IbcClient.connectWithSigner(
			chain.rpcAddr,
			signer,
			account.address,
			{
				prefix: chain.addrPrefix,
				gasPrice: chainGP,
			}
		);
		return client;
	}

	protected async relayPackets(
		link,
		relayHeights,
		options = { maxAgeDest: 86400, maxAgeSrc: 86400 }
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
			throw new Error("Error relaying packets");
		}
	}
}
