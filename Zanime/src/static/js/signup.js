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

// 获取按钮元素
const getVerificationCodeButton = document.getElementById('getVerificationCodeButton');
// 为按钮添加点击事件监听器
getVerificationCodeButton.addEventListener('click', function () {
    const email = document.getElementById('eMail').value;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (email && emailRegex.test(email)) {
        // 发送请求到后台的函数，使用 fetch API 发送 POST 请求
        sendVerificationCodeRequest(email);
        let countdown = 60;
        let timer;
        this.disabled = true;
        this.textContent = `${countdown}s`;
        timer = setInterval(() => {
            countdown--;
            if (countdown <= 0) {
                clearInterval(timer);
                this.disabled = false;
                this.textContent = 'Get Code';
            } else {
                this.textContent = `${countdown}s`;
            }
        }, 1000);
    } else {
        showError('请输入正确的邮箱地址')
    }
});

function sendVerificationCodeRequest(email) {
    fetch('/api/sendVerificationCode', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email: email })
    })
  .then(response => response.json())
  .then(data => {
        if (data.success) {
        } else {
            showError('验证码发送失败，请稍后重试')
        }
    })
  .catch(error => {
        console.error('Error:', error);
        showError('发生错误，请检查网络或联系管理员')
    });
}

document.getElementById('signupForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = new FormData(this);
    fetch('/api/signInfo', {
        method: 'POST',
        body: formData
    })
   .then(response => response.json())
   .then(data => {
        if (data.success) {
            window.location.href = '/login';
        } else {
            showError(data.message);
        }
    })
   .catch(error => {
        showError(error);
    });
});
