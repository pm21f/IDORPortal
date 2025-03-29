package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// User represents a simple user model
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email,omitempty"`
	Created  time.Time `json:"created"`
}

// Post represents a simple post model
type Post struct {
	ID      int       `json:"id"`
	UserID  int       `json:"user_id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

// Response represents a JSON response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// In-memory data store for demonstration
var (
	users = []User{
		{ID: 1, Username: "admin", Email: "admin@example.com", Created: time.Now().Add(-24 * time.Hour)},
		{ID: 2, Username: "user1", Email: "user1@example.com", Created: time.Now().Add(-12 * time.Hour)},
	}
	posts = []Post{
		{ID: 1, UserID: 1, Title: "Admin Post", Content: "This is a post by admin", Created: time.Now().Add(-12 * time.Hour)},
		{ID: 2, UserID: 2, Title: "User Post", Content: "This is a post by user1", Created: time.Now().Add(-6 * time.Hour)},
	}
)

// sendJSONResponse sends a JSON response
func sendJSONResponse(w http.ResponseWriter, success bool, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: success,
		Message: message,
		Data:    data,
	})
}

// IndexHandler handles the root path
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CycleSync - Home</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
				}
				.container {
					width: 80%;
					margin: auto;
					overflow: hidden;
				}
				header {
					background: #50b3a2;
					color: white;
					padding-top: 30px;
					min-height: 70px;
					border-bottom: #e8491d 3px solid;
				}
				header a {
					color: #ffffff;
					text-decoration: none;
					text-transform: uppercase;
					font-size: 16px;
				}
				header h1 {
					margin-bottom: 10px;
				}
				header .nav {
					margin-top: 10px;
					display: flex;
					justify-content: center;
				}
				header .nav a {
					margin: 0 15px;
					padding: 5px;
				}
				.highlight, .current a {
					font-weight: bold;
					background-color: #e8491d;
					border-radius: 3px;
					padding: 3px 7px;
				}
				.showcase {
					min-height: 400px;
					text-align: center;
					color: #333333;
					margin-top: 50px;
				}
				.showcase h1 {
					margin-top: 100px;
					font-size: 55px;
					margin-bottom: 10px;
				}
				.showcase p {
					font-size: 20px;
					margin-bottom: 20px;
				}
				.button {
					display: inline-block;
					font-size: 18px;
					text-decoration: none;
					color: #ffffff;
					background-color: #e8491d;
					padding: 10px 20px;
					border: none;
					border-radius: 3px;
					margin-top: 20px;
				}
				footer {
					padding: 20px;
					margin-top: 20px;
					color: #ffffff;
					background-color: #50b3a2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<header>
				<div class="container">
					<h1>CycleSync</h1>
					<div class="nav">
						<a href="/" class="highlight">Home</a>
						<a href="/login">Login</a>
						<a href="/signup">Signup</a>
					</div>
				</div>
			</header>

			<section class="showcase">
				<div class="container">
					<h1>Welcome to CycleSync</h1>
					<p>A platform for secure information sharing.</p>
					<p>Register an account to start posting!</p>
					<a href="/signup" class="button">Get Started</a>
				</div>
			</section>

			<footer>
				<p>CycleSync &copy; 2025</p>
			</footer>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginPageHandler renders the login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("login").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CycleSync - Login</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
				}
				.container {
					width: 80%;
					margin: auto;
					overflow: hidden;
				}
				header {
					background: #50b3a2;
					color: white;
					padding-top: 30px;
					min-height: 70px;
					border-bottom: #e8491d 3px solid;
				}
				header a {
					color: #ffffff;
					text-decoration: none;
					text-transform: uppercase;
					font-size: 16px;
				}
				header h1 {
					margin-bottom: 10px;
				}
				header .nav {
					margin-top: 10px;
					display: flex;
					justify-content: center;
				}
				header .nav a {
					margin: 0 15px;
					padding: 5px;
				}
				.highlight, .current a {
					font-weight: bold;
					background-color: #e8491d;
					border-radius: 3px;
					padding: 3px 7px;
				}
				.login-form {
					width: 400px;
					margin: 100px auto;
					padding: 20px;
					background: #ffffff;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.login-form h2 {
					text-align: center;
					margin-bottom: 20px;
				}
				.form-group {
					margin-bottom: 15px;
				}
				.form-group label {
					display: block;
					margin-bottom: 5px;
				}
				.form-group input {
					width: 100%;
					padding: 10px;
					border: 1px solid #ddd;
					border-radius: 3px;
				}
				.form-group input[type="submit"] {
					background-color: #50b3a2;
					color: white;
					cursor: pointer;
					font-size: 16px;
				}
				.form-group input[type="submit"]:hover {
					background-color: #429e8d;
				}
				.form-footer {
					text-align: center;
					margin-top: 20px;
				}
				.error-message {
					color: red;
					margin-bottom: 10px;
					text-align: center;
					display: none;
				}
				footer {
					padding: 20px;
					margin-top: 20px;
					color: #ffffff;
					background-color: #50b3a2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<header>
				<div class="container">
					<h1>CycleSync</h1>
					<div class="nav">
						<a href="/">Home</a>
						<a href="/login" class="highlight">Login</a>
						<a href="/signup">Signup</a>
					</div>
				</div>
			</header>

			<div class="container">
				<div class="login-form">
					<h2>Login</h2>
					<div id="error-message" class="error-message"></div>
					<form id="login-form">
						<div class="form-group">
							<label for="username">Username:</label>
							<input type="text" id="username" name="username" required>
						</div>
						<div class="form-group">
							<label for="password">Password:</label>
							<input type="password" id="password" name="password" required>
						</div>
						<div class="form-group">
							<input type="submit" value="Login">
						</div>
					</form>
					<div class="form-footer">
						<p>Don't have an account? <a href="/signup">Sign up</a></p>
					</div>
				</div>
			</div>

			<footer>
				<p>CycleSync &copy; 2025</p>
			</footer>

			<script>
				document.getElementById('login-form').addEventListener('submit', function(e) {
					e.preventDefault();
					
					const username = document.getElementById('username').value;
					const password = document.getElementById('password').value;
					
					// If you're the admin or user1 from our hardcoded data, allow login
					if ((username === 'admin' && password === 'password') || 
						(username === 'user1' && password === 'password')) {
						window.location.href = '/dashboard';
					} else {
						const errorDiv = document.getElementById('error-message');
						errorDiv.textContent = 'Invalid username or password';
						errorDiv.style.display = 'block';
					}
				});
			</script>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// SignupPageHandler renders the signup page
func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("signup").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CycleSync - Signup</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
				}
				.container {
					width: 80%;
					margin: auto;
					overflow: hidden;
				}
				header {
					background: #50b3a2;
					color: white;
					padding-top: 30px;
					min-height: 70px;
					border-bottom: #e8491d 3px solid;
				}
				header a {
					color: #ffffff;
					text-decoration: none;
					text-transform: uppercase;
					font-size: 16px;
				}
				header h1 {
					margin-bottom: 10px;
				}
				header .nav {
					margin-top: 10px;
					display: flex;
					justify-content: center;
				}
				header .nav a {
					margin: 0 15px;
					padding: 5px;
				}
				.highlight, .current a {
					font-weight: bold;
					background-color: #e8491d;
					border-radius: 3px;
					padding: 3px 7px;
				}
				.signup-form {
					width: 400px;
					margin: 100px auto;
					padding: 20px;
					background: #ffffff;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.signup-form h2 {
					text-align: center;
					margin-bottom: 20px;
				}
				.form-group {
					margin-bottom: 15px;
				}
				.form-group label {
					display: block;
					margin-bottom: 5px;
				}
				.form-group input {
					width: 100%;
					padding: 10px;
					border: 1px solid #ddd;
					border-radius: 3px;
				}
				.form-group input[type="submit"] {
					background-color: #50b3a2;
					color: white;
					cursor: pointer;
					font-size: 16px;
				}
				.form-group input[type="submit"]:hover {
					background-color: #429e8d;
				}
				.form-footer {
					text-align: center;
					margin-top: 20px;
				}
				.error-message {
					color: red;
					margin-bottom: 10px;
					text-align: center;
					display: none;
				}
				footer {
					padding: 20px;
					margin-top: 20px;
					color: #ffffff;
					background-color: #50b3a2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<header>
				<div class="container">
					<h1>CycleSync</h1>
					<div class="nav">
						<a href="/">Home</a>
						<a href="/login">Login</a>
						<a href="/signup" class="highlight">Signup</a>
					</div>
				</div>
			</header>

			<div class="container">
				<div class="signup-form">
					<h2>Create an Account</h2>
					<div id="error-message" class="error-message"></div>
					<form id="signup-form">
						<div class="form-group">
							<label for="username">Username:</label>
							<input type="text" id="username" name="username" required>
						</div>
						<div class="form-group">
							<label for="email">Email:</label>
							<input type="email" id="email" name="email" required>
						</div>
						<div class="form-group">
							<label for="password">Password:</label>
							<input type="password" id="password" name="password" required>
						</div>
						<div class="form-group">
							<label for="confirm-password">Confirm Password:</label>
							<input type="password" id="confirm-password" name="confirm-password" required>
						</div>
						<div class="form-group">
							<input type="submit" value="Sign Up">
						</div>
					</form>
					<div class="form-footer">
						<p>Already have an account? <a href="/login">Login</a></p>
					</div>
				</div>
			</div>

			<footer>
				<p>CycleSync &copy; 2025</p>
			</footer>

			<script>
				document.getElementById('signup-form').addEventListener('submit', function(e) {
					e.preventDefault();
					
					const username = document.getElementById('username').value;
					const email = document.getElementById('email').value;
					const password = document.getElementById('password').value;
					const confirmPassword = document.getElementById('confirm-password').value;
					
					const errorDiv = document.getElementById('error-message');
					
					// Simple validation
					if (password !== confirmPassword) {
						errorDiv.textContent = 'Passwords do not match';
						errorDiv.style.display = 'block';
						return;
					}
					
					// For demo, just simulate successful registration
					window.location.href = '/dashboard';
				});
			</script>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// DashboardHandler renders the dashboard page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("dashboard").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CycleSync - Dashboard</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
				}
				.container {
					width: 80%;
					margin: auto;
					overflow: hidden;
				}
				header {
					background: #50b3a2;
					color: white;
					padding-top: 30px;
					min-height: 70px;
					border-bottom: #e8491d 3px solid;
				}
				header a {
					color: #ffffff;
					text-decoration: none;
					text-transform: uppercase;
					font-size: 16px;
				}
				header h1 {
					margin-bottom: 10px;
				}
				header .nav {
					margin-top: 10px;
					display: flex;
					justify-content: center;
				}
				header .nav a {
					margin: 0 15px;
					padding: 5px;
				}
				.highlight, .current a {
					font-weight: bold;
					background-color: #e8491d;
					border-radius: 3px;
					padding: 3px 7px;
				}
				.dashboard {
					margin: 50px 0;
				}
				.dashboard h2 {
					color: #333;
				}
				.posts {
					margin-top: 30px;
				}
				.post {
					background: #fff;
					padding: 15px;
					margin-bottom: 15px;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.post h3 {
					margin-top: 0;
					color: #50b3a2;
				}
				.post .meta {
					color: #888;
					font-size: 0.8em;
					margin-bottom: 10px;
				}
				.post .content {
					line-height: 1.6;
				}
				.post .actions {
					margin-top: 10px;
				}
				.post .actions a {
					color: #50b3a2;
					margin-right: 10px;
					text-decoration: none;
				}
				.post .actions a:hover {
					text-decoration: underline;
				}
				.create-post {
					background: #fff;
					padding: 15px;
					margin-bottom: 30px;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.create-post h3 {
					margin-top: 0;
					color: #50b3a2;
				}
				.form-group {
					margin-bottom: 15px;
				}
				.form-group label {
					display: block;
					margin-bottom: 5px;
				}
				.form-group input, .form-group textarea {
					width: 100%;
					padding: 10px;
					border: 1px solid #ddd;
					border-radius: 3px;
				}
				.form-group textarea {
					height: 100px;
				}
				.form-group button {
					background-color: #50b3a2;
					color: white;
					border: none;
					padding: 10px 15px;
					border-radius: 3px;
					cursor: pointer;
					font-size: 16px;
				}
				.form-group button:hover {
					background-color: #429e8d;
				}
				.users {
					background: #fff;
					padding: 15px;
					margin-top: 30px;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.users h3 {
					margin-top: 0;
					color: #50b3a2;
				}
				.user-list {
					list-style: none;
					padding: 0;
				}
				.user-list li {
					padding: 10px;
					border-bottom: 1px solid #ddd;
				}
				.user-list li:last-child {
					border-bottom: none;
				}
				.user-list a {
					color: #50b3a2;
					text-decoration: none;
				}
				.user-list a:hover {
					text-decoration: underline;
				}
				footer {
					padding: 20px;
					margin-top: 20px;
					color: #ffffff;
					background-color: #50b3a2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<header>
				<div class="container">
					<h1>CycleSync</h1>
					<div class="nav">
						<a href="/">Home</a>
						<a href="/dashboard" class="highlight">Dashboard</a>
						<a href="/profile">Profile</a>
						<a href="#" id="logout-link">Logout</a>
					</div>
				</div>
			</header>

			<div class="container">
				<div class="dashboard">
					<h2>Dashboard</h2>
					
					<div class="create-post">
						<h3>Create a New Post</h3>
						<form id="post-form">
							<div class="form-group">
								<label for="title">Title:</label>
								<input type="text" id="title" name="title" required>
							</div>
							<div class="form-group">
								<label for="content">Content:</label>
								<textarea id="content" name="content" required></textarea>
							</div>
							<div class="form-group">
								<button type="submit">Post</button>
							</div>
						</form>
					</div>
					
					<div class="posts">
						<h3>Recent Posts</h3>
						<div id="posts-container">
							<!-- Posts will be loaded here -->
						</div>
					</div>
					
					<div class="users">
						<h3>Users</h3>
						<ul class="user-list" id="users-container">
							<!-- Users will be loaded here -->
						</ul>
					</div>
				</div>
			</div>

			<footer>
				<p>CycleSync &copy; 2025</p>
			</footer>

			<script>
				// This would be a client-side script to manage posts and users
				// For our simple demo, we'll just add some hardcoded content

				// Load posts
				function loadPosts() {
					const postsContainer = document.getElementById('posts-container');
					postsContainer.innerHTML = '';
					
					// Demo posts (these would normally come from the server)
					const posts = [
						{
							id: 1,
							user_id: 1,
							username: 'admin',
							title: 'Admin Post',
							content: 'This is a post by admin with some content.',
							created: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString()
						},
						{
							id: 2,
							user_id: 2,
							username: 'user1',
							title: 'User Post',
							content: 'This is a post by user1 with some different content.',
							created: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString()
						}
					];
					
					posts.forEach(post => {
						const postElement = document.createElement('div');
						postElement.className = 'post';
						postElement.innerHTML = `
							<h3>${escapeHtml(post.title)}</h3>
							<div class="meta">
								Posted by <a href="/profile?id=${post.user_id}">${escapeHtml(post.username)}</a> on ${formatDate(post.created)}
							</div>
							<div class="content">
								${escapeHtml(post.content)}
							</div>
							<div class="actions">
								<a href="#" onclick="editPost(${post.id})">Edit</a>
								<a href="#" onclick="deletePost(${post.id})">Delete</a>
							</div>
						`;
						postsContainer.appendChild(postElement);
					});
				}
				
				// Load users
				function loadUsers() {
					const usersContainer = document.getElementById('users-container');
					usersContainer.innerHTML = '';
					
					// Demo users (these would normally come from the server)
					const users = [
						{
							id: 1,
							username: 'admin'
						},
						{
							id: 2,
							username: 'user1'
						}
					];
					
					users.forEach(user => {
						const userElement = document.createElement('li');
						userElement.innerHTML = `
							<a href="/profile?id=${user.id}">${escapeHtml(user.username)}</a>
						`;
						usersContainer.appendChild(userElement);
					});
				}
				
				// Helper functions
				function formatDate(dateString) {
					const date = new Date(dateString);
					return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
				}
				
				function escapeHtml(str) {
					return str
						.replace(/&/g, "&amp;")
						.replace(/</g, "&lt;")
						.replace(/>/g, "&gt;")
						.replace(/"/g, "&quot;")
						.replace(/'/g, "&#039;");
				}
				
				// Post form handler
				document.getElementById('post-form').addEventListener('submit', function(e) {
					e.preventDefault();
					
					const title = document.getElementById('title').value;
					const content = document.getElementById('content').value;
					
					// In a real app, we'd send this to the server
					// For demo, just reload the posts
					
					// Clear form
					document.getElementById('title').value = '';
					document.getElementById('content').value = '';
					
					// Reload posts
					loadPosts();
					
					alert('Post created successfully!');
				});
				
				// Logout handler
				document.getElementById('logout-link').addEventListener('click', function(e) {
					e.preventDefault();
					
					// In a real app, we'd call the logout API
					// For demo, just redirect to home
					window.location.href = '/';
				});
				
				// Simulated actions
				function editPost(postId) {
					alert('Editing post ID: ' + postId);
					// In a real app, we'd open an edit form
				}
				
				function deletePost(postId) {
					if (confirm('Are you sure you want to delete this post?')) {
						alert('Post ID: ' + postId + ' deleted');
						// In a real app, we'd call the delete API and remove the post
						loadPosts();
					}
				}
				
				// Load data on page load
				document.addEventListener('DOMContentLoaded', function() {
					loadPosts();
					loadUsers();
				});
			</script>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ProfileHandler renders the profile page
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("profile").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CycleSync - Profile</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 0;
					padding: 0;
					background-color: #f4f4f4;
				}
				.container {
					width: 80%;
					margin: auto;
					overflow: hidden;
				}
				header {
					background: #50b3a2;
					color: white;
					padding-top: 30px;
					min-height: 70px;
					border-bottom: #e8491d 3px solid;
				}
				header a {
					color: #ffffff;
					text-decoration: none;
					text-transform: uppercase;
					font-size: 16px;
				}
				header h1 {
					margin-bottom: 10px;
				}
				header .nav {
					margin-top: 10px;
					display: flex;
					justify-content: center;
				}
				header .nav a {
					margin: 0 15px;
					padding: 5px;
				}
				.highlight, .current a {
					font-weight: bold;
					background-color: #e8491d;
					border-radius: 3px;
					padding: 3px 7px;
				}
				.profile {
					margin: 50px 0;
				}
				.profile h2 {
					color: #333;
				}
				.profile-info {
					background: #fff;
					padding: 20px;
					margin-bottom: 30px;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.profile-info h3 {
					margin-top: 0;
					color: #50b3a2;
				}
				.profile-info p {
					margin: 10px 0;
				}
				.profile-info .label {
					font-weight: bold;
					display: inline-block;
					width: 100px;
				}
				.posts {
					margin-top: 30px;
				}
				.post {
					background: #fff;
					padding: 15px;
					margin-bottom: 15px;
					border-radius: 5px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
				}
				.post h3 {
					margin-top: 0;
					color: #50b3a2;
				}
				.post .meta {
					color: #888;
					font-size: 0.8em;
					margin-bottom: 10px;
				}
				.post .content {
					line-height: 1.6;
				}
				.post .actions {
					margin-top: 10px;
				}
				.post .actions a {
					color: #50b3a2;
					margin-right: 10px;
					text-decoration: none;
				}
				.post .actions a:hover {
					text-decoration: underline;
				}
				.edit-profile {
					margin-top: 20px;
				}
				.edit-profile button {
					background-color: #50b3a2;
					color: white;
					border: none;
					padding: 10px 15px;
					border-radius: 3px;
					cursor: pointer;
				}
				.edit-profile button:hover {
					background-color: #429e8d;
				}
				footer {
					padding: 20px;
					margin-top: 20px;
					color: #ffffff;
					background-color: #50b3a2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<header>
				<div class="container">
					<h1>CycleSync</h1>
					<div class="nav">
						<a href="/">Home</a>
						<a href="/dashboard">Dashboard</a>
						<a href="/profile" class="highlight">Profile</a>
						<a href="#" id="logout-link">Logout</a>
					</div>
				</div>
			</header>

			<div class="container">
				<div class="profile">
					<h2>User Profile</h2>
					
					<div class="profile-info" id="profile-info">
						<!-- Profile info will be loaded here -->
					</div>
					
					<div class="posts">
						<h3>User's Posts</h3>
						<div id="posts-container">
							<!-- Posts will be loaded here -->
						</div>
					</div>
				</div>
			</div>

			<footer>
				<p>CycleSync &copy; 2025</p>
			</footer>

			<script>
				// In a real app, these would come from API calls
				
				// Get user ID from URL parameter
				const urlParams = new URLSearchParams(window.location.search);
				const userId = urlParams.get('id') || 1; // Default to user ID 1 if not specified
				
				// Load user profile
				function loadUserProfile(userId) {
					const profileContainer = document.getElementById('profile-info');
					
					// Demo user data (would come from server)
					const user = {
						id: userId,
						username: userId == 1 ? 'admin' : 'user1',
						email: userId == 1 ? 'admin@example.com' : 'user1@example.com',
						created: new Date(Date.now() - (userId == 1 ? 24 : 12) * 60 * 60 * 1000).toISOString()
					};
					
					profileContainer.innerHTML = `
						<h3>${escapeHtml(user.username)}'s Profile</h3>
						<p><span class="label">Username:</span> ${escapeHtml(user.username)}</p>
						<p><span class="label">Email:</span> ${escapeHtml(user.email)}</p>
						<p><span class="label">Joined:</span> ${formatDate(user.created)}</p>
						
						<!-- IDOR Vulnerability: Any logged-in user can see this button and edit any profile -->
						<div class="edit-profile">
							<button onclick="editProfile(${user.id})">Edit Profile</button>
						</div>
					`;
				}
				
				// Load user posts
				function loadUserPosts(userId) {
					const postsContainer = document.getElementById('posts-container');
					postsContainer.innerHTML = '';
					
					// Demo posts (would come from server)
					const posts = userId == 1 ? 
						[{
							id: 1,
							user_id: 1,
							username: 'admin',
							title: 'Admin Post',
							content: 'This is a post by admin with some content.',
							created: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString()
						}] : 
						[{
							id: 2,
							user_id: 2,
							username: 'user1',
							title: 'User Post',
							content: 'This is a post by user1 with some different content.',
							created: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString()
						}];
					
					if (posts.length === 0) {
						postsContainer.innerHTML = '<p>No posts yet.</p>';
						return;
					}
					
					posts.forEach(post => {
						const postElement = document.createElement('div');
						postElement.className = 'post';
						postElement.innerHTML = `
							<h3>${escapeHtml(post.title)}</h3>
							<div class="meta">
								Posted on ${formatDate(post.created)}
							</div>
							<div class="content">
								${escapeHtml(post.content)}
							</div>
							<!-- IDOR Vulnerability: Any logged-in user can see these buttons and modify any post -->
							<div class="actions">
								<a href="#" onclick="editPost(${post.id})">Edit</a>
								<a href="#" onclick="deletePost(${post.id})">Delete</a>
							</div>
						`;
						postsContainer.appendChild(postElement);
					});
				}
				
				// Helper functions
				function formatDate(dateString) {
					const date = new Date(dateString);
					return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
				}
				
				function escapeHtml(str) {
					return str
						.replace(/&/g, "&amp;")
						.replace(/</g, "&lt;")
						.replace(/>/g, "&gt;")
						.replace(/"/g, "&quot;")
						.replace(/'/g, "&#039;");
				}
				
				// Simulated actions (IDOR vulnerabilities)
				function editProfile(userId) {
					alert('Editing profile ID: ' + userId);
					// In a real app with proper security, this should check if the logged-in user
					// is the same as the profile being edited
				}
				
				function editPost(postId) {
					alert('Editing post ID: ' + postId);
					// In a real app with proper security, this should check if the logged-in user
					// is the owner of the post
				}
				
				function deletePost(postId) {
					if (confirm('Are you sure you want to delete this post?')) {
						alert('Post ID: ' + postId + ' deleted');
						// In a real app with proper security, this should check if the logged-in user
						// is the owner of the post
						loadUserPosts(userId);
					}
				}
				
				// Logout handler
				document.getElementById('logout-link').addEventListener('click', function(e) {
					e.preventDefault();
					window.location.href = '/';
				});
				
				// Load data on page load
				document.addEventListener('DOMContentLoaded', function() {
					loadUserProfile(userId);
					loadUserPosts(userId);
				});
			</script>
		</body>
		</html>
	`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// UsersHandler handles requests for all users
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all users
		sendJSONResponse(w, true, "", users, http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UserHandler handles requests for a specific user
// VULNERABLE TO IDOR: No authorization check on user access
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from path
	idStr := r.URL.Path[len("/api/user/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	// Find user by ID
	var user *User
	for _, u := range users {
		if u.ID == id {
			user = &u
			break
		}
	}

	if user == nil {
		sendJSONResponse(w, false, "User not found", nil, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		sendJSONResponse(w, true, "", user, http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PostsHandler handles requests for all posts
func PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all posts or filter by user ID
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr != "" {
			var userID int
			fmt.Sscanf(userIDStr, "%d", &userID)

			// Filter posts by user ID
			var userPosts []Post
			for _, post := range posts {
				if post.UserID == userID {
					userPosts = append(userPosts, post)
				}
			}
			sendJSONResponse(w, true, "", userPosts, http.StatusOK)
			return
		}

		// Get all posts
		sendJSONResponse(w, true, "", posts, http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PostHandler handles requests for a specific post
// VULNERABLE TO IDOR: No authorization check for viewing/modifying posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from path
	idStr := r.URL.Path[len("/api/post/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	// Find post by ID
	var post *Post
	for _, p := range posts {
		if p.ID == id {
			post = &p
			break
		}
	}

	if post == nil {
		sendJSONResponse(w, false, "Post not found", nil, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		sendJSONResponse(w, true, "", post, http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Starting CycleSync application...")

	// Static file server setup (not used in this simplified version)
	// fs := http.FileServer(http.Dir("static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main routes
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", LoginPageHandler)
	http.HandleFunc("/signup", SignupPageHandler)
	http.HandleFunc("/dashboard", DashboardHandler)
	http.HandleFunc("/profile", ProfileHandler)

	// API routes
	http.HandleFunc("/api/users", UsersHandler)
	http.HandleFunc("/api/user/", UserHandler) // Vulnerable to IDOR
	http.HandleFunc("/api/posts", PostsHandler)
	http.HandleFunc("/api/post/", PostHandler) // Vulnerable to IDOR

	// Serve on port 5000
	log.Println("Server starting on http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}