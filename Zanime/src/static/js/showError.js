export default function showError(message) {
    Swal.fire({
        icon: 'error',
        title: 'Oops...',
        text: message,
        customClass: {
            popup: 'swal2-modal swal2-show swal2-round-modal', // 自定义弹窗样式
            title: 'swal2-title swal2-round-title',            // 自定义标题样式
            content: 'swal2-content swal2-round-content',      // 自定义内容样式
            confirmButton: 'swal2-confirm swal2-round-confirm-button swal2-danger' // 自定义确认按钮样式
        },
        buttonsStyling: false, // 禁用默认样式，启用自定义样式
        confirmButtonText: 'OK' // 确认按钮文字
    });
}
