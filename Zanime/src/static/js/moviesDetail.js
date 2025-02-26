// 假设这是与后端获取评论数据的API端点相关的全局变量
const apiUrl = '/api/comments';

// 当前页码，初始化为1
let currentPage = 1;

let isLast = false;

// 从当前导航栏URL中获取videoId值的函数
function getVideoIdFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('videoId');
}
// 获取当前导航栏URL中的videoId值
const videoId = getVideoIdFromUrl();

// 创建一个全局Map，用于存储与commentId相关的信息
const commentIdMap = new Map();
const commentNum = new Map();

// 加载评论并展示的函数
async function loadComments() {
    const commentContainer = document.getElementById('comment-container');

    // 显示加载指示器
    const loadingIndicator = document.getElementById('loading-indicator');
    loadingIndicator.style.display = 'block';

    try {
        // 向后端发送请求获取下一批评论数据
        const response = await fetch(`${apiUrl}?videoId=${videoId}&currentPage=${currentPage}`);
        const newComments = await response.json();

        if (!newComments || newComments.length == 0) {
            updatePaginationButtonsDisabledState(0);
            isLast = true;
            return;
        }

        // 清空评论容器，用于重新加载当前页的评论
        commentContainer.innerHTML = '';

        // 将新获取的评论添加到评论容器中
        newComments.forEach((comment) => {
            createCommentsDiv(commentContainer, comment)
        });

        // 隐藏加载指示器
        loadingIndicator.style.display = 'none';

        if (newComments.length < 20) {
            isLast = true;
        }

        // 更新分页按钮的禁用状态以及当前页码显示
        updatePaginationButtonsDisabledState(newComments.length);
        updateCurrentPageNumDisplay();
    } catch (error) {
        console.error('加载评论时出错:', error);
        // 隐藏加载指示器
        loadingIndicator.style.display = 'none';
    }
}

function createCommentsDiv(commentContainer, comment) {
    const commentElement = createCommentElement(comment);
    commentContainer.appendChild(commentElement);

    // 创建一个外层div来包裹获取回复评论的按钮和展示回复评论的容器
    const replyWrapper = document.createElement('div');
    replyWrapper.classList.add('reply-wrapper');
    commentContainer.appendChild(replyWrapper);

    // 创建用于展示回复评论的容器
    const replyContainer = document.createElement('div');
    replyContainer.id = `reply-container-${comment.comment_id}`;
    replyContainer.classList.add('reply-container');
    replyWrapper.appendChild(replyContainer);

    // 添加间隔样式，这里使用 margin-bottom 来设置每条评论底部的间隔距离
    const marginBottom = '20px'; // 你可以根据需要调整这个值
    
    if (comment.reply_num > 0) {
        // 添加获取回复评论的按钮
        const getRepliesButton = document.createElement('button');
        getRepliesButton.textContent = '展开 '+comment.reply_num+' 回复';
        getRepliesButton.classList.add('get-replies-btn');
        getRepliesButton.dataset.commentId = comment.comment_id;
        replyWrapper.appendChild(getRepliesButton);
        
        // 初始化commentIdMap中对应评论ID的值，这里假设初始化为0，可以根据实际需求修改
        const commentIdStr = comment.comment_id.toString();
        commentIdMap.set(commentIdStr, 0);
        commentNum.set(commentIdStr, comment.reply_num)
    }
    
    replyWrapper.style.display = 'block';
    replyWrapper.style.marginBottom = marginBottom;
}

