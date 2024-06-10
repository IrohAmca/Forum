document.addEventListener("DOMContentLoaded", () => {
  fetch('/get-databyid')
      .then(response => response.json())
      .then(data => {
          if (data.success) {
              const postsContainer = document.getElementById('posts');
              data.data.forEach(post => {
                  const postElement = document.createElement('div');
                  postElement.classList.add('post');
                  postElement.innerHTML = `
                      <h2>Post ID: ${post.PostID}</h2>
                      <p>Thread ID: ${post.ThreadID}</p>
                      <p>Content: ${post.Content}</p>
                      <p>Created At: ${post.CreatedAt}</p>
                  `;
                  postsContainer.appendChild(postElement);
              });
          } else {
              alert('Failed to load posts: ' + data.message);
          }
      })
      .catch(error => console.error('Error fetching posts:', error));
});
$('#signOutButton').hide();
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
        $('#signInModal').modal('hide');
        alert(data.message);
        $(document).ready(function () {
          var userId = getCookie('user_id');
          //console.log(userId);
          if (userId) {
            $('#signInButton').hide();
            $('#signUpButton').hide();
            $('#signOutButton').show();
            $('#profileButton').show();
          }
        });
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
        $('#profileButton').hide();
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
