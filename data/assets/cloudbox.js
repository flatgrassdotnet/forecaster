function isScrolling(element) { // adapted from https://www.geeksforgeeks.org/check-whether-html-element-has-scrollbars-using-javascript/
	var res = !! element["scrollTop"]; 
	if (!res) { 
		element["scrollTop"] = 1; 
		res = !!element["scrollTop"]; 
		element["scrollTop"] = 0; 
	} 
	return res; 
}

var items = document.querySelectorAll(".item");
var metabox = document.querySelector(".metabox");
var metaboxTarget = metabox.querySelector(".metabox-target");
itemHover = function(e) {
	metabox.classList.add("active");
	if (metabox.dataset.itemid != this.dataset.itemid) {
		metabox.dataset.itemid = this.dataset.itemid;
		metaboxTarget.innerHTML = this.querySelector(".metabox-template").innerHTML;
		
		var metaDesc = metaboxTarget.querySelector(".meta-desc");
		while (isScrolling(metaDesc)) { // Make sure the description fits
			var scale = Number(metaDesc.dataset.scale);
			if (scale >= 14) {
				if (metaDesc.classList.contains("wrap")) {
					break;
				}
				metaDesc.classList.add("wrap");
				scale = -1;
			}
			scale++;
			metaDesc.dataset.scale = scale;
		}
	}
};
itemLeave = function(e) {
	metabox.classList.remove("active");
};
for (i=0; i<items.length; i++) {
	items[i].addEventListener("mouseover", itemHover);
	items[i].addEventListener("mouseout", itemLeave);
}


var content = document.querySelector(".content");
var bottomclouds = document.querySelector("#bottomclouds");
function checkScrollbar() {
	var scrollbarWidth = content.offsetWidth - content.clientWidth;
	bottomclouds.style.right = scrollbarWidth + "px";
	metabox.style.right = scrollbarWidth + "px";
}
window.addEventListener("resize", checkScrollbar);
checkScrollbar();


function cbalert(title, txt) {
	if (typeof txt === "undefined") {
		txt = title;
		title = "";
	}
	var dialog = document.querySelector("#alert");
	dialog.querySelector("h3").innerText = title;
	dialog.querySelector("p").innerText = txt;
	
	if (document.documentElement.classList.contains("no-dialog")) {
		cbalertfallbackModal();
	} else {
		dialog.showModal();
	}
}

function cbalertfallbackModal() {
	var dialog = document.querySelector("#alert");
	var backdrop = document.querySelector("#backdrop");
	dialog.setAttribute("open", "open");
	backdrop.setAttribute("open", "open");
}

function cbalertClose() {
	var dialog = document.querySelector("#alert");
	if (document.documentElement.classList.contains("no-dialog")) {
		var backdrop = document.querySelector("#backdrop");
		dialog.removeAttribute("open");
		backdrop.removeAttribute("open");
	} else {
		dialog.close();
	}
}


if (typeof HTMLDialogElement !== 'function') {
	// Fallback handling for Awesomium and GM12
	document.documentElement.classList.add("no-dialog");
	
	var backdrop = document.createElement("div");
	backdrop.setAttribute("id","backdrop");
	backdrop.addEventListener("click", cbalertClose);
	document.body.insertBefore(backdrop, document.querySelector("#alert"));
	
	document.addEventListener("keydown", function(e){
		if (event.keyCode == 27) {
			cbalertClose();
		}
	});
}