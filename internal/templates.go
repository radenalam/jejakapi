package internal

// layoutTemplate adalah template HTML untuk layout utama
const layoutTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - JejakAPI</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#3b82f6',
                        secondary: '#1e40af',
                        success: '#10b981',
                        warning: '#f59e0b',
                        danger: '#ef4444',
                        dark: '#1f2937'
                    }
                }
            }
        }
    </script>
</head>
<body class="bg-gray-50 min-h-screen">
    <!-- Navigation -->
    <nav class="bg-dark shadow-lg">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex items-center">
                    <h1 class="text-xl font-bold text-white flex items-center">
                        <span class="text-2xl mr-2">🔍</span>
                        JejakAPI
                    </h1>
                </div>
                <div class="flex items-center space-x-4">
                    <span class="text-gray-300 text-sm">Request & Response Monitoring</span>
                </div>
            </div>
        </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {{template "content" .}}
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t mt-12">
        <div class="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8 text-center text-gray-500 text-sm">
            <p>JejakAPI - Built with Go & Fiber</p>
        </div>
    </footer>
</body>
</html>`

// dashboardTemplate adalah template HTML untuk halaman dashboard
const dashboardTemplate = `{{define "content"}}
<!-- Controls Section -->
<div class="bg-white rounded-lg shadow-sm border mb-6">
    <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">Monitoring Controls</h2>
    </div>
    <div class="px-6 py-4">
        <div class="flex flex-wrap gap-3">
            <button onclick="toggleEnabled()" 
                    class="bg-primary hover:bg-secondary text-white px-4 py-2 rounded-md transition-colors duration-200 flex items-center">
                <span class="mr-2">⚡</span>
                Toggle Monitoring
            </button>
            <button onclick="refreshLogs()" 
                    class="bg-success hover:bg-green-600 text-white px-4 py-2 rounded-md transition-colors duration-200 flex items-center">
                <span class="mr-2">🔄</span>
                Refresh
            </button>
            <button onclick="clearLogs()" 
                    class="bg-danger hover:bg-red-600 text-white px-4 py-2 rounded-md transition-colors duration-200 flex items-center">
                <span class="mr-2">🗑️</span>
                Clear All Logs
            </button>
        </div>
    </div>
</div>

<!-- Stats Cards -->
<div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
    <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
            <div class="p-2 bg-blue-100 rounded-lg">
                <span class="text-2xl">📊</span>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Total Requests</p>
                <p class="text-2xl font-semibold text-gray-900" id="total-requests">-</p>
            </div>
        </div>
    </div>
    <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
            <div class="p-2 bg-green-100 rounded-lg">
                <span class="text-2xl">✅</span>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Success Rate</p>
                <p class="text-2xl font-semibold text-gray-900" id="success-rate">-</p>
            </div>
        </div>
    </div>
    <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
            <div class="p-2 bg-yellow-100 rounded-lg">
                <span class="text-2xl">⚡</span>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Avg Response</p>
                <p class="text-2xl font-semibold text-gray-900" id="avg-response">-</p>
            </div>
        </div>
    </div>
    <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
            <div class="p-2 bg-red-100 rounded-lg">
                <span class="text-2xl">❌</span>
            </div>
            <div class="ml-4">
                <p class="text-sm font-medium text-gray-500">Error Rate</p>
                <p class="text-2xl font-semibold text-gray-900" id="error-rate">-</p>
            </div>
        </div>
    </div>
</div>

<!-- Logs Section -->
<div class="bg-white rounded-lg shadow-sm border">
    <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-900">Request Logs</h2>
    </div>
    <div class="p-6">
        <div id="logs-container">
            <div class="text-center py-12">
                <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                <p class="mt-4 text-gray-600">Loading logs...</p>
            </div>
        </div>
    </div>
</div>

<script>
// Dashboard JavaScript untuk JejakAPI
let currentPage = 1;
const logsPerPage = 20;

// Utility functions
function formatTimestamp(timestamp) {
    return new Date(timestamp).toLocaleString();
}

function formatDuration(microseconds) {
    if (microseconds < 1000) {
        return microseconds + 'μs';
    } else if (microseconds < 1000000) {
        return Math.round(microseconds / 1000) + 'ms';
    } else {
        return Math.round(microseconds / 1000000) + 's';
    }
}

function getStatusClass(status) {
    if (status >= 200 && status < 300) return 'success';
    if (status >= 400) return 'error';
    if (status >= 300) return 'warning';
    return '';
}

