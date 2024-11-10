import { renderHeader } from '../base/header.js';
import { renderSidebar } from '../base/sidebar.js';
import { renderFooter } from '../base/footer.js';
import { authButton } from '../buttons/authBtn.js';

export async function renderAuth() {
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
                ${sidebarHtml} 
              </div>
              <div class="col-2"> 
                ${headerHtml}
                <div class="col-2 Content-flow"> 
                    <form id="registerForm" class="authform" method="POST" >
                        <div class="authinfobox">
                            <input type="text" id="nickname" name="username" class="authInput" placeholder="USERNAME"  required />
                        </div>
                        <div class="authinfobox">
                            <input type="email" id="email" name="email" class="authInput" placeholder="EMAIL"  required />
                        </div>
                        <div class="authinfobox">
                            <input type="text" id="first_name" name="first_name" class="authInput" placeholder="FIRST NAME"   required />
                        </div>
                        <div class="authinfobox">
                            <input type="text" id="last_name" name="last_name" class="authInput" placeholder="LAST NAME"   required />
                        </div>
                        <div class="authinfobox">
                            <input type="date" id="age" name="birthdate" class="authInput" placeholder="YOUR BIRTHDATE" required />
                        </div>
                        <div class="authinfobox">
                            <select id="gender" name="gender" class="authInput" required>
                                <option value="" disabled selected>YOUR GENDER</option>
                                <option value="male">Male</option>
                                <option value="female">Female</option>
                                <option value="other">Other</option>
                                <option value="prefer_not_to_say">Prefer not to say</option>
                            </select>
                        </div>
                        <div class="authinfobox">
                            <input type="password" id="password" name="password" class="authInput" placeholder="PASSWORD" minlength="8" maxlength="16" required />
                        </div>
                        
                        <div class="authbtnbox">
                            <input type="submit" class="authBnt" value="REGISTER" name="action">
                        </div>
                    </form>

                    <form class="authform" method="POST" action="">
                        <div class="authbtnbox2">
                            <div class="infostatic">Already have an account?</div>
                            <a href="/auth/login" style="color: var(--main-white) text-decoration: none; display: flex;">
                                <input type="submit" id="authButtonInput" class="inforender registerview" value="LOGIN" >
                            </a>
                        </div>
                    </form>
                </div>
              </div>
            </div>
            <div>${renderFooterHtml}</div>
          </div>
        `;

        // Add the event listener to the registration form
        const registerForm = document.getElementById('registerForm');
        if (registerForm) {
            registerForm.addEventListener('submit', authButton);
        }
    } catch (error) {
        console.error("Error rendering login page:", error);
        app.innerHTML = "<p>Error fetching login</p>";
    }

}
