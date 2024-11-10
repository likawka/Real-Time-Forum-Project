import { renderHeader } from '../base/header.js';
import { renderSidebar } from '../base/sidebar.js';
import { renderFooter } from '../base/footer.js';
import { fetchUserData } from '../../utils/api.js';
import { renderTags } from '../additional/tags.js';
import { timeAgo } from "../additional/time_count.js";
import { formatDate } from "../additional/formatDate.js";

export async function renderProfile(nickname) {
  const app = document.getElementById("app");

  // Load CSS dynamically
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/main.css";
  document.head.appendChild(link);

  let sidebarHtml = '';
  let headerHtml = '';
  let renderFooterHtml = '';

  try {
    await renderSidebar();
    sidebarHtml = app.innerHTML; 
    await renderHeader();
    headerHtml = app.innerHTML; 
    await renderFooter();
    renderFooterHtml = app.innerHTML; 
    await renderTags();
    const renderTagsHtml = app.innerHTML;

    const response = await fetchUserData(nickname, 'posts');
    const user = response.payload.user;
    const posts = response.payload.posts;
    const comments = response.payload.comments;

    const isMe = user.id === parseInt(localStorage.getItem('userId'));

    // Functions to render posts and comments
    function renderAllPosts(posts) {
      const container = document.createElement('div');
      container.className = 'post-container';

      posts.forEach(post => {
        const postBox = document.createElement('div');
        postBox.className = 'postBox';
        postBox.innerHTML = `
            <div class="questionbox">
                    <div class="infobox">
                        <div class="infoBoxRow">
                            <img class="mainIcon" src="/frontend/src/images/comments.svg" />
                            <div class="infostatic">${post.amount_of_comments}</div>
                        </div>
                        <div class="infoBoxRow">
                            <img class="mainIcon" src="/frontend/src/images/rate.svg" />
                            <div class="infostatic">${post.rate.rate}</div>
                        </div>
                    </div>
                    <div class="questContainer">
                        <div class="title click"><a href="/post/${post.id}">${post.title}</a></div>
                        <div class="questtextbox">${post.content}</div>
                        
                        <div class="addinfoframe-row">
                            <div class="tags">
                                ${post.categories
                .map(
                  (category) => `
                                    <div class="tag tagview click">${category.name}</div>
                                `
                )
                .join("")}
                            </div>
                            <div class="postedby">
                                <div class="tag">Posted by</div>
                                <div class="tag click">${post.nickname}</div>
                                <div class="inforender">${formatDate(post.created_at)}</div>
                                <img class="shareIcon click" src="/frontend/src/images/share.svg" alt="Copy Link">
                            </div>
                        </div>
                    </div>
            </div>
        `;
        container.appendChild(postBox);

        const line = document.createElement('div');
        line.className = 'line';
        container.appendChild(line);
      });

      return container.innerHTML;
    }

    function renderNoPosts() {
      return `
        <div class="addinfo noPostClass">
          <div class="headline">No</div>
          <div class="headline" style="color: var(--light-red);">Posts</div>
        </div>
      `;
    }

    // Determine what to render based on the posts and comments
    let contentHtml = '';
    if (!posts || posts.length === 0) {
      contentHtml = renderNoPosts();
    } else {
      contentHtml = renderAllPosts(posts);
    }

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
              <div class="addinfo">
                ${isMe ? "<div class=\"headline\">Hello  </div>" : "<div class=\"headline\">Hello from </div>" }
                <div class="headline" style="color: var(--light-red);">${user.nickname}</div>
              </div>
              <div class="addinfoframe1">
                <div class="addinfoframe">
                  <div class="addinfo">
                    <div class="infostatic">Enrolled</div>
                    <div class="inforender">${timeAgo(user.created_at)}</div>
                    <div class="infostatic">Posted</div>
                    <div class="inforender">${user.amount_of_posts} times</div>
                    <div class="infostatic">Commented</div>
                    <div class="inforender">${user.amount_of_comments} times</div>
                  </div>
                </div>
              </div>
              <div class="line"></div>
            <div class="addinfoframe-row">
              <div class="addinfo">
                  ${renderTagsHtml}
              </div>
              <div class="addinfo">
                  ${isMe ? " " : "<input type=\"submit\"  class=\"tagsort tagview click\" name=\"SortBy\" value=\"Text me\"> " }
              </div>
            </div>

              ${contentHtml}

            </div>
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
