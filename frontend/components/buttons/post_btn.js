// postbtn.js

import { API_PREFIX } from "../../config.js";

export async function postAnswer(event) {
  event.preventDefault(); 

  const content = document.getElementById("editorTextarea").value;
  const form = document.getElementById('editorForm');
  const postID = form.getAttribute('data-id');

  console.log("Button clicked, content:", content);
  console.log("Post ID:", postID);

  if (!content) {
    alert("Please enter your answer before submitting.");
    return;
  }

  try {
    const response = await fetch(`/${API_PREFIX}/posts/${postID}/comments`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ content }),
    });

    console.log("Response status:", response.status);

    if (response.ok) {
      console.log('Comment submitted:', content);
      // Optionally, refresh comments or give user feedback
    } else {
      const errorData = await response.json();
      console.error(`Error creating comment: ${response.status} - ${errorData.message}`);
      alert(`Error creating comment: ${errorData.message}`);
    }

  } catch (error) {
    console.error("Error during posting answer:", error);
    alert(`Error during posting answer: ${error.message}`);
  }
}
