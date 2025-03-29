document.addEventListener('DOMContentLoaded', function() {
    // Get DOM elements
    const usernameElement = document.getElementById('username');
    const postForm = document.getElementById('post-form');
    const postErrorMessage = document.getElementById('post-error-message');
    const myPostsList = document.getElementById('my-posts-list');
    const usersList = document.getElementById('users-list');
    
    // Get current user data
    fetch('/api/user/0', { // The backend will return the current user for ID 0
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Set username
            usernameElement.textContent = data.data.username;
            
            // Load user's posts
            loadUserPosts(data.data.id);
        } else {
            // Redirect to login if not authenticated
            window.location.href = '/login';
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
    
    // Load all users
    loadUsers();
    
    // Post form submission
    if (postForm) {
        postForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const title = document.getElementById('post-title').value;
            const content = document.getElementById('post-content').value;
            
            // Validate form
            if (!title || !content) {
                postErrorMessage.textContent = 'Please fill out all fields';
                postErrorMessage.classList.remove('hidden');
                return;
            }
            
            // Create post
            fetch('/api/posts', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    title: title,
                    content: content
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Clear form
                    postForm.reset();
                    
                    // Hide error message
                    postErrorMessage.classList.add('hidden');
                    
                    // Reload posts
                    loadUserPosts(data.data.user_id);
                } else {
                    // Display error message
                    postErrorMessage.textContent = data.message || 'Failed to create post';
                    postErrorMessage.classList.remove('hidden');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                postErrorMessage.textContent = 'An error occurred. Please try again.';
                postErrorMessage.classList.remove('hidden');
            });
        });
    }
    
    // Load current user's posts
    function loadUserPosts(userId) {
        fetch(`/api/posts?user_id=${userId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Display posts
                displayPosts(data.data, myPostsList);
            } else {
                myPostsList.innerHTML = `<p>Error loading posts: ${data.message}</p>`;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            myPostsList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        });
    }
    
    // Load all users
    function loadUsers() {
        fetch('/api/users', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Display users
                displayUsers(data.data);
            } else {
                usersList.innerHTML = `<p>Error loading users: ${data.message}</p>`;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            usersList.innerHTML = '<p>Error loading users. Please try again.</p>';
        });
    }
    
    // Display posts
    function displayPosts(posts, container) {
        if (!posts || posts.length === 0) {
            container.innerHTML = '<p>No posts found.</p>';
            return;
        }
        
        let html = '';
        posts.forEach(post => {
            html += `
                <div class="post-item" data-post-id="${post.id}">
                    <div class="post-header">
                        <span class="post-title">${escapeHtml(post.title)}</span>
                        <span class="post-meta">Post ID: ${post.id}</span>
                    </div>
                    <div class="post-content">${escapeHtml(post.content)}</div>
                    <div class="post-meta">Created: ${formatDate(post.created_at)}</div>
                    <div class="post-actions">
                        <button class="button button-small edit-post-btn">Edit</button>
                        <button class="button button-small button-danger delete-post-btn">Delete</button>
                    </div>
                </div>
            `;
        });
        
        container.innerHTML = html;
        
        // Add event listeners for edit and delete buttons
        container.querySelectorAll('.edit-post-btn').forEach(button => {
            button.addEventListener('click', function() {
                const postId = this.closest('.post-item').dataset.postId;
                editPost(postId);
            });
        });
        
        container.querySelectorAll('.delete-post-btn').forEach(button => {
            button.addEventListener('click', function() {
                const postId = this.closest('.post-item').dataset.postId;
                deletePost(postId);
            });
        });
    }
    
    // Display users
    function displayUsers(users) {
        if (!users || users.length === 0) {
            usersList.innerHTML = '<p>No users found.</p>';
            return;
        }
        
        let html = '';
        users.forEach(user => {
            html += `
                <div class="user-item">
                    <div class="user-info">
                        <span class="user-username">${escapeHtml(user.username)}</span>
                        <span class="user-email">${escapeHtml(user.email)}</span>
                    </div>
                    <button class="button button-small view-user-profile" data-user-id="${user.id}">View Profile</button>
                </div>
            `;
        });
        
        usersList.innerHTML = html;
        
        // Add event listeners for view profile buttons
        usersList.querySelectorAll('.view-user-profile').forEach(button => {
            button.addEventListener('click', function() {
                const userId = this.dataset.userId;
                // IDOR vulnerability: Direct navigation to user profile by ID
                window.location.href = `/profile?id=${userId}`;
            });
        });
    }
    
    // Edit post
    function editPost(postId) {
        fetch(`/api/post/${postId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const post = data.data;
                const newTitle = prompt('Edit title:', post.title);
                if (newTitle === null) return; // User cancelled
                
                const newContent = prompt('Edit content:', post.content);
                if (newContent === null) return; // User cancelled
                
                // Update post
                fetch(`/api/post/${postId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        title: newTitle,
                        content: newContent
                    })
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Reload posts to show updated content
                        loadUserPosts(data.data.user_id);
                    } else {
                        alert('Failed to update post: ' + (data.message || 'Unknown error'));
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('An error occurred while updating the post');
                });
            } else {
                alert('Failed to fetch post details: ' + (data.message || 'Unknown error'));
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('An error occurred while fetching post details');
        });
    }
    
    // Delete post
    function deletePost(postId) {
        if (confirm('Are you sure you want to delete this post?')) {
            fetch(`/api/post/${postId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Remove post from DOM
                    document.querySelector(`.post-item[data-post-id="${postId}"]`).remove();
                    
                    // Show message if no posts left
                    if (myPostsList.querySelectorAll('.post-item').length === 0) {
                        myPostsList.innerHTML = '<p>No posts found.</p>';
                    }
                } else {
                    alert('Failed to delete post: ' + (data.message || 'Unknown error'));
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred while deleting the post');
            });
        }
    }
    
    // Helper functions
    
    // Format date
    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString();
    }
    
    // Escape HTML to prevent XSS
    function escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
});
