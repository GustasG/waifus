const sidebar = document.getElementById("sidebar");
const hamburgerButton = document.getElementById("hamburger-button");
const closeButton = document.getElementById("close-button");

function openSidebar() {
  sidebar.classList.remove("-translate-x-full");
  sidebar.classList.add("translate-x-0");

  hamburgerButton.classList.add("hidden");
  closeButton.classList.remove("hidden");
}

function closeSidebar() {
  sidebar.classList.remove("translate-x-0");
  sidebar.classList.add("-translate-x-full");

  hamburgerButton.classList.remove("hidden");
  closeButton.classList.add("hidden");
}

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

hamburgerButton.addEventListener("click", () => {
  openSidebar();
});

closeButton.addEventListener("click", () => {
  closeSidebar();
});

document.addEventListener("htmx:beforeRequest", (e) => {
  changeActiveLink(e);
  blurImages();
  closeSidebar();

  document.title = e.detail.elt.textContent;
});