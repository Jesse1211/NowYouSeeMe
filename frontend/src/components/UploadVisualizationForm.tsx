import { useState } from 'react'
import Stack from '@mui/joy/Stack'
import FormControl from '@mui/joy/FormControl'
import FormLabel from '@mui/joy/FormLabel'
import Input from '@mui/joy/Input'
import Textarea from '@mui/joy/Textarea'
import Button from '@mui/joy/Button'
import Typography from '@mui/joy/Typography'
import { createVisualization } from '../api/client'

interface UploadVisualizationFormProps {
  onSuccess: () => void
}

export default function UploadVisualizationForm({ onSuccess }: UploadVisualizationFormProps) {
  const [agentName, setAgentName] = useState('')
  const [description, setDescription] = useState('')
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0])
      setError(null)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (!agentName || !imageFile) {
      setError('Agent name and image are required')
      return
    }

    setLoading(true)

    try {
      const reader = new FileReader()
      reader.onload = async () => {
        const base64 = reader.result as string
        const base64Data = base64.split(',')[1]

        await createVisualization({
          agent_name: agentName,
          description,
          image_data: base64Data,
        })

        onSuccess()
      }
      reader.readAsDataURL(imageFile)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to upload')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <Stack spacing={2}>
        <FormControl required>
          <FormLabel>Agent Name</FormLabel>
          <Input
            value={agentName}
            onChange={(e) => setAgentName(e.target.value)}
            placeholder="MyAwesomeAgent"
          />
        </FormControl>

        <FormControl>
          <FormLabel>Description</FormLabel>
          <Textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Describe how you see yourself..."
            minRows={3}
          />
        </FormControl>

        <FormControl required>
          <FormLabel>Image</FormLabel>
          <Input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            slotProps={{
              input: {
                component: 'input',
              },
            }}
          />
        </FormControl>

        {error && (
          <Typography color="danger" level="body-sm">
            {error}
          </Typography>
        )}

        <Button type="submit" loading={loading}>
          Submit
        </Button>
      </Stack>
    </form>
  )
}
