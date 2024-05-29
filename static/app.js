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
        $(document).ready(function() {
          var userId = getCookie('user_id');
          console.log(userId);
          if (userId) {
              $('#signInButton').hide();
              $('#signUpButton').hide();
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
      if (data.success) {
        $('#signUpModal').modal('hide');
        alert(data.message);
      } else {
        alert("Error signing up user: " + data.error);
      }
    })
    .catch(error => console.error('Error:', error));

});
function getCookie(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}


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
    newComment.innerHTML = '<p>' + commentText + '</p><div class="buttons"><button class="like-dislike-btn" onclick="likeComment(this)"><img src="dislike.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikeComment(this)"><img src="dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyComment(this)">Reply</button><button class="delete-btn" onclick="deleteComment(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a reply..."><button class="btn btn-primary" onclick="submitReply(this)">Submit</button></div>';
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
    newReply.innerHTML = '<p>' + replyText + '</p><div class="buttons"><button class="like-dislike-btn" onclick="likeComment(this)"><img src="like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikeComment(this)"><img src="dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyComment(this)">Reply</button><button class="delete-btn" onclick="deleteComment(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a reply..."><button class="btn btn-primary" onclick="submitReply(this)">Submit</button></div>';
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
    newPost.innerHTML = '<h2 class="blog-post-title">' + title + '</h2><p class="blog-post-meta">January 1, 2021 by <a href="#">Mark</a></p><p>' + content + '</p><hr><div class="buttons"><button class="like-dislike-btn" onclick="likePost(this)"><img src="like.png" alt="Like Icon">Like <span class="like-count">0</span></button><button class="like-dislike-btn" onclick="dislikePost(this)"><img src="dislike.png" alt="Dislike Icon">Dislike <span class="dislike-count">0</span></button><button class="reply-btn" onclick="replyPost(this)">Comment</button><button class="delete-btn" onclick="deletePost(this)">Sil</button></div><div class="reply-form" style="display:none;"><input type="text" class="form-control" placeholder="Write a comment..."><button class="btn btn-primary" onclick="submitComment(this)">Submit</button></div>';

    var postList = document.querySelector('.post-list');
    postList.prepend(newPost);
  };
});
document.getElementById('addPostButton').addEventListener('click', function () {
  addNewPost();
});


// filter

// Filtreleme Fonksiyonu
document.addEventListener('DOMContentLoaded', function () {
  // Kategori filtrelerini işle
  var categoryFilters = document.querySelectorAll('.category-filter');
  categoryFilters.forEach(function (filter) {
    filter.addEventListener('change', function () {
      // Seçilen kategorileri al
      var selectedCategories = Array.from(categoryFilters)
        .filter(function (checkbox) { return checkbox.checked; })
        .map(function (checkbox) { return checkbox.value; });
      console.log(selectedCategories); // Seçilen kategorileri konsola yazdır
      // Seçilen kategorilere göre işlem yapmak için bu bilgiyi kullanabilirsiniz
    });
  });

  // Likes ve Dates filtrelerini işle
  var filterTypeRadios = document.querySelectorAll('input[name="filterType"]');
  filterTypeRadios.forEach(function (radio) {
    radio.addEventListener('change', function () {
      var filterType = this.value; // Seçilen filtre türünü al
      console.log(filterType); // Seçilen filtreyi konsola yazdır
      // Seçilen filtreye göre işlem yapmak için bu bilgiyi kullanabilirsiniz
    });
  });
});



document.addEventListener('DOMContentLoaded', function () {
  // Handle sign-in form submission
  document.getElementById('loginForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // AJAX request for sign-in
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
          document.getElementById('signInModal').classList.remove('show');
          document.getElementById('user-content').style.display = 'block';
          document.getElementById('user-email').textContent = data.email;
        } else {
          alert('Invalid email or password');
        }
      })
      .catch(error => console.error('Error:', error));
  });

  // Handle sign-up form submission
  document.getElementById('signUpForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const name = document.getElementById('name').value;
    const email = document.getElementById('signUpEmail').value;
    const password = document.getElementById('signUpPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (password !== confirmPassword) {
      document.getElementById('passwordHelpBlock').innerText = "Passwords don't match!";
      return;
    }

    // AJAX request for sign-up
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
          document.getElementById('signUpModal').classList.remove('show');
          alert('Sign-up successful. Please sign in.');
        } else {
          alert('Error during sign-up: ' + data.message);
        }
      })
      .catch(error => console.error('Error:', error));
  });
});
