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
  return request(`/pass/file-list?f=${encodeURIComponent(f)}`)
}

export function deleteNode(path) {
  return request(`/pass/node-delete?f=${encodeURIComponent(path)}`)
}

export function addNode(path) {
  return request(`/pass/node-add?f=${encodeURIComponent(path)}`)
}

export function uploadFile(dirPath, file, onProgress) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    const form = new FormData()
    form.append('file', file)
    const target = dirPath || '/'
    xhr.open('POST', `/pass/file-upload?f=${encodeURIComponent(target)}`)
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

export function downloadUrl(path) {
  return `/pass/file-download?f=${encodeURIComponent(path)}`
}

export function staticFileUrl(path) {
  return `/files${path}`
}
