import { authHeaders, handleUnauthorized, withTokenQuery } from './auth'

async function request(url, options = {}) {
  const headers = {
    ...(options.headers || {}),
    ...authHeaders(),
  }
  const res = await fetch(url, { ...options, headers })
  const data = await res.json().catch(() => null)
  if (data && data.code === 401) {
    handleUnauthorized()
    throw new Error(data.msg || '请先登录')
  }
  if (!res.ok) {
    throw new Error((data && data.msg) || `请求失败 (${res.status})`)
  }
  if (data && typeof data.code !== 'undefined' && data.code !== 0 && data.code !== 200) {
    throw new Error(data.msg || '请求失败')
  }
  return data
}

export function fetchFileList(path = '/') {
  const f = path || '/'
  // 防浏览器缓存导致删除/上传后列表不更新
  return request(`/pass/file-list?f=${encodeURIComponent(f)}&_t=${Date.now()}`)
}

export function deleteNode(path) {
  return request(`/pass/node-delete?f=${encodeURIComponent(path)}&_t=${Date.now()}`)
}

export function addNode(path) {
  return request(`/pass/node-add?f=${encodeURIComponent(path)}`)
}

export function uploadFile(dirPath, file, onProgress, expireUnix = 0, saveAsName = '') {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const form = new FormData()
    const name = saveAsName || file.name
    form.append('file', file, name)
    const target = dirPath || '/'
    let url = `/pass/file-upload?f=${encodeURIComponent(target)}`
    if (expireUnix > 0) {
      url += `&expire=${expireUnix}`
    }
    xhr.open('POST', url)
    const headers = authHeaders()
    Object.keys(headers).forEach((k) => xhr.setRequestHeader(k, headers[k]))
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }
    xhr.onload = () => {
      try {
        const data = JSON.parse(xhr.responseText || '{}')
        if (data && data.code === 401) {
          handleUnauthorized()
          reject(new Error(data.msg || '请先登录'))
          return
        }
        if (xhr.status >= 200 && xhr.status < 300 && data && data.code === 0) {
          resolve(data)
        } else {
          reject(new Error((data && data.msg) || '上传失败'))
        }
      } catch (err) {
        reject(err)
      }
    }
    xhr.onerror = () => reject(new Error('上传失败'))
    xhr.send(form)
  })
}

export function setFileExpire(path, expireUnix = 0) {
  return request(
    `/pass/file-expire?f=${encodeURIComponent(path)}&expire=${expireUnix || 0}&_t=${Date.now()}`,
  )
}

export function fetchTextHistory() {
  return request(`/pass/text-history?_t=${Date.now()}`)
}

export function downloadUrl(path) {
  return withTokenQuery(`/pass/file-download?f=${encodeURIComponent(path)}`)
}

export function staticFileUrl(path) {
  return withTokenQuery(`/files${path}`)
}
