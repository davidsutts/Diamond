document.addEventListener("DOMContentLoaded", () => {
	document.getElementById("sub-prem").addEventListener("click", submit);
	document.getElementById("sub-super").addEventListener("click", submit);
});

function submit(e) {
	let tier = e.target.id.split("-", 2)[1];
	switch (tier) {
		case "prem":
			tier = "premium";
			break;
		case "super":
			tier = "super";
			break;
		default:
			break;
	}

	const data = {
		tier: tier
	};

	const query = new URLSearchParams(data).toString();

	window.location.href = "/checkout?" + query;
}