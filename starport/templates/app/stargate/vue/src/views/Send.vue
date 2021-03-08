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
			return this.$store.getters['common/wallet/address']
		}
	}
}
</script>
