import { API_PREFIX } from '../../config.js';

export async function authButton(event) {
    event.preventDefault(); // Prevent form submission

    const nickname = event.target.nickname.value;
    const email = event.target.email.value;
    const first_name = event.target.first_name.value;
    const last_name = event.target.last_name.value;
    const birthdate = event.target.birthdate.value; // Спочатку отримуємо дату в форматі YYYY-MM-DD
    const gender = event.target.gender.value;
    const password = event.target.password.value;

    // Перетворення дати з формату YYYY-MM-DD у формат DD.MM.YYYY
    let age;
    if (birthdate) {
        const dateParts = birthdate.split('-');
        age = `${dateParts[2]}.${dateParts[1]}.${dateParts[0]}`;
    } else {
        age = ''; // Обробка випадку, коли дата не вказана
    }

    console.log('Registering user:', JSON.stringify({ nickname, email, first_name, last_name, age, gender, password }));
    try {
        const response = await fetch(`/${API_PREFIX}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ nickname, email, first_name, last_name, age, gender, password })
        });

        if (response.ok) {
            alert('Registration successful!');
            window.location.reload(); // Refresh the page after registration
            window.location.href = '/';
        } else {
            const errorData = await response.json();
            console.error('Registration failed:', errorData);
            alert('Registration failed. Please check your credentials.');
        }
    } catch (error) {
        console.error('Error during registration:', error);
        alert('Error during registration. Please try again later.');
    }
}
