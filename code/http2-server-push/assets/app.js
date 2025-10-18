// HTTP/2 Server Push 示例 JavaScript
console.log('JavaScript 文件已加载（通过 HTTP/2 Server Push）');

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    console.log('页面 DOM 加载完成');

    // 检查是否通过 HTTP/2 加载
    if (performance && performance.getEntriesByType) {
        const resources = performance.getEntriesByType('resource');
        resources.forEach(resource => {
            if (resource.nextHopProtocol === 'h2') {
                console.log('通过 HTTP/2 加载:', resource.name);
            }
        });
    }

    // 显示加载时间
    if (performance && performance.timing) {
        const loadTime = performance.timing.loadEventEnd - performance.timing.navigationStart;
        console.log('页面加载时间:', loadTime, 'ms');
    }
});