// 根据评论数据创建评论元素的函数
function createCommentElement(comment) {
    const commentCard = document.createElement('div');
    if (comment.parent_comment_id.Valid) {
        commentCard.classList.add('review-card-2');
    } else {
        commentCard.classList.add('review-card');
    }
    
    const img = document.createElement('img');
    img.src = comment.avataru_url;
    img.alt = 'pic';
    const imgDiv = document.createElement('div');
    imgDiv.appendChild(img);
    
    const textBlock = document.createElement('div');
    textBlock.classList.add('textBlock');
    const dFlexDiv = document.createElement('div');
    dFlexDiv.classList.add('d-flex', 'align-items-end', 'justify-content-between', 'mb-16');
    
    // 创建一个新的div来包裹span和h6
    const infoDiv = document.createElement('div');
    infoDiv.classList.add('d-flex', 'align-items-end', 'justify-content-between')
    const span = document.createElement('span');
    span.classList.add('subtitle', 'light-gray');
    span.textContent = formatRelativeTime(comment.created_at);
    
    const h6 = document.createElement('h6');
    h6.classList.add('fw-500', 'light-gray');
    h6.textContent = comment.user_name;
    
    const infoDiv_FirstChild = document.createElement('div');
    // 将span和h6添加到新创建的div中
    infoDiv_FirstChild.appendChild(span);
    infoDiv_FirstChild.appendChild(h6);

    infoDiv.appendChild(infoDiv_FirstChild);

    if (comment.user_id === parseInt(localStorage.getItem('user_id'), 10)) {
        // 创建代表“我”的图标元素
        const meIcon = document.createElement('img');
        meIcon.src ='src/static/icons/me_2.png';
        meIcon.alt = '我';
        meIcon.classList.add('me-icon'); // 可以添加自定义的类名以便后续进行样式设置等操作

        // 将图标元素添加到合适的位置，这里假设添加到infoDiv中，紧挨着h6元素之后
        infoDiv.appendChild(meIcon);
    }

    // 检查是否有回复人名字，如果有则将评论者名字和回复人名字拼接展示
    if (comment.replied_name) {
        const h6_ = document.createElement('h6');
        h6_.classList.add('fw-500', 'replier');
        h6_.textContent = 'Reply:';
        infoDiv.appendChild(h6_)

        const h6_replier = document.createElement('h6');
        h6_replier.classList.add('fw-500', 'light-gray');
        h6_replier.textContent = comment.replied_name;
        infoDiv.appendChild(h6_replier)
        if (comment.replied_name == localStorage.getItem('user_name')) {
            // 创建代表“我”的图标元素
            const meIcon = document.createElement('img');
            meIcon.src ='src/static/icons/me_2.png';
            meIcon.alt = '我';
            meIcon.classList.add('me-icon'); // 可以添加自定义的类名以便后续进行样式设置等操作

            // 将图标元素添加到合适的位置，这里假设添加到infoDiv中，紧挨着h6元素之后
            infoDiv.appendChild(meIcon);
        }
    }
    
    const replyBtn = document.createElement('button');
    replyBtn.classList.add('reply-btn', 'fw-500');
    replyBtn.textContent = 'Reply';
    replyBtn.dataset.userName = comment.user_name
    replyBtn.dataset.commentId = comment.comment_id
    replyBtn.dataset.parentCommentId = comment.parent_comment_id.Int64
    replyBtn.dataset.userID = comment.user_id
    if (comment.parent_comment_id.Valid) {
        replyBtn.dataset.type = 2
    } else {
        replyBtn.dataset.type = 1
    }
    
    // 将包裹span和h6的div添加到dFlexDiv中
    dFlexDiv.appendChild(infoDiv)
    dFlexDiv.appendChild(replyBtn);
    
    const p = document.createElement('p');
    p.textContent = comment.content;
    p.classList.add('mb-24');
    
    textBlock.appendChild(dFlexDiv);
    textBlock.appendChild(p);
    
    commentCard.appendChild(imgDiv);
    commentCard.appendChild(textBlock);
    
    return commentCard;
}

