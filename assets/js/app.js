const sidebar = document.getElementById("sidebar");
const backdrop = document.getElementById("sidebar-backdrop");
const languageSearch = document.getElementById("language-search");
const lightbox = document.getElementById("lightbox");
const lightboxImg = document.getElementById("lightbox-img");
const lightboxCaption = document.getElementById("lightbox-caption");

let currentImageIndex = 0;
let currentImages = [];

function openSidebar() {
  sidebar.classList.remove("-translate-x-full");
  sidebar.classList.add("translate-x-0");
  backdrop.classList.add("active");
  const hamburgerButton = document.getElementById("hamburger-button");
  const closeButton = document.getElementById("close-button");
  if (hamburgerButton) hamburgerButton.classList.add("hidden");
  if (closeButton) closeButton.classList.remove("hidden");
}

function closeSidebar() {
  sidebar.classList.remove("translate-x-0");
  sidebar.classList.add("-translate-x-full");
  backdrop.classList.remove("active");
  const hamburgerButton = document.getElementById("hamburger-button");
  const closeButton = document.getElementById("close-button");
  if (hamburgerButton) hamburgerButton.classList.remove("hidden");
  if (closeButton) closeButton.classList.add("hidden");
}

function isSidebarOpen() {
  return sidebar.classList.contains("translate-x-0");
}

function updateActiveLink(href) {
  document.querySelectorAll("#language-list a").forEach((a) => {
    const isMatch = a.getAttribute("href") === href;
    if (isMatch) {
      a.classList.add("ring-2", "bg-pink-500");
      a.classList.remove("bg-pink-400", "hover:bg-pink-500");
    } else {
      a.classList.remove("ring-2", "bg-pink-500");
      a.classList.add("bg-pink-400", "hover:bg-pink-500");
    }
  });
}

function blurImages() {
  document
    .querySelectorAll("main img")
    .forEach((img) => img.classList.add("blur-sm"));
}

function unblurImages() {
  document
    .querySelectorAll("main img")
    .forEach((img) => img.classList.remove("blur-sm"));
}

function getFilename(path) {
  const name = path.split("/").pop();
  return decodeURIComponent(name).replace(/\.[^.]+$/, "");
}

function openLightbox(index) {
  if (!currentImages.length) return;
  currentImageIndex = index;
  updateLightboxImage();
  lightbox.classList.add("active");
  document.body.style.overflow = "hidden";
}

function closeLightbox() {
  lightbox.classList.remove("active");
  document.body.style.overflow = "";
}

function updateLightboxImage() {
  lightboxImg.src = currentImages[currentImageIndex];
  const filename = getFilename(currentImages[currentImageIndex]);

  lightboxCaption.textContent = currentImageIndex + 1 + " / " + currentImages.length + " — " + filename;
}

function prevImage() {
  currentImageIndex = (currentImageIndex - 1 + currentImages.length) % currentImages.length;
  updateLightboxImage();
}

function nextImage() {
  currentImageIndex = (currentImageIndex + 1) % currentImages.length;
  updateLightboxImage();
}

function filterLanguages(query) {
  const term = query.toLowerCase().trim();
  let visible = 0;

  document
    .querySelectorAll("#language-list li[data-language]")
    .forEach((li) => {
      const matches = li.dataset.language.toLowerCase().includes(term);
      li.style.display = matches ? "" : "none";
      if (matches) visible++;
    });

  const empty = document.getElementById("language-search-empty");
  if (empty) empty.classList.toggle("hidden", visible > 0);
}

function initPageContent() {
  const backToTop = document.getElementById("back-to-top");
  const scroller = document.querySelector("main .overflow-y-auto");
  const grid = document.getElementById("image-grid");

  if (grid) {
    const figures = grid.querySelectorAll("figure");
    currentImages = Array.from(figures).map((f) => f.dataset.image);
    figures.forEach((figure, index) => {
      figure.addEventListener("click", (e) => {
        e.preventDefault();
        openLightbox(index);
      });
    });
  } else {
    currentImages = [];
  }

  if (scroller) {
    scroller.addEventListener("scroll", () => {
      if (backToTop)
        backToTop.classList.toggle("hidden", scroller.scrollTop <= 400);
    });
  }

  if (backToTop) {
    backToTop.addEventListener("click", () => {
      if (scroller) scroller.scrollTo({ top: 0, behavior: "smooth" });
    });
  }
}

document.addEventListener("click", (e) => {
  if (e.target.closest("#hamburger-button")) {
    openSidebar();
    return;
  }

  if (e.target.closest("#close-button")) {
    closeSidebar();
    return;
  }

  if (isSidebarOpen() && !sidebar.contains(e.target)) {
    closeSidebar();
  }
});

if (backdrop) backdrop.addEventListener("click", closeSidebar);

if (languageSearch) {
  languageSearch.addEventListener("input", (e) =>
    filterLanguages(e.target.value),
  );
}

if (lightbox) {
  lightbox.addEventListener("click", (e) => {
    if (e.target === lightbox) closeLightbox();
  });
}

document
  .getElementById("lightbox-close")
  ?.addEventListener("click", closeLightbox);

document.getElementById("lightbox-prev")?.addEventListener("click", (e) => {
  e.stopPropagation();
  prevImage();
});
document.getElementById("lightbox-next")?.addEventListener("click", (e) => {
  e.stopPropagation();
  nextImage();
});

document.addEventListener("keydown", (e) => {
  if (!lightbox.classList.contains("active")) return;
  if (e.key === "Escape") closeLightbox();
  if (e.key === "ArrowLeft") prevImage();
  if (e.key === "ArrowRight") nextImage();
});

document.addEventListener("htmx:beforeRequest", (e) => {
  const elt = e.detail.elt;
  if (elt.tagName === "A" && elt.closest("#language-list")) {
    updateActiveLink(elt.getAttribute("href"));
    blurImages();
    closeSidebar();
    if (languageSearch) {
      languageSearch.value = "";
      filterLanguages("");
    }
    document.title = elt.textContent.trim() + " — Waifus";
  }
});

document.addEventListener("htmx:afterSwap", () => {
  unblurImages();
  const scroller = document.querySelector("main .overflow-y-auto");
  if (scroller) scroller.scrollTo({ top: 0, behavior: "auto" });
  const backToTop = document.getElementById("back-to-top");
  if (backToTop) backToTop.classList.add("hidden");
  initPageContent();
});

initPageContent();