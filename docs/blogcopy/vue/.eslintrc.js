module.exports = {
	root: true,
	env: {
		node: true
	},
	extends: [
		'plugin:vue/vue3-essential',
		'eslint:recommended',
		'plugin:prettier/recommended',
		'@vue/prettier'
	],
	parserOptions: {
		parser: 'babel-eslint'
	},
	plugins: ['prettier'],
	rules: {
		'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
		'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
		'no-unused-vars': 'off',
		'vue/component-name-in-template-casing': ['error', 'PascalCase']
	}
}
