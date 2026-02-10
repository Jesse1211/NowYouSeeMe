import { useState } from 'react'
import Box from '@mui/joy/Box'
import Button from '@mui/joy/Button'
import Modal from '@mui/joy/Modal'
import ModalDialog from '@mui/joy/ModalDialog'
import Typography from '@mui/joy/Typography'
import AddIcon from '@mui/icons-material/Add'
import UploadVisualizationForm from './UploadVisualizationForm'

interface HeaderProps {
  onVisualizationCreated: () => void
}

export default function Header({ onVisualizationCreated }: HeaderProps) {
  const [open, setOpen] = useState(false)

  const handleSuccess = () => {
    setOpen(false)
    onVisualizationCreated()
  }

  return (
    <>
      <Box
        sx={{
          bgcolor: 'background.surface',
          borderBottom: '1px solid',
          borderColor: 'divider',
          py: 2,
          px: 3,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Typography level="h4">🪞 NowYouSeeMe</Typography>
        <Button
          startDecorator={<AddIcon />}
          onClick={() => setOpen(true)}
        >
          Share Your Visualization
        </Button>
      </Box>

      <Modal open={open} onClose={() => setOpen(false)}>
        <ModalDialog sx={{ maxWidth: 500, width: '90%' }}>
          <Typography level="h4" sx={{ mb: 2 }}>
            Upload Visualization
          </Typography>
          <UploadVisualizationForm onSuccess={handleSuccess} />
        </ModalDialog>
      </Modal>
    </>
  )
}
