import { renderHeader } from '../base/header.js';
import { renderSidebar } from '../base/sidebar.js';
import { renderFooter } from '../base/footer.js';
import { loginButton } from '../buttons/login_btn.js';

export async function renderLogin() {
    const app = document.getElementById("app");

    // Load CSS dynamically
    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.href = "/frontend/src/css/authourisation.css";
    document.head.appendChild(link);
    const link2 = document.createElement("link");
    link2.rel = "stylesheet";
    link2.href = "/frontend/src/css/main.css";
    document.head.appendChild(link2);

    let sidebarHtml = '';
    let headerHtml = '';
    let renderFooterHtml = '';

    try {
        await renderSidebar();
        sidebarHtml = app.innerHTML; 
        await renderHeader();
        headerHtml = app.innerHTML; 
        await renderFooter();
        renderFooterHtml = app.innerHTML; 

        // Set the app innerHTML after all rendering is complete
        app.innerHTML = `
          <div class="mainBody">
            <div class="mainPage">
              <div class="col-1">  
                ${sidebarHtml} </div>
              <div class="col-2"> 
                ${headerHtml}
                <div class="col-2 Content-flow"> 
                    <form id="loginForm" class="authform" method="POST" action="">
                        <div class="authinfobox">
                            <input type="email" name="email" class="authInput" placeholder="EMAIL" 
                                required />
                        </div>
                        <div class="authinfobox">
                            <input type="password" name="password" class="authInput" placeholder="PASSWORD" minlength="8" maxlength="16"
                                required />
                        </div>
                        <div class="authbtnbox">
                            <input type="submit" class="authBnt" value="LOGIN" name="action">
                        </div>
                    </form>
                    <form class="authform" method="POST" action="">
                <div class="authbtnbox2">
                    <div class="infostatic">Don’t have an account?</div>
                    <a href="/auth/register" style="color: var(--main-white) text-decoration: none; display: flex;"> <input type="submit" class="inforender registerview" value="REGISTER" name="authStatus"> </a> 
                </div>
                </div>
              </div>
            </div>
            <div>${renderFooterHtml}</div>
          </div>
        `;

        // Now that the HTML is set, add the event listener to the form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', loginButton);
        }
    } catch (error) {
        console.error("Error rendering login page:", error);
        app.innerHTML = "<p>Error fetching login</p>";
    }
}
