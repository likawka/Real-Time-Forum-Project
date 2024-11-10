import { renderHome } from './components/pages/home.js';
import { renderPost } from './components/pages/post.js';
import { renderLogin } from './components/pages/login.js';
import { renderAuth } from './components/pages/auth.js';
import { createPostPage } from './components/pages/create_post.js';
import { renderProfile } from './components/pages/profile.js';
import { renderSearch } from './components/pages/renderSearch.js';
import { renderChatsPage } from './components/chats/chatPage.js';


const routes = {
  '/': renderHome,
  '/post/:id': renderPost,
  '/create': createPostPage,
  '/auth/login': renderLogin,
  '/auth/register': renderAuth,
  '/users/:userid': renderProfile,
  '/search/:query': renderSearch,
  '/search': renderSearch,
  '/chats/:chatid': renderChatsPage,
};

export function router() {
  const path = window.location.pathname;
  const query = new URLSearchParams(window.location.search);
  const searchQuery = query.get('search');

  if (searchQuery) {
    renderSearch(searchQuery);
    return;
  }

  const route = Object.keys(routes).find(route => {
    const regex = new RegExp(`^${route.replace(/:\w+/g, '\\w+')}$`);
    return regex.test(path);
  });

  if (route) {
    const regex = new RegExp(route.replace(/:\w+/g, '(\\w+)'));
    const values = path.match(regex).slice(1);
    routes[route](...values);
  } else {
    renderNotFound();
  }

  // Handle browser navigation
  window.addEventListener('popstate', router);

  // Handle link clicks
  document.body.addEventListener('click', (event) => {
    const target = event.target.closest('a');
    if (target) {
      event.preventDefault();
      history.pushState(null, null, target.href);
      router();
    }
  });
}

function renderNotFound() {
  const app = document.getElementById('app');
  app.innerHTML = '<h1>404 Not Found</h1>';
}
