import { fetchPosts } from "../../utils/api.js";
import { likePost, dislikePost } from "../buttons/like_dislike.js";
import { renderBody } from "./body.js";
import { renderTags } from "../additional/tags.js";
import { timeAgo } from "../additional/time_count.js";

export async function renderAllPosts(sortBy = "Newest") {
  
  const app = document.getElementById("app");

  try {
    await renderTags(sortBy);
    

    const response = await fetchPosts();
    let posts = response.payload.posts;
    const pagination = response.pagination;


    if (sortBy === "Newest") {
      posts.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    } else if (sortBy === "Oldest") {
      posts.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
    } else if (sortBy === "Popular") {
      posts.sort((a, b) => b.rate.rate - a.rate.rate);
    }

    // Function to format date
    function formatDate(dateString) {
      const date = new Date(dateString);
      const day = String(date.getDate()).padStart(2, "0");
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const year = date.getFullYear();
      return `${day}.${month}.${year}`;
    }

    app.innerHTML = `
      <div class="addinfo">
          <div class="headline"> All Posts </div>
      </div>
      <div class="addinfoframe">
        <div class="addinfo">
          <div class="infostatic">Amount</div>
          <div class="inforender">${pagination.total_count}</div>
        </div>
      </div>
      <div class="line"></div>
      <div class="addinfoframe">
        <div class="addinfo">
          ${app.innerHTML}
        </div>
      </div>
      
      ${posts
        .map(
          (post) => `
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
          `
        )
        .join("")}
    `;
  } catch (error) {
    console.error("Error fetching posts:", error);
    app.innerHTML += `<p>Error fetching posts: ${error.message}</p>`;
  }
}
