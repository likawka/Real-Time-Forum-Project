import { fetchPosts } from "../../utils/api.js";
import { renderSearch } from "../pages/renderSearch.js"; // Import renderSearch function

export async function renderHeader() {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/header.css";
  document.head.appendChild(link);

  try {
    // Fetch authentication status
    const response = await fetchPosts(); // Assuming fetchPosts() fetches authentication status
    const { authenticated } = response;
    const userNickname = localStorage.getItem('userNickname');

    // Determine the text to display based on authentication status
    const buttonText = authenticated ? `
    <a href="/users/${userNickname}" class="headerBnt click">
          <img class="mainIcon" src="/frontend/src/images/person.svg" />
          PROFILE
        </a>` : 
        `<a href="/auth/login" class="headerBnt click">
          <img class="mainIcon" src="/frontend/src/images/person.svg" />
          LOG IN
        </a>`;

    app.innerHTML = `
      <header>
        <form id="searchForm" action="#" method="GET" class="headerSearchBarFrame">
          <img class="mainIcon" src="/frontend/src/images/search.svg" />
          <div class="headerSearchBar">
            <input type="search" id="search" name="search" class="searchInput" placeholder="Search..." />
          </div>
        </form>
        <div id="searchHints" class="searchHints">
          <div class="searchHintItem static">
            <div class="searchHintItem static">[tag] search within a tag</div>
          </div>
        </div>
      
          ${buttonText}
        
      </header>
    `;

    document.getElementById('searchForm').addEventListener('submit', async (event) => {
      event.preventDefault(); // Prevent the default form submission
      const query = document.getElementById('search').value.trim();
      await renderSearch(query); // Call the search function
    });

  } catch (error) {
    console.error("Error:", error);
    app.innerHTML = "<p>Error fetching header</p>";
  }
}
