var app = new Vue({
    el: '#app',
    data: {
      users: []
    },
    mounted: function() {
      var xhr = new XMLHttpRequest();
      xhr.open('GET', 'http://localhost:8080/users');
      xhr.onload = function() {
        if (xhr.status === 200) {
          app.users = JSON.parse(xhr.responseText);
        } else {
          console.log('Ошибка загрузки данных');
        }
      };
      xhr.send();
    }
  });
  