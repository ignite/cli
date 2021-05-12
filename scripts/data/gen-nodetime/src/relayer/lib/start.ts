import { readOrCreateConfig, writeConfig } from "./persistence";
import { getFullPath } from "./chain";
import { getLink } from "./link";
export const startMethod = "start";

interface Response {}
const relayers = new Map();

// start starts relaying ibc packets for requested paths.
export async function start(paths: string[]): Promise<Response> {
	const config = readOrCreateConfig();
	if (config.paths) {
		for (let pathName of paths) {
			const path = getFullPath(pathName);
			if (path.path.isLinked) {
				const link = await getLink(path);
				relayers.set(
					pathName,
					setInterval(async () => {
						let heights = config.paths.find((x) => x.path.id == pathName)
							.relayerData;
						let newHeights = await relayPackets(link, heights);
						config.paths.find(
							(x) => x.path.id == pathName
						).relayerData = newHeights;
						writeConfig(config);
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

async function relayPackets(
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
