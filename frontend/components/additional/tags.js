import { fetchPosts } from "../../utils/api.js";

export async function renderTags(sortBy = "Newest") {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/main.css";
  document.head.appendChild(link);

  try {
    // Render sorting buttons
    app.innerHTML = `
      <div class="addinfoframe">
      </div>
    `;

    // Attach event listeners to buttons
    document.querySelectorAll(".tagsort").forEach(button => {
      button.addEventListener("click", async (event) => {
        const selectedSortBy = event.target.getAttribute("data-sortby");
        await renderAllPosts(selectedSortBy); // Re-render posts with the selected sort option
      });
    });

  } catch (error) {
    console.error("Error:", error);
    app.innerHTML = "<p>Error fetching tags</p>";
  }
}
