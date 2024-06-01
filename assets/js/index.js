function changeActiveLink(e) {
  document.querySelectorAll(".active-link").forEach(element => {
    element.classList.remove("active-link");
  });

  e.detail.elt.classList.add("active-link");
}

function blurImages() {
  document.querySelectorAll("main img").forEach(element => {
    element.classList.add("blur-sm");
  });
}

document.addEventListener("htmx:beforeRequest", (e) => {
  changeActiveLink(e);
  blurImages();

  document.title = e.detail.elt.textContent;
});