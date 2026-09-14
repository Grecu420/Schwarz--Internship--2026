export const uploadImageToCloudinary = async (file: File, presetName: string): Promise<string> => {
  const cloudName = import.meta.env.VITE_CLOUDINARY_CLOUD_NAME
  
  if (!cloudName) {
    throw new Error("Cloudinary cloud name is missing in .env")
  }

  const formData = new FormData()
  formData.append('file', file)
  formData.append('upload_preset', presetName) 

  try {
    const response = await fetch(
      `https://api.cloudinary.com/v1_1/${cloudName}/image/upload`,
      {
        method: 'POST',
        body: formData,
      }
    )

    if (!response.ok) {
      throw new Error('Failed to upload image to Cloudinary')
    }

    const data = await response.json()
    return data.secure_url 
    
  } catch (error) {
    console.error('Cloudinary upload error:', error)
    throw error
  }
}