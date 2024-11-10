import { fetchPost } from "../../utils/api.js";
import { likeComment, dislikeComment } from "../buttons/like_dislike.js";

export async function renderComments(postID) {
    const app = document.getElementById("app");

    try {
        const payload = await fetchPost(postID);
        const comments = payload.payload.comments;

        if (comments === null) {
            return '<div class="numAnswers">Be the legend who drops the first comment on this post!</div> <div class="line"></div>';
        } else {
            return comments.map(comment => `
                <div class="questionbox" id="comment-${comment.id}">
                    <div class="ratebox">
                        <button onclick="likeComment('${comment.id}', '${postID}')" style="background:none; border:none;">
                            <img src="/frontend/src/images/like.svg" class="rateicon click" alt="Like">
                        </button>
                        <div class="rateNum" id="rateNumCom-${comment.id}">${comment.rate.rate}</div>
                        <button onclick="dislikeComment('${comment.id}', '${postID}')" style="background:none; border:none; transform: rotate(180deg)">
                            <img src="/frontend/src/images/like.svg" class="rateicon click" alt="Dislike">
                        </button>
                    </div>
                    <div class="questContainer">
                        <div class="questtextbox" id="contentComm-${comment.id}" >${comment.content}</div>
                        <div class="addinfoframe-row">
                            <div class="postedby">
                                <div class="tag">Posted by</div>
                                <div class="tag click" onclick="window.location.href = '/users/${comment.nickname}';" >${comment.nickname}</div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="line"></div>
            `).join("");
        }
    } catch (error) {
        console.error("Error fetching comments:", error);
        return "<p>Error fetching comments</p>";
    }
}
