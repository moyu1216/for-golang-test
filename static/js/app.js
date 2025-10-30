const API_BASE = '/api/tasks';

let tasks = [];
let selectedTasks = [];

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', () => {
    loadTasks();
    loadStats();
});

// 加载任务列表
async function loadTasks() {
    try {
        const status = document.getElementById('statusFilter').value;
        const priority = document.getElementById('priorityFilter').value;
        const sortBy = document.getElementById('sortBy').value;

        let url = API_BASE;
        const params = new URLSearchParams();
        if (status) params.append('status', status);
        if (priority) params.append('priority', priority);
        if (sortBy) params.append('sort_by', sortBy);
        if (params.toString()) url += '?' + params.toString();

        const response = await fetch(url);
        if (!response.ok) throw new Error('加载任务失败');
        
        tasks = await response.json();
        renderTasks();
        selectedTasks = [];
    } catch (error) {
        alert('加载任务失败: ' + error.message);
    }
}

// 渲染任务列表
function renderTasks() {
    const container = document.getElementById('tasksList');
    if (tasks.length === 0) {
        container.innerHTML = '<div class="empty-state"><h3>暂无任务</h3><p>点击"新建任务"按钮开始创建</p></div>';
        return;
    }

    container.innerHTML = tasks.map(task => `
        <div class="task-item">
            <input type="checkbox" value="${task.id}" onchange="toggleTaskSelection(${task.id})">
            <div class="task-content">
                <div class="task-title">${escapeHtml(task.title)}</div>
                ${task.description ? `<div class="task-description">${escapeHtml(task.description)}</div>` : ''}
                <div class="task-meta">
                    <span class="task-badge badge-status-${task.status}">${getStatusText(task.status)}</span>
                    <span class="task-badge badge-priority-${task.priority}">${getPriorityText(task.priority)}</span>
                    <span>创建: ${formatDate(task.created_at)}</span>
                    ${task.due_date ? `<span>截止: ${formatDate(task.due_date)}</span>` : ''}
                </div>
            </div>
            <div class="task-actions">
                <button onclick="editTask(${task.id})" class="btn-secondary">编辑</button>
                <button onclick="deleteTask(${task.id})" class="btn-danger">删除</button>
            </div>
        </div>
    `).join('');
}

// 加载统计信息
async function loadStats() {
    try {
        const response = await fetch(API_BASE + '/stats');
        if (!response.ok) throw new Error('加载统计失败');
        
        const stats = await response.json();
        const statsContainer = document.getElementById('stats');
        statsContainer.innerHTML = Object.entries(stats).map(([status, count]) => `
            <div class="stat-item">
                <strong>${count}</strong>
                <div>${getStatusText(status)}</div>
            </div>
        `).join('');
    } catch (error) {
        console.error('加载统计失败:', error);
    }
}

// 搜索任务
async function searchTasks() {
    const query = document.getElementById('searchInput').value.trim();
    if (!query) {
        loadTasks();
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/search?q=${encodeURIComponent(query)}`);
        if (!response.ok) throw new Error('搜索失败');
        
        tasks = await response.json();
        renderTasks();
    } catch (error) {
        alert('搜索失败: ' + error.message);
    }
}

// 显示创建表单
function showCreateForm() {
    document.getElementById('formTitle').textContent = '新建任务';
    document.getElementById('taskForm').style.display = 'block';
    document.getElementById('taskFormElement').reset();
    document.getElementById('taskId').value = '';
    document.getElementById('taskStatus').value = 'pending';
    document.getElementById('taskPriority').value = 'medium';
}

// 编辑任务
async function editTask(id) {
    try {
        const response = await fetch(`${API_BASE}/${id}`);
        if (!response.ok) throw new Error('加载任务失败');
        
        const task = await response.json();
        document.getElementById('formTitle').textContent = '编辑任务';
        document.getElementById('taskId').value = task.id;
        document.getElementById('taskTitle').value = task.title;
        document.getElementById('taskDescription').value = task.description || '';
        document.getElementById('taskStatus').value = task.status;
        document.getElementById('taskPriority').value = task.priority;
        if (task.due_date) {
            const dueDate = new Date(task.due_date);
            document.getElementById('taskDueDate').value = formatDateTimeLocal(dueDate);
        } else {
            document.getElementById('taskDueDate').value = '';
        }
        document.getElementById('taskForm').style.display = 'block';
    } catch (error) {
        alert('加载任务失败: ' + error.message);
    }
}

// 保存任务
async function saveTask(event) {
    event.preventDefault();
    
    const id = document.getElementById('taskId').value;
    const task = {
        title: document.getElementById('taskTitle').value,
        description: document.getElementById('taskDescription').value,
        status: document.getElementById('taskStatus').value,
        priority: document.getElementById('taskPriority').value,
    };

    const dueDate = document.getElementById('taskDueDate').value;
    if (dueDate) {
        task.due_date = new Date(dueDate).toISOString();
    }

    try {
        const url = id ? `${API_BASE}/${id}` : API_BASE;
        const method = id ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(task)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || '保存失败');
        }

        hideForm();
        loadTasks();
        loadStats();
    } catch (error) {
        alert('保存失败: ' + error.message);
    }
}

// 删除任务
async function deleteTask(id) {
    if (!confirm('确定要删除这个任务吗？')) return;

    try {
        const response = await fetch(`${API_BASE}/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('删除失败');
        
        loadTasks();
        loadStats();
    } catch (error) {
        alert('删除失败: ' + error.message);
    }
}

// 批量删除
async function batchDelete() {
    if (selectedTasks.length === 0) {
        alert('请先选择要删除的任务');
        return;
    }

    if (!confirm(`确定要删除选中的 ${selectedTasks.length} 个任务吗？`)) return;

    try {
        const response = await fetch(`${API_BASE}/batch/delete`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: selectedTasks })
        });

        if (!response.ok) throw new Error('批量删除失败');
        
        loadTasks();
        loadStats();
    } catch (error) {
        alert('批量删除失败: ' + error.message);
    }
}

// 批量更新状态
async function batchUpdateStatus() {
    if (selectedTasks.length === 0) {
        alert('请先选择要更新的任务');
        return;
    }

    const status = prompt('请输入新状态 (pending/in_progress/completed):');
    if (!status) return;

    try {
        const response = await fetch(`${API_BASE}/batch/status`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: selectedTasks, status })
        });

        if (!response.ok) throw new Error('批量更新失败');
        
        loadTasks();
        loadStats();
    } catch (error) {
        alert('批量更新失败: ' + error.message);
    }
}

// 切换任务选择
function toggleTaskSelection(id) {
    const index = selectedTasks.indexOf(id);
    if (index > -1) {
        selectedTasks.splice(index, 1);
    } else {
        selectedTasks.push(id);
    }
}

// 隐藏表单
function hideForm() {
    document.getElementById('taskForm').style.display = 'none';
}

// 工具函数
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN');
}

function formatDateTimeLocal(date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
}

function getStatusText(status) {
    const map = {
        'pending': '待处理',
        'in_progress': '进行中',
        'completed': '已完成'
    };
    return map[status] || status;
}

function getPriorityText(priority) {
    const map = {
        'low': '低',
        'medium': '中',
        'high': '高'
    };
    return map[priority] || priority;
}