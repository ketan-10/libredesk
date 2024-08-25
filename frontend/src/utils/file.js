export function formatBytes(bytes) {
  if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(2) + ' KB'
  } else {
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
  }
}

export function getThumbFilepath(filepath) {
  const urlParts = filepath.split('/')
  const filename = urlParts.pop()
  const thumbFilename = `thumb_${filename}`
  urlParts.push(thumbFilename)
  return urlParts.join('/')
}
