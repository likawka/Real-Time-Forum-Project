import { API_PREFIX } from '../../config.js';

export default async function handleLogout(event) {
    console.log("Logout button clicked"); // Log to console to verify function is called

    event.preventDefault(); // Prevent default behavior

    try {
        const response = await fetch(`/${API_PREFIX}/auth/logout`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        if (response.ok) {
            const responseData = await response.json();
            localStorage.removeItem('userId');
            localStorage.removeItem('userNickname');
            alert('Logout successful!');
            window.location.href = '/'; // Redirect to home after logout
        } else {
            const errorData = await response.json();
            console.error('Logout failed:', errorData);
        }
    } catch (error) {
        console.error('Error during logout:', error);
        alert('Error during logout. Please try again later.');
    }
}

// Make the handleLogout function globally accessible
window.handleLogout = handleLogout;