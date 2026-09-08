async function request(url, options = {}) {
  const res = await fetch(url, options)
  const data = await res.json().catch(() => null)
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

export function uploadFile(dirPath, file, onProgress, expireUnix = 0) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const form = new FormData()
    form.append('file', file)
    const target = dirPath || '/'
    let url = `/pass/file-upload?f=${encodeURIComponent(target)}`
    if (expireUnix > 0) {
      url += `&expire=${expireUnix}`
    }
    xhr.open('POST', url)
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    }
    xhr.onload = () => {
      try {
        const data = JSON.parse(xhr.responseText || '{}')
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

export function downloadUrl(path) {
  return `/pass/file-download?f=${encodeURIComponent(path)}`
}

export function staticFileUrl(path) {
  return `/files${path}`
}
