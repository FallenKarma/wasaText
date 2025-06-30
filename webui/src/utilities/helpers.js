export function getFullPhotoUrl(relativePath) {
  const backendBaseUrl = import.meta.env.VITE_PHOTO_SERVER_URL
  if (relativePath && !relativePath.startsWith('/')) {
    relativePath = '/' + relativePath
  }
  return `${backendBaseUrl}${relativePath}`
}
