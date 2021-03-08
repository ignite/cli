<template>
	<div v-if="initialized">
		<SpWallet ref="wallet" v-on:dropdown-opened="$refs.menu.closeDropdown()" />
		<SpLayout>
			<template v-slot:sidebar>
				<SpSidebar
					v-on:sidebar-open="sidebarOpen = true"
					v-on:sidebar-close="sidebarOpen = false"
				>
					<template v-slot:header>
						<SpLogo />
					</template>
					<template v-slot:default>
						<SpLinkIcon link="/" text="Dashboard" icon="Dashboard" />
						<SpLinkIcon link="/modules" text="Modules" icon="Modules" />
						<SpLinkIcon
							link="/transactions"
							text="Transactions"
							icon="Transactions"
						/>
						<SpLinkIcon link="/types" text="Custom Type" icon="Form" />
						<div class="sp-dash"></div>
						<SpLinkIcon link="/settings" text="Settings" icon="Settings" />
						<SpLinkIcon link="/docs" text="Documentation" icon="Docs" />
					</template>
					<template v-slot:footer>
						<SpStatusAPI :showText="sidebarOpen" />
						<SpStatusRPC :showText="sidebarOpen" />
						<SpStatusWS :showText="sidebarOpen" />
						<div class="sp-text">Build: v0.3.8</div>
					</template>
				</SpSidebar>
			</template>
			<template v-slot:content>
				<router-view />
			</template>
		</SpLayout>
	</div>
</template>

<style>
body {
	margin: 0;
}
</style>

<script>
import './scss/app.scss'
import '@starport/vue/lib/starport-vue.css'

export default {
	data() {
		return {
			initialized: false,
			sidebarOpen: true
		}
	},
	computed: {
		hasWallet() {
			return this.$store.hasModule([ 'common', 'wallet'])
		}
	},
	async created() {
		await this.$store.dispatch('common/env/init')
		this.initialized = true
	},
	errorCaptured(err) {
		console.log(err)
		return false
	}
}
</script>
