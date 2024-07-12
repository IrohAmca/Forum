checkToken();
var token;
document.addEventListener('DOMContentLoaded', function () {
  getAllPosts();

  document.getElementById('search-form').addEventListener('submit', function (event) {
    event.preventDefault();
    getAllPosts();
  });
});

function checkToken() {
  var cookie = getCookie('cookie');
  fetch('/check-token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ cookie: cookie })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        console.log(data.username);
        token = data.token;
        const initial = data.username.charAt(0).toUpperCase();
        const profileIconHTML = `
          <a href="/profile/${data.username}" class="profile-icon">${initial}</a>`;
        document.getElementById('profileIconContainer').innerHTML = profileIconHTML;
        $('#profileIconContainer').show();
        $('#signInButton').hide();
        $('#signUpButton').hide();
        $('#signOutButton').show();
        $('#postForm').show();
      } else {
        $('#postForm').show();
        $('#signOutButton').hide();
        $('#signInButton').show();
        $('#signUpButton').show();
      }
    })
    .catch(error => console.error('Error:', error));
}
function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function deleteCookie(name) {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}
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
function loginWithGithub() {
  window.location.href = "/auth/github";
}

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
        location.reload();
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
        var postElement = document.querySelector(`[data-post-id="${PostID}"]`);
        if (postElement) {
          postElement.remove();
        }
      } else {
        alert("Error deleting post: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
}

function getDeletePostButtonHtml(postToken, PostID) {
  if (token == postToken) {
    return '<button class="delete-btn" onclick="deletePost(\'' + PostID + '\')"><img src="../png/delete.png" alt="Delete Icon"></button>';
  }
  return '';
}

function ld_submit(PostID, isLike) {
  fetch('/ld_post', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ PostID: PostID, isLike: isLike })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        getAllPosts();
      } else {
        alert("Please log in to the website first!!!");
      }
    })
    .catch(error => console.error('Error:', error));
  //location.reload();
}



function ld_comment_submit(CommentID, isLike) {
  fetch('/ld_comment', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ CommentID: CommentID, isLike: isLike })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        getAllPosts();
      } else {
        alert("Please log in to the website first!!!");
      }
    })
    .catch(error => console.error('Error:', error));
}


window.writeComment = function (button) {
  if (getCookie('cookie') === undefined) {
    alert('You must be logged in to comment');
    return;
  }
  var replyForm = button.closest('.post').querySelector('.reply-form');
  replyForm.style.display = 'block';
};

function DeleteComment(CommentID) {
  fetch('/delete-comment', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ CommentID: CommentID })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        location.reload();
      } else {
        alert("Error deleting comment: " + data.message);
      }
    }
    )
    .catch(error => console.error('Error:', error));
}

function getDeleteCommentButtonHtml(commentToken, CommentID) {
  if (token == commentToken) {
    return '<button class="delete-btn" onclick="DeleteComment(\'' + CommentID + '\')"><img src="../png/delete.png" alt="Delete Icon"></button>';
  }
  return '';
}

window.submitComment = function (button) {
  var replyForm = button.closest('.reply-form');
  var commentText = replyForm.querySelector('input').value;
  var postId = replyForm.closest('.post').dataset.postId;
  let selectedCategories = [];
  let checkboxes = document.querySelectorAll('input[name="category"]:checked');

  checkboxes.forEach((checkbox) => {
    selectedCategories.push(checkbox.value);
  });

  fetch('/create-comment', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ postId: postId, comment: commentText, categories: selectedCategories })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
      } else {
        alert("You cannot leave empty comments!!!");
      }
    })
    .catch(error => console.error('Error:', error));
  location.reload();
};

