import { API_PREFIX } from "../../config.js";

export async function createPostButton(event) {
  event.preventDefault(); // Prevent default form s               ubmission

  const form = event.target;
  const title = form.title.value;
  const content = form.content.value;
  const categories = form.categories.value;
  


  try {
    const response = await fetch(`/${API_PREFIX}/posts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title, content, categories }),
    });

    console.log("Response status:", response.status);

    if (response.ok) {
      console.log('Post submitted:', content);
      window.location.href = '/';
    } else {
      const errorData = await response.json();
      console.error(`Error creating comment: ${response.status} - ${errorData.message}`);
      // alert(`Error creating comment: ${errorData.message}`);
    }
    alert("Please enter your Post before submitting.");
    
  } catch (error) {
    console.error("Error during posting Post:", error);
    // alert(`Error during posting Post: ${error.message}`);
  }
}
