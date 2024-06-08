const sidebar = document.getElementById("sidebar");
const hamburgerButton = document.getElementById("hamburger-button");

function openSidebar() {
  sidebar.classList.remove("-translate-x-full");
  sidebar.classList.add("translate-x-0");
}

function closeSidebar() {
  sidebar.classList.remove("translate-x-0");
  sidebar.classList.add("-translate-x-full");
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
  if (sidebar.classList.contains("-translate-x-full")) {
    openSidebar();
  } else {
    closeSidebar();
  }
});

document.addEventListener("htmx:beforeRequest", (e) => {
  changeActiveLink(e);
  blurImages();
  closeSidebar();

  document.title = e.detail.elt.textContent;
});