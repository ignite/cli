<template>
	<div class="container">
		<div class="SpWalletInfo" v-if="wallet">
			<h3 class="SpWalletInfoTitle">WALLET INFORMATION: '{{ wallet.name }}'</h3>
			<hr />
			<p><strong>MNEMONIC</strong></p>
			<div class="SpWalletInfoMnemonic">
				<label for="show_mnemonic" v-if="!showMnemonic">Show mnemonic? </label>
				<span v-if="showMnemonic"> {{ wallet.mnemonic }} </span><br />
				<input type="checkbox" id="show_mnemonic" v-model="showMnemonic" />
			</div>
			<hr />
			<p><strong>ACCOUNTS</strong></p>
			<table class="SpTable">
				<thead>
					<tr>
						<th>ACCOUNT</th>
						<th>PATH</th>
						<th class="SpWideCell">PRIVATE KEY</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="account in wallet.accounts" v-bind:key="account.address">
						<td>
							{{ account.address }}
						</td>
						<td>
							{{ '' + wallet.HDpath + account.pathIncrement }}
						</td>
						<td>
							<label
								:for="account.address + '_key'"
								v-if="!viewTogglers[account.address]"
								>Show Private Key?
							</label>
							<span v-if="viewTogglers[account.address]">
								{{ pKeys[account.address] }} </span
							><br />
							<input
								type="checkbox"
								:id="account.address + '_key'"
								v-model="viewTogglers[account.address]"
							/>
						</td>
					</tr>
				</tbody>
			</table>

			<button @click="downloadBackup()" class="SpButton">
				<div class="SpButtonText">DOWNLOAD BACKUP</div>
			</button>
		</div>
	</div>
</template>

<script>
import moment from 'moment'
import { saveAs } from 'file-saver'
import CryptoJS from 'crypto-js'
import {
	Bip39,
	EnglishMnemonic,
	stringToPath,
	Slip10,
	Slip10Curve
} from '@cosmjs/crypto'
import { keyToWif } from '@starport/vuex'

export default {
	name: 'Wallet',
	data() {
		return {
			viewTogglers: {},
			pKeys: {},
			showMnemonic: false
		}
	},
	computed: {
		wallet() {
			return this.$store.getters['common/wallet/wallet']
		}
	},
	async created() {
		for (let account of this.wallet.accounts) {
			this.pKeys[account.address] = await this.getPrivateKey(
				account.pathIncrement
			)
		}
	},
	methods: {
		downloadBackup() {
			const backup = CryptoJS.AES.encrypt(
				JSON.stringify(this.wallet),
				this.wallet.password
			)

			const blob = new Blob([backup.toString()], {
				type: 'application/octet-stream; charset=us-ascii'
			})
			saveAs(blob, this.backupName())
		},
		getPrivateKey(pathIncrement) {
			const mnemonicChecked = new EnglishMnemonic(this.wallet.mnemonic)
			return Bip39.mnemonicToSeed(mnemonicChecked).then((seed) => {
				const hdPath = stringToPath(this.wallet.HDpath + pathIncrement)
				const { privkey } = Slip10.derivePath(
					Slip10Curve.Secp256k1,
					seed,
					hdPath
				)
				return keyToWif(privkey)
			})
		},
		backupName() {
			return (
				this.wallet.name + '_Backup_' + moment().format('YYYY-MM-DD') + '.bin'
			)
		}
	}
}
</script>
