import { renderHeader } from './header.js';
import { renderSidebar } from './sidebar.js';
import { renderAllPosts } from './posts.js';
import { renderFooter } from './footer.js';


export async function renderBody() {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/main.css";
  document.head.appendChild(link);

  let sidebarHtml = '';
  let headerHtml = '';
  let renderAllPostsHtml = '';
  let renderFooterHtml = '';

  try {
    await renderSidebar();
    sidebarHtml = app.innerHTML;
    await renderHeader();
    headerHtml = app.innerHTML;
    await renderAllPosts();
    renderAllPostsHtml = app.innerHTML;
    await renderFooter();
    renderFooterHtml = app.innerHTML;



    // Set the app innerHTML after all rendering is complete
    app.innerHTML = `
      <div class="mainBody">
        <div class="mainPage">
          <div class="col-1">  
            ${sidebarHtml} </div>
          <div class="col-2"> 
            ${headerHtml}
            <div class="col-2 Content-flow"> ${renderAllPostsHtml} </div>
          </div>
        </div>
        <div>${renderFooterHtml}</div>
      </div>
    `;
  } catch (error) {
    console.error("Error:", error);
    app.innerHTML = "<p>Error fetching body</p>";
  }
}
