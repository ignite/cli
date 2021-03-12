<template>
	<div>
		<div class="container">
			<SpTokenSend :address="address" />
			<SpTransferList :address="address" />
			<!-- this line is used by starport scaffolding # 4 -->
		</div>
	</div>
</template>

<style scoped>
.container {
	max-width: 800px;
	margin: 2rem auto;
	padding: 1rem;
}
</style>
<script>
export default {
	name: 'Index',
	computed: {
		address() {
			return this.$store.getters['common/wallet/address']
		}
	}
}
</script>
