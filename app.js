document.getElementById('loginForm').addEventListener('submit', function(event) {
  event.preventDefault(); // Formun varsayılan submit davranışını engelle

  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

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
      document.getElementById('user-email').textContent = email;
      $('#signInModal').modal('hide'); // Modal'ı kapat
      document.getElementById('user-content').style.display = 'block';
    } else {
      alert('Giriş başarısız: ' + data.message);
    }
  })
  .catch(error => {
    console.error('Error:', error);
    alert('Sunucu hatası. Lütfen daha sonra tekrar deneyin.');
  });
});

document.getElementById('signUpForm').addEventListener('submit', function(event) {
  event.preventDefault(); // Formun varsayılan submit davranışını engelle

  const name = document.getElementById('name').value;
  const email = document.getElementById('signUpEmail').value;
  const password = document.getElementById('signUpPassword').value;

  fetch('/signup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name: name, email: email, password: password })
  })
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      $('#signUpModal').modal('hide'); // Modal'ı kapat
      alert('Kayıt başarılı! Şimdi giriş yapabilirsiniz.');
    } else {
      alert('Kayıt başarısız: ' + data.message);
    }
  })
  .catch(error => {
    console.error('Error:', error);
    alert('Sunucu hatası. Lütfen daha sonra tekrar deneyin.');
  });
});





document.addEventListener('DOMContentLoaded', function() {
  const postContainer = document.getElementById('postContainer');
  const newPostForm = document.getElementById('newPostForm');

  newPostForm.addEventListener('submit', function(event) {
    event.preventDefault();
    const postContent = document.getElementById('postContent').value;
    addPost(postContent);
    newPostForm.reset();
  });

  function addPost(content) {
    const postId = Date.now();
    const postDiv = document.createElement('div');
    postDiv.className = 'post';
    postDiv.id = `post-${postId}`;
    postDiv.innerHTML = `
      <button class="delete-btn" onclick="deletePost(${postId})">&times;</button>
      <p>${content}</p>
      <div class="buttons">
        <button class="btn btn-success" onclick="likePost(${postId})">Like <span id="like-count-${postId}">0</span></button>
        <button class="btn btn-danger" onclick="dislikePost(${postId})">Dislike <span id="dislike-count-${postId}">0</span></button>
        <button class="btn btn-secondary" onclick="toggleCommentForm(${postId})">Comment</button>
      </div>
      <div id="comments-${postId}" class="mt-3"></div>
      <div id="commentForm-${postId}" class="comment-form" style="display: none;">
        <textarea class="form-control" rows="2" placeholder="Write a comment..." id="commentContent-${postId}"></textarea>
        <button class="btn btn-primary mt-2" onclick="addComment(${postId})">Add Comment</button>
      </div>
    `;
    postContainer.appendChild(postDiv);
  }

  window.toggleCommentForm = function(postId) {
    const commentForm = document.getElementById(`commentForm-${postId}`);
    if (commentForm.style.display === 'none') {
      commentForm.style.display = 'block';
    } else {
      commentForm.style.display = 'none';
    }
  }

  window.addComment = function(postId) {
    const commentContent = document.getElementById(`commentContent-${postId}`).value;
    const commentsDiv = document.getElementById(`comments-${postId}`);
    const commentDiv = document.createElement('div');
    commentDiv.className = 'comment';
    commentDiv.innerHTML = `
      <button class="delete-btn" onclick="deleteComment(${postId}, this)">&times;</button>
      <p>${commentContent}</p>
    `;
    commentsDiv.appendChild(commentDiv);
    document.getElementById(`commentForm-${postId}`).style.display = 'none';
    document.getElementById(`commentContent-${postId}`).value = ''; // Clear comment textarea
  }

  window.likePost = function(postId) {
    const likeCount = document.getElementById(`like-count-${postId}`);
    likeCount.textContent = parseInt(likeCount.textContent) + 1;
  }

  window.dislikePost = function(postId) {
    const dislikeCount = document.getElementById(`dislike-count-${postId}`);
    dislikeCount.textContent = parseInt(dislikeCount.textContent) + 1;
  }

  window.deletePost = function(postId) {
    const postDiv = document.getElementById(`post-${postId}`);
    postDiv.remove();
  }

  window.deleteComment = function(postId, commentElement) {
    const commentDiv = commentElement.parentElement;
    commentDiv.remove();
  }
});