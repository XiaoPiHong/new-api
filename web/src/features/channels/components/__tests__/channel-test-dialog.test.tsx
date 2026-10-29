import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createInstance } from 'i18next'
import { useEffect } from 'react'
import { I18nextProvider } from 'react-i18next'
import { afterEach, expect, it, vi } from 'vitest'

import { channelSchema, type Channel } from '../../types'
import { ChannelsProvider, useChannels } from '../channels-provider'
import { ChannelTestDialog } from '../dialogs/channel-test-dialog'

const i18n = createInstance()
await i18n.init({
  lng: 'en',
  resources: { en: { translation: {} } },
  initAsync: false,
})
const client = new QueryClient({
  defaultOptions: { queries: { retry: false } },
})

afterEach(() => {
  cleanup()
  client.clear()
})

function TestDialog({ channel }: { channel: Channel }) {
  const { setCurrentRow } = useChannels()
  useEffect(() => setCurrentRow(channel), [channel, setCurrentRow])
  return <ChannelTestDialog open onOpenChange={vi.fn()} />
}

it('turns off streaming for image editing and leaves it off when returning to chat', async () => {
  const user = userEvent.setup()
  const channel = channelSchema.parse({
    id: 1,
    type: 1,
    key: '',
    status: 1,
    name: 'image-edit-channel',
    created_time: 0,
    test_time: 0,
    response_time: 0,
    balance_updated_time: 0,
    models: 'image-model',
  })
  render(
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={client}>
        <ChannelsProvider>
          <TestDialog channel={channel} />
        </ChannelsProvider>
      </QueryClientProvider>
    </I18nextProvider>
  )

  const streamToggle = screen.getByRole('switch', { name: 'Stream Mode' })
  await user.click(streamToggle)
  expect(streamToggle).toBeChecked()

  await user.click(screen.getByRole('combobox', { name: 'Endpoint Type' }))
  await user.click(
    await screen.findByRole('option', { name: 'Image Edit (/v1/images/edits)' })
  )
  expect(streamToggle).toHaveAttribute('aria-disabled', 'true')
  expect(streamToggle).not.toBeChecked()

  await user.click(screen.getByRole('combobox', { name: 'Endpoint Type' }))
  await user.click(
    await screen.findByRole('option', { name: 'OpenAI (/v1/chat/completions)' })
  )
  expect(streamToggle).not.toHaveAttribute('aria-disabled', 'true')
  expect(streamToggle).not.toBeChecked()
})
