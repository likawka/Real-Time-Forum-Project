import { fetchAllUsers, fetchChats, fetchChatsHash } from "../../utils/api.js";
import { timeAgo } from "../additional/time_count.js";
import { renderHeader } from "../base/header.js";
import { renderSidebar } from "../base/sidebar.js";
import { renderFooter } from "../base/footer.js";
import { API_PREFIX } from '../../config.js';

let conn; // WebSocket connection
let roomHash; // Define roomHash

export async function renderChatsPage() {
  const app = document.getElementById("app");

  // Add stylesheets
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = "/frontend/src/css/chats.css";
  document.head.appendChild(link);

  const link2 = document.createElement("link");
  link2.rel = "stylesheet";
  link2.href = "/frontend/src/css/main.css";
  document.head.appendChild(link2);

  let sidebarHtml = "";
  let headerHtml = "";
  let renderFooterHtml = "";

  // Render sidebar, header, and footer
  await renderSidebar();
  sidebarHtml = app.innerHTML;
  await renderHeader();
  headerHtml = app.innerHTML;
  await renderFooter();
  renderFooterHtml = app.innerHTML;

  const { chatPartnerId } = getCurrentUserId();

  // Fetch all users
  const usersResponse = await fetchAllUsers();
  if (!usersResponse || !usersResponse.payload || !usersResponse.payload.users) {
    console.error("Failed to fetch users");
    return;
  }

  const users = usersResponse.payload.users;
  const chatPartner = users.find(user => user.id === parseInt(chatPartnerId, 10));

  // Get or create the chat room and get the roomHash
  const { currentUserId } = getCurrentUserId();
  const chatRoomResult = await getOrCreateChatRoom(chatPartnerId, currentUserId);

  if (!chatRoomResult.success) {
    console.error("Failed to get or create chat room:", chatRoomResult.message);
    return;
  }

  roomHash = chatRoomResult.roomHash;

  // Fetch chats using roomHash
  const chatsMessages = await fetchChatsHash(roomHash);
  const messagesDB = chatsMessages.payload.messages || []; // Ensure messagesDB is always an array

  // Generate messages HTML or a placeholder if there are no messages
  const messagesHtml = messagesDB.length > 0 ? messagesDB.map(
    (messagesData) => `
    <div class="message">
      <span>${messagesData.sender.nickname}: ${messagesData.message}</span>
      <span class="timestamp">${formatDateTime(messagesData.created_at)}</span>
    </div>`
  ).join("") : '<div class="message-placeholder"> </div>';

  app.innerHTML = `
    <div class="mainBody">
      <div class="mainPage">
        <div class="col-1">${sidebarHtml}</div>
        <div class="col-2"> 
          ${headerHtml}
          <div class="col-2 Content-flow"> 
            <div class="chats-box">
              <div>
                <div class="addinfo">
                  <div class="headline">To:</div>
                  <div class="headline" style="color: var(--light-red);">${chatPartner.nickname}</div>
                </div>
                <div class="line"></div>
              </div>
              <div class="chat-container">
                <div class="date-separator hidden">Today</div>
                <div id="msgcontainer" class="message-row">
                  ${messagesHtml}
                  <div class="message hidden"></div>
                </div>
                <form id="chatroom-message">
                  <div class="chat-container-row">
                    <input id="message" class="chat-input" type="text" placeholder="Type your message...">
                    <button id="send-button" class="send-button">Send</button>
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

  const ChatUserId = getCurrentUserId();
  document.getElementById("chatroom-message").onsubmit = sendMessage;

  if (window["WebSocket"]) {
    console.log("Browser supports WebSocket.");
    setupWebSocket(ChatUserId);
  } else {
    alert("Your browser does not support WebSockets.");
  }

  const msgContainer = document.getElementById("msgcontainer");
  const messageDiv = msgContainer.querySelector('.message');

  if (messageDiv && !messageDiv.innerHTML.trim()) {
    messageDiv.classList.add('hidden');
  }
}

function formatDateTime(isoDate) {
  const dateObj = new Date(isoDate);
  const hours = String(dateObj.getHours()).padStart(2, '0');
  const minutes = String(dateObj.getMinutes()).padStart(2, '0');
  const day = String(dateObj.getDate()).padStart(2, '0');
  const month = String(dateObj.getMonth() + 1).padStart(2, '0');
  const year = dateObj.getFullYear();
  return `${hours}:${minutes} ${day}.${month}.${year}`;
}


function getCurrentUserId() {
  const url_string = window.location.href;  
  const url = new URL(url_string);        
  const pathSegments = url.pathname.split('/');
  const chatPartnerId = pathSegments[pathSegments.length - 1];
  
  // Assuming you store the current user's ID in localStorage
  const currentUserId = localStorage.getItem('userId');
  
  return { currentUserId, chatPartnerId };
}

function setupWebSocket(ChatUserId) {
  conn = new WebSocket(`ws://localhost:8080/api/ws`);

  conn.onopen = function () {
    console.log("WebSocket connection established.");
    const { currentUserId, chatPartnerId } = getCurrentUserId();
    getOrCreateChatRoom(chatPartnerId, currentUserId);
  };

  conn.onmessage = function (event) {
    console.log("WebSocket message received:", event.data);
    const message = JSON.parse(event.data);
    handleMessage(message);
  };

  conn.onerror = function (error) {
    console.log("WebSocket error: " + error.message);
  };

  conn.onclose = function (event) {
    console.log("WebSocket connection closed: ", event);
  };
}


