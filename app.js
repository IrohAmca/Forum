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

  fetch('/sign-up', {
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



//forum kısmı -- comment, like- dislike


document.addEventListener('DOMContentLoaded', function () {
  function updateCount(button, countClass, increment) {
    var countElement = button.querySelector(countClass);
    var count = parseInt(countElement.textContent);
    countElement.textContent = count + increment;
  }

  window.likePost = function (button) {
    updateCount(button, '.like-count', 1);
  };

  window.dislikePost = function (button) {
    updateCount(button, '.dislike-count', 1);
  };

  window.replyPost = function (button) {
    var replyForm = button.closest('.post').querySelector('.reply-form');
    replyForm.style.display = 'block';
  };

  window.submitComment = function (button) {
    var replyForm = button.closest('.reply-form');
    var commentText = replyForm.querySelector('input').value;
    var newComment = document.createElement('div');
    newComment.classList.add('comment');
    newComment.innerHTML = '<p>' + commentText + '</p><div class="buttons"><button class="like-dislike-btn" onclick="likeComment(this)"><img src="/png/dislike.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikeComment(this)"><img src="/png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyComment(this)">Reply</button><button class="delete-btn" onclick="deleteComment(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a reply..."><button class="btn btn-primary" onclick="submitReply(this)">Submit</button></div>';
    replyForm.insertAdjacentElement('afterend', newComment);
    replyForm.style.display = 'none';
  };

  window.likeComment = function (button) {
    updateCount(button, '.like-count', 1);
  };

  window.dislikeComment = function (button) {
    updateCount(button, '.dislike-count', 1);
  };

  window.replyComment = function (button) {
    var comment = button.closest('.comment');
    var replyForm = comment.querySelector('.reply-form');
    replyForm.style.display = 'block';
  };

  window.submitReply = function (button) {
    var replyForm = button.closest('.reply-form');
    var replyText = replyForm.querySelector('input').value;
    var newReply = document.createElement('div');
    newReply.classList.add('comment');
    newReply.innerHTML = '<p>' + replyText + '</p><div class="buttons"><button class="like-dislike-btn" onclick="likeComment(this)"><img src="/png/like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikeComment(this)"><img src="/png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyComment(this)">Reply</button><button class="delete-btn" onclick="deleteComment(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a reply..."><button class="btn btn-primary" onclick="submitReply(this)">Submit</button></div>';
    replyForm.insertAdjacentElement('afterend', newReply);
    replyForm.style.display = 'none';
  };

  window.deletePost = function (button) {
    var post = button.closest('.post');
    post.remove();
  };

  window.deleteComment = function (button) {
    var comment = button.closest('.comment');
    comment.remove();
  };

  window.createPost = function () {
    var title = document.getElementById('new-post-title').value;
    var content = document.getElementById('new-post-content').value;

    var newPost = document.createElement('article');
    newPost.classList.add('post');
    newPost.innerHTML = '<h2 class="blog-post-title">' + title + '</h2><p class="blog-post-meta">January 1, 2021 by <a href="#">Mark</a></p><p>' + content + '</p><hr><div class="buttons"><button class="like-dislike-btn" onclick="likePost(this)"><img src="/png/like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikePost(this)"><img src="/png/dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyPost(this)">Comment</button><button class="delete-btn" onclick="deletePost(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a comment..."><button class="btn btn-primary" onclick="submitComment(this)">Submit</button></div>';

    var postList = document.querySelector('.post-list');
    postList.prepend(newPost);
  };
});
