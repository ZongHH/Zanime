export default async function Logout() {
    document.getElementById('logoutLink').addEventListener('click', function (event) {
        event.preventDefault();
        fetch('/api/logout', {
            method: 'GET'
        })
        .then(() => {
            localStorage.clear();
            window.location.href = '/login';
        })
        .catch(error => {
            showError('Error during logout:' + error);
        });
    });
}