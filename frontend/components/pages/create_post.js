import { renderHeader } from "../base/header.js";
import { renderSidebar } from "../base/sidebar.js";
import { renderFooter } from "../base/footer.js";
import { createPostButton } from '../buttons/сreatePostBtn.js';

export async function createPostPage() {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/main.css";
  document.head.appendChild(link);

  const link2 = document.createElement("link");
  link2.rel = "stylesheet";
  link2.href = "/frontend/src/css/createpost.css";
  document.head.appendChild(link2);

  let sidebarHtml = "";
  let headerHtml = "";
  let renderFooterHtml = "";

  try {
    await renderSidebar();
    sidebarHtml = app.innerHTML;
    await renderHeader();
    headerHtml = app.innerHTML;
    await renderFooter();
    renderFooterHtml = app.innerHTML;

    // Set the app innerHTML after all rendering is complete
    app.innerHTML = `
      <div class="mainBody">
        <div class="mainPage">
          <div class="col-1">  
            ${sidebarHtml} 
          </div>
          <div class="col-2"> 
            ${headerHtml}
            <div class="col-2 Content-flow"> 
              <form id="createPostForm" method="POST" action="">
                <div class="inputbox">
                  <input class="inputview inputbox" placeholder="Input title..." name="title" required minlength="5" maxlength="100" id="title">
                </div>

                <textarea id="ContenPostInput" name="content" class="editor-textarea" placeholder="Type your text here..." required></textarea>

                <div class="inputaTagbox">
                  <input class="inputview inputaTagbox" placeholder="Input tags..." name="categories" required minlength="2" maxlength="50" id="categories"">
                </div>
                <div class="rowBox">
                  <input type="submit" value="POST" class="click headerBnt" >
                </div>
              </form>
            
            </div>
          </div>
        </div>
        <div>${renderFooterHtml}</div>
      </div>
    `;

    const createPost = document.getElementById("createPostForm");
    if (createPost) {
      createPost.addEventListener('submit', createPostButton);
    } else {
      console.error("Form element 'createPostForm' not found in the DOM.");
    }

    
  } catch (error) {
    console.error("Error:", error);
    app.innerHTML = "<p>Error fetching body</p>";
  }
}
