import { readOrCreateConfig, writeConfig } from "./persistence";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { GasPrice } from "@cosmjs/stargate";
import { Tendermint34Client } from "@cosmjs/tendermint-rpc";
import { stringToPath } from "@cosmjs/crypto";
import { Link, IbcClient } from "@confio/relayer/build";
import { Coin } from "@cosmjs/stargate";
import { getFullPath, FullPath } from "./chain";
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
					throw new Error("Could not link path: " + pathName);
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
	let chainA = config.chains.find((x) => x.id == path.src.chainID);
	let chainB = config.chains.find((x) => x.id == path.dst.chainID);
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
	let configPath = config.paths.find((x) => x.path.id == path.id).path;
	configPath.src.channelID = channels.src.channelId;
	configPath.dst.channelID = channels.dest.channelId;
	configPath.isLinked = true;
	writeConfig(config);
}
