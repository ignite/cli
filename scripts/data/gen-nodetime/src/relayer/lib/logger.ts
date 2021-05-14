type LogMethod = (message: string) => Logger;

interface Logger {
	error: LogMethod;
	warn: LogMethod;
	info: LogMethod;
	verbose: LogMethod;
	debug: LogMethod;
}

export default class ConsoleLogger {
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
