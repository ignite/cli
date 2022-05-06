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
		this.error = () => {
			return this;
		};
		this.warn = () => {
			return this;
		};
		this.info = (msg) => {
			if (msg.indexOf('Relay') == 0 && msg.indexOf('Relay 0') == -1) {
				console.log(msg);
			}
			return this;
		};
		this.verbose = () => {
			return this;
		};
		this.debug = () => {
			return this;
		};
	}
}
