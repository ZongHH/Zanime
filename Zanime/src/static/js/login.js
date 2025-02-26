function showError(message) {
    Swal.fire({
        icon: 'error',
        title: 'Oops...',
        text: message,
        customClass: {
            popup: 'swal2-modal swal2-show swal2-round-modal',
            title: 'swal2-title swal2-round-title',
            content: 'swal2-content swal2-round-content',
            confirmButton: 'swal2-confirm swal2-round-confirm-button'
        },
        buttonsStyling: false,
        confirmButtonText: 'OK',
        confirmButtonClass: 'swal2-confirm swal2-styled swal2-round-confirm-button swal2-danger'
    });
}

export function login(stroe, email, password) {
    // 创建一个新的 FormData 实例
    const formData = new FormData();

    // 将 email 和 password 添加到 FormData 中
    formData.append('email', email);
    formData.append('password', password);
    fetch('/api/loginInfo', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // window.location.href = '/';
            // this.reset();
            localStorage.setItem('user_name', data.message.user_name)
            localStorage.setItem('user_id', data.message.user_id)
            localStorage.setItem('email', data.message.email)
            localStorage.setItem('gender', data.message.gender)
            localStorage.setItem('avataru_url', data.message.avataru_url)

            //在这一块映射connectWebSocket过来进行websocket连接
            stroe.dispatch('websocket/connectWebSocket', data.pushUrl); // 调用 Vuex action 来连接 WebSocket
        } else {
            showError(data.message)
        }

    })
    .catch(error => {
        showError(error);
    });
}