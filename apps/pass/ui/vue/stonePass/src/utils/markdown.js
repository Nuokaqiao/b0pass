import { marked } from 'marked'
import DOMPurify from 'dompurify'

marked.setOptions({
  gfm: true,
  breaks: true,
})

const PURIFY = {
  USE_PROFILES: { html: true },
  ALLOWED_TAGS: [
    'p',
    'br',
    'strong',
    'b',
    'em',
    'i',
    'del',
    's',
    'code',
    'pre',
    'a',
    'ul',
    'ol',
    'li',
    'blockquote',
    'hr',
    'h1',
    'h2',
    'h3',
    'h4',
  ],
  ALLOWED_ATTR: ['href', 'title', 'target', 'rel'],
  ALLOW_DATA_ATTR: false,
}

/** 简易 Markdown → 安全 HTML */
export function renderMarkdown(src) {
  const text = String(src || '')
  if (!text.trim()) return ''
  const dirty = marked.parse(text, { async: false })
  const clean = DOMPurify.sanitize(dirty, PURIFY)
  // 外链统一新窗口打开
  return clean.replace(/<a /g, '<a target="_blank" rel="noopener noreferrer" ')
}
