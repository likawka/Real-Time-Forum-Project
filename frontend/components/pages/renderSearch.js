// File: renderSearch.js

import { fetchPostsSearch , fetchPosts } from "../../utils/api.js"; // Import fetchPostsSearch function
import { searchPosts } from "../buttons/search.js"; // Import the search function
import { renderHeader } from "../base/header.js";
import { renderSidebar } from "../base/sidebar.js";
import { renderFooter } from "../base/footer.js";
import { timeAgo } from "../additional/time_count.js";
import { renderTags } from "../additional/tags.js";
import { formatDate } from "../additional/formatDate.js";

export async function renderSearch(query) {
  const app = document.getElementById("app");

  try {
    await renderTags();
    const renderTagsHtml = app.innerHTML;
    await renderSidebar();
    const sidebarHtml = app.innerHTML;
    await renderHeader();
    const headerHtml = app.innerHTML;
    await renderFooter();
    const footerHtml = app.innerHTML;

    const response = await fetchPosts();
    const postsAll = response.payload.posts;
    var pagination = response.pagination;

    // Perform the search
    const posts = await searchPosts(query);

    // Generate the HTML for the posts
    const postsHtml = posts.map(post => `
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
                              <div class="tag click" onclick="window.location.href = '/users/${post.nickname}';" >${post.nickname}</div>
                              <div class="inforender">${formatDate(post.created_at)}</div>
                              <img class="shareIcon click" src="/frontend/src/images/share.svg" alt="Copy Link">
                          </div>
                      </div>
                  </div>
          </div>
          <div class="line"></div>
    `).join('');

    const postHtml = `
      <div class="mainBody">
        <div class="mainPage">
          <div class="col-1">${sidebarHtml}</div>
          <div class="col-2">
            ${headerHtml}
            <div class="col-2 Content-flow">

              ${postsHtml ? `
                <div class="addinfo">
                  <div class="headline"> Founds by  </div>
                  <div class="headline" style="color: var(--light-red);">"${query}"</div>
                </div>
                <div class="line"></div>`
              : ''}

              ${postsHtml || '<div class="headline"> No posts found matching your search. </div>'}
              
            </div>
          </div>
        </div>
        <div>${footerHtml}</div>
      </div>
    `;

    app.innerHTML = postHtml;

    // Optionally, update the URL with the search query
    history.pushState(null, null, `/?search=${encodeURIComponent(query)}`);

  } catch (error) {
    console.error("Error fetching posts:", error);
    app.innerHTML += `<p>Error fetching posts: ${error.message}</p>`;
  }
}
