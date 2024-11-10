import { API_PREFIX } from '../../config.js';

export async function loginButton(event) {
    event.preventDefault(); // Prevent form submission

    const email = event.target.email.value;
    const password = event.target.password.value;

    if (!email || !password) {
        alert('Please enter both email and password.');
        return;
    }

    try {
        const response = await fetch(`/${API_PREFIX}/auth/login`, {  // Corrected the URL
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })  // Keep credentials in the body
        });
        
        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('userId', data.user.id);
            localStorage.setItem('userNickname', data.user.nickname);
            alert('Login successful!');
            window.location.href = '/'; // Redirect to home page after login
        } else {
            const errorData = await response.json();
            console.error('Login failed:', errorData);
            alert('Login failed. Please check your credentials.');
        }
    } catch (error) {
        console.error('Error during login:', error);
        alert('Error during login. Please try again later.');
    }
}
