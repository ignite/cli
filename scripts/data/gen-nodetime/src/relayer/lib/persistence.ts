import os from "os";
import yaml from "js-yaml";
import fs from "fs";
import { Bip39, Random } from "@cosmjs/crypto";

export function generateMnemonic(): string {
	return Bip39.encode(Random.getBytes(16)).toString();
}
const homedir = os.homedir();

export function createConfigFolder() {
	const configPath = homedir + "/.ts-relayer";
	try {
		if (!fs.existsSync(configPath)) {
			fs.mkdirSync(configPath);
		}
	} catch (e) {
		throw new Error("Could not create config folder: " + e);
	}
}

export function readOrCreateConfig() {
	createConfigFolder();
	try {
		if (fs.existsSync(homedir + "/.ts-relayer/config.yaml")) {
			let configFile = fs.readFileSync(
				homedir + "/.ts-relayer/config.yaml",
				"utf8"
			);
			return yaml.load(configFile);
		} else {
			let config = {
				mnemonic: Bip39.encode(Random.getBytes(32)).toString(),
			};
			let configFile = yaml.dump(config);
			fs.writeFileSync(
				homedir + "/.ts-relayer/config.yaml",
				configFile,
				"utf8"
			);
			return config;
		}
	} catch (e) {
		throw new Error("Failed reading config: " + e);
	}
}
export function writeConfig(config) {
	try {
		let configFile = yaml.dump(config);
		fs.writeFileSync(homedir + "/.ts-relayer/config.yaml", configFile, "utf8");
	} catch (e) {
		throw new Error("Failed writing  config: " + e);
	}
}
