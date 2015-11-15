var baseNavHeight = 4; //em
var animationTime = 300; //ms

document.addEventListener('load', load());

function load(){
	
	/*
	animating nav bar toggle
	*/

	var button = document.querySelector('.navbar__button');
	var collapse = document.querySelector('.navbar__collapse');
	var collapseList = document.querySelector('.navbar__collapse>.navbar__item-list');
	var navbar = document.querySelector('.navbar');

	button.addEventListener('click', function(){
		navbarToggle();
	});

	/*
	client side validation
	*/

	var form = document.querySelector("#p2d-form");
	var name = document.querySelector("#p2d-name");
	var sequence = document.querySelector("#p2d-sequence");
	var range = document.querySelector("#p2d-aa-range");
	var features = document.querySelectorAll("input[name=\"Features\"]");
	var featuresLabel = document.querySelector("#p2d-features-label");
	var seqPattern = new RegExp("^([AaCcDdEeFfGgHhIiKkLlMmNnPpQqRrSsTtUuVvWwYy]+)$");
	var rangePattern = new RegExp("^([0-9]+)-([0-9]+)$")

	form.addEventListener('submit',function(e){
		validate(e);
	});

	function navbarToggle(){

		switch (navbar.style.height != false){

			case true:
				navbar.style.height = '';
				break;

			case false:
				var itemCount = collapseList.childElementCount;
				navbar.style.height = String((itemCount + 1) * baseNavHeight) + 'em';
				break;
		}
	}

	function validate(e){

		e.stopPropagation();
		e.preventDefault();

		if (sequence.value.includes("\n")){
			sequence.value = sequence.value.split("\n").join("");
		}

		if (sequence.value.includes("\r")){
			sequence.value = sequence.value.split("\r").join("");
		}


		if (!name.value){
			name.parentElement.classList.add("invalid-name");
		}

		else {
			name.parentElement.classList.remove("invalid-name");
		}

		if (!sequence.value.match(seqPattern)){
			sequence.parentElement.classList.add("invalid-sequence");
		}

		else {
			sequence.parentElement.classList.remove("invalid-sequence");
		}

		if (!range.value.match(rangePattern)){
			range.parentElement.classList.add("invalid-range");
		}

		else {
			range.parentElement.classList.remove("invalid-range");
		}

		var rangeLen = range.value.split("-");
		rangeLen = rangeLen[1] - rangeLen[0] + 1;
		var seqLen = sequence.value.length;

		if (rangeLen != seqLen){
			range.parentElement.classList.add("invalid-range-length");
		}

		else{
			range.parentElement.classList.remove("invalid-range-length");
		}

		var checked = false;

		Object.keys(features).forEach(function(index){
			checked = checked || features[index].checked
		});

		if (!checked) {
			featuresLabel.classList.add("invalid-features");
		}

		else {
			featuresLabel.classList.remove("invalid-features");
		}

		if (sequence.value.match(seqPattern)
			&& range.value.match(rangePattern)
			&& (rangeLen == seqLen)
			&& checked
			&& name.value){
			form.submit();
		}
	}

	/*
		auto-add protein range
	*/

	var fullLengthBtn = document.querySelector("#p2d-full-length");

	fullLengthBtn.addEventListener("click", function(){
		getSeqLength();
	});

	function getSeqLength(){

		//first make sure sequence exists and is valid
		if (!sequence.value.match(seqPattern)){
			sequence.parentElement.classList.add("invalid-sequence");
			return;
		}

		else {
			sequence.parentElement.classList.remove("invalid-sequence");
		}

		length = sequence.value.length;
		range.value = "1-" + String(length);
	}

}