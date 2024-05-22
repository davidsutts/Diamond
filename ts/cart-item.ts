import { html, css, CSSResult, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

declare module "./modules/twlit.js" {
	export const TWStyles: CSSResult;
}

import { TWStyles } from "./modules/twlit.js";

@customElement('cart-item')
export class cartItem extends LitElement {

	static styles = TWStyles;

	@property({ type: Number })
	price?: number;

	@property({ type: String })
	tier?: string;

	render() {
		return this.price && this.tier
			? html`
			<div class="rounded-md pb-5 w-full flex gap-5">
				<div class="flex flex-col gap-1">
					<h2 class="font-bold capitalize">${this.tier ? this.tier + " Subscription" : "loading"}</h2>
					<div>${this.price ? this.price + "/month" : "loading"}</div>
				</div>
			</div>
		`
			: html`
			<div class="h-20"></div>
		`
	}

}