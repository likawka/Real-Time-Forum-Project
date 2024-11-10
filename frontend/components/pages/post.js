// File: renderPost.js

import { fetchPost } from "../../utils/api.js";
import { likePost, dislikePost } from "../buttons/like_dislike.js";
import { postAnswer } from "../buttons/post_btn.js";
import { renderHeader } from "../base/header.js";
import { renderSidebar } from "../base/sidebar.js";
import { renderFooter } from "../base/footer.js";
import { timeAgo } from "../additional/time_count.js";
import { renderTags } from "../additional/tags.js";
import { renderComments } from "../base/comments.js";

export async function renderPost(postID) {
  const app = document.getElementById("app");

  try {
    await renderSidebar();
    const sidebarHtml = app.innerHTML;
    await renderHeader();
    const headerHtml = app.innerHTML;
    await renderTags();
    const renderTagsHtml = app.innerHTML;
    await renderFooter();
    const renderFooterHtml = app.innerHTML;

    const response = await fetchPost(postID);
    const post = response.payload.post;
    const pagination = response.pagination;

    const postHtml = `
      <div class="mainBody">
        <div class="mainPage">
          <div class="col-1">${sidebarHtml}</div>
          <div class="col-2"> 
            ${headerHtml}
            <div class="col-2 Content-flow"> 
              <div class="questionbox-column">
                <div class="questContainer">
                  <div class="headline">${post.title}</div>
                  <div class="addinfoframe-row">
                    <div class="addinfo">
                      <div class="infostatic">Asked</div>
                      <div class="inforender">${timeAgo(post.created_at)}</div>
                    </div>
                    <div class="shareFrame">
                      <img class="shareIcon click" onclick="copyLink('/post/${post.id}')" src="/frontend/src/images/share.svg" alt="Copy Link">
                    </div>
                  </div>
                </div>
                <div class="line"></div>
                <div class="questionbox">
                  <div class="ratebox">
                    <button onclick="likePost('${post.id}')" style="background:none; border:none;">
                      <img src="/frontend/src/images/like.svg" class="rateicon click" alt="Like">
                    </button>
                    <div class="rateNum" id="rateNum-${post.id}">${post.rate.rate}</div>
                    <button onclick="dislikePost('${post.id}')" style="background:none; border:none; transform: rotate(180deg)">
                      <img src="/frontend/src/images/like.svg" class="rateicon click" alt="Dislike">
                    </button>
                  </div>
                  <div class="questContainer">
                    <div class="questtextbox">${post.content}</div>
                    <div class="addinfoframe-row">
                      <div class="tags">
                        ${post.categories.map(category => `<div class="tag tagview click">${category.name}</div>`).join("")}
                      </div>
                      <div class="postedby">
                        <div class="tag">Posted by</div>
                        <div class="tag click" onclick="window.location.href = '/users/${post.nickname}';">${post.nickname}</div>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="addinfoframe-row">
                  <div class="numAnswers">${pagination.total_count} Answers</div>
                  <div class="addinfo">${renderTagsHtml}</div>
                </div>
                <div class="line"></div>
                ${await renderComments(postID)} <!-- Wait for comments to be rendered -->
                <div class="headline">Your Answer</div>
                <div class="editorForm">
                  <form id="editorForm" class="editorForm" method="POST" data-id="${post.id}">
                    <textarea id="editorTextarea" class="editor-textarea" placeholder="Type your text here..."></textarea>
                    <div class="rowBox">
                      <input type="submit" id="postButton" value="POST" class="click headerBnt">
                    </div>
                  </form>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div>${renderFooterHtml}</div>
      </div>
    `;

    app.innerHTML = postHtml;

    // Adding event listener after content is rendered
    const editorForm = document.getElementById('editorForm');
    if (editorForm) {
      editorForm.addEventListener('submit', postAnswer);
    } else {
      console.error("Form element 'editorForm' not found in the DOM.");
    }

  } catch (error) {
    console.error("Error fetching posts:", error);
    app.innerHTML += `<p>Error fetching posts: ${error.message}</p>`;
  }

}
