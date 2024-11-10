import { fetchAllUsers, fetchChats, fetchPosts } from "../../utils/api.js";
import { timeAgo } from "../additional/time_count.js";

export async function renderUsers() {
    const app = document.getElementById("app");

    try {
        const responseUser = await fetchAllUsers();
        const chats = await fetchChats();
        console.log(chats);
        const users = responseUser.payload.users;

        const response = await fetchPosts();
        const { authenticated } = response;

        // Initialize chatsHtML here
        let chatsHtML = '';

        if (authenticated) {
            // Define chatsHtML if authenticated
            chatsHtML = (user) => `
                <div class="user" data-user-id="${user.id}" onclick="window.location.href = '/chats/${user.id}';">
            `;
        } else {
            // Define chatsHtML if not authenticated
            chatsHtML = (user) => `
                <div class="user" data-user-id="${user.id}" onclick="window.location.href = '/';">
            `;
        }

        // Render users
        app.innerHTML = `
            <div class="sidebar">
                <div id="users-list">
                    ${users
                        .map(
                            (user) => `
                            <div class="sidebar-chat">
                                ${chatsHtML(user)}
                                <div class="click">
                                    ${user.nickname}
                                    <div class="addinfoframe-row">
                                        <div class="online">Online</div>
                                    </div>
                                </div>
                            </div>
                            </div>
                        `
                        )
                        .join("")}
                </div>
            </div>
        `;

    } catch (error) {
        console.error("Error fetching Users:", error);
        app.innerHTML = "<p>Error fetching Users</p>";
    }
}
