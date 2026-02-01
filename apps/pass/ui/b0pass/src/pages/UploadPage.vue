<template>
  <div class="upload-shell">
    <header class="upload-header">
      <div class="logo-mark">TeamTransfer</div>
      <a class="text-link" href="/">← 回到首页</a>
    </header>
    <main class="upload-main">
      <div class="upload-card">
        <h2>上传文件</h2>
        <p class="upload-legend">{{ uploadLegend }}</p>
        <fieldset class="upload-fieldset">
          <div class="upload-layout">
            <div
              class="upload-drag-area"
              :class="{ 'is-active': dragActive }"
              @click="triggerSelect"
              @dragenter.prevent="handleDragOver"
              @dragover.prevent="handleDragOver"
              @dragleave.prevent="handleDragLeave"
              @drop="handleDrop"
            >
              <input
                ref="fileInput"
                class="upload-hidden-input"
                type="file"
                multiple
                @change="handleFileSelect"
              />
              <div class="upload-drag-content">
                <p class="upload-drag-icon">⬆</p>
                <p class="upload-drag-text">点击或者拖拽文件到此区域</p>
                <p class="upload-drag-hint">支持大文件 · 多选 · 拖拽</p>
                <button type="button" class="outline-btn" :disabled="isUploading">
                  选择文件
                </button>
              </div>
            </div>

            <div class="upload-table-wrapper">
              <table class="upload-table">
                <colgroup>
                  <col>
                  <col width="90">
                  <col width="180">
                  <col width="160">
                </colgroup>
                <thead>
                  <tr>
                    <th>文件名</th>
                    <th>大小</th>
                    <th>进度</th>
                    <th>提示</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in files" :key="item.id">
                    <td>{{ item.file.name }}</td>
                    <td>{{ item.readableSize }}</td>
                    <td>
                      <div class="progress-holder">
                        <div class="progress-bar" :style="{ width: item.progress + '%' }"></div>
                      </div>
                      <span class="progress-value">{{ item.progress }}%</span>
                    </td>
                    <td>
                      <div class="status-wrap">
                        <span class="status-chip" :class="item.status">{{ getStatusLabel(item.status) }}</span>
                        <div class="status-actions">
                          <button
                            type="button"
                            class="tiny-btn"
                            :disabled="item.status === 'uploading'"
                            @click="removeFile(item.id)"
                          >
                            移除
                          </button>
                          <button
                            v-if="item.status === 'error'"
                            type="button"
                            class="tiny-btn warning"
                            @click="reuploadFile(item)"
                          >
                            重传
                          </button>
                        </div>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="!files.length">
                    <td colspan="4" class="empty-state">还没有选择要上传的文件</td>
                  </tr>
                </tbody>
              </table>
              <div class="upload-footer">
                <button
                  type="button"
                  class="primary-btn"
                  :disabled="isUploading || !files.length"
                  @click="startUpload"
                >
                  提交上传
                </button>
                <button type="button" class="ghost-btn" :disabled="isUploading || !files.length" @click="clearFiles">
                  清空列表
                </button>
                <p class="upload-feedback" v-if="uploadFeedback">{{ uploadFeedback }}</p>
              </div>
              <p class="upload-instructions">
                * 只需与 B0Pass 使用同一 WIFI · * 充分发挥局域网带宽 · * 多文件上传推荐使用桌面浏览器
              </p>
            </div>
          </div>
        </fieldset>
      </div>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';

const files = ref([]);
const dragActive = ref(false);
const fileInput = ref(null);
const isUploading = ref(false);
const uploadFeedback = ref('');
const targetPath = ref('/');
const baseUploadUrl = '/pass/file-upload';

const normalizePath = (value) => {
  let normalized = value ?? '/';
  normalized = normalized.replace(/\\/g, '/').replace(/\/{2,}/g, '/');
  if (!normalized.startsWith('/')) {
    normalized = `/${normalized}`;
  }
  return normalized;
};

const uploadLegend = computed(() => `上传到 ${targetPath.value} 目录`);

onMounted(() => {
  try {
    const url = new URL(window.location.href);
    const fParam = url.searchParams.get('f');
    targetPath.value = fParam ? normalizePath(decodeURI(fParam)) : '/';
  } catch (error) {
    targetPath.value = '/';
  }
});

