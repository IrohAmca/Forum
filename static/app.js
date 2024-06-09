function CheckToken() {
  var token = getCookie('token');
  fetch('/check-token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ token: token })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        $('#signInButton').hide();
        $('#signUpButton').hide();
        $('#signOutButton').show();
      }
      else {
        $('#signOutButton').hide();
        $('#signInButton').show();
        $('#signUpButton').show();
      }
    })
    .catch(error => console.error('Error:', error));
}
CheckToken();
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
        location.reload();
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
      if (data.success) {
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
        alert(data.message);
        location.reload();
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

function deletePost(PostID) {
  fetch('/delete-post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ PostID: PostID })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        location.reload();
      } else {
        alert("Error deleting post: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
}

function getDeleteButtonHtml(postToken, PostID) {
  var token = getCookie('token');

  if (token == postToken) {
    return '<button class="delete-btn" onclick="deletePost(\'' + PostID + '\')">Delete</button>';
  }
  return '';
}
window.writeComment = function (button) {
  var replyForm = button.closest('.post').querySelector('.reply-form');
  replyForm.style.display = 'block';
};

window.submitComment = function (button) {
  var replyForm = button.closest('.reply-form');
  var commentText = replyForm.querySelector('input').value;
  var newComment = document.createElement('div');
  var postId = replyForm.closest('.post').dataset.postId;

  fetch('/create-comment', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ postId: postId, comment: commentText })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        alert(data.message);
      } else {
        alert("Error creating comment: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
  newComment.classList.add('comment');
  newComment.innerHTML = '<p>' + commentText + '</p><div class="buttons"><button class="like-dislike-btn" onclick="likeComment(this)"><img src="/png/dislike.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikeComment(this)"><img src="../png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyComment(this)">Reply</button>' + getDeleteButtonHtml(post.UserToken, post.PostID) + '</div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a reply..."><button class="btn btn-primary" onclick="submitReply(this)">Submit</button></div>';
  replyForm.insertAdjacentElement('afterend', newComment);
  replyForm.style.display = 'none';
};

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
          newPost.innerHTML = '<h2 class="blog-post-title">'
            + post.Title +
            '</h2><p class="blog-post-meta">'
            + post.CreatedAt +
            ' by <a href="#">'
            + post.Username +
            '</a></p><p>' +
            post.Content +
            '</p><hr><div class="buttons"><button class="like-dislike-btn" onclick="likePost(this)"><img src="../png/like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikePost(this)"><img src="../png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="writeComment(this)">Comment</button>'
            + getDeleteButtonHtml(post.UserToken, post.PostID) +
            '</div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a comment..."><button class="btn btn-primary" onclick="submitComment(this)">Submit</button></div>';

          newPost.dataset.postId = post.PostID;

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