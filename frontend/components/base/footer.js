// import { fetchPosts } from "../../utils/api.js";

export async function renderFooter() {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/header.css";
  document.head.appendChild(link);

  try {
    app.innerHTML = `
      <footer>
          <div>  </div>
      </footer>
      `;

    function isPageFullyScrolled() {
      const pageHeight = document.documentElement.scrollHeight;
      const windowHeight = window.innerHeight;
      const scrollY = window.scrollY || window.pageYOffset;
      return pageHeight - windowHeight - scrollY <= 0;
    }

    // Функція для показу або приховування футера в залежності від прокрутки
    function toggleFooterVisibility() {
      const footer = document.querySelector("footer");

      if (isPageFullyScrolled()) {
        footer.classList.remove("hiddenElement"); // Показати футер
      } else {
        footer.classList.add("hiddenElement"); // Приховати футер
      }
    }
    // Викликати toggleFooterVisibility() при прокрутці сторінки
    window.addEventListener("scroll", toggleFooterVisibility);

    // Викликати toggleFooterVisibility() при завантаженні сторінки, щоб ініціювати перевірку на початку
    window.addEventListener("load", toggleFooterVisibility);
  } catch (error) {
    console.error("Error:", error);
    app.innerHTML = "<p>Error fetching footer</p>";
  }
}
