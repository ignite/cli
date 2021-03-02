<template>
	<div>
		<div class="container">
			<SpBankBalances :address="address" :refresh="true" v-if="address" />
		</div>
	</div>
</template>

<script>
export default {
	computed: {
		address() {
			return this.$route.params.address
		}
	}
}
</script>
