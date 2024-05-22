document.addEventListener("DOMContentLoaded", () => {
	document.getElementById("sub-prem").addEventListener("click", submit);
	document.getElementById("sub-super").addEventListener("click", submit);
});

function submit(e) {
	let tier = e.target.id.split("-", 2)[1];

	console.log(tier);

	const data = {
		tier: tier
	};

	const query = new URLSearchParams(data).toString();

	window.location.href = "/?" + query;
}