const humanFileSize = (size) => {
  if (size >= 1024 * 1024 * 1024) {
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)}G`;
  }
  if (size >= 1024 * 1024) {
    return `${(size / (1024 * 1024)).toFixed(1)}M`;
  }
  if (size >= 1024) {
    return `${(size / 1024).toFixed(1)}K`;
  }
  return `${size}B`;
};

const getStatusLabel = (status) => {
  const map = {
    ready: '等待上传',
    uploading: '上传中',
    done: '上传完成',
    error: '上传失败'
  };
  return map[status] ?? '未知状态';
};

const addFiles = (incoming) => {
  const source = Array.from(incoming ?? []);
  if (!source.length) {
    return;
  }
  const extra = source.map((file) => ({
    id: `${file.name}-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    file,
    readableSize: humanFileSize(file.size),
    progress: 0,
    status: 'ready'
  }));
  files.value = [...files.value, ...extra];
  uploadFeedback.value = '';
};

const handleFileSelect = (event) => {
  addFiles(event.target.files);
  if (fileInput.value) {
    fileInput.value.value = '';
  }
};

const triggerSelect = () => {
  fileInput.value?.click();
};

const handleDragOver = (event) => {
  event.preventDefault();
  dragActive.value = true;
};

const handleDragLeave = () => {
  dragActive.value = false;
};

const handleDrop = (event) => {
  event.preventDefault();
  dragActive.value = false;
  addFiles(event.dataTransfer?.files);
};

const removeFile = (id) => {
  files.value = files.value.filter((item) => item.id !== id);
};

const clearFiles = () => {
  files.value = [];
  uploadFeedback.value = '';
};

const uploadFileItem = (item) =>
  new Promise((resolve, reject) => {
    const form = new FormData();
    form.append('file', item.file);
    const xhr = new XMLHttpRequest();
    const uploadTarget = `${baseUploadUrl}?f=${encodeURIComponent(targetPath.value)}`;

    xhr.open('POST', uploadTarget);
    xhr.upload.onprogress = (progressEvent) => {
      if (progressEvent.lengthComputable) {
        item.progress = Math.round((progressEvent.loaded / progressEvent.total) * 100);
      }
    };
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 400) {
        item.progress = 100;
        item.status = 'done';
        resolve(xhr.response);
      } else {
        item.status = 'error';
        reject(new Error('上传失败'));
      }
    };
    xhr.onerror = () => {
      item.status = 'error';
      reject(new Error('上传失败'));
    };
    xhr.send(form);
  });

const startUpload = async () => {
  if (isUploading.value || !files.value.length) {
    return;
  }
  uploadFeedback.value = '';
  isUploading.value = true;
  for (const item of files.value) {
    if (item.status === 'done') {
      continue;
    }
    item.status = 'uploading';
    try {
      await uploadFileItem(item);
    } catch (error) {
      uploadFeedback.value = '部分文件上传失败，可重试';
    }
  }
  if (files.value.every((item) => item.status === 'done')) {
    uploadFeedback.value = '所有文件已上传';
  }
  isUploading.value = false;
};

const reuploadFile = async (targetItem) => {
  if (isUploading.value) {
    return;
  }
  uploadFeedback.value = '';
  targetItem.status = 'uploading';
  try {
    await uploadFileItem(targetItem);
    uploadFeedback.value = '文件重传完成';
  } catch (error) {
    uploadFeedback.value = '重传失败，请检查网络';
  }
};
</script setup>

<style scoped>
.upload-shell {
  min-height: 100vh;
  background: linear-gradient(180deg, #eef2ff 0%, #f9fbff 100%);
  display: flex;
  flex-direction: column;
  color: #0f172a;
}

.upload-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 3rem;
  background: #ffffff;
  box-shadow: 0 1px 30px rgba(15, 23, 42, 0.08);
}

.logo-mark {
  font-weight: 700;
  font-size: 1.25rem;
  letter-spacing: 0.05em;
}

.text-link {
  font-size: 0.95rem;
  color: #2563eb;
  text-decoration: none;
}

