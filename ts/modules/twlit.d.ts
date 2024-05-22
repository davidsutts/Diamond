import { CSSResult } from "lit";
import { TWStyles } from "./twlit.js";

declare module "twstyles" {
	const TWStyles: CSSResult;
	export default TWStyles;
}