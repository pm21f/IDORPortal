document.addEventListener('DOMContentLoaded', function() {
    // Get DOM elements
    const userIdElement = document.getElementById('user-id');
    const usernameElement = document.getElementById('username');
    const emailElement = document.getElementById('email');
    const createdAtElement = document.getElementById('created-at');
    const currentUserIdElement = document.getElementById('current-user-id');
    const editUsernameInput = document.getElementById('edit-username');
    const editEmailInput = document.getElementById('edit-email');
    const profileForm = document.getElementById('profile-form');
    const profileErrorMessage = document.getElementById('profile-error-message');
    const deleteAccountBtn = document.getElementById('delete-account-btn');
    const userPostsList = document.getElementById('user-posts-list');
    const viewUserBtn = document.getElementById('view-user-btn');
    
    // Get user ID from URL query parameter if any
    const urlParams = new URLSearchParams(window.location.search);
    let userIdToView = urlParams.get('id');
    
    // If no user ID in URL, get current user
    if (!userIdToView) {
        userIdToView = 0; // The backend will return the current user for ID 0
    }
    
    // Load user profile
    loadUserProfile(userIdToView);
    
    // Load user posts
    loadUserPosts(userIdToView);
    
    // Profile form submission
    if (profileForm) {
        profileForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            const username = editUsernameInput.value;
            const email = editEmailInput.value;
            
            // Validate form
            if (!username || !email) {
                profileErrorMessage.textContent = 'Please fill out all fields';
                profileErrorMessage.classList.remove('hidden');
                return;
            }
            
            // Validate email format
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(email)) {
                profileErrorMessage.textContent = 'Please enter a valid email address';
                profileErrorMessage.classList.remove('hidden');
                return;
            }
            
            // Update profile
            // IDOR vulnerability: Can update any user's profile by changing the ID
            fetch(`/api/user/${userIdToView}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    username: username,
                    email: email
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Update displayed profile info
                    usernameElement.textContent = data.data.username;
                    emailElement.textContent = data.data.email;
                    
                    // Show success message
                    profileErrorMessage.textContent = 'Profile updated successfully!';
                    profileErrorMessage.style.backgroundColor = '#d4edda';
                    profileErrorMessage.style.color = '#155724';
                    profileErrorMessage.classList.remove('hidden');
                    
                    // Hide success message after 3 seconds
                    setTimeout(function() {
                        profileErrorMessage.classList.add('hidden');
                    }, 3000);
                } else {
                    // Display error message
                    profileErrorMessage.textContent = data.message || 'Failed to update profile';
                    profileErrorMessage.style.backgroundColor = '#f8d7da';
                    profileErrorMessage.style.color = '#721c24';
                    profileErrorMessage.classList.remove('hidden');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                profileErrorMessage.textContent = 'An error occurred. Please try again.';
                profileErrorMessage.style.backgroundColor = '#f8d7da';
                profileErrorMessage.style.color = '#721c24';
                profileErrorMessage.classList.remove('hidden');
            });
        });
    }
    
    // Delete account button
    if (deleteAccountBtn) {
        deleteAccountBtn.addEventListener('click', function() {
            if (confirm('Are you sure you want to delete this account? This action cannot be undone.')) {
                // IDOR vulnerability: Can delete any user's account by changing the ID
                fetch(`/api/user/${userIdToView}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Redirect to login page
                        window.location.href = '/login';
                    } else {
                        alert('Failed to delete account: ' + (data.message || 'Unknown error'));
                    }
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('An error occurred while deleting the account');
                });
            }
        });
    }
    
    // View user button
    if (viewUserBtn) {
        viewUserBtn.addEventListener('click', function() {
            const testUserId = document.getElementById('test-user-id').value;
            if (testUserId) {
                // IDOR vulnerability demonstration
                window.location.href = `/profile?id=${testUserId}`;
            }
        });
    }
    
    // Load user profile
    function loadUserProfile(userId) {
        fetch(`/api/user/${userId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const user = data.data;
                
                // Set profile information
                userIdElement.textContent = user.id;
                usernameElement.textContent = user.username;
                emailElement.textContent = user.email;
                createdAtElement.textContent = formatDate(user.created_at);
                
                // Set form values
                editUsernameInput.value = user.username;
                editEmailInput.value = user.email;
                
                // Set current user ID for vulnerability testing
                if (currentUserIdElement) {
                    currentUserIdElement.textContent = user.id;
                }
                
                // Display warning if viewing another user's profile
                if (userId != 0 && userId != user.id) {
                    const warningDiv = document.createElement('div');
                    warningDiv.className = 'warning';
                    warningDiv.innerHTML = `
                        <h3>⚠️ IDOR Vulnerability Detected</h3>
                        <p>You are currently viewing another user's profile (ID: ${user.id}).</p>
                        <p>This demonstrates an Insecure Direct Object Reference vulnerability.</p>
                    `;
                    
                    // Insert warning at the top of the profile container
                    const profileContainer = document.getElementById('profile-container');
                    profileContainer.insertBefore(warningDiv, profileContainer.firstChild);
                }
                
            } else {
                // Redirect to login if not authenticated or user not found
                window.location.href = '/login';
            }
        })
        .catch(error => {
            console.error('Error:', error);
            profileErrorMessage.textContent = 'An error occurred while loading the profile';
            profileErrorMessage.classList.remove('hidden');
        });
    }
    
    // Load user posts
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
                displayPosts(data.data);
            } else {
                userPostsList.innerHTML = `<p>Error loading posts: ${data.message}</p>`;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            userPostsList.innerHTML = '<p>Error loading posts. Please try again.</p>';
        });
    }
    
    // Display posts
    function displayPosts(posts) {
        if (!posts || posts.length === 0) {
            userPostsList.innerHTML = '<p>No posts found.</p>';
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
        
        userPostsList.innerHTML = html;
        
        // Add event listeners for edit and delete buttons
        userPostsList.querySelectorAll('.edit-post-btn').forEach(button => {
            button.addEventListener('click', function() {
                const postId = this.closest('.post-item').dataset.postId;
                editPost(postId);
            });
        });
        
        userPostsList.querySelectorAll('.delete-post-btn').forEach(button => {
            button.addEventListener('click', function() {
                const postId = this.closest('.post-item').dataset.postId;
                deletePost(postId);
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
                
                // Update post - IDOR vulnerability
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
                        loadUserPosts(userIdToView);
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
            // IDOR vulnerability
            fetch(`/api/post/${postId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Reload posts
                    loadUserPosts(userIdToView);
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
