# Forum Project README

Welcome to the Forum Project! This README provides an overview of the project's structure and functionalities. The goal is to build a feature-rich forum using a combination of SQLite, Golang, JavaScript, HTML, and CSS. This project involves implementing registration and login systems, post and comment management, and private messaging. Below are detailed instructions on how to set up and contribute to the project.

## Project Structure

1. **SQLite** : Database management
2. **Golang** : Backend services and WebSocket handling
3. **JavaScript** : Frontend event handling and WebSocket communication
4. **HTML** : Single-page layout
5. **CSS** : Styling and layout customization

## Features

### 1. Registration and Login

* **Registration** :
* Users must provide the following information:
  * Nickname
  * Age
  * Gender
  * First Name
  * Last Name
  * E-mail
  * Password
* **Login** :
* Users can log in using either their nickname or e-mail combined with the password.
* **Logout** :
* Users can log out from any page within the forum.

### 2. Posts and Comments

* **Creating Posts** :
* Users can create new posts with categories.
* **Commenting on Posts** :
* Users can comment on posts. Comments are visible only when a user clicks on a post.
* **Feed Display** :
* Posts are displayed in a feed format.

### 3. Private Messages

* **Chat Section** :
* Displays online/offline users, organized by the last message sent.
* For new users, organize alphabetically.
* **Sending Messages** :
* Users can send messages to those who are online.
* A chat section should always be visible.
* **Viewing and Loading Messages** :
* When clicking on a user, past messages are loaded.
* Implement throttling/debouncing to manage scrolling.
* **Message Format** :
* Each message displays:
  * Date sent
  * Sender's username
* Messages should be updated in real-time using WebSockets.

## Setup Instructions

### Prerequisites

* Go (Golang) installed
* SQLite database setup
* A modern web browser

### Installation

1. **Clone the Repository** :
2. **Run the Application** :

* Start the Go server:
  ```
  bash scripts/start.sh
  ```
* Open your web browser and navigate to `http://localhost:8080`

## Users

Login: TestUser1@gmail.com
Password: !QAZ@WSX3edc

Login: TestUser2@gmail.com
Password: !QAZ@WSX3edc

## Development

* **HTML** : Use a single HTML file to manage the layout.
* **CSS** : Customize styles in the CSS file to improve UI/UX.
* **JavaScript** : Handle all frontend events and WebSocket interactions in the JavaScript file.
* **Golang** : Implement backend logic and WebSocket handling in the Go files.

## Created by

izinko, lya