.upload-main {
  flex: 1;
  padding: 3rem;
  display: flex;
  justify-content: center;
}

.upload-card {
  width: min(1024px, 100%);
  background: #fff;
  border-radius: 1.5rem;
  padding: 2.5rem;
  box-shadow: 0 25px 40px rgba(15, 23, 42, 0.14);
  border: 1px solid rgba(37, 99, 235, 0.1);
}

.upload-card h2 {
  text-align: center;
  margin-bottom: 0.25rem;
  font-size: 1.75rem;
}

.upload-legend {
  text-align: center;
  margin: 0.5rem 0 1.5rem;
  color: #475569;
  font-size: 0.95rem;
}

.upload-fieldset {
  border: none;
  margin: 0;
  padding: 0;
}

.upload-layout {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.75rem;
}

.upload-drag-area {
  border: 2px dashed #cbd5f5;
  border-radius: 1.25rem;
  padding: 2.4rem 1.25rem;
  text-align: center;
  cursor: pointer;
  background: #f7f8fc;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.upload-drag-area.is-active {
  border-color: #2563eb;
  background: #e9f0ff;
}

.upload-drag-area:hover {
  border-color: #2563eb;
}

.upload-hidden-input {
  display: none;
}

.upload-drag-icon {
  font-size: 2rem;
  margin: 0;
}

.upload-drag-text {
  margin: 0.4rem 0;
  font-weight: 500;
}

.upload-drag-hint {
  margin: 0;
  font-size: 0.85rem;
  color: #475569;
}

.outline-btn {
  border: 1px solid #2563eb;
  background: transparent;
  border-radius: 999px;
  padding: 0.5rem 1.75rem;
  color: #2563eb;
  font-size: 0.9rem;
  margin-top: 0.75rem;
  cursor: pointer;
}

.outline-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.upload-table-wrapper {
  background: #fefefe;
  border-radius: 1.1rem;
  padding: 1rem;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.upload-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.upload-table th,
.upload-table td {
  padding: 0.75rem 0.5rem;
  text-align: left;
  border-bottom: 1px solid #e2e8f0;
}

.upload-table thead th {
  font-weight: 600;
  color: #1d1f3d;
}

.progress-holder {
  position: relative;
  height: 6px;
  background: #e2e8f0;
  border-radius: 999px;
  overflow: hidden;
  margin-bottom: 0.3rem;
}

.progress-bar {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, #2563eb, #4f46e5);
  width: 0%;
  transition: width 0.2s ease;
}

.progress-value {
  font-size: 0.75rem;
  color: #475569;
}

.status-wrap {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.status-chip {
  display: inline-flex;
  align-items: center;
  padding: 0.2rem 0.8rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 600;
  color: #1d1f3d;
  background: #e2e8f0;
  width: fit-content;
}

.status-chip.uploading {
  background: #e0f2ff;
  color: #0f62ff;
}

.status-chip.done {
  background: #ecfdf4;
  color: #059669;
}

.status-chip.error {
  background: #fee2e2;
  color: #b91c1c;
}

.status-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.tiny-btn {
  border: none;
  border-radius: 999px;
  padding: 0.2rem 0.95rem;
  font-size: 0.75rem;
  cursor: pointer;
  background: #e2e8f0;
  color: #1d1f3d;
}

.tiny-btn.warning {
  background: #fee2e2;
  color: #b91c1c;
}

.tiny-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.empty-state {
  text-align: center;
  color: #94a3b8;
  padding: 1.5rem 0;
}

.upload-footer {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  align-items: center;
  justify-content: flex-end;
}

.upload-footer .ghost-btn {
  min-width: 120px;
}

.upload-feedback {
  font-size: 0.85rem;
  color: #475569;
  margin: 0;
}

.upload-instructions {
  margin: 0;
  font-size: 0.85rem;
  color: #475569;
  text-align: center;
  line-height: 1.6;
}

@media (max-width: 720px) {
  .upload-header {
    padding: 1rem;
    flex-direction: column;
    gap: 0.5rem;
  }

  .upload-main {
    padding: 2rem 1rem;
  }

  .upload-card {
    padding: 1.5rem;
  }

  .upload-layout {
    gap: 1rem;
  }
}
</style>