function getMethodColor(method) {
    const colors = {
        'GET': 'bg-green-500',
        'POST': 'bg-blue-500',
        'PUT': 'bg-yellow-500',
        'DELETE': 'bg-red-500',
        'PATCH': 'bg-purple-500'
    };
    return colors[method] || 'bg-gray-500';
}

function getStatusColor(status) {
    if (status >= 200 && status < 300) return 'bg-green-500';
    if (status >= 400) return 'bg-red-500';
    if (status >= 300) return 'bg-yellow-500';
    return 'bg-gray-500';
}

function formatJSON(str) {
    try {
        return JSON.stringify(JSON.parse(str), null, 2);
    } catch (e) {
        return str;
    }
}

function toggleLogDetails(id) {
    const details = document.getElementById('details-' + id);
    const arrow = document.getElementById('arrow-' + id);
    
    if (details.classList.contains('hidden')) {
        details.classList.remove('hidden');
        arrow.innerHTML = '▼';
    } else {
        details.classList.add('hidden');
        arrow.innerHTML = '▶';
    }
}

// Main functions
async function loadLogs(page = 1) {
    try {
        const response = await fetch('/jejakapi/api/logs?page=' + page + '&limit=' + logsPerPage);
        const data = await response.json();
        
        let html = '';
        
        if (data.logs && data.logs.length > 0) {
            // Update stats
            updateStats(data);
            
            data.logs.forEach(log => {
                const methodColor = getMethodColor(log.method);
                const statusColor = getStatusColor(log.status_code);
                
                html += '<div class="border border-gray-200 rounded-lg mb-4 overflow-hidden">';
                html += '<div class="bg-gray-50 px-4 py-3 cursor-pointer hover:bg-gray-100 transition-colors duration-200" onclick="toggleLogDetails(\'' + log.id + '\')">';
                html += '<div class="flex items-center justify-between">';
                html += '<div class="flex items-center space-x-3">';
                html += '<span id="arrow-' + log.id + '" class="text-gray-400 text-sm">▶</span>';
                html += '<span class="' + methodColor + ' text-white px-2 py-1 rounded text-xs font-semibold">' + log.method + '</span>';
                html += '<span class="text-gray-900 font-medium">' + log.url + '</span>';
                html += '</div>';
                html += '<div class="flex items-center space-x-3 text-sm">';
                html += '<span class="' + statusColor + ' text-white px-2 py-1 rounded text-xs font-semibold">' + log.status_code + '</span>';
                html += '<span class="text-gray-600">' + formatDuration(log.duration) + '</span>';
                html += '<span class="text-gray-500">' + formatTimestamp(log.created_at) + '</span>';
                html += '</div>';
                html += '</div>';
                html += '</div>';
                
                html += '<div id="details-' + log.id + '" class="hidden bg-white border-t border-gray-200">';
                html += '<div class="p-4 space-y-4">';
                
                // Request Headers
                html += '<div>';
                html += '<h4 class="font-semibold text-gray-900 mb-2">Request Headers</h4>';
                html += '<div class="bg-gray-50 rounded-md p-3 max-h-48 overflow-y-auto">';
                html += '<pre class="text-xs text-gray-700">' + formatJSON(JSON.stringify(log.headers)) + '</pre>';
                html += '</div>';
                html += '</div>';
                
                // Request Body
                if (log.body) {
                    html += '<div>';
                    html += '<h4 class="font-semibold text-gray-900 mb-2">Request Body</h4>';
                    html += '<div class="bg-gray-50 rounded-md p-3 max-h-48 overflow-y-auto">';
                    html += '<pre class="text-xs text-gray-700">' + formatJSON(log.body) + '</pre>';
                    html += '</div>';
                    html += '</div>';
                }
                
                // Response Headers
                html += '<div>';
                html += '<h4 class="font-semibold text-gray-900 mb-2">Response Headers</h4>';
                html += '<div class="bg-gray-50 rounded-md p-3 max-h-48 overflow-y-auto">';
                html += '<pre class="text-xs text-gray-700">' + formatJSON(JSON.stringify(log.response_headers)) + '</pre>';
                html += '</div>';
                html += '</div>';
                
                // Response Body
                if (log.response_body) {
                    html += '<div>';
                    html += '<h4 class="font-semibold text-gray-900 mb-2">Response Body</h4>';
                    html += '<div class="bg-gray-50 rounded-md p-3 max-h-48 overflow-y-auto">';
                    html += '<pre class="text-xs text-gray-700">' + formatJSON(log.response_body) + '</pre>';
                    html += '</div>';
                    html += '</div>';
                }
                
                // Client Info
                html += '<div>';
                html += '<h4 class="font-semibold text-gray-900 mb-2">Client Information</h4>';
                html += '<div class="bg-gray-50 rounded-md p-3">';
                html += '<div class="text-sm text-gray-700">';
                html += '<div><strong>IP:</strong> ' + log.ip + '</div>';
                html += '<div><strong>User Agent:</strong> ' + log.user_agent + '</div>';
                html += '</div>';
                html += '</div>';
                html += '</div>';
                
                html += '</div>';
                html += '</div>';
                html += '</div>';
            });
            
            // Add pagination
            const totalPages = Math.ceil(data.total / logsPerPage);
            if (totalPages > 1) {
                html += '<div class="flex justify-center mt-6 space-x-2">';
                for (let i = 1; i <= totalPages; i++) {
                    const activeClass = i === page ? 'bg-primary text-white' : 'bg-white text-gray-700 hover:bg-gray-50';
                    html += '<button onclick="loadLogs(' + i + ')" class="' + activeClass + ' px-3 py-2 border border-gray-300 rounded-md text-sm font-medium transition-colors duration-200">' + i + '</button>';
                }
                html += '</div>';
            }
        } else {
            html = '<div class="text-center py-12">';
            html += '<div class="text-gray-400 text-6xl mb-4">📭</div>';
            html += '<p class="text-gray-500 text-lg">No logs found</p>';
            html += '<p class="text-gray-400 text-sm mt-2">Logs will appear here when requests are made</p>';
            html += '</div>';
        }
        
        document.getElementById('logs-container').innerHTML = html;
        currentPage = page;
    } catch (error) {
        console.error('Error loading logs:', error);
        document.getElementById('logs-container').innerHTML = '<div class="text-center py-12"><div class="text-red-400 text-6xl mb-4">⚠️</div><p class="text-red-500 text-lg">Error loading logs</p><button onclick="loadLogs()" class="mt-4 bg-primary text-white px-4 py-2 rounded-md">Try Again</button></div>';
    }
}

