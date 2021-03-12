# Run a blockchain

<!-- what goes here? -->
(this old content seems to introduce a directory structure, do we want to talk about that structure? no edits yet)
When you start the blockchain app with the `starport serve`, the blockchain folder user homefolder `~/.myappd` (the name of your app with a `d` for `daemon` attached) and initiate your blockchain with the genesis file, located under `~/.myappd/config`. The second folder you can find in the `~/.myappd` folder is `data` - this is where the blockchain will write the consecutive blocks and transactions. The other folder created is the `~/.myappcli` folder, which contains a configuration file for your current command line interface, such as `chain-id`, output parameters such as `json` or `indent` mode. If you want to make sure all of your data from the blockchain setup is deleted, make sure to remove the `~/.myappd` and `~/.myappcli` folder.