async function createChatRoom(user1_id, user2_id) {
  try {
    // Ensure that user1_id and user2_id are integers
    user1_id = parseInt(user1_id, 10);
    user2_id = parseInt(user2_id, 10);

    console.log(`Attempting to create chat room for users: ${user1_id} and ${user2_id}`);
    
    const response = await fetch(`/api/chats`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ user1_id, user2_id }),  // These will now be integers
    });

    // The next line contains a typo, correcting it to 'user1_id'
    console.log(JSON.stringify({ user1_id, user2_id }));

    console.log("Response status:", response.status);

    if (response.ok) {
      const data = await response.json();
      console.log('Chat created:', data);
      return { success: true, data };
    } else {
      const errorData = await response.json();
      console.error(`Error creating chat: ${response.status} - ${errorData.message}`);
      console.error('Full error response:', errorData);
      return { success: false, message: errorData.message };
    }
  } catch (error) {
    console.error("Error during chat creation:", error);
    return { success: false, message: error.message };
  }
}


function joinRoom(roomHash) {  
  if (conn && conn.readyState === WebSocket.OPEN) {
    const msg = {
      type: "join_room",
      payload: {
        roomHash: roomHash,
      },
    };
    conn.send(JSON.stringify(msg));
  }
}

async function getOrCreateChatRoom(ChatUserId, userId) {
  try {
    console.log(`Getting or creating chat room for users: ${ChatUserId} and ${userId}`);
    
    const response = await fetchChats();
    const chatsBL = response.payload.chats;

    console.log('All chats:', response);

    const user1IdStr = String(ChatUserId);
    const user2IdStr = String(userId);

    let chatRoom = null;

    for (let i = 0; i < chatsBL.length; i++) {
      const chat = chatsBL[i];

      if (
        (String(chat.user1_id) === user1IdStr && String(chat.user2_id) === user2IdStr) ||
        (String(chat.user1_id) === user2IdStr && String(chat.user2_id) === user1IdStr)
      ) {
        console.log('Matching chat found:', chat);
        chatRoom = chat;
        break;
      }
    }
    
    if (chatRoom && chatRoom.chatHash) {
      roomHash = chatRoom.chatHash;  // Save roomHash globally
      console.log('Existing chat room found:', chatRoom.chatHash);
      
      joinRoom(roomHash);
      return { success: true, message: "Joined existing chat room", roomHash: chatRoom.chatHash };
    } else {
      console.log('No existing chat room found, creating a new one');
      const result = await createChatRoom(userId, ChatUserId);

      if (result.success) {
        roomHash = result.data.chatHash;  // Save roomHash globally
        joinRoom(roomHash);
        return { success: true, message: "Created and joined new chat room", roomHash: result.data.chatHash };
      } else {
        return { success: false, message: result.message };
      }
    }
  } catch (error) {
    console.error("Error during chat room check or creation:", error);
    return { success: false, message: error.message };
  }
}






export function sendMessage(event) {
  event.preventDefault();
  const newmessage = document.getElementById("message").value;

  if (newmessage && conn && conn.readyState === WebSocket.OPEN) {
    const msg = {
      type: "message",
      payload: {
        message: newmessage,
        roomHash,  // Use the global roomHash
      },
    };

    console.log("Sending message:", msg , "roomHash:", roomHash);
    conn.send(JSON.stringify(msg));
    document.getElementById("message").value = "";
  }
}


function handleMessage(message) {
  const msgContainer = document.getElementById("msgcontainer");
  const emptyMessageDiv = msgContainer.querySelector('.message');

  switch (message.type) {
    case "message":
      if (emptyMessageDiv && !emptyMessageDiv.innerHTML.trim()) {
        emptyMessageDiv.classList.add("hidden");
      }

      const msg = message.payload;
      const msgDiv = document.createElement("div");
      msgDiv.className = "message";
      msgDiv.innerHTML = `
        <span>${msg.sender.nickname}: ${msg.message}</span>
        <span class="timestamp">${new Date(msg.created_at).toLocaleTimeString()}</span>
      `;
      msgContainer.appendChild(msgDiv);
      break;

    case "active_users":
      console.log("Active users: ", message.payload.users);
      break;

    case "typing":
      console.log(`${message.payload.sender.username} is typing...`);
      break;

    case "error":
      console.error("Server error:", message);
      break;

    default:
      console.log("Unknown message type:", message.type);
  }
}