function getAllPosts() {
  let selectedCategories = [];
  let checkboxes = document.querySelectorAll('input[name="category-filter"]:checked');

  checkboxes.forEach((checkbox) => {
    selectedCategories.push(checkbox.value);
  });

  title = document.getElementById('keyword').value;
  short_type = document.getElementById('sort-by').value;
  fetch('/get-posts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ categories: selectedCategories, title: title, short_type: short_type })
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        var postList = document.querySelector('.post-list');
        postList.innerHTML = '';
        var posts = data.posts;
        posts.forEach(post => {
          console.log(post.Image)
          var newPost = document.createElement('article');
          newPost.classList.add('post');
          newPost.innerHTML = `<h2 class="blog-post-title">
            ${post.Title}
          </h2>
          <p class="blog-post-meta">
            ${post.CreatedAt}
            <a href="/profile/${post.Username}">
              ${post.Username}
            </a>
          </p>
          ${post.Image ? `<p><img src="${post.Image}" alt="Post Image" style="width: 40%; height: 50%;"></p>` : ''}
          <p>
            <div class="post-categories">${post.Categories}</div>
            ${post.Content}
          </p>
          <hr>
          <div class="buttons">
            <button class="like-dislike-btn" onclick="ld_submit('${post.PostID}', true)">
              <img src="../png/like.png" alt="Like Icon">
              <span class="like-count">${post.LikeCounter}</span>
            </button>
            <button class="like-dislike-btn" onclick="ld_submit('${post.PostID}', false)">
              <img src="../png/dislike.png" alt="Dislike Icon">
              
              <span class="dislike-count">${post.DislikeCounter}</span>
            </button>
            <button class="reply-btn" onclick="writeComment(this)"><img src="../png/comment.png" alt="Comment Icon"></button>
            ${getDeletePostButtonHtml(post.UserToken, post.PostID)}

    
          </div>

          <div class="reply-form" style="display:none;">
            <input type="text" class="form-control" placeholder="Write a comment...">
            <button class="btn btn-primary" onclick="submitComment(this)">Submit</button>
          </div>`;
          newPost.dataset.postId = post.PostID;
          var postList = document.querySelector('.post-list');
          postList.prepend(newPost);
          var comments = post.Comment;
          comments.forEach(comment => {
            var newComment = document.createElement('div');
            newComment.classList.add('comment');
            newComment.innerHTML = '<p class="blog-post-meta">'
              + comment.CreatedAt +
              ' by <a href="#">'
              + comment.Username +
              '</a></p><p>'
              + comment.Content +
              '</p><hr><div class="buttons"><button class="like-dislike-btn" onclick="ld_comment_submit(\'' + comment.CommentID + '\', true)"><img src="../png/like.png" alt="Like Icon"> <span class="like-count">' + comment.LikeCounter + '</span></button><button class="like-dislike-btn" onclick="ld_comment_submit(\'' + comment.CommentID + '\', false)"><img src="../png/dislike.png" alt="Dislike Icon"> <span class="dislike-count">' + comment.DislikeCounter + '</span></button>'
              + getDeleteCommentButtonHtml(post.UserToken, comment.CommentID) +
              '</div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a comment..."><button class="btn btn-primary" onclick="submitComment(this)">Submit</button></div>';

            newPost.appendChild(newComment);
          });
        });
      } else {
        alert("Error getting posts: " + data.message);
      }
    })
    .catch(error => console.error('Error:', error));
}

document.getElementById('postForm').addEventListener('submit', function (event) {
  event.preventDefault();
  var title = document.getElementById('postTitle').value;
  var content = document.getElementById('postContent').value;
  var selectedCategories = [];
  var checkboxes = document.querySelectorAll('input[name="category"]:checked');
  checkboxes.forEach((checkbox) => {
    selectedCategories.push(checkbox.value);
  });

  var fileInput = document.getElementById('image');
  var image = fileInput.files[0];
  var formData = new FormData();

  formData.append('title', title);
  formData.append('content', content);
  selectedCategories.forEach((category) => {
    formData.append('categories', category);
  });
  formData.append('image', image);

  fetch('/create-post', {
    method: 'POST',
    body: formData
  })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        location.reload();
      } else {
        alert(data.message);
      }
    })
    .catch(error => console.error('Error:', error));
});

function loginWithGoogle() {
  window.location.href = "/auth/google";
}

function loginWithGithub() {
  window.location.href = "/auth/github";
}

function loginWithFacebook() {
  window.location.href = "/auth/facebook";
}

document.getElementById("image").addEventListener("change", function (event) {
  var file = event.target.files[0];
  if (file) {
    var reader = new FileReader();
    reader.onload = function (e) {
      var img = document.getElementById("preview");
      img.src = e.target.result;
      img.style.display = "block";
      document.getElementById("cancel").style.display = "inline";
    }
    reader.readAsDataURL(file);
  }
});

document.getElementById("cancel").addEventListener("click", function () {
  document.getElementById("image").value = "";
  var img = document.getElementById("preview");
  img.src = "";
  img.style.display = "none";
  this.style.display = "none";
});

document.getElementById("uploadForm").addEventListener("submit", function (event) {
  event.preventDefault();

  var formData = new FormData(this);
  var file = document.getElementById("image").files[0];

  if (file.size > 20 * 1024 * 1024) {
    document.getElementById("message").innerText = "File size exceeds 20 MB";
    return;
  }

  fetch("/upload", {
    method: "POST",
    body: formData
  })
    .then(response => response.json())
    .then(data => {
      if (data.error) {
        document.getElementById("message").innerText = data.error;
      } else {
        document.getElementById("message").innerText = data.message;
      }
    })
    .catch(error => {
      console.error("Error:", error);
      document.getElementById("message").innerText = "An error occurred";
    });
});
