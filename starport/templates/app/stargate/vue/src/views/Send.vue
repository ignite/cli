<template>
	<div>
		<div class="container">
			<SpBankBalances :address="address" :refresh="true" />
			<SpTokenSend :address="address" />
		</div>
	</div>
</template>

<script>
export default {
	computed: {
		address() {
			return this.$store.getters['chain/common/wallet/address']
		}
	}
}
</script>
