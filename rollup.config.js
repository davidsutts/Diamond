import typescript from '@rollup/plugin-typescript';
import { nodeResolve } from '@rollup/plugin-node-resolve';

export default {
	input: 'ts/cart-item.ts',
	output: {
		dir: 'static/build',
		format: 'esm'
	},
	plugins: [typescript(), nodeResolve()]
};