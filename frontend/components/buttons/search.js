// File: search.js

import { fetchPostsSearch } from "../../utils/api.js"; // Import fetchPostsSearch function

export async function searchPosts(query) {
  try {
    // Fetch all posts using fetchPostsSearch
    const response = await fetchPostsSearch();
    const posts = response.payload.posts;

    // Filter posts based on the search query
    const filteredPosts = posts.filter(post => {
      const searchInContent = post.content.toLowerCase().includes(query.toLowerCase());
      const searchInTitle = post.title.toLowerCase().includes(query.toLowerCase());
      const searchInCategories = post.categories.some(category =>
        category.name.toLowerCase().includes(query.toLowerCase())
      );

      return searchInContent || searchInTitle || searchInCategories;
    });

    return filteredPosts;
  } catch (error) {
    console.error("Error searching posts:", error);
    return [];
  }
}
