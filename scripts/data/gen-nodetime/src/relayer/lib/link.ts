import { readOrCreateConfig, writeConfig } from "./persistence";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { GasPrice } from "@cosmjs/stargate";
import { stringToPath } from "@cosmjs/crypto";
import { Link, IbcClient } from "@confio/relayer/build";
import { getFullPath, FullPath, getDefaultAccountBalance } from "./chain";
import { orderFromJSON } from "@confio/relayer/build/codec/ibc/core/channel/v1/channel";

export const linkMethod = "link";

interface Response {
	// linkedPaths is a list of paths that are linked when link() called.
	linkedPaths: string[];

	// alreadyLinkedPaths is a list of paths that already linked on chain.
	alreadyLinkedPaths: string[];
}

// link connects src and dst chains by their paths on chain with ibc txs.
export async function link(paths: string[]): Promise<Response> {
	const config = readOrCreateConfig();
	if (config.paths) {
		let response = {
			linkedPaths: [],
			alreadyLinkedPaths: [],
		};
		for (let pathName of paths) {
			const path = getFullPath(pathName);
			if (path.path.isLinked) {
				response.alreadyLinkedPaths.push(pathName);
			} else {
				try {
					await createLink(path);
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

async function createLink({ path, options }: FullPath) {
	const config = readOrCreateConfig();
	let chainA = config.chains.find((x) => x.chainId == path.src.chainID);
	let chainB = config.chains.find((x) => x.chainId == path.dst.chainID);
	let chainABalances = await getDefaultAccountBalance([path.src.chainID]);
	let chainBBalances = await getDefaultAccountBalance([path.dst.chainID]);
	let chainAGP = GasPrice.fromString(chainA.gasPrice);
	let chainBGP = GasPrice.fromString(chainB.gasPrice);
	if (!chainABalances.find((x) => x.denom == chainAGP.denom)) {
		throw new Error(
			"Not enough balance available on '" +
				path.src.chainID +
				"'. You need at least " +
				chainAGP.amount.toFloatApproximation() * 2256000 +
				chainAGP.denom
		);
	}
	if (!chainBBalances.find((x) => x.denom == chainBGP.denom)) {
		throw new Error(
			"Not enough balance available on '" +
				path.dst.chainID +
				"'. You need at least " +
				chainBGP.amount.toFloatApproximation() * 2256000 +
				chainBGP.denom
		);
	}

	if (
		chainABalances.find(
			(x) =>
				x.denom == chainAGP.denom &&
				parseInt(x.amount) < chainAGP.amount.toFloatApproximation() * 2256000
		)
	) {
		throw new Error(
			"Not enough balance available on '" +
				path.src.chainID +
				"'. You need at least " +
				chainAGP.amount.toFloatApproximation() * 2256000 +
				chainAGP.denom
		);
	}
	if (
		chainBBalances.find(
			(x) =>
				x.denom == chainBGP.denom &&
				parseInt(x.amount) < chainBGP.amount.toFloatApproximation() * 2256000
		)
	) {
		throw new Error(
			"Not enough balance available on '" +
				path.dst.chainID +
				"'. You need at least " +
				chainBGP.amount.toFloatApproximation() * 2256000 +
				chainBGP.denom
		);
	}
	let signerA = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
		hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
		prefix: chainA.addrPrefix,
	});
	// Create a signer from chain B config
	let signerB = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
		hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
		prefix: chainB.addrPrefix,
	});
	// get addresses
	const [accountA] = await signerA.getAccounts();
	const [accountB] = await signerB.getAccounts();
	// Create IBC Client for chain A

	const clientA = await IbcClient.connectWithSigner(
		chainA.rpcAddr,
		signerA,
		accountA.address,
		{
			prefix: chainA.addrPrefix,
			gasPrice: GasPrice.fromString(chainA.gasPrice),
		}
	);
	// Create IBC Client for chain B

	const clientB = await IbcClient.connectWithSigner(
		chainB.rpcAddr,
		signerB,
		accountB.address,
		{
			prefix: chainB.addrPrefix,
			gasPrice: GasPrice.fromString(chainB.gasPrice),
		}
	);
	const link = await Link.createWithNewConnections(clientA, clientB);
	const channels = await link.createChannel(
		"A",
		options.sourcePort,
		options.targetPort,
		orderFromJSON(options.ordering),
		options.targetVersion
	);
	let configPath = config.paths.find((x) => x.path.id == path.id);
	configPath.path.src.channelID = channels.src.channelId;
	configPath.path.dst.channelID = channels.dest.channelId;
	configPath.path.isLinked = true;
	configPath.connections = {
		srcConnection: link.endA.connectionID,
		destConnection: link.endB.connectionID,
	};
	configPath.relayerData = null;
	writeConfig(config);
}
export async function getLink({ path, connections }: FullPath) {
	const config = readOrCreateConfig();
	let chainA = config.chains.find((x) => x.chainId == path.src.chainID);
	let chainB = config.chains.find((x) => x.chainId == path.dst.chainID);
	let signerA = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
		hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
		prefix: chainA.addrPrefix,
	});
	// Create a signer from chain B config
	let signerB = await DirectSecp256k1HdWallet.fromMnemonic(config.mnemonic, {
		hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
		prefix: chainB.addrPrefix,
	});
	// get addresses
	const [accountA] = await signerA.getAccounts();
	const [accountB] = await signerB.getAccounts();
	// Create IBC Client for chain A

	const clientA = await IbcClient.connectWithSigner(
		chainA.rpcAddr,
		signerA,
		accountA.address,
		{
			prefix: chainA.addrPrefix,
			gasPrice: GasPrice.fromString(chainA.gasPrice),
		}
	);
	// Create IBC Client for chain B

	const clientB = await IbcClient.connectWithSigner(
		chainB.rpcAddr,
		signerB,
		accountB.address,
		{
			prefix: chainB.addrPrefix,
			gasPrice: GasPrice.fromString(chainB.gasPrice),
		}
	);
	const link = Link.createWithExistingConnections(
		clientA,
		clientB,
		connections.srcConnection,
		connections.destConnection
	);
	return link;
}
