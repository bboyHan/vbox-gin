
// 设置缓存控制头部
export const setCacheControl = (seconds) => { // 设置缓存时间为seconds秒
    const headers = new Headers({
        'Cache-Control': `max-age=${seconds}`
    })
    document.defaultView.fetch('/', { headers })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error, status = ${response.status}`);
            }
            // 处理正常的响应
            return response.json();
        })
        .then(data => {
            // 处理返回的数据
        })
        .catch(error => {
            // 处理网络错误或 HTTP 错误
            console.error(error);
        });
}