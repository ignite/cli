type LogMethod = (message: string) => Logger;

interface Logger {
	error: LogMethod;
	warn: LogMethod;
	info: LogMethod;
	verbose: LogMethod;
	debug: LogMethod;
}
export enum LogLevels {
	ERROR = 0,
	WARN = 1,
	INFO = 2,
	VERBOSE = 3,
	DEBUG = 4
}
export default class ConsoleLogger {
	public readonly error: LogMethod;
	public readonly warn: LogMethod;
	public readonly info: LogMethod;
	public readonly verbose: LogMethod;
	public readonly debug: LogMethod;

	constructor(logLevel:LogLevels) {
		this.error = (msg) => {
			if(logLevel>=LogLevels.ERROR) {
				console.log(msg);
			}
			return this;
		};
		this.warn = (msg) => {
			if(logLevel>=LogLevels.WARN) {
				console.log(msg);
			}
			return this;
		};
		this.info = (msg) => {
			if(logLevel>=LogLevels.INFO) {
				if (msg.indexOf('Relay') == 0 && msg.indexOf('Relay 0') == -1) {
					console.log(msg);
				}
			}
			return this;
		};
		this.verbose = (msg) => {
			if(logLevel>=LogLevels.VERBOSE) {
				console.log(msg);
			}
			return this;
		};
		this.debug = (msg) => {
			if(logLevel>=LogLevels.DEBUG) {
				console.log(msg);
			}
			return this;
		};
	}
}