function updateStats(data) {
    if (!data.logs) return;
    
    const totalRequests = data.total || 0;
    const logs = data.logs;
    
    let successCount = 0;
    let errorCount = 0;
    let totalDuration = 0;
    
    logs.forEach(log => {
        if (log.status_code >= 200 && log.status_code < 300) {
            successCount++;
        } else if (log.status_code >= 400) {
            errorCount++;
        }
        totalDuration += log.duration || 0;
    });
    
    const successRate = logs.length > 0 ? Math.round((successCount / logs.length) * 100) : 0;
    const errorRate = logs.length > 0 ? Math.round((errorCount / logs.length) * 100) : 0;
    const avgResponse = logs.length > 0 ? formatDuration(totalDuration / logs.length) : '0ms';
    
    document.getElementById('total-requests').textContent = totalRequests;
    document.getElementById('success-rate').textContent = successRate + '%';
    document.getElementById('error-rate').textContent = errorRate + '%';
    document.getElementById('avg-response').textContent = avgResponse;
}

async function clearLogs() {
    if (confirm('Are you sure you want to clear all logs?')) {
        try {
            await fetch('/jejakapi/api/logs', { method: 'DELETE' });
            loadLogs(1);
        } catch (error) {
            console.error('Error clearing logs:', error);
            alert('Error clearing logs');
        }
    }
}

async function toggleEnabled() {
    try {
        await fetch('/jejakapi/api/toggle', { method: 'POST' });
        showNotification('Monitoring toggled successfully', 'success');
    } catch (error) {
        console.error('Error toggling monitoring:', error);
        showNotification('Error toggling monitoring', 'error');
    }
}

function refreshLogs() {
    loadLogs(currentPage);
}

function showNotification(message, type = 'info') {
    const colors = {
        'success': 'bg-green-500',
        'error': 'bg-red-500',
        'info': 'bg-blue-500'
    };
    
    const notification = document.createElement('div');
    notification.className = 'fixed top-4 right-4 ' + colors[type] + ' text-white px-4 py-2 rounded-md shadow-lg z-50 transition-opacity duration-300';
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.opacity = '0';
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    loadLogs();
    
    // Auto refresh every 10 seconds
    setInterval(function() {
        if (currentPage === 1) {
            loadLogs(1);
        }
    }, 10000);
});
</script>
{{end}}`
