var baseNavHeight = 4; //em
var animationTime = 300; //ms

document.addEventListener('load', load());

function load(){
	
	var button = document.querySelector('.navbar__button');
	var collapse = document.querySelector('.navbar__collapse');
	var collapseList = document.querySelector('.navbar__collapse>.navbar__item-list');
	var navbar = document.querySelector('.navbar');

	button.addEventListener('click', function(){
		navbarToggle();
	});

	function navbarToggle(){
		
		//console.log(navbar.style.height);

		switch (navbar.style.height != false){

			case true:

				console.log('hi');

				navbar.style.height = '';

				break;

			case false:

				var itemCount = collapseList.childElementCount;
				navbar.style.height = String((itemCount + 1) * baseNavHeight) + 'em';
				break;
		}
	}
}