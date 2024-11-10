import { fetchPosts, fetchUserData } from "../../utils/api.js";
import {renderUsers} from '../additional/users.js';
import handleLogout from '../buttons/logout.js';

export async function renderSidebar() {
    const app = document.getElementById("app");

    // Load CSS dynamically

    const link2 = document.createElement("link");
    link2.rel = "stylesheet";
    link2.href = "/frontend/src/css/sidebar.css";
    document.head.appendChild(link2);

    let usersHtml = '';

    try {
        await renderUsers();
        usersHtml = app.innerHTML;

        const response = await fetchPosts();
        const { authenticated } = response;
        const userNickname = localStorage.getItem('userNickname');



        app.innerHTML = `
            <div class="sideBar">
                <div class="logoSideBarFrame click" onclick="window.location.href = '/';">UNDER pressure</div>
                <div class="sideBarMenuFrame">
                    <div class="sideBarMenuBtn click" onclick="window.location.href = '/';">
                        <img class="mainIcon" src="/frontend/src/images/home.svg" />
                        Home
                    </div>
                    <div  class="hiddenElement" id="createPost" onclick="window.location.href = '/create';"> 
                            <div class="sideBarMenuBtn click">
                                <img class="mainIcon click" src="/frontend/src/images/createPost.svg">
                                <div class="textview click">Create Post</div>
                            </div>
                    </div>
                    <div class="hiddenElement" id="logoutBtn" onclick="handleLogout(event)">
                        <div style="width: 100%;" id="logoutForm">
                            <div class="sideBarMenuBtn click">
                                <img class="mainIcon click" src="/frontend/src/images/logOut.svg">
                                <div class="textview click">Log Out</div>
                            </div>
                        </div>
                    </div>

                    <br/><br/>

                    ${usersHtml}
                </div>
                </div>
           
        `;

        const logoutBtn = document.getElementById("logoutBtn");
        if (authenticated) {
            logoutBtn.classList.remove("hiddenElement");
        } else {
            logoutBtn.classList.add("hiddenElement");
        }
        
        const createPost = document.getElementById("createPost");
        if (authenticated) {
            createPost.classList.remove("hiddenElement");
        } else {
            createPost.classList.add("hiddenElement");
        }
    } catch (error) {
        console.error("Error fetching sidebar:", error);
        app.innerHTML = `<p>Error fetching sidebar: ${error.message}</p>`;
    }
}

