import { API_PREFIX } from '../../config.js';

let clickTimeouts = {}; // Store timeouts for each comment to prevent rapid clicking

export async function ratePost(postID, status) {
    try {
        const response = await fetch(`/${API_PREFIX}/rate`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ post_id: Number(postID), status })
        });

        if (response.ok) {
            const data = await response.json();
            console.log('Post rating updated:', data);

            const rateElement = document.getElementById(`rateNum-${postID}`);
            if (rateElement) {
                rateElement.textContent = data.payload.rate.rate;
            } else {
                console.error(`Post rating element with ID rateNum-${postID} not found.`);
            }
        } else {
            const errorData = await response.json();
            console.error(`Error updating post rating: ${response.status} - ${errorData.message}`);
        }
    } catch (error) {
        console.error('Network or other error:', error);
    }
}

export async function likePost(postID) {
    await ratePost(postID, 'up');
}

export async function dislikePost(postID) {
    await ratePost(postID, 'down');
}

export async function rateComment(commentID, postID, status) {
    try {
        const response = await fetch(`/${API_PREFIX}/rate`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ comment_id: Number(commentID), post_id: Number(postID), status })
        });
        
        console.log(`Sending status: ${status}`);

        if (response.ok) {
            const data = await response.json();
            console.log('Comment rating updated:', data);

            const rateElement = document.getElementById(`rateNumCom-${commentID}`);
            if (rateElement) {
                rateElement.textContent = data.payload.rate.rate; // Update the rating
            } else {
                console.error(`Comment rating element with ID rateNumCom-${commentID} not found.`);
            }
        } else {
            const errorData = await response.json();
            console.error(`Error updating comment rating: ${response.status} - ${errorData.message}`);
        }
    } catch (error) {
        console.error('Network or other error:', error);
    }
}

function throttleClick(commentID, callback) {
    if (!clickTimeouts[commentID]) {
        callback();
        clickTimeouts[commentID] = setTimeout(() => {
            clearTimeout(clickTimeouts[commentID]);
            delete clickTimeouts[commentID];
        }, 2000); // 2 seconds cooldown
    }
}

export function likeComment(commentID, postID) {
    throttleClick(commentID, () => rateComment(commentID, postID, 'up'));
}

export function dislikeComment(commentID, postID) {
    throttleClick(commentID, () => rateComment(commentID, postID, 'down'));
}

window.likePost = likePost;
window.dislikePost = dislikePost;
window.likeComment = likeComment;
window.dislikeComment = dislikeComment;
