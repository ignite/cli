<template>
	<div>
		<div class="container">
			<SpBlockDisplayFull :block="block" v-if="block" />
		</div>
	</div>
</template>

<script>
import axios from 'axios'

export default {
	data() {
		return {
			block: null
		}
	},
	async created() {
		const blockDetails = await axios.get(
			this.$store.getters['common/env/apiTendermint'] +
				'/block?height=' +
				this.$route.params.block
		)

		const txDecoded = blockDetails.data.result.block.data.txs.map(async (x) => {
			const dec = await this.$store.getters[
				'common/env/apiClient'
			].decodeTx(x)
			return dec
		})
		const txs = await Promise.all(txDecoded)
		this.block = {
			height: blockDetails.data.result.block.header.height,
			timestamp: blockDetails.data.result.block.header.time,
			hash: blockDetails.data.result.block_id.hash,
			details: blockDetails.data.result.block,
			txDecoded: [...txs]
		}
	}
}
</script>
