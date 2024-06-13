
function CheckCookie(){
  var user_id =getCookie("user_id")
  if (user_id) {
    $('#signInButton').hide();
    $('#signUpButton').hide();
    $('#signOutButton').show();
  }else{
    $('#signInButton').show();
    $('#signUpButton').show();
    $('#signOutButton').hide();
  }
}
CheckCookie();
document.getElementById('loginForm').addEventListener('submit', function (event) {
  event.preventDefault();

  var email = document.getElementById('loginEmail').value;
  var password = document.getElementById('loginPassword').value;

  fetch('/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email: email, password: password })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        alert(data.message);
        //location.reload();
      } else {
        alert('Error logging in user: ' + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
});

document.getElementById('signUpForm').addEventListener('submit', function (event) {
  event.preventDefault();

  const username = document.getElementById('signUpUsername').value;
  const email = document.getElementById('signUpEmail').value;
  const password = document.getElementById('signUpPassword').value;
  const confirmPassword = document.getElementById('confirmSignUpPassword').value;
  const passwordHelp = document.getElementById('passwordHelpBlock');

  if (password !== confirmPassword) {
    passwordHelp.textContent = 'Passwords do not match.';
    passwordHelp.style.color = 'red';
    return;
  } else {
    passwordHelp.textContent = '';
  }


  fetch('/sign-up', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username: username, email: email, password: password })
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      if (data.success) { // true
        $('#signUpModal').modal('hide');
        alert(data.message);
      } else {
        alert("Error signing up user: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));

});
document.getElementById('signOutButton').addEventListener('click', function () {
  fetch('/sign-out', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        $('#signOutButton').hide(); // Gizleme ve görüntüleme işlemlerini dinamik olarak cookie üzerinden kontrol eden bir fonksiyon yazılabilir.
        $('#signInButton').show();
        $('#signUpButton').show();
        alert(data.message);
      } else {
        alert("Error signing out: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
});

function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function deleteCookie(name) {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

// sign in - sign out ends

//user profile

document.getElementById('loginForm').addEventListener('submit', function (event) {
  event.preventDefault();
  const username = document.getElementById('loginEmail').value;
  showProfileIcon(username);
  $('#signInModal').modal('hide');
});

document.getElementById('registerForm').addEventListener('submit', function (event) {
  event.preventDefault();
  const username = document.getElementById('registerUsername').value;
  showProfileIcon(username);
  $('#signUpModal').modal('hide');
});

document.getElementById('signOutButton').addEventListener('click', function () {
  resetProfileIcon();
});

// Function to show profile icon and make it clickable
function showProfileIcon(username) {
  const initial = username.charAt(0).toUpperCase();
  const profileIconHTML = `
    <a href="userprofile.html" class="profile-icon">${initial}</a>
  `;
  document.getElementById('profileIconContainer').innerHTML = profileIconHTML;
  document.getElementById('signInButton').style.display = 'none';
  document.getElementById('signUpButton').style.display = 'none';
  document.getElementById('signOutButton').style.display = 'block';
}

// Function to reset profile icon to initial state
function resetProfileIcon() {
  document.getElementById('profileIconContainer').innerHTML = '';
  document.getElementById('signInButton').style.display = 'block';
  document.getElementById('signUpButton').style.display = 'block';
  document.getElementById('signOutButton').style.display = 'none';
}

//filter

const posts = [
  { id: 1, content: "Post about JavaScript", date: "2024-06-01", likes: 10 },
  { id: 2, content: "Learning Python", date: "2024-06-02", likes: 20 },
  { id: 3, content: "CSS Flexbox Guide", date: "2024-06-03", likes: 5 },
  { id: 4, content: "HTML Basics", date: "2024-06-04", likes: 15 },
  { id: 5, content: "Advanced React", date: "2024-06-05", likes: 8 }
];

function filterPosts() {
  const keyword = document.getElementById('keyword').value.toLowerCase();
  const sortBy = document.getElementById('sort-by').value;

  let filteredPosts = posts.filter(post => 
    post.content.toLowerCase().includes(keyword)
  );

  if (sortBy === 'date-asc') {
    filteredPosts.sort((a, b) => new Date(a.date) - new Date(b.date));
  } else if (sortBy === 'date-desc') {
    filteredPosts.sort((a, b) => new Date(b.date) - new Date(a.date));
  } else if (sortBy === 'likes-asc') {
    filteredPosts.sort((a, b) => a.likes - b.likes);
  } else if (sortBy === 'likes-desc') {
    filteredPosts.sort((a, b) => b.likes - a.likes);
  }

  renderPosts(filteredPosts);
}

function renderPosts(posts) {
  const postsList = document.getElementById('posts');
  postsList.innerHTML = '';
  
  posts.forEach(post => {
    const listItem = document.createElement('li');
    listItem.className = 'list-group-item';
    listItem.textContent = `${post.content} - Date: ${post.date} - Likes: ${post.likes}`;
    postsList.appendChild(listItem);
  });
}

// Initial render
renderPosts(posts);

// post - forum starts

function getAllPosts() {
  fetch('/get-posts', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        var posts = data.posts;
        posts.forEach(post => {
          var newPost = document.createElement('article');
          newPost.classList.add('post');
          newPost.innerHTML = '<h2 class="blog-post-title">' + post.Title + '</h2><p class="blog-post-meta">'+post.CreatedAt +' by <a href="#">' + post.Username + '</a></p><p>' + post.Content + '</p><hr><div class="buttons"><button class="like-dislike-btn" onclick="likePost(this)"><img src="../png/like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikePost(this)"><img src="../png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyPost(this)">Comment</button><button class="delete-btn" onclick="deletePost(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a comment..."><button class="btn btn-primary" onclick="submitComment(this)">Submit</button></div>';

          var postList = document.querySelector('.post-list');
          postList.prepend(newPost);
        });
      } else {
        alert("Error getting posts: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
}

getAllPosts();

document.getElementById('postForm').addEventListener('submit', function (event) {
  event.preventDefault();

  var title = document.getElementById('postTitle').value;
  var content = document.getElementById('postContent').value;

  fetch('/create-post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ title: title, content: content })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        location.reload();
      } else {
        alert("Error creating post: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
});



 