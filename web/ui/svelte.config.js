import adapter from '@sveltejs/adapter-auto';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter(),
		alias: {
			$components: 'src/lib/components',
			$utils: 'src/lib/utils',
			'$utils.js': 'src/lib/utils.ts'
		}
	}
};

export default config;
