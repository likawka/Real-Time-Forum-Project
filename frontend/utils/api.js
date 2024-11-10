import { API_PREFIX } from '../config.js';

export async function fetchPosts() {
    const response = await fetch(`/${API_PREFIX}/posts`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
}

// utils/api.js

export async function fetchPostsSearch() {
    try {
        const response = await fetch(`/${API_PREFIX}/posts`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        }); 
      const data = await response.json();
      return data; 
      
    } catch (error) {
      console.error("Error fetching posts:", error);
      throw error;
    }
  }
  

export async function fetchPost(postID) {
    const response = await fetch(`/${API_PREFIX}/posts/${postID}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
}

export async function fetchUserData(nickname, type) {
    const response = await fetch(`/${API_PREFIX}/users/${nickname}/${type}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
  }
  
  export async function fetchAllUsers() {
    const response = await fetch(`/${API_PREFIX}/users`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
  }

  export async function fetchChats() {
    const response = await fetch(`/${API_PREFIX}/chats`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
  }
  
  export async function fetchChatsHash(chat_hash) {
    const response = await fetch(`/${API_PREFIX}/chats/${chat_hash}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });
    return response.json();
  }
  
  export async function fetchChatsArr() {
    try {
      const response = await fetch(`/${API_PREFIX}/chats`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      const data = await response.json();
      if (data.status === 'success' && data.payload && Array.isArray(data.payload.chats)) {
        return data.payload.chats;
      } else {
        throw new Error('Invalid chats data structure');
      }
    } catch (error) {
      console.error('Error fetching chats:', error);
      return [];
    }
  }