function formatRelativeTime(timeString) {
    const now = new Date();
    const commentTime = new Date(timeString);

    const diffInMilliseconds = now - commentTime;

    const seconds = Math.floor(diffInMilliseconds / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    const months = Math.floor(days / 30.44); // 平均每月按30.44天计算
    const years = Math.floor(months / 12);

    if (years > 0) {
        return `${years}年前`;
    } else if (months > 0) {
        return `${months}个月前`;
    } else if (days > 0) {
        return `${days}天前`;
    } else if (hours > 0) {
        return `${hours}小时前`;
    } else if (minutes > 0) {
        return `${minutes}分钟前`;
    } else if (seconds > 0) {
        return `${seconds}秒前`;
    } else {
        return `刚刚`;
    }
}

// 更新分页按钮的禁用状态的函数
function updatePaginationButtonsDisabledState(loadedCommentsCount) {
    const prevButton = document.getElementById('prevButton');
    const nextButton = document.getElementById('nextButton');

    prevButton.disabled = currentPage === 1;
    nextButton.disabled = loadedCommentsCount < 20;
}

// 更新当前页码显示的函数
function updateCurrentPageNumDisplay() {
    const currentPageNumSpan = document.getElementById('currentPageNum');
    currentPageNumSpan.textContent = currentPage;
}

function NowDataTime() {
    const currentTime = new Date();
    const year = currentTime.getFullYear();
    const month = String(currentTime.getMonth() + 1).padStart(2, '0');
    const day = String(currentTime.getDate()).padStart(2, '0');
    const hours = String(currentTime.getHours()).padStart(2, '0');
    const minutes = String(currentTime.getMinutes()).padStart(2, '0');
    const seconds = String(currentTime.getSeconds()).padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

// 处理按钮点击事件的函数
function handleButtonClicks(e) {
    if (e.target && e.target.classList.contains('get-replies-btn')) {
        const commentId = e.target.dataset.commentId;
        const button = e.target;

        // 根据按钮当前文本判断是展开还是折叠操作
        if (button.textContent !== '收起') {
            loadReplyComments(commentId);
        } else if (button.textContent === '收起') {
            commentIdMap.set(commentId, 0);
            collapseReplyComments(commentId);
        }
    } else if (e.target && e.target.classList.contains('post-btn')) {
        // 处理点击Post按钮的情况
        e.preventDefault();
        const commentInput = document.getElementById('message');
        const commentContent = commentInput.value;
        commentInput.value = '';
        postComment(commentContent)

    } else if (e.target && e.target.classList.contains('reply-btn')) {
        const commentElement = e.target.closest('.textBlock');
        let parentCommentId = e.target.dataset.commentId
        const repliedName = e.target.dataset.userName
        const repliedID = e.target.dataset.userID

        if (parseInt(e.target.dataset.type, 10) === 2) {
            parentCommentId = e.target.dataset.parentCommentId
        }

        // 创建包含输入框和发送按钮的div容器
        const replyDiv = document.createElement('div');
        replyDiv.classList.add('reply-comment-input');

        // 创建输入框元素
        const inputBox = document.createElement('input');
        inputBox.type = 'text';
        inputBox.placeholder = '请输入回复内容';
        inputBox.classList.add('reply-input');

        // 创建发送按钮元素
        const sendButton = document.createElement('button');
        sendButton.textContent = '发送';
        sendButton.type = 'button';
        sendButton.classList.add('send-btn');

        // 将输入框和发送按钮添加到新创建的div容器中
        replyDiv.appendChild(inputBox);
        replyDiv.appendChild(sendButton);

        // 将包含输入框和发送按钮的div容器添加到评论元素底部
        commentElement.appendChild(replyDiv);

        // 监听整个文档的点击事件
        document.addEventListener('click', function (e) {
            // 先判断replyDiv是否存在
            if (replyDiv) {
                // 判断点击事件的目标元素是否不在replyDiv及其子元素内
                if (!replyDiv.contains(e.target)) {
                    // 额外添加判断：检查replyDiv是否是commentElement的子节点
                    if (commentElement.contains(replyDiv)) {
                        // 如果是，则移除replyDiv
                        commentElement.removeChild(replyDiv);
                    }
                }
            }
        });

        // 给发送按钮添加点击事件处理函数
        sendButton.addEventListener('click', function () {
            const commentContent = inputBox.value;
            const parentComment = {
                "Int64": parseInt(parentCommentId, 10),
                "Valid": true,
            }
            let comment;
            if (parseInt(e.target.dataset.type, 10) === 1) {
                comment = sendComment(commentContent, repliedID, parentComment)
            } else {
                comment = sendComment(commentContent, repliedID, parentComment, repliedName)
            }
            // 发送成功后，清空输入框并移除输入框和发送按钮元素
            inputBox.value = '';
            commentElement.removeChild(replyDiv);

            createTempComment(comment, e.target.dataset.type, commentElement.closest('.review-card, .review-card-2'));

        });
    }
}

function createTempComment(comment, type, fatherElement) {
    const commentElement = createCommentElement(comment)
    if (type == 1) {
        const replyWrapper = fatherElement.nextElementSibling;
        if (replyWrapper) {
            replyWrapper.insertBefore(commentElement, replyWrapper.firstChild);
        }
    } else if (type == 2) {
        fatherElement.insertAdjacentElement('afterend', commentElement)
    } else if (type == 0) {
        // 获取评论容器元素
        const commentContainer = document.getElementById('comment-container');

        // 创建一个外层div来包裹获取回复评论的按钮和展示回复评论的容器
        const replyWrapper = document.createElement('div');
        replyWrapper.classList.add('reply-wrapper');
        commentContainer.insertBefore(replyWrapper, commentContainer.firstChild);

        // 创建用于展示回复评论的容器
        const replyContainer = document.createElement('div');
        replyContainer.id = `reply-container-${comment.comment_id}`;
        replyContainer.classList.add('reply-container');
        replyWrapper.appendChild(replyContainer);

        // 添加间隔样式，这里使用 margin-bottom 来设置每条评论底部的间隔距离
        const marginBottom = '20px'; // 你可以根据需要调整这个值
        
        if (comment.reply_num > 0) {
            // 添加获取回复评论的按钮
            const getRepliesButton = document.createElement('button');
            getRepliesButton.textContent = '展开 '+comment.reply_num+' 回复';
            getRepliesButton.classList.add('get-replies-btn');
            getRepliesButton.dataset.commentId = comment.comment_id;
            replyWrapper.appendChild(getRepliesButton);
            replyWrapper.style.display = 'block';
            
            // 初始化commentIdMap中对应评论ID的值，这里假设初始化为0，可以根据实际需求修改
            const commentIdStr = comment.comment_id.toString();
            commentIdMap.set(commentIdStr, 0);
            commentNum.set(commentIdStr, comment.reply_num)
        }
        replyWrapper.style.marginBottom = marginBottom;

        // 将新创建的包含评论及回复相关元素的commentElement插入到评论容器的第一位
        commentContainer.insertBefore(commentElement, commentContainer.firstChild);
    }
}

function sendComment(commentContent, repliedID, parentCommentId, repliedName) {
    // 如果没有接收到repliedName参数，将其默认为空字符串
    if (typeof repliedName === 'undefined') {
        repliedName = '';
    }
    if (typeof parentCommentId === 'undefined') {
        parentCommentId = {
            "Int64": 0,
            "Valid": false,
        };
    }
    if (typeof repliedID === 'undefined') {
        repliedID = 0;
    }
    if (commentContent) {
        const createdAt = NowDataTime();
        const avataruUrl = localStorage.getItem('avataru_url')
        const userName = localStorage.getItem('user_name')
        const userID = localStorage.getItem('user_id')
        // 构造要发送到后端的数据对象
        const comment = {
            "comment_id": virtual_id--,
            "video_id": parseInt(videoId, 10),
            "user_name": userName,
            "content": commentContent,
            "user_id": parseInt(userID, 10),
            "reply_num": 0,
            "parent_comment_id": parentCommentId,
            "created_at": createdAt,
            // "user_id": "123", 后台可以从cookies里取到
            "avataru_url": avataruUrl,
            "replied_name": repliedName,
            "replied_id": parseInt(repliedID, 10),
        };

        // 向Go后台发送POST请求
        fetch('/api/submit-comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(comment)
        })
        .then(response => response.json())
        .then(data => {
            if (!data.success) {
                showError(data.message)
            }
        })
        .catch(error => {
            showError(error);
        });

        return comment;
    }
    return null
}

//虚拟commentID
let virtual_id = -1

async function postComment(commentContent) {
    try {
        const comment = sendComment(commentContent)
        console.log(comment)
        if (!comment) return;
        createTempComment(comment, 0)
    } catch (error) {
        console.error('提交评论时出错:', error);
    }
}

// 加载指定评论的回复评论函数
async function loadReplyComments(commentId) {
    const replyContainer = document.getElementById(`reply-container-${commentId}`);

    try {
        commentIdMap.set(commentId, (commentIdMap.get(commentId) || 0) + 1);

        const response = await fetch(`${apiUrl}?videoId=${videoId}&parent_comment_id=${commentId}&currentPage=${commentIdMap.get(commentId)}`);
        const replyComments = await response.json();

        // 将新获取的回复评论添加到回复评论容器中
        replyComments.forEach((replyComment) => {
            const replyCommentElement = createCommentElement(replyComment);
            replyContainer.appendChild(replyCommentElement);
        });

        // 如果获取到的回复评论数量小于限制数量，说明已经获取完所有回复评论，修改按钮文本为折叠评论
        if (commentIdMap.get(commentId)*5 >= commentNum.get(commentId)) {
            const getRepliesButton = document.querySelector(`.get-replies-btn[data-comment-id="${commentId}"]`);
            getRepliesButton.textContent = '收起';
        } else {
            const getRepliesButton = document.querySelector(`.get-replies-btn[data-comment-id="${commentId}"]`);
            getRepliesButton.textContent = '展开更多回复';
        }
    } catch (error) {
        console.error('加载回复评论时出错:', error);
    }
}

// 折叠指定评论的单独回复评论函数
function collapseReplyComments(commentId) {
    const replyContainer = document.getElementById(`reply-container-${commentId}`);
    replyContainer.innerHTML = '';
    const getRepliesButton = document.querySelector(`.get-replies-btn[data-comment-id="${commentId}"]`);
    getRepliesButton.textContent = '展开 '+commentNum.get(commentId)+' 回复';
}

function showError(message) {
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

export default async function Movies() {
    // 给上一页按钮添加点击事件处理函数
    document.getElementById('prevButton').addEventListener('click', function () {
        if (currentPage > 1) {
            currentPage--;
            isLast = false;
            loadComments();
            // 页面滚动到指定位置
            const targetElement = document.querySelector('.add-comment');
            targetElement.scrollIntoView({ behavior: 'smooth' });
        }
    });

    // 给下一页按钮添加点击事件处理函数
    document.getElementById('nextButton').addEventListener('click', function () {
        if (isLast == false) {
            currentPage++;
            loadComments();
            // 页面滚动到指定位置
            const targetElement = document.querySelector('.add-comment');
            targetElement.scrollIntoView({ behavior: 'smooth' });
        }
    });
    // 页面加载时绑定事件委托，处理获取回复评论按钮的点击事件
    window.addEventListener('load', function () {
        document.addEventListener('click', handleButtonClicks);
    });
    loadComments()
}