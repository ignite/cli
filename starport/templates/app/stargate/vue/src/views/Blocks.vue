<template>
	<div class="container">
		<table class="SpTable SpBlocksTable">
			<thead>
				<tr>
					<th><strong>HEIGHT</strong></th>
					<th><strong>HASH</strong></th>
					<th><strong>TIME</strong></th>
					<th><strong>TXs</strong></th>
				</tr>
			</thead>
			<tbody>
				<SpBlockDisplayLine
					:block="block"
					tsFormat="YYYY-MM-DD HH:mm:ss"
					v-for="block in blocks"
					v-bind:key="block.hash"
				/>
			</tbody>
		</table>
		<div class="SpPagination">
			<div class="SpPaginationTitle">PAGES</div>
			<router-link to="/blocks/1" v-if="page >= 2">
				<button class="SpButton">
					<div class="SpButtonText">&lt;&lt;</div>
				</button>
			</router-link>
			<router-link :to="'/blocks/' + (page - 1)" v-if="page >= 2">
				<button class="SpButton">
					<div class="SpButtonText">&lt;</div>
				</button>
			</router-link>
			<span class="SpPaginationItem active">
				{{ page }}
			</span>
			<router-link :to="'/blocks/' + (page + 1)" v-if="page < pages">
				<button class="SpButton">
					<div class="SpButtonText">&gt;</div>
				</button>
			</router-link>
			<router-link :to="'/blocks/' + pages" v-if="page < pages">
				<button class="SpButton">
					<div class="SpButtonText">&gt;&gt;</div>
				</button>
			</router-link>
		</div>
	</div>
</template>

<script>
import axios from 'axios'

export default {
	data() {
		return {
			blocks: [],

			pages: 1
		}
	},

	computed: {
		page() {
			return parseInt(this.$route.params.page) || 1
		}
	},
	watch: {
		page: async function (newPage) {
			this.blocks = []
			const chain = await axios.get(
				this.$store.getters['common/env/apiTendermint'] +
					'/blockchain?minHeight=1&maxHeight=20'
			)
			const lowest = parseInt(
				chain.data.result.block_metas[chain.data.result.block_metas.length - 1]
					.header.height
			)

			const highest = parseInt(chain.data.result.last_height)
			this.pages = Math.ceil((highest - lowest + 1) / 20)
			let from
			if (highest + 1 - this.page * 20 >= 1) {
				from = highest + 1 - this.page * 20
			} else {
				from = 1
			}
			const page = await axios.get(
				this.$store.getters['common/env/apiTendermint'] +
					'/blockchain?minHeight=' +
					from +
					'&maxHeight=' +
					(highest - (newPage - 1) * 20)
			)
			for (let block_meta of page.data.result.block_metas) {
				const block = {
					height: block_meta.header.height,
					timestamp: block_meta.header.time,
					hash: block_meta.block_id.hash,
					details: { num_txs: block_meta.num_txs }
				}
				this.blocks.push(block)
			}
		}
	},
	async mounted() {
		const chain = await axios.get(
			this.$store.getters['common/env/apiTendermint'] +
				'/blockchain?minHeight=1&maxHeight=20'
		)
		const lowest = parseInt(
			chain.data.result.block_metas[chain.data.result.block_metas.length - 1]
				.header.height
		)

		const highest = parseInt(chain.data.result.last_height)
		this.pages = Math.ceil((highest - lowest + 1) / 20)
		let from
		if (highest + 1 - this.page * 20 >= 1) {
			from = highest + 1 - this.page * 20
		} else {
			from = 1
		}
		const page = await axios.get(
			this.$store.getters['common/env/apiTendermint'] +
				'/blockchain?minHeight=' +
				from +
				'&maxHeight=' +
				(highest - (this.page - 1) * 20)
		)
		for (let block_meta of page.data.result.block_metas) {
			const block = {
				height: block_meta.header.height,
				timestamp: block_meta.header.time,
				hash: block_meta.block_id.hash,
				details: { num_txs: block_meta.num_txs }
			}
			this.blocks.push(block)
		}
	}
}
</script>
